package jobs

import (
	jobs "encoder/application/jobs/accounts"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type JobLoan struct {
	ID    string `json:"job_id" valid:"uuid" bson:"_id,omitempty"`
	Loans []Loan `json:"Loans" bson:"Loans" valid:"-"`
}

func NewJobLoan(accounts *[]jobs.Account) (*JobLoan, error) {

	job := JobLoan{
		ID: uuid.NewV4().String(),
	}

	job.prepare(accounts)

	err := job.Validate()

	if err != nil {
		return nil, err
	}

	return &job, err
}

func (job *JobLoan) prepare(accounts *[]jobs.Account) {
	loans := make([]Loan, 0, 10)

	for i := 0; i < len(*accounts); i++ {
		loans = append(loans, NewLoan((*accounts)[i].AccountNumber))
	}

	job.Loans = loans
}

func (job *JobLoan) Validate() error {
	_, err := govalidator.ValidateStruct(job)

	if err != nil {
		return err
	}

	return nil
}
