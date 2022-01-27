package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	// CollectionAPI collection interface
	CollectionAPI interface {
		InsertOne(ctx context.Context, document interface{},
			opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
		UpdateOne(ctx context.Context, filter interface{}, update interface{},
			opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
		UpdateMany(ctx context.Context, filter interface{}, update interface{},
			opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
		FindOne(ctx context.Context, filter interface{},
			opts ...*options.FindOneOptions) *mongo.SingleResult
	}
)
