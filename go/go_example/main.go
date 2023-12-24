package main

import (
	"go_example/controllers"
	// "go_example/initializers"
	"go_example/middleware"
	"go_example/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	// initializers.LoadEnvVariables()
	// initializers.ConnnectToDb()
	// initializers.SyncDb()
}

func main() {
	// services.SendMail()
	ginInitializers()
}

func ginInitializers() {
	r := gin.Default()

	r.Use(cors.Default())

	// Or, enable CORS for specific origins
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"*"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept"},
	}))

	r.Static("/assets", "./assets")
	r.Static("/templates", "./templates")

	r.POST("/uploadSingleFile", services.UploadSingleFile)
	r.POST("/uploadMultiFile", services.UploadMultiFile)
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.POST("/verifyCaptcha", controllers.VerifyCaptcha)
	r.Run() // listen and serve on 0.
}
