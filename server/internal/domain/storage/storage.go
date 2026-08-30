package storage

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Hạn mức dung lượng tối đa (Ví dụ: 5MB)
const MaxFileSize = 5 * 1024 * 1024

// Danh sách các MIME Type cho phép upload
var AllowedMimeTypes = []string{
	"image/jpeg",
	"image/png",
	"image/webp",
	"image/heic",
}

type File struct {
	Key         string
	ContentType string // Định dạng do client truyền lên hoặc tự động nhận diện
	Data        []byte // Dữ liệu nhị phân (viết hoa để Export)
}

type StorageProvider interface {
	Upload(ctx context.Context, file *File) (string, error)
	Delete(ctx context.Context, fileURL string) error
}

type UploadUsecase interface {
    ProcessUploadFile(ctx context.Context, file *File) (string, error)
    DeleteFile(ctx context.Context, fileURL string) error
}
// Validate kiểm tra Key, Dung lượng và Định dạng file
func (f *File) Validate() error {
	// 1. Kiểm tra Key
	if strings.TrimSpace(f.Key) == "" {
		return errors.New("key không được để trống")
	}

	// 2. Kiểm tra Dung lượng File (Size)
	if len(f.Data) == 0 {
		return errors.New("file không được để trống (size = 0)")
	}

	if len(f.Data) > MaxFileSize {
		return fmt.Errorf("dung lượng file vượt quá giới hạn %dMB", MaxFileSize/(1024*1024))
	}

	// 3. Magic Bytes Check - Kiểm tra định dạng thực tế của dữ liệu nhị phân
	detectedType := http.DetectContentType(f.Data)

	if (detectedType == "application/octet-stream" || isAllowedType(f.ContentType)) && IsHEIC(f.Data) {
		detectedType = "image/heic"
	}

	// 4. Kiểm tra định dạng thực tế nhận diện được có nằm trong danh sách cho phép không
	if !isAllowedType(detectedType) {
		return fmt.Errorf("dữ liệu file thực tế (%s) không đúng định dạng ảnh cho phép", detectedType)
	}

	// 5. Cập nhật lại ContentType chuẩn cho struct File
	f.ContentType = detectedType

	return nil
}

// Helper kiểm tra Content-Type có nằm trong danh sách cho phép không
func isAllowedType(contentType string) bool {
	for _, allowed := range AllowedMimeTypes {
		if strings.HasPrefix(contentType, allowed) {
			return true
		}
	}
	return false
}

// Helper kiểm tra Magic Bytes thủ công cho file HEIC/HEIF từ iPhone
func IsHEIC(data []byte) bool {
	if len(data) < 12 {
		return false
	}
	// File HEIC/HEIF luôn chứa chuỗi "ftyp" ở vị trí byte từ 4 đến 8
	if string(data[4:8]) == "ftyp" {
		brand := string(data[8:12])
		// Các Major Brand phổ biến của chuẩn HEIC/HEIF trên iOS
		switch brand {
		case "heic", "heix", "hevc", "mif1", "msf1":
			return true
		}
	}
	return false
}