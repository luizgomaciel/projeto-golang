package jobs

import (
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Loan struct {
	ID                   string    `bson:"_id,omitempty" valid:"uuid"`
	AccountNumber        string    `bson:"account_number" valid:"notnull"`
	Type                 string    `bson:"type" valid:"notnull"`
	NumberOfInstallments int       `bson:"number_installments" valid:"notnull"`
	Value                float64   `bson:"value" valid:"notnull"`
	Total                float64   `bson:"total_installments" valid:"notnull"`
	CreatedAt            time.Time `bson:"created_at" valid:"-"`
}

func NewLoan(accountNumber string) Loan {
	value := 100.10
	numberOfInstallments := 10
	total := value * float64(numberOfInstallments)

	ac := Loan{
		ID:                   uuid.NewV4().String(),
		AccountNumber:        accountNumber,
		Type:                 "LOAN_" + strconv.Itoa(1000),
		NumberOfInstallments: numberOfInstallments,
		Value:                value,
		Total:                total,
		CreatedAt:            time.Now(),
	}

	return ac
}
