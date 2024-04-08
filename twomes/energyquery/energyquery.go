package energyquery

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes/account"
	"github.com/energietransitie/twomes-backoffice-api/twomes/energyquerytype"
	"github.com/energietransitie/twomes-backoffice-api/twomes/upload"
)

// An EnergyQuery is the connection between an account and an upload
type EnergyQuery struct {
	ID              uint                            `json:"id"`
	EnergyQueryType energyquerytype.EnergyQueryType `json:"energy_query_type_id"`
	Account         account.Account                 `json:"account_id"`
	Uploads         []upload.Upload                 `json:"uploads,omitempty"`
	ActivatedAt     *time.Time                      `json:"activated_at"`
}

// Create a new EnergyQuery.
func MakeEnergyQuery(energyQueryType energyquerytype.EnergyQueryType, account account.Account, uploads []upload.Upload) EnergyQuery {
	now := time.Now().UTC()
	return EnergyQuery{
		EnergyQueryType: energyQueryType,
		Account:         account,
		Uploads:         uploads,
		ActivatedAt:     &now,
	}
}
