package upload

import (
	"context"
	"errors"
	"spendsnap-backend/internal/domain/storage"
	"spendsnap-backend/pkg/utils"
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

func (uc *CreateUploadUsecase) ProcessUploadFile(ctx context.Context, file *storage.File) (string, error) {
	if uc == nil || uc.storage == nil {
		return "", errors.New("storage provider not initialized")
	}
	if file == nil {
		return "", errors.New("file không được để trống")
	}
	err := file.Validate()
	if err != nil {
		return "", err
	}

	fileConverted, newContentType, err := media.ToWebP(file.Data, 80)
	if err != nil {
		return "", err
	}

	newImageID := utils.NewID()
	newKey := "images/" + newImageID + ".webp"
	fileToUpload := &storage.File{
		Key:         newKey,
		ContentType: newContentType,
		Data:        fileConverted,
	}
	return uc.storage.Upload(ctx,fileToUpload)
}

func (uc *CreateUploadUsecase) DeleteFile(ctx context.Context, fileURL string) error {
	if uc == nil || uc.storage == nil {
		return errors.New("storage provider not initialized")
	}
	if fileURL == "" {
		return nil
	}
	return uc.storage.Delete(ctx, fileURL)
}