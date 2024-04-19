package service

import "github.com/TEDxITS/website-backend-2024/repository"

type (
	StorageService interface {
		GetMainEventPaymentFile(string) ([]byte, error)
	}

	storageService struct {
		bucketRepo repository.BucketRepository
	}
)

func NewStorageService(bRepo repository.BucketRepository) StorageService {
	return &storageService{
		bucketRepo: bRepo,
	}
}

func (s *storageService) GetMainEventPaymentFile(id string) ([]byte, error) {
	return s.bucketRepo.DownloadFile("main-event", id)
}
