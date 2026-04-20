package http

import (
	"net/http"

	"purchase-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc *usecase.PurchaseUsecase
}

func NewHandler(uc *usecase.PurchaseUsecase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) CreatePurchase(c *gin.Context) {
	var req struct {
		CustomerID string `json:"customer_id"`
		ItemName   string `json:"item_name"`
		Amount     int64  `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p, err := h.uc.CreatePurchase(req.CustomerID, req.ItemName, req.Amount)
	if err != nil {
		if err.Error() == "payment service unavailable" {
			c.JSON(http.StatusServiceUnavailable, p)
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *Handler) GetPurchase(c *gin.Context) {
	id := c.Param("id")
	p, err := h.uc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *Handler) CancelPurchase(c *gin.Context) {
	id := c.Param("id")
	p, err := h.uc.CancelPurchase(id)
	if err != nil {
		if err.Error() == "paid orders cannot be cancelled" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "order not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *Handler) GetOrders(c *gin.Context) {
	customerID := c.Query("customer_id")
	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer_id is required"})
		return
	}
	orders, err := h.uc.GetByCustomerID(customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/orders", h.CreatePurchase)
	r.GET("/orders", h.GetOrders)
	r.GET("/orders/:id", h.GetPurchase)
	r.PATCH("/orders/:id/cancel", h.CancelPurchase)
}
