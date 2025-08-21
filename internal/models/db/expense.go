package dbmodels

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
)

// TODO comments
type ExpenseDBResponse struct {
	InstanceID string    `db:"instance_id"`
	Amount     float64   `db:"amount"`
	Date       time.Time `db:"date"`
}

// TODO comments
func (e ExpenseDBResponse) ToExpenseDTOResponse() *dto.ExpenseDTOResponse {
	return &dto.ExpenseDTOResponse{
		InstanceID: e.InstanceID,
		Amount:     e.Amount,
		Date:       e.Date,
	}
}
