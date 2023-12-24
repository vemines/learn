package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

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
