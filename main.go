package main

import (
	"github.com/Dariocent/matchedbetting-clone-go/betting"
	"github.com/Dariocent/matchedbetting-clone-go/controllers"
	"github.com/Dariocent/matchedbetting-clone-go/initializers"
	"github.com/Dariocent/matchedbetting-clone-go/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoanEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.StaticFile("/", "./templates/index.html")
	r.GET("/betting", middleware.RequireAuth, betting.Betting)
	r.Run()
}
