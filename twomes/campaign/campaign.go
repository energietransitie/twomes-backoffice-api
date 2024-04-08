package campaign

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes/app"
	"github.com/energietransitie/twomes-backoffice-api/twomes/datasourcelist"
)

// A campaign is a timeframe where we gather measurements with a specific goal.
type Campaign struct {
	ID             uint                          `json:"id"`
	Name           string                        `json:"name"`
	App            app.App                       `json:"app"`
	InfoURL        string                        `json:"info_url"`
	StartTime      *time.Time                    `json:"start_time,omitempty"`
	EndTime        *time.Time                    `json:"end_time,omitempty"`
	DataSourceList datasourcelist.DataSourceList `json:"data_sources_list"`
}

// Create a new Campaign.
func MakeCampaign(name string, app app.App, infoURL string, startTime, endTime *time.Time, dataSourceList datasourcelist.DataSourceList) Campaign {
	return Campaign{
		Name:           name,
		App:            app,
		InfoURL:        infoURL,
		StartTime:      startTime,
		EndTime:        endTime,
		DataSourceList: dataSourceList,
	}
}
