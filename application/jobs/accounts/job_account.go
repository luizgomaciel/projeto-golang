package jobs

import (
	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type JobAccount struct {
	ID              string    `json:"job_id" valid:"uuid" bson:"_id,omitempty"`
	TypeAccount     int       `json:"type_account" bson:"type_account" valid:"-"`
	NumberOfAccount int       `json:"number_of_account" bson:"number_of_account" valid:"-"`
	Accounts        []Account `json:"accounts" bson:"accounts" valid:"-"`
}

func NewJobAccount(typeAccount int, number int) (*JobAccount, error) {
	job := &JobAccount{
		ID:              uuid.NewV4().String(),
		TypeAccount:     typeAccount,
		NumberOfAccount: number,
	}

	job.prepare(typeAccount, number)

	err := job.Validate()

	if err != nil {
		return nil, err
	}

	return job, err
}

func (job *JobAccount) prepare(typeAccount int, number int) {
	accounts := make([]Account, 0, 10)

	for i := 0; i < number; i++ {
		accounts = append(accounts, *NewAccount(typeAccount))
	}

	job.Accounts = accounts
}

func (job *JobAccount) Validate() error {
	_, err := govalidator.ValidateStruct(job)

	if err != nil {
		return err
	}

	return nil
}
