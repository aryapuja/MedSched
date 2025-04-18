package middleware

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Struct untuk claims JWT
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

var secretKey string

// Ambil secretKey
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	secretKey = os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY is not set in .env file")
	}
}

// GenerateJWT
func GenerateJWT(username string, role string) (string, error) {
	claims := Claims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // Token kadaluarsa dalam 24 jam
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Memverifikasi JWT di setiap request
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mendapatkan token dari header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(400, gin.H{"error": "Authorization token is required"})
			c.Abort()
			return
		}

		// Menghapus "Bearer " dari token
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Verifikasi token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Menyimpan klaim di context
		claims, _ := token.Claims.(*Claims)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// Role-Based Auth Middleware (tes)
func RoleRequired(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleFromToken := c.MustGet("role").(string)

		// Cek role dari token
		if roleFromToken != role && role != "any" {
			c.JSON(403, gin.H{"error": "Forbidden: You do not have the required role"})
			c.Abort()
			return
		}

		c.Next()
	}
}