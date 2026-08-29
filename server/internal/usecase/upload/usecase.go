package upload

import (
	"context"
	"spendsnap-backend/internal/domain/storage"
	"spendsnap-backend/pkg/utils/media"
)


type CreateUploadUsecase struct {
	storage storage.StorageProvider
}

func NewCreateUploadUsecase(storage storage.StorageProvider) *CreateUploadUsecase {
	return &CreateUploadUsecase{
		storage: storage,
	}
}

func (uc *CreateUploadUsecase) Upload(ctx context.Context, file *storage.File) (string, error) {
	err := file.Validate()
	if err != nil {
		return "", err
	}

	fileConverted, newContentType, err := media.ToWebP(file.Data, 80)
	if err != nil {
		return "", err
	}


	fileToUpload := &storage.File{
		Key:         file.Key,
		ContentType: newContentType,
		Data:        fileConverted,
	}
	return uc.storage.Upload(ctx, fileToUpload.Key, fileToUpload.Data, fileToUpload.ContentType)
}