package upload

import (
	"spendsnap-backend/internal/domain/storage"
	"spendsnap-backend/internal/usecase/upload"
)

type Dependencies struct {
	Usecase *upload.CreateUploadUsecase
	Handler *UploadHandler
}

func NewDependencies(storageProvider storage.StorageProvider) *Dependencies {
	uc := upload.NewCreateUploadUsecase(storageProvider)
	handler := NewUploadHandler(uc)

	return &Dependencies{
		Usecase: uc,
		Handler: handler,
	}
}
