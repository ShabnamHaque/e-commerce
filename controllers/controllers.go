package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ShabnamHaque/database"
	"github.com/ShabnamHaque/models"
	generate "github.com/ShabnamHaque/tokens"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

var Validate *validator.Validate
var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Users")

func Login(userCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		var foundUser models.User
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or password incorrect"})
			return
		}
		PasswordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)

		defer cancel()

		if !PasswordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}
		token, refreshToken, _ := generate.TokenGenerator(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_Name, foundUser.User_ID)
		defer cancel()

		generate.UpdateAllTokens(token, refreshToken, foundUser.User_ID)
		c.JSON(http.StatusFound, foundUser)
	}
}

func Signup(userCollection *mongo.Collection) gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//we are validating the data fields from the user struct.
		validationErr := Validate.Struct(user) //from validator package - checks if key - values are intact with the struct defined
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}
		var count, err = userCollection.CountDocuments(ctx, bson.M{"email": user.Email}) //count users with same email
		defer cancel()

		if err != nil {
			log.Panic(err)
			ErrorData := gin.H{"error": "Error occurred while checking for email"}
			c.JSON(http.StatusInternalServerError, ErrorData)
		}
		if count > 0 { //because signup we need to maintain both email and phone as unique for each user..
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "This email or/and phone number already exists"})
		}

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone}) //count documents with the entered phone we want unique phone and email
		defer cancel()
		if err != nil {
			log.Panic(err)
			ErrorData := gin.H{"error": "Error occurred while checking for Phone Number"}
			c.JSON(http.StatusInternalServerError, ErrorData)

		}
		if count > 0 { //because signup we need to maintain both email and phone as unique for each user..
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "This email or/and phone number already exists"})
			return
		}
		password := HashPassword(*user.Password)
		user.Password = &password
		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		token, refreshToken, _ := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshToken

		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0) //using slices
		user.Order_Status = make([]models.Order, 0)

		_, inserterr := userCollection.InsertOne(ctx, user)

		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User did not get inserted"})
			return
		}

		defer cancel()

		c.JSON(http.StatusCreated, "Successfully signed in!")
	}
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPass string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPass))
	valid := true
	msg := ""

	if err != nil {
		msg = "Login or password is incorrect"
		valid = false
	}
	return valid, msg
}
func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchProducts []models.Product
		queryParam := c.Query("name")
		if queryParam == "" {
			log.Println("query empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid search"})
			c.Abort()
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		searchquerydb, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Something went wrong while fetching the product data")
			return
		}
		err = searchquerydb.All(ctx, &searchProducts) //converts bson data back to json golang understandable data
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer searchquerydb.Close(ctx)

		if err := searchquerydb.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid req")
			return
		}
		defer cancel()
		c.IndentedJSON(200, searchProducts)
	}
}

func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productList []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Something went wrong")
			return
		}
		err = cursor.All(ctx, &productList) //converts bson data back to json golang understandable data

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)

		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer cancel()
		c.IndentedJSON(200, productList)
	}
}
