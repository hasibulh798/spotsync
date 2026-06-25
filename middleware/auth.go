package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"spotsync/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWTMiddleware validates the Bearer token in Authorization header and injects user claims.
func JWTMiddleware(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return utils.SendError(c, http.StatusUnauthorized, "Unauthorized", "Missing authorization header")
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return utils.SendError(c, http.StatusUnauthorized, "Unauthorized", "Invalid authorization header format")
			}

			tokenString := parts[1]
			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				return utils.SendError(c, http.StatusUnauthorized, "Unauthorized", "Invalid or expired token")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return utils.SendError(c, http.StatusUnauthorized, "Unauthorized", "Failed to parse token claims")
			}

			// Extract id
			var userID uint
			if idVal, ok := claims["id"]; ok {
				if floatVal, ok := idVal.(float64); ok {
					userID = uint(floatVal)
				} else if intVal, ok := idVal.(int); ok {
					userID = uint(intVal)
				}
			}

			// Extract role
			var role string
			if roleVal, ok := claims["role"]; ok {
				role, _ = roleVal.(string)
			}

			if userID == 0 || role == "" {
				return utils.SendError(c, http.StatusUnauthorized, "Unauthorized", "Invalid token claims payload")
			}

			// Set claims in context
			c.Set("user_id", userID)
			c.Set("role", role)

			return next(c)
		}
	}
}

// RoleMiddleware restricts access based on user roles.
func RoleMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			roleVal := c.Get("role")
			role, ok := roleVal.(string)
			if !ok || role == "" {
				return utils.SendError(c, http.StatusForbidden, "Forbidden", "Insufficient permissions: role not found")
			}

			// Check if role is in allowedRoles
			for _, allowed := range allowedRoles {
				if role == allowed {
					return next(c)
				}
			}

			return utils.SendError(c, http.StatusForbidden, "Forbidden", "Insufficient permissions: access denied for this role")
		}
	}
}
