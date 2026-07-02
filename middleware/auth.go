package middleware

import (
	"net/http"
	"ppdb-be/config"
	"ppdb-be/models"
	"ppdb-be/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.SendError(c, http.StatusUnauthorized, utils.MsgUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.SendError(c, http.StatusUnauthorized, utils.MsgUnauthorized, "Authorization header must be Bearer token")
			c.Abort()
			return
		}

		tokenStr := parts[1]
		claims, err := utils.ValidateToken(tokenStr)
		if err != nil {
			utils.SendError(c, http.StatusUnauthorized, utils.MsgUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		var session models.UserSession
		if err := config.DB.WithContext(c).Where("access_token = ? AND is_active = ? AND is_revoked = ?", tokenStr, true, false).First(&session).Error; err != nil {
			utils.SendError(c, http.StatusUnauthorized, utils.MsgUnauthorized, "Sesi Anda telah berakhir atau login di perangkat lain, silakan login kembali")
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userName", claims.UserName)
		c.Set("userEmail", claims.Email)
		c.Set("sessionID", session.SessionID)

		c.Next()
	}
}
