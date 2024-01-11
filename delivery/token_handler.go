package delivery

import (
	"net/http"
	"project/models/request"
	"project/service"
	"project/utils"
	"project/validator"

	"github.com/gin-gonic/gin"
)

type TokenHandler struct {
	tokenSvc service.TokenService
}

func NewTokenHandler(tokenSvc service.TokenService) *TokenHandler {
	return &TokenHandler{tokenSvc}
}

func ConfigureTokenHandler(router *gin.RouterGroup, tokenSvc service.TokenService) {
	TokenHandler := NewTokenHandler(tokenSvc)
	addTokenHandlerRoutes(router, TokenHandler)
}

func addTokenHandlerRoutes(router *gin.RouterGroup, TokenHandler *TokenHandler) {
	router.POST("/renew", TokenHandler.RenewAccessToken)
}

func (h *TokenHandler) RenewAccessToken(c *gin.Context) {
	ctx := c.Request.Context()
	var renewAccessTokenRequest request.RenewAccessTokenRequest
	if err := c.ShouldBindJSON(&renewAccessTokenRequest); err != nil {
		validationErrs := validator.GetValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrs})
		return
	}

	renewAccessTokenResponse, err := h.tokenSvc.RenewAccessToken(ctx, &renewAccessTokenRequest)
	if err != nil {
		statusCode := utils.GetHTTPStatusCode(err)
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, renewAccessTokenResponse)
}
