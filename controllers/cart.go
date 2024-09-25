package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/ShabnamHaque/database"
	"github.com/ShabnamHaque/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {

	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	} //returns the address.
}
func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Print(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "Successfully added to the cart.")

	}
}
func RemoveItem(app *Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Print(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(200, "Successfully removed Item from Cart")
	}
}

func BuyFromCart(app *Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryID := c.Query("id")

		if userQueryID == "" {
			log.Panic("user id empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "successfully placed the order")
	}
}
func GetItemsFromCart(app *Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Context-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid id"})
			c.Abort()
			return
		}
		usert_id, _ := primitive.ObjectIDFromHex(user_id)
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var filledCart models.User
		err := UserCollection.FindOne(ctx, bson.M{"_id": usert_id}).Decode(&filledCart)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "not found")
			return
		}

		///	filter :=

		// filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: usert_id}}}}
		// unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
		// grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}} // ***

		filter_match := bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "_id", Value: usert_id},
			}},
		}

		unwind := bson.D{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$usercart"},
			}},
		}

		grouping := bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$_id"},
				{Key: "total", Value: bson.D{
					{Key: "$sum", Value: "$usercart.price"},
				}},
			}},
		}

		// pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})

		// if err != nil {
		// 	log.Println(err)
		// }

		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{
			filter_match,
			unwind,
			grouping,
		})

		if err != nil {
			log.Fatal(err)
		}

		var listing []bson.M
		if err = pointcursor.All(ctx, &listing); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		for _, json := range listing {
			c.IndentedJSON(200, json["total"])
			c.IndentedJSON(200, filledCart.UserCart)
		}
		ctx.Done()
	}
}
func InstantBuy(app *Application) gin.HandlerFunc {

	return func(c *gin.Context) {

		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Print(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.InstantBuy(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(200, "Successfully placed the order.")
	}

}
