package main

import (
	"log"
	"os"

	"github.com/ShabnamHaque/ecommerce/controllers"
	"github.com/ShabnamHaque/ecommerce/database"
	"github.com/ShabnamHaque/ecommerce/middleware"
	"github.com/ShabnamHaque/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT") //start the server
	if port == "" {
		port = "8080"
	}
	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))
	// two tables Products and Users in the Database

	router := gin.New() //creates a new gin router
	router.Use(gin.Logger()) // Use Logger middleware for logging requests

	routes.UserRoutes(router) //sends them to routes.go
	router.Use(middleware.Authentication())  // authentication

	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())
	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())

	log.Fatal(router.Run(":" + port))

}
