package services

import (
	jobs "encoder/application/jobs/accounts"
	jobs_loans "encoder/application/jobs/loan"
	"encoder/model/repository"
)

type LoanService struct {
	LoanRepository repository.LoanRepositoryDb
}

func NewLoanService() LoanService {
	return LoanService{}
}

func (s *LoanService) InsertLoan(accounts *[]jobs.Account) (*[]jobs.Account, error) {
	loans, err := jobs_loans.NewJobLoan(accounts)
	if err != nil {
		return nil, err
	}

	for _, loan := range loans.Loans {
		if _, err := s.LoanRepository.Insert(loan); err != nil {
			return nil, err
		}
	}

	return accounts, nil
}

func (s *LoanService) FindAllLoansByAccountNumber(accountNumber string) (*jobs_loans.Loan, error) {
	return s.LoanRepository.FindAllByAccountNumber(accountNumber)
}
