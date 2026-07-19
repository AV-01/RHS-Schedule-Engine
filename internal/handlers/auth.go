package handlers

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/AV-01/RHS-Schedule-Engine/internal/db"
	"github.com/AV-01/RHS-Schedule-Engine/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login godoc
//
//	@Summary		Authenticate and get a JWT
//	@Description	Login with username (firstname.lastname) and your student ID as the password. Returns a JWT valid for 24 hours. For public demo access, skip this endpoint and use the token "dmeo-key" directly in the Authorization header.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		LoginRequest	true	"Login credentials"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Router			/api/v1/auth/login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password fields are required"})
		return
	}

	var userID, username string
	err := db.DB.QueryRow(`SELECT id::text, username FROM students WHERE username = $1 AND student_id = $2`,
		req.Username, req.Password).Scan(&userID, &username)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	claims := &middleware.Claims{
		UserID:   userID,
		Username: username,
		IsDemo:   false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to gen token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   signed,
		"message": "login successful, token expires in 24 hours",
	})
}
