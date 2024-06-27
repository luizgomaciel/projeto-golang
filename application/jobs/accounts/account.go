package jobs

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Account struct {
	ID            string    `json:"job_id" bson:"_id,omitempty" valid:"uuid"`
	AccountNumber string    `json:"account_number" bson:"account_number" valid:"notnull"`
	Type          string    `json:"type" bson:"type" valid:"notnull"`
	Name          string    `json:"name" bson:"name" valid:"notnull"`
	Active        bool      `json:"active" bson:"active" valid:"-"`
	Blocked       bool      `json:"blocked" bson:"blocked" valid:"-"`
	Origin        string    `json:"origin" bson:"origin" valid:"notnull"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at" valid:"-"`
}

func NewAccount(typeNumber int) *Account {
	ac := Account{
		ID:            uuid.NewV4().String(),
		AccountNumber: strconv.FormatInt(rand.Int63n(10000000000), 10),
		Type:          "CONTA_CORRENTE_" + strconv.Itoa(typeNumber),
		Name:          "Luiz Gustavo",
		Active:        true,
		Blocked:       false,
		Origin:        "site",
		CreatedAt:     time.Now(),
	}

	return &ac
}

func (account *Account) Validate() error {
	_, err := govalidator.ValidateStruct(account)

	if err != nil {
		return err
	}

	return nil
}
