package db

import (
	"context"

	"github.com/GrigoryNazarov96/hotel-reservation.git/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const roomColl string = "rooms"

type RoomStore interface {
	CreateRoom(context.Context, *types.Room) (*types.Room, error)
	GetRooms(context.Context, bson.M) ([]*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore, dbname string) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(dbname).Collection(roomColl),
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) CreateRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.Id = res.InsertedID.(primitive.ObjectID)
	u_query := bson.M{"$push": bson.M{"rooms": room.Id}}
	if err := s.HotelStore.UpdateHotel(ctx, room.HotelID, u_query); err != nil {
		return nil, err
	}
	return room, nil
}

func (s *MongoRoomStore) GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error) {
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var rooms []*types.Room
	if err := cur.All(ctx, &rooms); err != nil {
		return []*types.Room{}, err
	}
	return rooms, err
}
