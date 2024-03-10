package services

type FineService interface {
}

type fineService struct {
}

func NewFineService() FineService {
	return &fineService{}
}
