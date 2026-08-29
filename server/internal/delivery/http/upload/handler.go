package upload

import (
	"io"
	"net/http"

	"spendsnap-backend/internal/domain/storage"
	usecase "spendsnap-backend/internal/usecase/upload"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	usecase *usecase.CreateUploadUsecase
}

type FileUploadResponse struct {
	FileURL string `json:"file_url"`
}

func NewUploadHandler(usecase *usecase.CreateUploadUsecase) *UploadHandler {
	return &UploadHandler{usecase: usecase}
}

func (h *UploadHandler) UploadFile(c *gin.Context) {
	// 1. Nhận file nhị phân trực tiếp từ Form với key là "file"
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vui lòng chọn file để upload"})
		return
	}

	// 2. Lấy Content-Type tự động từ File (image/png, image/jpeg...)
	contentType := fileHeader.Header.Get("Content-Type")

	// 3. Mở luồng đọc file
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể mở file"})
		return
	}
	defer file.Close()

	// 4. Đọc dữ liệu file thành []byte
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể đọc dữ liệu file"})
		return
	}

	// 5. Tên file lấy trực tiếp từ file gốc hoặc Param URL
	filename := "uploads/" + c.Param("filename")
	if filename == "" {
		filename = fileHeader.Filename
	}

	fileToUpload := &storage.File{
		Key:         filename,
		ContentType: contentType,
		Data:        fileBytes,
	}


	// 6. Upload lên R2
	res, err := h.usecase.Upload(c.Request.Context(), fileToUpload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, FileUploadResponse{FileURL: res})
}

func RegisterUploadRoutes(router *gin.Engine, uploadHandler *UploadHandler) {
	uploadGroup := router.Group("/api/v1/upload")
	{
		uploadGroup.POST("/:filename", uploadHandler.UploadFile)
	}
}