package database

import (
	"context"

	"github.com/gsiragusa/short-to-me/config"
	"github.com/gsiragusa/short-to-me/shortener"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const CollShortUrls = "short_urls"

type Client struct {
	mc     *mongo.Client
	db     *mongo.Database
	config *config.AppConfig
}

func NewMongoClient(config *config.AppConfig) (*Client, error) {
	// connect to MongoDB
	clientOptions := options.Client().ApplyURI(config.MongoUri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return &Client{
		mc:     client,
		config: config,
		db:     client.Database(config.MongoDbName),
	}, nil
}

// TODO: ensure index on url
func (c *Client) FindUrl(ctx context.Context, url string) (*shortener.ModelShorten, error) {
	u := &shortener.ModelShorten{}
	collection := c.db.Collection(CollShortUrls)
	if err := collection.FindOne(ctx, bson.M{"url": url}).Decode(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (c *Client) FindById(ctx context.Context, id string) (*shortener.ModelShorten, error) {
	u := &shortener.ModelShorten{}
	collection := c.db.Collection(CollShortUrls)
	if err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (c *Client) StoreUrl(ctx context.Context, document interface{}) error {
	collection := c.db.Collection(CollShortUrls)
	_, err := collection.InsertOne(ctx, document)
	return err
}

func (c *Client) DeleteById(ctx context.Context, id string) error {
	collection := c.db.Collection(CollShortUrls)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (c *Client) IncrementCount(ctx context.Context, id string) (*shortener.ModelShorten, error) {
	u := &shortener.ModelShorten{}
	collection := c.db.Collection(CollShortUrls)
	increment := bson.M{"$inc": bson.M{"count": 1}}
	if err := collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, increment).Decode(u); err != nil {
		return nil, err
	}
	return u, nil
}
