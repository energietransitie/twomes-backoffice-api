package enelogic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes"
)

// Client is the http client used to make requests to enelogic.
var client = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{

		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     10,
		IdleConnTimeout:     time.Second * 30,
	},
}

const (
	baseURL                    = "https://enelogic.com/api"
	endpointMeasuringPoints    = "/measuringpoints"
	endpointDatapointsMonths   = "/measuringpoints/{{.MeasuringPointID}}/datapoint/months/{{.From}}/{{.To}}"
	endpointDatapointsDays     = "/measuringpoints/{{.MeasuringPointID}}/datapoint/days/{{.From}}/{{.To}}"
	endpointDatapointsInterval = "/measuringpoints/{{.MeasuringPointID}}/datapoints/{{.From}}/{{.To}}"

	enelogicTimeFormat          = `"2006-01-02 15:04:05"`
	enelogicDefaultTimeLocation = "Europe/Amsterdam"
)

var (
	ErrNoData = errors.New("no data from enelogic")
)

// EnelogicTime is a custom time type for enelogic.
// It is used to parse the time format used by enelogic.
type EnelogicTime struct {
	time.Time
	LocationName *string
}

func (t *EnelogicTime) UnmarshalJSON(b []byte) error {
	if t.LocationName == nil {
		location := enelogicDefaultTimeLocation
		t.LocationName = &location
	}

	loc, err := time.LoadLocation(*t.LocationName)
	if err != nil {
		return err
	}

	parsed, err := time.ParseInLocation(enelogicTimeFormat, string(b), loc)
	if err != nil {
		return err
	}

	t.Time = parsed
	return nil
}

func (t EnelogicTime) MarshalJSON() ([]byte, error) {
	enelogicFormattedTime := t.Time.Format(enelogicTimeFormat)
	return []byte(enelogicFormattedTime), nil
}

func (t EnelogicTime) String() string {
	return t.Time.Format(time.RFC3339)
}

// UnitType is the type of the unit.
// 0 = electricity.
// 1 = gas.
type UnitType int

const (
	unitTypeElectricity UnitType = 0
	unitTypeGas         UnitType = 1
)

// Rate is the type of the Rate, as defined by enelogic.
type Rate int

const (
	g_use_cum__m3     Rate = 180
	e_use_lo_cum__kWh Rate = 181
	e_use_hi_cum__kWh Rate = 182
	e_ret_lo_cum__kWh Rate = 281
	e_ret_hi_cum__kWh Rate = 282
)

func (r Rate) String() string {
	propertyNames := map[Rate]string{
		g_use_cum__m3:     "g_use_cum__m3",
		e_use_lo_cum__kWh: "e_use_lo_cum__kWh",
		e_use_hi_cum__kWh: "e_use_hi_cum__kWh",
		e_ret_lo_cum__kWh: "e_ret_lo_cum__kWh",
		e_ret_hi_cum__kWh: "e_ret_hi_cum__kWh",
	}

	return propertyNames[r]
}

type MeasuringsPointsResponse []MeasuringPoint

type MeasuringPoint struct {
	Timezone string       `json:"timezone"`
	ID       int          `json:"id"`
	UnitType UnitType     `json:"unitType"`
	DayMin   EnelogicTime `json:"dayMin"`
	DayMax   EnelogicTime `json:"dayMax"`
	MonthMin EnelogicTime `json:"monthMin"`
	MonthMax EnelogicTime `json:"monthMax"`
	YearMin  EnelogicTime `json:"yearMin"`
	YearMax  EnelogicTime `json:"yearMax"`
	Active   bool         `json:"active"`
}

func (m *MeasuringPoint) UnmarshalJSON(b []byte) error {
	// Set all EnelogicTime fields' LocationName to the timezone of the measuring point using a pointer.
	m.DayMin.LocationName = &m.Timezone
	m.DayMax.LocationName = &m.Timezone
	m.MonthMin.LocationName = &m.Timezone
	m.MonthMax.LocationName = &m.Timezone
	m.YearMin.LocationName = &m.Timezone
	m.YearMax.LocationName = &m.Timezone

	type localMeasuringPoint *MeasuringPoint
	return json.Unmarshal(b, localMeasuringPoint(m))
}

