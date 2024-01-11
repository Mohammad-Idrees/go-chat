package middleware

import (
	"fmt"
	"net/http"
	"project/constants"
	"project/service"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenSvc service.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(constants.HeaderAuthorization)
		if len(authHeader) == 0 {
			err := constants.ErrEmptyAuthHeader
			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := constants.ErrInvalidAuthHeader
			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != constants.BearerAuthorizationType {
			err := fmt.Errorf("unsupported authorization type %s", authType)
			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}

		token := fields[1]
		// TODO: only allow verification of access token
		claims, err := tokenSvc.VerifyToken(c, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}

		c.Set(constants.JWTClaims, claims)
		c.Next()
	}
}
