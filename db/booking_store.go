package db

import (
	"context"

	"github.com/GrigoryNazarov96/hotel-reservation.git/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const bookingColl string = "bookings"

type BookingStore interface {
	CreateBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(context.Context, bson.M) ([]*types.Booking, error)
	GetBookingByID(context.Context, string) (*types.Booking, error)
	UpdateBooking(context.Context, string, types.UpdateBookingDTO) error
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	BookingStore
}

func NewMongoBookingStore(c *mongo.Client, dbname string) *MongoBookingStore {
	return &MongoBookingStore{
		client: c,
		coll:   c.Database(dbname).Collection(bookingColl),
	}
}

func (s *MongoBookingStore) CreateBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	cur, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.Id = cur.InsertedID.(primitive.ObjectID)
	return booking, nil
}

func (s *MongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error) {
	var bookings []*types.Booking
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err := cur.All(ctx, &bookings); err != nil {
		return []*types.Booking{}, err
	}
	return bookings, nil
}

func (s *MongoBookingStore) GetBookingByID(ctx context.Context, id string) (*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var booking *types.Booking
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking); err != nil {
		return nil, err
	}
	return booking, err
}

func (s *MongoBookingStore) UpdateBooking(ctx context.Context, id string, dto types.UpdateBookingDTO) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	u_query := bson.M{
		"$set": dto.ToBSONM(),
	}
	if _, err := s.coll.UpdateByID(ctx, oid, u_query); err != nil {
		return err
	}
	return nil
}
