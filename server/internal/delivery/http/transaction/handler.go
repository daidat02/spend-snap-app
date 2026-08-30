package transaction

import (
	"net/http"
	"time"

	domain "spendsnap-backend/internal/domain/transaction"
	usecase "spendsnap-backend/internal/usecase/transaction"

	"github.com/gin-gonic/gin"
)

type CreateTransactionHandler struct {
	usecase *usecase.CreateTransactionUsecase
}

func NewCreateTransactionHandler(usecase *usecase.CreateTransactionUsecase) *CreateTransactionHandler {
	return &CreateTransactionHandler{
		usecase: usecase,
	}
}

func (h *CreateTransactionHandler) CreateTransaction(c *gin.Context) {
	var req domain.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	uid := userID.(string)

	var txDate time.Time
	if req.TransactionDate != nil && *req.TransactionDate != "" {
		parsed, err := time.Parse(time.RFC3339, *req.TransactionDate)
		if err != nil {
			parsed2, err2 := time.Parse("2006-01-02", *req.TransactionDate)
			if err2 != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "transaction_date phải là RFC3339 hoặc YYYY-MM-DD"})
				return
			}
			txDate = parsed2
		} else {
			txDate = parsed
		}
	}

	transaction := &domain.Transaction{
		UserID:          uid,
		CategoryID:      req.CategoryID,
		Amount:          req.Amount,
		Type:            req.Type,
		Note:            req.Note,
		TransactionDate: txDate,
	}

	resp, err := h.usecase.CreateTransaction(c.Request.Context(), transaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}
