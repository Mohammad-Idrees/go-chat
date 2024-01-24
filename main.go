package main

import (
	"context"
	"fmt"
	"log"
	"project/chat"
	"project/config"
	db "project/db/sqlc"
	"project/delivery"
	"project/middleware"
	"project/models"
	service "project/service/impl"
	"project/validator"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	var wg sync.WaitGroup
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("failed loading config file", err)
		return
	}

	redis, err := newRedis(config)
	if err != nil {
		log.Fatalln("failed connecting to redis", err)
		return
	}

	database, err := newDatabase(config)
	if err != nil {
		log.Fatalln("failed connecting to database", err)
		return
	}
	defer database.ConnPool.Close()

	err = runMigration(config)
	if err != nil {
		log.Fatalln("failed to run migration ", err)
		return
	}

	router := gin.Default()
	err = validator.Init()
	if err != nil {
		log.Println(err)
	}

	// Init Hub
	hub := chat.InitHub(&wg, config, redis.Client)

	repository := db.NewRepository(database)
	tokenService := service.ConfigureTokenService(config, repository)
	userService := service.ConfigureUserService(config, repository, tokenService)

	authMiddleware := middleware.AuthMiddleware(tokenService)

	delivery.ConfigureTokenHandler(&router.RouterGroup, tokenService)
	delivery.ConfigureUserHandler(&router.RouterGroup, authMiddleware, userService)
	delivery.ConfigureWSHandler(&router.RouterGroup, hub, userService, repository)

	err = router.Run(config.Server.Address)
	if err != nil {
		log.Fatalln("failed to start server")
		return
	}
}

func newRedis(config *config.StartupConfig) (*models.Redis, error) {
	addr := fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port)
	rds := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.Redis.Password,
		DB:       0, // Default DB
	})

	_, err := rds.Ping(context.Background()).Result()
	if err != nil {
		log.Println("failed connecting to redis")
		return nil, err
	}

	return &models.Redis{
		Client: rds,
	}, nil
}

func runMigration(cfg *config.StartupConfig) error {
	dbSource := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", cfg.Database.Type, cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	migrationURL := cfg.Migration.MigrationURL
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Println("failed to create migration instance", dbSource, err)
		return err
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Println("failed to run migrate up", err)
		return err
	}

	log.Println("db migration successful!!")
	return nil
}

func newDatabase(config *config.StartupConfig) (*models.Database, error) {
	cfg := config.Database
	connString := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", cfg.Type, cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	connPool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Println("failed connecting to postgres", err)
		return nil, err
	}

	log.Println("connected to postgres!!")
	return &models.Database{ConnPool: connPool}, nil
}
