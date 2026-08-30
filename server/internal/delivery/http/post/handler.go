package post

import (
	"net/http"
	postDomain "spendsnap-backend/internal/domain/post"
	"spendsnap-backend/internal/usecase/post"
	"spendsnap-backend/pkg/utils/media"

	transactionDomain "spendsnap-backend/internal/domain/transaction"

	"github.com/gin-gonic/gin"
)

type CreatePostHandler struct{
	usecase *post.CreatePostUsecase
}

func NewCreatePostHandler(usecase *post.CreatePostUsecase) *CreatePostHandler {
	return &CreatePostHandler{
		usecase: usecase,
	}
}

func (h *CreatePostHandler) CreatePost(c *gin.Context){
	var req postDomain.CreatePostRequest
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	uid := userID.(string)

	if err := c.ShouldBind(&req); err != nil{
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return  
    }
	fileToUpload, err := media.ExtractFileFromRequest(c, "file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Lỗi khi trích xuất file": err.Error()})
		return
	}

	post := &postDomain.Post{
		UserID: uid,
		Caption: req.Caption,
		Visibility: req.Visibility,
		LocationName: req.LocationName,
	}

	var transaction *transactionDomain.Transaction
	// Chỉ tạo transaction nếu client gửi amount + category_id (cho phép post không kèm chi tiêu)
	if req.Amount != nil && req.CategoryID != nil && *req.CategoryID != "" {
		if *req.Amount <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "amount phải lớn hơn 0"})
			return
		}
		transaction = &transactionDomain.Transaction{
			UserID: uid,
			Amount: *req.Amount,
			CategoryID: *req.CategoryID,
			Note: req.Note,
		}
		if req.Type != nil && *req.Type != "" {
			transaction.Type = *req.Type
		}
	}

	postCreated, err := h.usecase.CreatePost(c.Request.Context(), post, fileToUpload, transaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, postCreated)
}

