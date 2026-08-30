package upload

import (
	"net/http"

	usecase "spendsnap-backend/internal/usecase/upload"
	"spendsnap-backend/pkg/utils/media"

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
	fileToUpload, err := media.ExtractFileFromRequest(c, "file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.usecase.ProcessUploadFile(c.Request.Context(), fileToUpload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, FileUploadResponse{FileURL: res})
}

