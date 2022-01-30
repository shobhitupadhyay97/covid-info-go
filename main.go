package main

import (
	"context"
	"covid-info-go/handlers"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	c      *mongo.Client
	db     *mongo.Database
	col    *mongo.Collection
	client *redis.Client
	err    error
)

func init() {
	godotenv.Load()
	mongoOption := options.Client().ApplyURI(os.Getenv("DB_URI"))
	c, err = mongo.Connect(context.Background(), mongoOption)
	if err != nil {
		log.Fatalf("Can not connect to DB: %v", err)
	}
	db = c.Database(os.Getenv("DB_NAME"))
	col = db.Collection("state_data")
	client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URI"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

}

func main() {
	godotenv.Load()
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	h := handlers.CovidInfoHandler{Col: col, Redis: client}

	e.GET("/knockknock", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello there!!")
	})

	e.POST("/covid-info", h.InsertCovidInfo)
	e.POST("/covid-info-goroutine", h.InsertCovidInfoV2)

	e.GET("/covid-info", h.GetCovidInfo)

	e.Logger.Print("Listening to port %s", os.Getenv("PORT"))
	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))))
}
