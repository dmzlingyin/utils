package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindOne 查找单条数据
func FindOne[T any](ctx context.Context, collection *mongo.Collection, filter any) (result T, err error) {
	if filter == nil {
		filter = bson.M{}
	}
	err = collection.FindOne(ctx, filter).Decode(&result)
	return
}

// FetchByPage 分页获取数据
func FetchByPage[T any](ctx context.Context, collection *mongo.Collection, pageIndex, pageSize int64, filter, sorter any) (results T, err error) {
	var (
		cur  *mongo.Cursor
		opts = options.Find().SetSkip(pageSize * (pageIndex - 1)).SetLimit(pageSize)
	)
	if filter == nil {
		filter = bson.M{}
	}
	if sorter != nil {
		opts.SetSort(sorter)
	}
	cur, err = collection.Find(ctx, filter, opts)
	if err != nil {
		return
	}
	defer cur.Close(ctx)

	err = cur.All(ctx, &results)
	return
}

// Fetch 获取全部数据
func Fetch[T any](ctx context.Context, collection *mongo.Collection, filter, sorter any) (results T, err error) {
	var (
		cur  *mongo.Cursor
		opts = options.Find()
	)
	if filter == nil {
		filter = bson.M{}
	}
	if sorter != nil {
		opts.SetSort(sorter)
	}
	cur, err = collection.Find(ctx, filter, opts)
	if err != nil {
		return
	}
	defer cur.Close(ctx)

	err = cur.All(ctx, &results)
	return
}
