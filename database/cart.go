package database

import (
	"context"
	"errors"
	"log"

	"github.com/ShabnamHaque/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

var (
	ErrCantFindProduct    = errors.New("cant find product")
	ErrCantDecodeProducts = errors.New("cant find the products")
	ErrCantUpdateUser     = errors.New("cant update user")
	ErrCantRemoveItem     = errors.New("cant remove this item from cart")
	ErrCantBuyCartItem    = errors.New("cant buy this item")
	ErrUserIdNotValid     = errors.New("user is not valid")
)

func AddProductToCart(ctx context.Context, productCollection *mongo.Collection, userCollection *mongo.Collection, productId primitive.ObjectID, userId string) error {
	searchfromdb, err := productCollection.Find(ctx, bson.M{"_id": productId})
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}
	var productCart []models.ProductUser
	err = searchfromdb.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	}

	usert_id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		return ErrUserIdNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
	update := bson.D{primitive.E{Key: "$push", Value: bson.D{Key: "usercart", Value: bson.D{Key: "$each", Value: productCart}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrCantUpdateUser
	}
	return nil
}
func RemoveCartItem(ctx context.Context, productCollection *mongo.Collection, userCollection *mongo.Collection, productId primitive.ObjectID, userId string) error {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		return ErrUserIdNotValid
	}
	filter := bson.D(primitive.E{Key: "_id", Value: id})
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productId}}} //pull away that one product from usercart
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantRemoveItem
	}
	return nil
}
func BuyItemFromCart() gin.HandlerFunc {}
func InstantBuy() gin.HandlerFunc      {}
