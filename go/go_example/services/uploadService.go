package services

import (
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadSingleFile(c *gin.Context) {
	// Get file in multipart/form-data
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get file" + err.Error(),
		})
		return
	}
	// Generate new file name with hash (time + random float)
	fileName := generateName(parseFileExt(file.Filename))
	// Create path where to save file
	filePath := os.Getenv("UPLOADPATH") + fileName
	// Save file
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to save file" + err.Error(),
		})
		return
	}
	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully: " + filePath[1:],
	})
}
func UploadMultiFile(c *gin.Context) {
	// Get files in multipart/form-data
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get file" + err.Error(),
		})
		return
	}
	files := form.File["files"]
	fileMapSlice := []map[string]string{}

	for _, file := range files {
		// Generate new file name with hash (time + random float)
		fileName := generateName(parseFileExt(file.Filename))
		// Create path where to save file
		filePath := os.Getenv("UPLOADPATH") + fileName
		// Save file
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to save file" + err.Error(),
			})
			return
		}
		// Add to slice to display at response
		fileMap := map[string]string{
			"fileName": file.Filename,
			"savePath": filePath[1:],
		}
		fileMapSlice = append(fileMapSlice, fileMap)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Files uploaded successfully",
		"data":    fileMapSlice,
	})
}

func parseFileExt(fileName string) (ext string) {
	fileNameParts := strings.Split(fileName, ".")
	if len(fileNameParts) > 1 {
		ext = "." + fileNameParts[len(fileNameParts)-1]
	} else {
		ext = ""
	}
	return
}

func generateName(ext string) (name string) {
	randomFloat := rand.Float64()
	floatString := strconv.FormatFloat(randomFloat, 'f', -1, 64)
	hash := sha1.New()
	hash.Write([]byte(time.Now().Format(time.RFC3339Nano) + floatString))
	name = hex.EncodeToString(hash.Sum(nil)) + ext
	return
}
