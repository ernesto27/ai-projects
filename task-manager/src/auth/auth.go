package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ernesto/task-manager/src/config"
	"github.com/ernesto/task-manager/src/models"
	"github.com/ernesto/task-manager/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// TokenResponse represents the authentication token response
type TokenResponse struct {
	Token   string `json:"token"`
	Type    string `json:"type"`
	Expires int64  `json:"expires"`
}

// Register handles user registration
func Register(c *gin.Context) {
	var input models.UserRegistrationInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Sanitize inputs
	input.Name = utils.SanitizeInput(input.Name)
	input.Email = utils.SanitizeInput(input.Email)

	// Validate email format
	if !utils.ValidateEmail(input.Email) {
		c.JSON(400, gin.H{"error": "Invalid email format"})
		return
	}

	// Validate password strength
	if valid, message := utils.ValidatePassword(input.Password); !valid {
		c.JSON(400, gin.H{"error": message})
		return
	}

	// Check if user with the same email already exists
	var existingUser models.User
	if result := config.DB.Where("email = ?", input.Email).First(&existingUser); result.RowsAffected > 0 {
		c.JSON(400, gin.H{"error": "Email already in use"})
		return
	}

	// Create user
	user := models.User{
		Name:  input.Name,
		Email: input.Email,
	}

	// Hash the password
	if err := user.HashPassword(input.Password); err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	// Save user to database
	if result := config.DB.Create(&user); result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token for the new user
	token, err := GenerateToken(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set token as HTTP-only cookie
	secure := false
	if os.Getenv("GIN_MODE") == "release" {
		secure = true
	}

	c.SetCookie(
		"auth_token",       // name
		token.Token,        // value
		int(token.Expires), // max age in seconds
		"/",                // path
		"",                 // domain
		secure,             // secure (HTTPS only)
		true,               // httpOnly (prevents JavaScript access)
	)

	c.JSON(201, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
		"auth": token,
	})
}

// Login handles user authentication
func Login(c *gin.Context) {
	var input models.UserLoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Sanitize input
	input.Email = utils.SanitizeInput(input.Email)

	// Validate email format
	if !utils.ValidateEmail(input.Email) {
		c.JSON(400, gin.H{"error": "Invalid email format"})
		return
	}

	// Find user by email
	var user models.User
	if result := config.DB.Where("email = ?", input.Email).First(&user); result.Error != nil {
		// Use generic error message for security (don't reveal if email exists)
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	// Verify password
	if err := user.CheckPassword(input.Password); err != nil {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := GenerateToken(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set token as HTTP-only cookie
	// Using HTTP-only prevents JavaScript from accessing the cookie (more secure)
	// Secure should be true in production when using HTTPS
	secure := false
	if os.Getenv("GIN_MODE") == "release" {
		secure = true
	}

	c.SetCookie(
		"auth_token",       // name
		token.Token,        // value
		int(token.Expires), // max age in seconds
		"/",                // path
		"",                 // domain
		secure,             // secure (HTTPS only)
		true,               // httpOnly (prevents JavaScript access)
	)

	c.JSON(200, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
		"auth": token,
	})
}

// GetCurrentUser returns the currently authenticated user
func GetCurrentUser(c *gin.Context) (*models.User, error) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		return nil, errors.New("user ID not found in context")
	}

	var user models.User
	if result := config.DB.First(&user, userID); result.Error != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

// GenerateToken creates a new JWT token for a user
func GenerateToken(userID uint) (TokenResponse, error) {
	// Set expiration time
	expirationTime := time.Now().Add(24 * time.Hour) // 24 hours

	// Get JWT secret from environment or use default
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")

	// Create claims with user ID
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate signed token
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return TokenResponse{}, err
	}

	return TokenResponse{
		Token:   tokenString,
		Type:    "Bearer",
		Expires: expirationTime.Unix(),
	}, nil
}

// ValidateToken validates the provided JWT token
func ValidateToken(tokenString string) (uint, error) {
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	// Validate claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return 0, errors.New("token expired")
		}

		// Get user ID from claims
		userID := uint(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, errors.New("invalid token")
}

// AuthMiddleware is a Gin middleware to validate JWT tokens
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header required"})
			return
		}

		// Check if header has correct format
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid authorization format"})
			return
		}

		// Extract token
		tokenString := authHeader[7:]

		// Validate token
		userID, err := ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Set user ID in context
		c.Set("userID", userID)

		c.Next()
	}
}

// AuthMiddlewareWeb is a middleware for web pages that checks for authentication via cookie or Authorization header
func AuthMiddlewareWeb() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// Check cookie first
		cookie, err := c.Cookie("auth_token")
		if err == nil && cookie != "" {
			tokenString = cookie
		} else {
			// Try to get token from Authorization header
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenString = authHeader[7:]
			} else {
				// Try to get token from query parameter (for HTMX requests)
				tokenString = c.Query("token")
			}
		}

		// If no token found, redirect to login
		if tokenString == "" {
			c.Redirect(302, "/login")
			c.Abort()
			return
		}

		// Validate token
		userID, err := ValidateToken(tokenString)
		if err != nil {
			c.Redirect(302, "/login")
			c.Abort()
			return
		}

		// Get user from database
		var user models.User
		if result := config.DB.First(&user, userID); result.Error != nil {
			c.Redirect(302, "/login")
			c.Abort()
			return
		}

		// Set user in context
		c.Set("userID", userID)
		c.Set("user", user)

		c.Next()
	}
}

// getEnv retrieves an environment variable or returns a default value if not set
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
