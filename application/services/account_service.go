package services

import (
	"context"
	jobs "encoder/application/jobs/accounts"
	"encoder/application/producer"
	"encoder/graph/model"
	"encoder/grpc"
	"encoder/model/repository"
	"strings"
)

type AcccountService struct {
	grpc.UnimplementedAccountServiceRequestServer
	AccountRepository repository.AccountRepositoryDb
}

type AccountDTO struct {
	AccountNumber string
	AccountType   string
	CustomerName  string
	LoanProducts  []*ProductsDTO
}

type ProductsDTO struct {
	ID                 string
	LoanType           string
	NumberInstallments int
	ValueInstallments  float64
	TotalInstalments   float64
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

func (s *AcccountService) findAvailableAccounts() ([]*AccountDTO, error) {
	accountsPtr, errorAccount := s.AccountRepository.FindAllAccounts()
	if errorAccount != nil {
		return nil, errorAccount
	}

	accounts := *accountsPtr
	loanService := NewLoanService()
	loanService.LoanRepository.Db = s.AccountRepository.Db

	consolidated := make([]*AccountDTO, 0)

	for _, ac := range accounts {
		if strings.TrimSpace(ac.AccountNumber) == "" {
			continue
		}

		loan, errorLoan := loanService.FindAllLoansByAccountNumber(ac.AccountNumber)
		if errorLoan != nil {
			continue
		}

		lo := make([]*ProductsDTO, 0)
		lo = append(lo, &ProductsDTO{
			ID:                 loan.ID,
			LoanType:           loan.Type,
			NumberInstallments: loan.NumberOfInstallments,
			ValueInstallments:  loan.Value,
			TotalInstalments:   loan.Total,
		})

		data := AccountDTO{
			AccountNumber: ac.AccountNumber,
			AccountType:   ac.Type,
			CustomerName:  ac.Name,
			LoanProducts:  lo,
		}

		consolidated = append(consolidated, &data)
	}

	return consolidated, nil
}

func (s AcccountService) mustEmbedUnimplementedAccountServiceRequestServer() {
	panic("unimplemented")
}

func (s AcccountService) CreateAccounts(ctx context.Context, request *grpc.CreateRequest) (*grpc.CreateResponse, error) {
	message := producer.Message{
		TypeAccount: int(request.GetTypeAccount()),
		Qtd:         int(request.GetQuantity()),
		Products:    request.Products,
	}

	resp, err := producer.ProduceMessage(message)
	if err != nil {
		return nil, err
	}

	return &grpc.CreateResponse{
		IsStarted: resp,
	}, nil
}

func (s AcccountService) FindAccounts(ctx context.Context, em *grpc.Empty) (*grpc.FindAccountsResponse, error) {
	dto, err := s.findAvailableAccounts()
	if err != nil {
		return nil, err
	}

	consolidated := make([]*grpc.Account, 0)

	for _, ac := range dto {
		lo := make([]*grpc.Product, 0)
		lo = append(lo, &grpc.Product{
			Id:                 ac.LoanProducts[0].ID,
			LoanType:           ac.LoanProducts[0].LoanType,
			NumberInstallments: int32(ac.LoanProducts[0].NumberInstallments),
			ValueInstallments:  ac.LoanProducts[0].ValueInstallments,
			TotalInstallments:  ac.LoanProducts[0].TotalInstalments,
		})

		data := grpc.Account{
			AccountNumber: ac.AccountNumber,
			AccountType:   ac.AccountType,
			CustomerName:  ac.CustomerName,
			LoanProducts:  lo,
		}

		consolidated = append(consolidated, &data)
	}

	return &grpc.FindAccountsResponse{
		Accounts: consolidated,
	}, nil
}

func (s *AcccountService) FindAvailableAccounts() ([]*model.Account, error) {
	dto, err := s.findAvailableAccounts()
	if err != nil {
		return nil, err
	}

	consolidated := make([]*model.Account, 0)

	for _, ac := range dto {
		lo := make([]*model.Products, 0)
		lo = append(lo, &model.Products{
			ID:                 ac.LoanProducts[0].ID,
			LoanType:           ac.LoanProducts[0].LoanType,
			NumberInstallments: ac.LoanProducts[0].NumberInstallments,
			ValueInstallments:  ac.LoanProducts[0].ValueInstallments,
			TotalInstalments:   ac.LoanProducts[0].TotalInstalments,
		})

		data := model.Account{
			AccountNumber: ac.AccountNumber,
			AccountType:   ac.AccountType,
			CustomerName:  ac.CustomerName,
			LoanProducts:  lo,
		}

		consolidated = append(consolidated, &data)
	}

	return consolidated, nil
}
