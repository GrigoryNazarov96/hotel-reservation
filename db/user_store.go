package db

import (
	"context"

	"github.com/GrigoryNazarov96/hotel-reservation.git/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userColl string = "users"

type UserStore interface {
	Dropper
	GetUserByID(context.Context, string) (*types.User, error)
	GetUserByEmail(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	CreateUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(context.Context, string, types.UpdateUserDTO) (*types.User, error)
	Drop(context.Context) error
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(c *mongo.Client, dbname string) *MongoUserStore {
	return &MongoUserStore{
		client: c,
		coll:   c.Database(dbname).Collection(userColl),
	}
}

func (s *MongoUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []*types.User
	if err := cur.All(ctx, &users); err != nil {
		return []*types.User{}, nil
	}
	return users, nil
}

func (s *MongoUserStore) CreateUser(ctx context.Context, user *types.User) (*types.User, error) {
	cur, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.Id = cur.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, id string, dto types.UpdateUserDTO) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	u_query := bson.D{
		{Key: "$set", Value: dto.ToBSONM()},
	}
	_, err = s.coll.UpdateByID(ctx, oid, u_query)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *MongoUserStore) Drop(ctx context.Context) error {
	if err := s.coll.Drop(ctx); err != nil {
		return err
	}
	return nil
}
