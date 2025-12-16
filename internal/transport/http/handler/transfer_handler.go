package handler

import (
	"Gobank/internal/service"
	"Gobank/internal/token"
	"Gobank/internal/transport/http/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransferHandler struct {
	transferService service.TransferService
}

func NewTransferHandler(transferService service.TransferService) *TransferHandler {
	return &TransferHandler{transferService: transferService}
}
func (h *TransferHandler) CreateTransfer(c *gin.Context) {
	var req dto.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authPayload := c.MustGet("authorization_payload").(*token.TokenPayload)
	idempotencyKey := c.GetHeader("Idempotency-Key")
	if idempotencyKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "idempotency key is required"})
		return
	}
	transfer, err := h.transferService.CreateTransfer(c, &req, authPayload.UserID, idempotencyKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, transfer)
}
