package http

import (
	"net/http"

	"transaction-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase *usecase.TransactionUsecase
}

func NewHandler(uc *usecase.TransactionUsecase) *Handler {
	return &Handler{usecase: uc}
}

type createRequest struct {
	PurchaseID string `json:"purchase_id"`
	Amount     int64  `json:"amount"`
}

func (h *Handler) CreateTransaction(c *gin.Context) {
	var req createRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.usecase.CreateTransaction(req.PurchaseID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tx)
}

func (h *Handler) GetTransaction(c *gin.Context) {
	purchaseID := c.Param("purchase_id")

	tx, err := h.usecase.GetByPurchaseID(purchaseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, tx)
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/transactions", h.CreateTransaction)
	r.GET("/transactions/:purchase_id", h.GetTransaction)
}
