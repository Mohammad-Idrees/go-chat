package delivery

import (
	"net/http"
	"project/constants"
	"project/models"
	"project/models/request"
	"project/service"
	"project/utils"
	"project/validator"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc service.UserService
}

func NewUserHandler(userSvc service.UserService) *UserHandler {
	return &UserHandler{userSvc}
}

func ConfigureUserHandler(router *gin.RouterGroup, authMiddleware gin.HandlerFunc, userSvc service.UserService) {
	userHandler := NewUserHandler(userSvc)
	addUserHandlerRoutes(router, authMiddleware, userHandler)
}

func addUserHandlerRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc, userHandler *UserHandler) {
	router.POST("/signup", userHandler.CreateUser)
	router.POST("/login", userHandler.LoginUser)
	router.POST("/logout", userHandler.LogoutUser)
	router.GET("/users/:email", authMiddleware, userHandler.GetUserByEmail)
	router.GET("/users", userHandler.GetUsers)
	router.POST("/users", userHandler.CreateUser)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()
	var createUserRequest request.CreateUserRequest
	if err := c.ShouldBindJSON(&createUserRequest); err != nil {
		validationErrs := validator.GetValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrs})
		return
	}

	user, err := h.userSvc.CreateUser(ctx, &createUserRequest)
	if err != nil {
		statusCode := utils.GetHTTPStatusCode(err)
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	ctx := c.Request.Context()
	var loginUserRequest request.LoginUserRequest
	loginUserRequest.UserAgent = c.Request.UserAgent()
	loginUserRequest.ClientIp = c.ClientIP()
	if err := c.ShouldBindJSON(&loginUserRequest); err != nil {
		validationErrs := validator.GetValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrs})
		return
	}

	loginUserResponse, err := h.userSvc.LoginUser(ctx, &loginUserRequest)
	if err != nil {
		statusCode := utils.GetHTTPStatusCode(err)
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, loginUserResponse)
}

func (h *UserHandler) LogoutUser(c *gin.Context) {
	ctx := c.Request.Context()
	var logoutUserRequest request.LogoutUserRequest
	if err := c.ShouldBindJSON(&logoutUserRequest); err != nil {
		validationErrs := validator.GetValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrs})
		return
	}

	err := h.userSvc.LogoutUser(ctx, &logoutUserRequest)
	if err != nil {
		statusCode := utils.GetHTTPStatusCode(err)
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "logged out successfully")
}

func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.GetUserByEmailRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userSvc.GetUserByEmail(ctx, &req)
	if err != nil {
		statusCode := utils.GetHTTPStatusCode(err)
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	jwtClaims := c.MustGet(constants.JWTClaims).(*models.JWTClaims)
	if user.Email != jwtClaims.Email {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrAccessDenied.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	ctx := c.Request.Context()
	users, err := h.userSvc.GetUsers(ctx)
	if err != nil {
		statusCode := utils.GetHTTPStatusCode(err)
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
