package mongodb

import (
	"context"
	"time"

	"github.com/gofrs/uuid"

	"github.com/nairobi-gophers/fupisha/internal/encoding"
	"github.com/nairobi-gophers/fupisha/internal/store/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

type userStore struct {
	db  *mongo.Database
	ctx context.Context
}

//New creates a new user document
func (s userStore) New(name, email, password string) (primitive.ObjectID, error) {

	tkn := encoding.GenUniqueID()

	var insertedID primitive.ObjectID //zero value

	user := model.User{
		ID:                  primitive.NewObjectID(),
		Name:                name,
		Email:               email,
		Password:            password,
		VerificationToken:   tkn,
		VerificationExpires: time.Now().Add(time.Minute * 60), //60 mins
		CreatedAt:           time.Now(),
	}

	if err := user.HashPassword(); err != nil {
		return insertedID, err
	}

	result, err := s.db.Collection("users").InsertOne(s.ctx, user)
	if err != nil {
		return insertedID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

//Get finds a user by id
func (s userStore) Get(id string) (model.User, error) {

	user := model.User{}

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	if err := s.db.Collection("users").FindOne(s.ctx, bson.M{"_id": docID}).Decode(&user); err != nil {
		return user, err
	}

	return user, nil
}

//GetByEmail retrieve an existing user with the given email
func (s userStore) GetByEmail(email string) (model.User, error) {
	user := model.User{}

	if err := s.db.Collection("users").FindOne(s.ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return user, err
	}

	return user, nil
}

//SetAPIKey sets the api key for the given user id.
func (s userStore) SetAPIKey(id string, key uuid.UUID) (model.User, error) {

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.User{}, err
	}

	filter := bson.M{"_id": docID}
	update := bson.M{
		"$set": model.User{APIKey: key},
	}

	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	result := s.db.Collection("users").FindOneAndUpdate(s.ctx, filter, update, &opt)

	if result.Err() != nil {
		return model.User{}, result.Err()
	}

	doc := model.User{}
	decodeErr := result.Decode(&doc)

	return doc, decodeErr
}
