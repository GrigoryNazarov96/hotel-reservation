package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserId   primitive.ObjectID `bson:"userId,omitempty" json:"userId,omitempty"`
	RoomId   primitive.ObjectID `bson:"roomId,omitempty" json:"roomId,omitempty"`
	FromDate time.Time          `bson:"fromDate" json:"fromDate"`
	TillDate time.Time          `bson:"tillDate" json:"tillDate"`
	Guests   int                `bson:"guests" json:"guests"`
}

type BookingDTO struct {
	FromDate time.Time `json:"fromDate"`
	TillDate time.Time `json:"tillDate"`
	Guests   int       `json:"guests"`
}
