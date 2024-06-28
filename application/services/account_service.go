package services

import (
	jobs "encoder/application/jobs/accounts"
	"encoder/model/repository"
)

type AcccountService struct {
	AccountRepository repository.AccountRepositoryDb
}

func NewAccountService() AcccountService {
	return AcccountService{}
}

func (s *AcccountService) InsertAccount(typeAccount int, number int) error {
	accounts, err := jobs.NewJobAccount(typeAccount, number)
	if err != nil {
		return err
	}

	for _, ac := range accounts.Accounts {
		if _, err := s.AccountRepository.Insert(&ac); err != nil {
			return err
		}
	}

	return nil
}
