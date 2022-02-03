package userprofile

import (
	"context"
	"golang-vscode-setup/service/userprofile/defined"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Collection struct {
	UserId     string    `bson:"_id,omitempty"`
	FirstName  string    `bson:"first_name"`
	LastName   string    `bson:"last_name"`
	Otp        int8      `bson:"otp"`
	OtpCreated time.Time `bson:"otp_created"`
	FcmToken   string    `bson:"fcm_token"`
}

func NewUserProfile(userProfile defined.UserProfile) *Collection {
	return &Collection{
		UserId:     userProfile.UserId,
		FirstName:  userProfile.FirstName,
		LastName:   userProfile.LastName,
		Otp:        userProfile.Otp,
		OtpCreated: userProfile.OtpCreated,
		FcmToken:   userProfile.FcmToken,
	}
}

func (collection *Collection) Map() defined.UserProfile {
	var data defined.UserProfile
	data.UserId = collection.UserId
	data.FirstName = collection.FirstName
	data.LastName = collection.LastName
	data.Otp = collection.Otp
	data.OtpCreated = collection.OtpCreated
	data.FcmToken = collection.FcmToken
	return data
}

type Repository struct {
	db *mongo.Collection
}

func NewRepository(db *mongo.Database) (*Repository, error) {
	repository := Repository{db.Collection("user_profile")}
	return &repository, nil
}

type IRepository interface {
	Create(userProfile defined.UserProfile) (*defined.UserProfile, error)
	FindById(userId string) (*defined.UserProfile, error)
	UpdateById(userId string, userProfile defined.UserProfile) (*defined.UserProfile, error)
	DeleteById(userId string) error
}

func (repository Repository) Create(userProfile defined.UserProfile) (*defined.UserProfile, error) {
	data := NewUserProfile(userProfile)
	_, err := repository.db.InsertOne(context.TODO(), data)
	if err != nil {
		return nil, err
	}
	result := data.Map()
	return &result, nil
}

func (repository Repository) FindById(userId string) (*defined.UserProfile, error) {
	var data Collection
	err := repository.db.FindOne(context.TODO(), bson.M{"_id": userId}).Decode(&data)
	if err != nil {
		return nil, err
	}
	result := data.Map()
	return &result, nil
}

func (repository Repository) UpdateById(userId string, userProfile defined.UserProfile) (*defined.UserProfile, error) {
	result, err := repository.db.UpdateByID(context.TODO(), userId, bson.M{"$set": userProfile})
	if result.MatchedCount == 0 {
		return nil, mongo.ErrNoDocuments
	}
	if err != nil {
		return nil, err
	}
	data, _ := repository.FindById(userId)
	return data, nil
}

func (repository Repository) DeleteById(userId string) error {
	result, err := repository.db.DeleteOne(context.TODO(), bson.M{"_id": userId})
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	if err != nil {
		return err
	}
	return nil
}
