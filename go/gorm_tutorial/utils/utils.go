package utils

import (
	"gorm_tutorial/models"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func CreateToken(db *gorm.DB, userID uint) (string, error) {
	// Generate a new JWT token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "my-api",
		"sub": userID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	// Sign the token.
	signedToken, err := token.SignedString([]byte("my-secret-key"))
	if err != nil {
		return "", err
	}

	// Create a new token record in the database.
	tokenRecord := models.Token{
		UserId:    userID,
		Token:     signedToken,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	err = db.Create(tokenRecord).Error
	if err != nil {
		return "", err
	}

	return tokenRecord.Token, nil
}

func VerifyToken(db *gorm.DB, tokenString string) (bool, error) {
	// Parse the token.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the token signature.
		return []byte("my-secret-key"), nil
	})
	if err != nil {
		return false, err
	}

	// Check if the token is expired.
	claims := token.Claims.(jwt.MapClaims)
	exp := claims["exp"].(float64)
	nowUnix := float64(time.Now().Unix())
	if nowUnix > exp {
		return false, nil
	}

	// Check if the token exists in the database.
	var tokenRecord models.Token
	err = db.First(&tokenRecord, "token = ?", tokenString).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
