package services

import (
	jobs "encoder/application/jobs/accounts"
	"encoder/graph/model"
	"encoder/model/repository"
	"strings"
)

type AcccountService struct {
	AccountRepository repository.AccountRepositoryDb
}

func NewAccountService() AcccountService {
	return AcccountService{}
}

func (s *AcccountService) InsertAccount(typeAccount int, number int) (*[]jobs.Account, error) {
	accounts, err := jobs.NewJobAccount(typeAccount, number)
	if err != nil {
		return nil, err
	}

	for _, ac := range accounts.Accounts {
		if _, err := s.AccountRepository.Insert(&ac); err != nil {
			return nil, err
		}
	}

	return &accounts.Accounts, nil
}

func (s *AcccountService) FindAvailableAccounts() ([]*model.Account, error) {
	accountsPtr, errorAccount := s.AccountRepository.FindAllAccounts()
	if errorAccount != nil {
		return nil, errorAccount
	}

	accounts := *accountsPtr
	loanService := NewLoanService()
	loanService.LoanRepository.Db = s.AccountRepository.Db

	consolidated := make([]*model.Account, 0)

	for _, ac := range accounts {
		if strings.TrimSpace(ac.AccountNumber) == "" {
			continue
		}

		loan, errorLoan := loanService.FindAllLoansByAccountNumber(ac.AccountNumber)
		if errorLoan != nil {
			continue
		}

		lo := make([]*model.Products, 0)
		lo = append(lo, &model.Products{
			ID:                 loan.ID,
			LoanType:           loan.Type,
			NumberInstallments: loan.NumberOfInstallments,
			ValueInstallments:  loan.Value,
			TotalInstalments:   loan.Total,
		})

		data := model.Account{
			AccountNumber: ac.AccountNumber,
			AccountType:   ac.Type,
			CustomerName:  ac.Name,
			LoanProducts:  lo,
		}

		consolidated = append(consolidated, &data)
	}

	return consolidated, nil
}
