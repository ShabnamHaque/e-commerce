package tokens

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ShabnamHaque/database"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// contains logic to create and refresh tokens also to validate jwt authentciation
type signedDetails struct {
	Email      string
	First_Name string
	Last_Name  string
	Uid        string
	jwt.StandardClaims
}

var UserData *mongo.Collection = database.UserData(database.Client, "Users")
var SECRET_KEY = os.Getenv("SECRET_KEY")

func TokenGenerator(email string, firstName string, lastName string, uid string) (signedtoken string, signedrefreshtoken string, err error) {
	claims := &signedDetails{
		Email:      email,
		First_Name: firstName,
		Last_Name:  lastName,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(), //valid for 24 hrs
		},
	}

	refreshclaims := &signedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}
	refreshtoken, err := jwt.NewWithClaims(jwt.SigningMethodHS384, refreshclaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}
	return token, refreshtoken, err
}
func ValidateToken(signedtoken string) (claims *signedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedtoken, &signedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*signedDetails)
	if ok == false {
		msg = "token is invalid"
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
		return
	}
	return claims, msg
}
func UpdateAllTokens(signedtoken string, signedrefreshtoken string, userid string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var updatedobj primitive.D
	updatedobj = append(updatedobj, bson.E{Key: "token", Value: signedtoken})
	updatedobj = append(updatedobj, bson.E{Key: "refresh_token", Value: signedrefreshtoken})
	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updatedobj = append(updatedobj, bson.E{Key: "updated_at", Value: updated_at})
	upsert := true
	filter := bson.M{"user_id": userid}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := UserData.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: updatedobj},
	},
		&opt)

	defer cancel()
	if err != nil {
		log.Panic(err)
		return
	}
}
