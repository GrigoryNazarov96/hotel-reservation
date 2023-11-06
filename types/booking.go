package types

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserId    primitive.ObjectID `bson:"userId,omitempty" json:"userId,omitempty"`
	RoomId    primitive.ObjectID `bson:"roomId,omitempty" json:"roomId,omitempty"`
	FromDate  time.Time          `bson:"fromDate" json:"fromDate"`
	TillDate  time.Time          `bson:"tillDate" json:"tillDate"`
	Guests    int                `bson:"guests" json:"guests"`
	Cancelled bool               `bson:"cancelled,omitempty" json:"cancelled,omitempty"`
}

type BookingDTO struct {
	FromDate time.Time `json:"fromDate"`
	TillDate time.Time `json:"tillDate"`
	Guests   int       `json:"guests"`
}

type UpdateBookingDTO struct {
	Cancelled bool `json:"cancelled,omitempty"`
}

func (d UpdateBookingDTO) ToBSONM() bson.M {
	m := bson.M{}
	m["cancelled"] = d.Cancelled
	return m
}

func (d BookingDTO) Validate() error {
	now := time.Now()
	if now.After(d.FromDate) {
		return fmt.Errorf("you can not book a room for past date")
	}
	if d.FromDate.After(d.TillDate) {
		return fmt.Errorf("check in date should be after check out date")
	}
	return nil
}
