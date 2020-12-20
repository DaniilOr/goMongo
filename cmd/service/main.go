package main

import (
	"context"
	"github.com/DaniilOr/goMongo/cmd/service/app"
	"github.com/DaniilOr/goMongo/pkg/cache"
	"github.com/DaniilOr/goMongo/pkg/payments"
	"github.com/DaniilOr/goMongo/pkg/security"
	"github.com/go-chi/chi"
	"github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net"
	"net/http"
	"os"
)

const (
	defaultPort        = "9999"
	defaultHost        = "0.0.0.0"
	defaultClientsDB   = "db"
	defaultClients8DSN = "postgres://app:pass@localhost:5432/" + defaultClientsDB
	defaultMongoDB     = "predictions"
	defaultMongoDSN    = "mongodb://app:pass@localhost:27017/" + defaultMongoDB
	defaultCacheDSN = "redis://localhost:6379/0"
)

func main() {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}

	clientsDsn, ok := os.LookupEnv("CLIENTS_DSN")
	if !ok {
		clientsDsn = defaultClients8DSN
	}
	mongoDsn, ok := os.LookupEnv("Mongo_DSN")
	if !ok {
		mongoDsn = defaultMongoDSN
	}
	mongoDB, ok := os.LookupEnv("Mongo_DB")
	if !ok {
		mongoDB = defaultMongoDB
	}
	cacheDSN, ok := os.LookupEnv("APP_CACHE_DSN")
	if !ok {
		cacheDSN = defaultCacheDSN
	}
	if err := execute(net.JoinHostPort(host, port), clientsDsn, mongoDsn, mongoDB, cacheDSN); err != nil {
		os.Exit(1)
	}
}

func execute(addr string, clientsDsn string, mongoDsn string, mongDB string, cacheDSN string) error {
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, clientsDsn)
	if err != nil {
		log.Print(err)
		return err
	}
	cachePool := &redis.Pool{
		DialContext: func(ctx context.Context) (redis.Conn, error) {
			return redis.DialURL(cacheDSN)
		},
	}
	defer pool.Close()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDsn))
	if err != nil {
		log.Print(err)
		return err
	}

	database := client.Database(mongDB)
	cacheSvc := cache.Service{cachePool}
	securitySvc := security.NewService(pool)
	paymentsSvc := payments.NewService(database)
	router := chi.NewRouter()
	application := app.NewServer(securitySvc, paymentsSvc, router, &cacheSvc)
	err = application.Init()
	if err != nil {
		log.Print(err)
		return err
	}

	server := &http.Server{
		Addr:    addr,
		Handler: application,
	}
	return server.ListenAndServe()
}
