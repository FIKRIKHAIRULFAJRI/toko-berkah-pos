package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/toko-berkah/pos-backend/repositories"
	"github.com/toko-berkah/pos-backend/utils"
	"gorm.io/gorm"
)

// AuthMiddleware validates JWT token and PasswordVersion
func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, utils.NewErrorResponse(utils.NewUnauthorizedError()))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, utils.NewErrorResponse(utils.NewUnauthorizedError()))
			c.Abort()
			return
		}

		claims, err := utils.VerifyToken(parts[1])
		if err != nil {
			c.JSON(401, utils.NewErrorResponse(utils.NewUnauthorizedError()))
			c.Abort()
			return
		}

		// Verify PasswordVersion matches current database version
		userRepo := repositories.NewUserRepository(db)
		user, err := userRepo.GetUserByID(claims.UserID)
		if err != nil {
			c.JSON(401, utils.NewErrorResponse(utils.NewUnauthorizedError()))
			c.Abort()
			return
		}

		// Check if token version matches current password version
		if claims.PasswordVersion != user.PasswordVersion {
			c.JSON(401, utils.NewErrorResponse(
				utils.NewCustomError("UNAUTHORIZED", "Token has been invalidated. Please login again."),
			))
			c.Abort()
			return
		}

		// Store claims in context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("mustChangePassword", user.MustChangePassword)

		c.Next()
	}
}

// RequirePasswordChange enforces password change for first-time users
// Apply this to endpoints that should be accessible even with MustChangePassword flag
func RequirePasswordChange() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow only /auth/change-password endpoint when MustChangePassword is true
		mustChange, exists := c.Get("mustChangePassword")
		if exists && mustChange.(bool) {
			// Only allow change-password endpoint
			if c.Request.URL.Path != "/api/v1/auth/change-password" {
				c.JSON(403, utils.NewErrorResponse(
					utils.NewCustomError("FORBIDDEN", "You must change your password before accessing other endpoints"),
				))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// RequireRole checks if user has required role (kept for future use, not used in MVP)
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(403, utils.NewErrorResponse(utils.NewForbiddenError()))
			c.Abort()
			return
		}

		hasRole := false
		for _, role := range roles {
			if userRole == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(403, utils.NewErrorResponse(utils.NewForbiddenError()))
			c.Abort()
			return
		}

		c.Next()
	}
}

// CORSMiddleware handles CORS
func CORSMiddleware(frontendURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", frontendURL)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// ErrorHandlerMiddleware handles panics and errors
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(500, utils.NewErrorResponse(
					utils.NewInternalError("Internal server error"),
				))
			}
		}()
		c.Next()
	}
}
