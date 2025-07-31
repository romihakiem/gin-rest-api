package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"gin-rest-api/config"
	"gin-rest-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Middleware struct {
	cfg *config.Config
	rds *redis.Client
}

func NewMiddleware(cfg *config.Config, rds *redis.Client) *Middleware {
	return &Middleware{cfg, rds}
}

func (m *Middleware) CheckAuth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is missing"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	authToken := strings.Split(authHeader, " ")
	if len(authToken) != 2 || authToken[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token format"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userId, err := utils.ValidateToken(m.cfg, authToken[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	_, err = m.rds.Get(c.Request.Context(), strconv.Itoa(int(userId))).Result()
	if err == redis.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	} else if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("userAuth", userId)
	c.Next()
}
