package mongo

import (
	"context"
	"github.com/dmzlingyin/utils/config"
	"github.com/dmzlingyin/utils/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"
)

func NewDatabase() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.GetString("mongo.uri")))
	if err != nil {
		panic(err)
	}
	return client.Database(config.GetString("mongo.database"))
}

func NewScope(db *mongo.Database) Scope {
	w := &scope{Database: db}
	// MongoDB 只有在 v4 之后、并且必须是 ReplicaSet 模式下才支持事务，为便于测试增加此开关
	w.txEnabled = config.GetBool("mongo.tx_enabled")
	if w.txEnabled {
		log.Info("tx_enabled: ", w.txEnabled)
	}
	return w
}

type Scope interface {
	Transact(ctx context.Context, fn func(c context.Context) error) error
}

type scope struct {
	*mongo.Database
	txEnabled bool
}

func (d *scope) Transact(ctx context.Context, fn func(c context.Context) error) error {
	if !d.txEnabled {
		return fn(ctx)
	}

	session, err := d.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	opts := options.Transaction().SetWriteConcern(writeconcern.Majority()).SetReadConcern(readconcern.Snapshot())
	_, err = session.WithTransaction(ctx, func(sc mongo.SessionContext) (any, error) {
		return nil, fn(sc)
	}, opts)
	return err
}