type DatapointsResponse []DataPoint

type DataPoint struct {
	Quantity float64      `json:"quantity"`
	Rate     Rate         `json:"rate"`
	Date     EnelogicTime `json:"date"`
}

type RequestTime struct {
	time.Time
}

func (t RequestTime) String() string {
	return t.Format("2006-01-02")
}

type RequestArgs struct {
	MeasuringPointID int
	From             RequestTime
	To               RequestTime
}

func newRequestArgs(measuringPointID int, from, to time.Time) RequestArgs {
	return RequestArgs{
		MeasuringPointID: measuringPointID,
		From:             RequestTime{from},
		To:               RequestTime{to},
	}
}

type Timespan struct {
	From time.Time
	To   time.Time
}

type DownloadArgs struct {
	RequestMonthsDatapoints   Timespan
	RequestDaysDatapoints     Timespan
	RequestIntervalDatapoints Timespan
}

// Download downloads the data from enelogic.
// A slice of measurements is returned, which can be saved to the database.
func Download(ctx context.Context, token string, args DownloadArgs) ([]twomes.Measurement, error) {
	var measurements []twomes.Measurement

	measuringPoints, err := getMeasuringPoints(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("error getting measuring points: %w", err)
	}

	if len(measuringPoints) == 0 {
		return nil, ErrNoData
	}

	for _, measuringPoint := range measuringPoints {
		args := newRequestArgs(measuringPoint.ID, args.RequestMonthsDatapoints.From, args.RequestMonthsDatapoints.To)
		datapoints, err := getDatapointsMonths(ctx, token, args)
		if err != nil {
			return nil, fmt.Errorf("error getting datapoints: %w", err)
		}

		measurements = append(measurements, parseDatapoints(datapoints)...)
	}

	return measurements, nil
}

func parseDatapoints(datapoints DatapointsResponse) []twomes.Measurement {
	var measurements []twomes.Measurement

	for _, datapoint := range datapoints {
		measurement := twomes.Measurement{
			Property: twomes.Property{
				Name: datapoint.Rate.String(),
			},
			Time:  twomes.Time(datapoint.Date.Time),
			Value: fmt.Sprintf("%.3f", datapoint.Quantity),
		}

		measurements = append(measurements, measurement)
	}

	return measurements
}

// GetMeasuringPoints returns the measuring points for the account with the given token.
func getMeasuringPoints(ctx context.Context, token string) (MeasuringsPointsResponse, error) {
	requestURL := baseURL + endpointMeasuringPoints

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request to Enelogic: %w", err)
	}

	var response MeasuringsPointsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error decoding response to json: %w", err)
	}
	defer resp.Body.Close()

	return response, nil
}

// GetDatapointsMonths returns the data points for the account with the given token.
func getDatapointsMonths(ctx context.Context, token string, args RequestArgs) (DatapointsResponse, error) {
	requestUrl, err := getRequestURL(endpointDatapointsMonths, args)
	if err != nil {
		return nil, fmt.Errorf("error getting request url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request to Enelogic: %w", err)
	}

	var response DatapointsResponse
	json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error decoding response to json: %w", err)
	}

	return response, nil
}

// GetDatapointsDays returns the data points for the account with the given token.
func getDatapointsDays(ctx context.Context, token string, args RequestArgs) (DatapointsResponse, error) {
	return nil, nil
}

// GetDatapointsInterval returns the data points for the account with the given token.
func getDatapointsInterval(ctx context.Context, token string, args RequestArgs) (DatapointsResponse, error) {
	return nil, nil
}

func getRequestURL(endpoint string, args RequestArgs) (string, error) {
	requestURL := strings.Builder{}

	t, err := template.New("url").Parse(baseURL + endpoint)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	err = t.Execute(&requestURL, args)
	if err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return requestURL.String(), nil
}
