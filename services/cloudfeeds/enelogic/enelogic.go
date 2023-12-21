package enelogic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes"
	"github.com/sirupsen/logrus"
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

	Day = time.Hour * 24
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
	UnitTypeElectricity UnitType = 0
	UnitTypeGas         UnitType = 1
)

func (u UnitType) String() string {
	switch u {
	case UnitTypeElectricity:
		return "electricity"
	case UnitTypeGas:
		return "gas"
	}

	return "unknown"
}

// Rate is the type of the Rate, as defined by enelogic.
type Rate int

// Rate constants as defined by enelogic.
// See https://enelogic.com/nl/blog/slimme-meter-data-exporteren for more information.
const (
	RateUsageTotal  Rate = 180
	RateUsageLow    Rate = 181
	RateUsageHigh   Rate = 182
	RateReturnTotal Rate = 280
	RateReturnLow   Rate = 281
	RateReturnHigh  Rate = 282
)

func (r Rate) Parse(unit UnitType) string {
	propertyNames := map[Rate]string{
		RateUsageTotal:  "use_cum",
		RateUsageLow:    "use_lo_cum",
		RateUsageHigh:   "use_hi_cum",
		RateReturnTotal: "ret_cum",
		RateReturnLow:   "ret_lo_cum",
		RateReturnHigh:  "ret_hi_cum",
	}

	switch unit {
	case UnitTypeElectricity:
		return "e_" + propertyNames[r] + "__kWh"
	case UnitTypeGas:
		return "g_" + propertyNames[r] + "__m3"
	}

	return propertyNames[r]
}

type Quantity float64

func (q *Quantity) UnmarshalJSON(b []byte) error {
	var quantity float64
	err := json.Unmarshal(b, &quantity)
	if err != nil {
		// Try to unmarshal as string.
		var quantity string
		err := json.Unmarshal(b, &quantity)
		if err != nil {
			return err
		}

		q64, err := strconv.ParseFloat(quantity, 64)
		if err != nil {
			return err
		}
		*q = Quantity(q64)
		return nil
	}

	*q = Quantity(quantity)
	return nil
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
	Quantity Quantity     `json:"quantity"`
	Rate     Rate         `json:"rate"`
	Date     EnelogicTime `json:"date"`
	Datetime EnelogicTime `json:"datetime"`
}

func (d DataPoint) Parse(unit UnitType) twomes.Measurement {
	t := twomes.Time(d.Date.Time)
	if time.Time(t).IsZero() {
		t = twomes.Time(d.Datetime.Time)
	}

	return twomes.Measurement{
		Property: twomes.Property{
			Name: d.Rate.Parse(unit),
		},
		Time:  t,
		Value: fmt.Sprintf("%.3f", d.Quantity),
	}
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
//
// StartPeriod is the start of the period from which data should be downloaded.
func Download(ctx context.Context, token string, startPeriod time.Time) ([]twomes.Measurement, error) {
	var measurements []twomes.Measurement

	if dateEqual(startPeriod, time.Now()) {
		return nil, ErrNoData
	}

	measuringPoints, err := getMeasuringPoints(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("error getting measuring points: %w", err)
	}

	if len(measuringPoints) == 0 {
		return nil, ErrNoData
	}

	for _, measuringPoint := range measuringPoints {

		// Get month datapoints.
		logrus.Infoln("downloading", measuringPoint.UnitType.String(), "month datapoints from", RequestTime{startPeriod}, "to", RequestTime{time.Now()})

		args := newRequestArgs(measuringPoint.ID, startPeriod, time.Now())
		datapoints, err := getDatapoints(ctx, token, endpointDatapointsMonths, args)
		if err != nil {
			return nil, fmt.Errorf("error getting month datapoints: %w", err)
		}
		measurements = append(measurements, parseDatapoints(datapoints, measuringPoint.UnitType)...)

		// Get day datapoints.
		{
			// Shadow startPeriod to avoid changing the original value. This is needed for the next iteration.
			var startPeriod time.Time
			if time.Since(startPeriod) > Day*40 {
				// Set startPeriod to 40 days ago, if the real startPeriod is more than 40 days ago.
				startPeriod = time.Now().Add(-Day * 40)
			}

			logrus.Infoln("downloading", measuringPoint.UnitType.String(), "day datapoints from", RequestTime{startPeriod}, "to", RequestTime{time.Now()})

			args = newRequestArgs(measuringPoint.ID, startPeriod, time.Now())
			datapoints, err = getDatapoints(ctx, token, endpointDatapointsDays, args)
			if err != nil {
				return nil, fmt.Errorf("error getting day datapoints: %w", err)
			}
			measurements = append(measurements, parseDatapoints(datapoints, measuringPoint.UnitType)...)
		}

		// Get interval datapoints.
		{
			// Shadow startPeriod to avoid changing the original value. This is needed for the next iteration.
			var startPeriod time.Time
			if time.Since(startPeriod) > Day*10 {
				// Set startPeriod to 10 days ago, if the real startPeriod is more than 10 days ago.
				startPeriod = time.Now().Add(-Day * 10)
			}

			logrus.Infoln("downloading", measuringPoint.UnitType.String(), "interval datapoints from", RequestTime{startPeriod}, "to", RequestTime{time.Now()})

			for _, day := range splitDays(startPeriod, time.Now()) {
				args = newRequestArgs(measuringPoint.ID, day.Start, day.End)
				datapoints, err = getDatapoints(ctx, token, endpointDatapointsInterval, args)
				if err != nil {
					return nil, fmt.Errorf("error getting interval datapoints: %w", err)
				}
				measurements = append(measurements, parseDatapoints(datapoints, measuringPoint.UnitType)...)
			}
		}
	}

	return measurements, nil
}

func parseDatapoints(datapoints DatapointsResponse, unit UnitType) []twomes.Measurement {
	var measurements []twomes.Measurement

	for _, datapoint := range datapoints {
		measurements = append(measurements, datapoint.Parse(unit))
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

// GetDatapoints returns the data points for the account with the given token.
func getDatapoints(ctx context.Context, token string, endpoint string, args RequestArgs) (DatapointsResponse, error) {
	requestUrl, err := getRequestURL(endpoint, args)
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

func dateEqual(a, b time.Time) bool {
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

type day struct {
	Start time.Time
	End   time.Time
}

func splitDays(from, to time.Time) []day {
	var days []day

	for from.Before(to) && !dateEqual(from, to) {
		to := from.AddDate(0, 0, 1)
		days = append(days, day{
			Start: from,
			End:   to,
		})
		from = to
	}

	return days
}
