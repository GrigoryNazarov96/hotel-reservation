package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type RoomType int

const (
	_ RoomType = iota
	Single
	Double
	Suite
	Deluxe
)

type Room struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type      RoomType           `bson:"type" json:"type"`
	BasePrice float64            `bson:"basePrice" json:"basePrice"`
	HotelID   primitive.ObjectID `bson:"hotelID" json:"hotelID"`
}
