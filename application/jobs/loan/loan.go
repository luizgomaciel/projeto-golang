package jobs

import (
	"encoder/application/utils"
	rand_math "math/rand"
	"math/rand/v2"
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
	min := 5.0
	max := 1000.0
	value := min + rand.Float64()*(max-min)

	numberOfInstallments := rand_math.Intn(24)
	total := value * float64(numberOfInstallments)

	ac := Loan{
		ID:                   uuid.NewV4().String(),
		AccountNumber:        accountNumber,
		Type:                 "LOAN_" + strconv.Itoa(1000),
		NumberOfInstallments: numberOfInstallments,
		Value:                utils.RoundToTwoDecimalPlaces(value),
		Total:                utils.RoundToTwoDecimalPlaces(total),
		CreatedAt:            time.Now(),
	}

	return ac
}
