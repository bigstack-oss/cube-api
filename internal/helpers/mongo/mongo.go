package mongo

import (
	"context"
	"fmt"
	"time"

	log "go-micro.dev/v5/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	CreateRecordIfNotExist = options.Update().SetUpsert(true)
)

type Client interface {
	Database(string, ...*options.DatabaseOptions) *mongo.Database
	StartSession(...*options.SessionOptions) (mongo.Session, error)
	Disconnect(context.Context) error
}

type DBClient interface {
	Collection(string, ...*options.CollectionOptions) *mongo.Collection
	ListCollectionNames(context.Context, interface{}, ...*options.ListCollectionsOptions) ([]string, error)
}

type TxnClient interface {
	WithTransaction(context.Context, func(mongo.SessionContext) (interface{}, error), ...*options.TransactionOptions) (interface{}, error)
	EndSession(context.Context)
}

type CollClient interface {
	Find(context.Context, interface{}, ...*options.FindOptions) (*mongo.Cursor, error)
	FindOne(context.Context, interface{}, ...*options.FindOneOptions) *mongo.SingleResult
	FindOneAndUpdate(context.Context, interface{}, interface{}, ...*options.FindOneAndUpdateOptions) *mongo.SingleResult
	FindOneAndDelete(context.Context, interface{}, ...*options.FindOneAndDeleteOptions) *mongo.SingleResult
	InsertOne(context.Context, interface{}, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	DeleteOne(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteMany(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	UpdateOne(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateMany(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	CountDocuments(context.Context, interface{}, ...*options.CountOptions) (int64, error)
}

type CursorClient interface {
	All(context.Context, interface{}) error
	Next(context.Context) bool
}

type Helper struct {
	Client
	Config
}

type Config struct {
	Sockets    string `validate:"required"`
	Auth       `validate:"required"`
	ReplicaSet string
	Connect    string

	Database    string
	Collection  string
	Databases   map[string]string
	Collections map[string]string

	Fetch
}

type Auth struct {
	Enable   bool
	Source   string
	Username string
	Password string
}

type Fetch struct {
	Interval time.Duration
	Number   int
	Retry    int
}

func NewHelper(conf Config) *Helper {
	h := Helper{Config: conf}
	h.SetMongoClient()
	return &h
}

func NewDefaultConf(db string) Config {
	return Config{
		Sockets:    "mongodb://0.0.0.0:27019/?directConnection=true",
		Database:   db,
		ReplicaSet: "rs0",
		// Auth: Auth{
		// 	Enable:   true,
		// 	Username: "root",
		// 	Password: "example",
		// },
		Fetch: Fetch{
			Retry: 3,
		},
	}
}

func (h *Helper) NewDBCli(db string) (DBClient, error) {
	if db == "" {
		return nil, fmt.Errorf(
			"db is nil. value: db(%s)",
			db,
		)
	}

	return h.Client.Database(db), nil
}

func (h *Helper) NewCollCli(db, coll string) (CollClient, error) {
	if db == "" || coll == "" {
		return nil, fmt.Errorf(
			"db or coll is nil or both are nil. values: db(%s); coll(%s)",
			db,
			coll,
		)
	}
	dbCli := h.Client.Database(db)

	return dbCli.Collection(coll), nil
}

func (h *Helper) NewTxnCli() (TxnClient, error) {
	s, err := h.Client.StartSession()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (h *Helper) GetQueryCursor(db, coll string, query bson.M, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	c, err := h.NewCollCli(db, coll)
	if err != nil {
		return nil, err
	}

	cursor, err := c.Find(context.Background(), query, opts...)
	if err != nil {
		return nil, err
	}

	return cursor, nil
}

func (h *Helper) Get(db, coll string, filter bson.M) (*mongo.SingleResult, error) {
	c, err := h.NewCollCli(db, coll)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := c.FindOne(ctx, filter)
	return result, nil
}

func (h *Helper) GetCount(db, coll string, filter bson.M) (int64, error) {
	c, err := h.NewCollCli(db, coll)
	if err != nil {
		return 0, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	count, err := c.CountDocuments(ctx, filter)
	defer cancel()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (h *Helper) SetMongoClient() {
	opt := options.Client()
	opt.ApplyURI(h.Sockets)

	if h.Auth.Enable {
		opt.Auth = &options.Credential{
			AuthSource: h.Auth.Source,
			Username:   h.Auth.Username,
			Password:   h.Auth.Password,
		}
	}

	if h.ReplicaSet != "" {
		opt.ReplicaSet = &h.ReplicaSet
	}

	var err error
	h.Client, err = mongo.Connect(context.Background(), opt)
	if err != nil {
		log.Errorf("err of connect mongo: %s", err.Error())
	}
}

func (h *Helper) Insert(db, coll string, data interface{}) error {
	c, err := h.NewCollCli(db, coll)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = c.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (h *Helper) UpdateOne(db, coll string, filter interface{}, data interface{}, opts ...*options.UpdateOptions) error {
	c, err := h.NewCollCli(db, coll)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = c.UpdateOne(ctx, filter, data, opts...)
	if err != nil {
		return err
	}

	return nil
}

func (h *Helper) UpdateMany(db, coll string, filter interface{}, data interface{}) error {
	c, err := h.NewCollCli(db, coll)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = c.UpdateMany(ctx, filter, data)
	if err != nil {
		return err
	}

	return nil
}

func (h *Helper) DeleteOne(db, coll string, filter interface{}) error {
	c, err := h.NewCollCli(db, coll)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = c.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (h *Helper) DeleteAll(db, coll string, filter interface{}) error {
	c, err := h.NewCollCli(db, coll)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = c.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (h *Helper) GetAllCollections(db string) ([]string, error) {
	dbCli, err := h.NewDBCli(db)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collections, err := dbCli.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	return collections, nil
}

func (h *Helper) Close() {
	err := h.Client.Disconnect(context.Background())
	if err != nil {
		log.Errorf("failed to close mongo connection: %s", err.Error())
	}
}
