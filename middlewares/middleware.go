package middlewares

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"order-service/clients"
	"order-service/common/response"
	"order-service/config"
	"order-service/constants"
	errConstants "order-service/constants/error"
	"strings"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func HandlePanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("Recovered from panic: %v", r)
				c.JSON(http.StatusInternalServerError, response.Response{
					Status:  constants.Error,
					Message: errConstants.ErrInternalServer.Error(),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

func RateLimiter(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusTooManyRequests, response.Response{
				Status:  constants.Error,
				Message: errConstants.ErrTooManyRequest.Error(),
			})
			c.Abort()
		}
		c.Next()
	}
}

func extractBearerToken(token string) string {
	arrayToken := strings.Split(token, " ")

	if len(arrayToken) == 2 {
		return arrayToken[1]
	}

	return ""
}

func responseUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, response.Response{
		Status:  constants.Error,
		Message: message,
	})
	c.Abort()
}

func validateAPIKey(c *gin.Context) error {
	apiKey := c.GetHeader(constants.XApiKey)
	requestAt := c.GetHeader(constants.XRequestAt)
	serviceName := c.GetHeader(constants.XServiceName)
	signatureKey := config.Config.SignatureKey

	validateKey := fmt.Sprintf("%s:%s:%s", serviceName, signatureKey, requestAt)
	hash := sha256.New()
	hash.Write([]byte(validateKey))
	resultHash := hex.EncodeToString(hash.Sum(nil))
	logrus.Infof("API Key: %s", resultHash)

	if apiKey != resultHash {
		logrus.Info("Api key didn't same like result hash")
		return errConstants.ErrUnauthorized
	}

	return nil
}

func contains(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}

	return false
}

func CheckRole(roles []string, clients clients.IClientRegistry) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := clients.GetUser().GetUserByToken(c.Request.Context())
		if err != nil {
			responseUnauthorized(c, errConstants.ErrUnauthorized.Error())
			return
		}

		if !contains(roles, user.Role) {
			responseUnauthorized(c, errConstants.ErrUnauthorized.Error())
			return
		}
		userLogin := c.Request.WithContext(context.WithValue(c.Request.Context(), constants.User, user))
		c.Request = userLogin
		c.Next()
	}
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {

		var err error
		token := c.GetHeader(constants.Authorization)
		if token == "" {
			responseUnauthorized(c, errConstants.ErrUnauthorized.Error())
			return
		}

		err = validateAPIKey(c)
		if err != nil {
			responseUnauthorized(c, err.Error())
			return
		}

		tokenString := extractBearerToken(token)
		tokenUser := c.Request.WithContext(context.WithValue(c.Request.Context(), constants.Token, tokenString))
		c.Request = tokenUser
		c.Next()

		c.Next()
	}
}

func AuthenticateWithoutToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := validateAPIKey(c)
		if err != nil {
			responseUnauthorized(c, err.Error())
			return
		}
		c.Next()
	}
}
