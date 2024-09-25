package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	First_Name      *string            `json:"first_name" validate="required,min=2,max=20"`
	Last_Name       *string            `json:"last_name"  validate="required,min=2,max=20"`
	Email           *string            `json:"email"      validate="email,required"`
	Password        *string            `json:"password"   validate="required,min=6,max=20"`
	Phone           *string            `json:"phone"      validate="required"`
	Token           *string            `json:"token"`
	Refresh_Token   *string            `json:"refresh_token"`
	Created_At      time.Time          `json:"created_at"`
	Updated_At      time.Time          `json:"updated_at"`
	User_ID         string             `json:"user_id"`
	UserCart        []ProductUser      `json:"usercart" bson:"usercart"`
	Address_Details []Address          `json:"address" bson:"address"`
	Order_Status    []Order            `json:"orders" bson:"orders"`
}
type Product struct {
	Product_ID   primitive.ObjectID `bson:"_id"`
	Product_Name *string            `json:"product_name"`
	Price        *uint64            `json:"price"`
	Image        *string            `json:"image"`
	Rating       *uint64            `json:"rating"`
}
type ProductUser struct {
	Product_ID   primitive.ObjectID `bson:"_id"`
	Product_Name *string            `json:"product_name"`
	Price        *uint64            `json:"price"`
	Image        *string            `json:"image"`
	Rating       *uint64            `json:"rating"`
}
type Address struct {
	House      *string            `json:"house_name"`
	Street     *string            `json:"street"`
	Address_ID primitive.ObjectID `bson:"_id"`
	City       *string            `json:"city"`
	Pincode    *string            `json:"pin_code"`
}
type Order struct {
	Order_ID       primitive.ObjectID `bson:"_id"`
	Order_Cart     []ProductUser      `json:"order_cart" bson:"order_cart"`
	Ordered_At     time.Time          `json:"ordered_at"`
	Price          *uint64            `json:"price"`
	Discount       *uint64            `json:"discount"`
	Payment_Method Payment            `json:"payment_method"`
}
type Payment struct {
	Digital bool `json:"digital"`
	COD     bool `json:"cod"`
}
