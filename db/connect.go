package db

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UrlShortnerCollection struct {
	Id       string `json:"id"`
	LongUrl  string `json:"long_url" bson:"long_url"`   // long url
	ShortUrl string `json:"short_url" bson:"short_url"` // short url endpoint
}

type UrlShortner struct {
	Collection            *mongo.Collection
	UrlShortnerCollection *UrlShortnerCollection
}

type UrlShortnerInterface interface {
	Insert(url string) (string, error)
	Find(shorUrl string) (string, error)
}

func UrlShortnerImpl() UrlShortnerInterface {
	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(context, options.Client().ApplyURI(os.Getenv("MONGO_URL")))
	if err != nil {
		panic(err)
	}

	err = client.Ping(context, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("MongoDB connected successfully")

	db := client.Database(os.Getenv("MONGO_DB_NAME"))
	collection := db.Collection(os.Getenv("MONGO_COLLECTION_NAME"))

	return &UrlShortner{
		Collection: collection,
	}
}

func (u *UrlShortner) Insert(url string) (string, error) {
	shortUrl := strings.Split(uuid.NewMD5(uuid.New(), []byte(url)).String(), "-")

	result, err := u.Find(shortUrl[len(shortUrl)-1])
	if err != nil {
		return "", err
	}
	if result != "" {
		return u.Insert(url)
	}

	u.UrlShortnerCollection = &UrlShortnerCollection{
		Id:       uuid.New().String(),
		LongUrl:  url,
		ShortUrl: shortUrl[len(shortUrl)-1],
	}

	_, err = u.Collection.InsertOne(context.Background(), u.UrlShortnerCollection)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", os.Getenv("DOMAIN_NAME"), shortUrl[len(shortUrl)-1]), nil
}

func (u *UrlShortner) Find(shortUrl string) (string, error) {

	var found UrlShortnerCollection
	result := u.Collection.FindOne(context.Background(), bson.M{"short_url": strings.Trim(shortUrl, "/")})
	if result.Err() == mongo.ErrNoDocuments {
		return "", nil
	}
	if result.Err() != nil {
		return "", result.Err()
	}
	if err := result.Decode(&found); err != nil {
		return "", err
	}
	return found.LongUrl, nil
}
