package media

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"spendsnap-backend/internal/domain/storage"

	"github.com/chai2010/webp"
	"github.com/jdeng/goheif"   // Decoder cho file HEIC/HEIF từ iPhone
	_ "golang.org/x/image/webp" // Đăng ký decoder khi file đầu vào là WebP
)

// ToWebP nén mảng bytes ảnh gốc (JPG, PNG, WEBP, HEIC) sang WebP
func ToWebP(input []byte, quality float32) ([]byte, string, error) {
	var img image.Image

	var err error

	// 1. Kiểm tra nếu là file HEIC/HEIF từ iPhone -> Dùng goheif để decode
	if storage.IsHEIC(input) {
		img, err = goheif.Decode(bytes.NewReader(input))
		if err != nil {
			return nil, "", fmt.Errorf("lỗi decode file HEIC từ iPhone: %w", err)
		}
	} else {
		// 2. Với JPG, PNG, WEBP -> Dùng image.Decode chuẩn của Go
		img, _, err = image.Decode(bytes.NewReader(input))
		if err != nil {
			return nil, "", fmt.Errorf("lỗi decode ảnh: %w", err)
		}
	}
	// Cắt ảnh về tỷ lệ 1:1 chính giữa (Crop Square) trước khi encode sang WebP
	img = CropSquare(img)
	// 3. Encode image.Image thu được sang WebP
	var buf bytes.Buffer
	err = webp.Encode(&buf, img, &webp.Options{Quality: quality})
	if err != nil {
		return nil, "", fmt.Errorf("lỗi encode webp: %w", err)
	}

	return buf.Bytes(), "image/webp", nil
}

