package media

import (
	"fmt"
	"io"
	"spendsnap-backend/internal/domain/storage"

	"github.com/gin-gonic/gin"
)


func ExtractFileFromRequest(c *gin.Context, formKey string) (*storage.File, error) {
    // 1. Lấy file từ request
    fileHeader, err := c.FormFile(formKey)
    if err != nil {
        return nil, fmt.Errorf("không tìm thấy file trong form với key '%s'", formKey)
    }

    // 2. Mở file
    file, err := fileHeader.Open()
    if err != nil {
        return nil, fmt.Errorf("không thể đọc nội dung file: %w", err)
    }
    defer file.Close()

    // 3. Đọc ra mảng byte
    fileBytes, err := io.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("lỗi trong quá trình đọc dữ liệu file: %w", err)
    }

    // 4. Lấy content-type và tên
    contentType := fileHeader.Header.Get("Content-Type")
    
    // (Tên file có thể nhận từ param hoặc lấy tên gốc)
    filename := c.Param("filename") 
    if filename == "" {
        filename = fileHeader.Filename
    }

    return &storage.File{
        Key:         filename,
        ContentType: contentType,
        Data:        fileBytes,
    }, nil
}