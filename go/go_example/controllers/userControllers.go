package controllers

import (
	"encoding/json"
	"fmt"
	"go_example/initializers"
	"go_example/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	// Get email/password of req body
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	if body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email or password are empty",
		})
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}
	// Create user
	user := models.User{Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&user) // pass pointer of data to Create

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user" + result.Error.Error(),
		})
		return
	}
	// Respond
	c.JSON(http.StatusOK, gin.H{})
}

func Login(c *gin.Context) {
	// Get email/password of req body
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	if body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email or password are empty",
		})
		return
	}

	// Look up requested user
	var user models.User
	initializers.DB.First(&user, "email = ? ", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email not exist",
		})
		return
	}

	// check hash password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong password",
		})
		return
	}

	// generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SERCRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token" + err.Error(),
		})
		return
	}

	// Set cookie

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 2600*24*30, "", "", false, true)

	// Response
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}

func VerifyCaptcha(c *gin.Context) {
	// Get the SecretKey and Response from the request body
	var body struct {
		SecretKey string `json:"secret"`
		Response  string `json:"response"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	if body.SecretKey == "" || body.Response == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "SecretKey or Response are empty",
		})
		return
	}

	// Verify the reCAPTCHA token with Google's API
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://www.google.com/recaptcha/api/siteverify", strings.NewReader(fmt.Sprintf("secret=%s&response=%s", body.SecretKey, body.Response)))
	if err != nil {
		// Handle request creation error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create verification request",
		})
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		// Handle verification request error
		c.JSON(http.StatusBadRequest, gin.H{
			"success":     false,
			"error-codes": []string{"Failed to verify captcha"},
		})
		return
	}

	defer resp.Body.Close()

	var responseData struct {
		Success     bool     `json:"success"`
		ErrorCodes  []string `json:"error-codes"`
		ChallengeTs string   `json:"challenge_ts"`
		Hostname    string   `json:"hostname"`
		Score       float64  `json:"score"`
		Action      string   `json:"action"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		// Handle response parsing error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse verification response",
		})
		return
	}

	// Send the verification result back to the client
	if responseData.Success {
		c.JSON(http.StatusOK, gin.H{
			"success":      true,
			"score":        responseData.Score,
			"challenge_ts": responseData.ChallengeTs,
			"hostname":     responseData.Hostname,
			"action":       responseData.Action,
			// Add other fields as needed
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":     false,
			"error-codes": responseData.ErrorCodes,
		})
	}
}