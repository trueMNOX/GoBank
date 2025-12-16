package handler

import (
	"Gobank/internal/service"
	"Gobank/internal/token"
	"Gobank/internal/transport/http/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accountService service.AccountService
}

func NewAccountHandler(accountService service.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req dto.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authPayload := c.MustGet("authorization_payload").(*token.TokenPayload)
	account, err := h.accountService.CreateAccount(c, &req, authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, account)
}
func (h *AccountHandler) GetAccount(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}
	authPayload := c.MustGet("authorization_payload").(*token.TokenPayload)
	accoount, err := h.accountService.GetAccountByID(c, id, authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, accoount)
}
func (h *AccountHandler) ListAccount(c *gin.Context) {
	var req dto.ListAccountsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authPayload := c.MustGet("authorization_payload").(*token.TokenPayload)
	accounts, err := h.accountService.ListAccounts(c, authPayload.UserID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, accounts)
}
