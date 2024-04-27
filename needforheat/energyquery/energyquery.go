package energyquery

import (
	"time"

	"github.com/energietransitie/needforheat-server-api/needforheat/energyquerytype"
	"github.com/energietransitie/needforheat-server-api/needforheat/upload"
)

// An EnergyQuery is collects measurements in a subject's account.
type EnergyQuery struct {
	ID              uint                            `json:"id"`
	EnergyQueryType energyquerytype.EnergyQueryType `json:"energy_query_type"`
	AccountID       uint                            `json:"account_id"` // This can be removed if a device uses JWT's too.
	ActivatedAt     *time.Time                      `json:"activated_at"`
	Uploads         []upload.Upload                 `json:"uploads,omitempty"`
}

// Create a new EnergyQuery, this should be with uploads
func MakeEnergyQuery(energyQueryType energyquerytype.EnergyQueryType, accountID uint, uploads []upload.Upload) EnergyQuery {
	now := time.Now().UTC()
	return EnergyQuery{
		EnergyQueryType: energyQueryType,
		AccountID:       accountID,
		ActivatedAt:     &now,
		Uploads:         uploads,
	}
}
