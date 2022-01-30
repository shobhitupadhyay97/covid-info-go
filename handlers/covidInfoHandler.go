package handlers

import (
	"context"
	"covid-info-go/database"
	"covid-info-go/external"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
)

type CovidStateInfo struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" redis:"na_idme"`
	Name       string             `json:"state_code" bson:"state_code" redis:"state_code"`
	ActiveCase int                `json:"active_case" bson:"active_case" redis:"active_case"`
	LastSync   time.Time          `json:"sync_time" bson:"sync_time" redis:"sync_time"`
	TotalCases int                `json:"total_case" bson:"total_case" redis:"total_case"`
}

type CovidInfoHandler struct {
	Col   database.CollectionAPI
	Redis *redis.Client
}

func (h *CovidInfoHandler) InsertCovidInfo(c echo.Context) error {
	covidInfo, totalCase := external.GetCovidData()
	today := time.Now()
	for _, stateData := range covidInfo {
		filter := bson.D{{Key: "state_code", Value: stateData.StateCode}}
		update := bson.D{{Key: "$set", Value: bson.D{
			{Key: "active_case", Value: stateData.ActiveCase},
			{Key: "sync_time", Value: today},
			{Key: "total_case", Value: totalCase}}}}
		opts := options.Update().SetUpsert(true)
		_, err := h.Col.UpdateOne(context.TODO(), filter, update, opts)
		if err != nil {
			fmt.Println(err)
		}
	}
	return c.JSON(http.StatusCreated, covidInfo)
}

func (h *CovidInfoHandler) InsertCovidInfoV2(c echo.Context) error {
	covidInfo, totalCase := external.GetCovidData()
	today := time.Now()
	g, _ := errgroup.WithContext(context.Background())
	for _, stateData := range covidInfo {
		filter := bson.D{{Key: "state_code", Value: stateData.StateCode}}
		update := bson.D{{Key: "$set", Value: bson.D{
			{Key: "active_case", Value: stateData.ActiveCase},
			{Key: "sync_time", Value: today},
			{Key: "total_case", Value: totalCase}}}}
		opts := options.Update().SetUpsert(true)
		g.Go(func() error {
			_, err := h.Col.UpdateOne(context.TODO(), filter, update, opts)
			if err != nil {
				fmt.Println(err)
			}
			return nil
		})
	}
	g.Wait()
	return c.JSON(http.StatusCreated, covidInfo)
}

func getStateCode(ctx context.Context, lat string, long string, redisClient redis.Client) string {
	latLongKey := lat + long + "STATE_MAP"
	stateCode, err := redisClient.Get(ctx, latLongKey).Result()
	if err == redis.Nil {
		locationInfo := external.GetLoacationInfo(lat, long)
		err := redisClient.Set(ctx, latLongKey, locationInfo.StateCode, 30*time.Minute).Err()
		if err != nil {
			log.Fatalf("Can not connect to redis: %v", err)
		}
		stateCode = locationInfo.StateCode
	} else if err != nil {
		log.Fatalf("Error while getting lat, long to state code mapping %v", err)
	}
	return stateCode
}

func getStateCovidData(ctx context.Context, stateCode string, redisClient redis.Client, mongoCollection database.CollectionAPI) map[string]interface{} {
	fieldsMap := make(map[string]interface{})
	val, err := redisClient.HGetAll(ctx, stateCode).Result()
	if err != nil {
		log.Fatalf("Error while getting state covid data from redis %v", err)
	}
	_, exists := val["state_code"]
	if exists {
		fieldsMap["_id"] = val["_id"]
		fieldsMap["state_code"] = val["state_code"]
		fieldsMap["active_case"] = val["active_case"]
		fieldsMap["sync_time"] = val["sync_time"]
		fieldsMap["total_case"] = val["total_case"]
		return fieldsMap
	}
	var stateInfo CovidStateInfo
	filter := bson.D{{Key: "state_code", Value: stateCode}}
	if err := mongoCollection.FindOne(ctx, filter).Decode(&stateInfo); err != nil {
		log.Fatal(err)
	}
	fieldsMap["_id"] = stateInfo.ID.Hex()
	fieldsMap["state_code"] = stateInfo.Name
	fieldsMap["active_case"] = stateInfo.ActiveCase
	fieldsMap["sync_time"] = stateInfo.LastSync.String()
	fieldsMap["total_case"] = stateInfo.TotalCases
	err = redisClient.HMSet(context.Background(), stateInfo.Name, fieldsMap).Err()
	if err != nil {
		log.Fatalf("Can not connect to redis: %v", err)
	}
	return fieldsMap
}

func (h *CovidInfoHandler) GetCovidInfo(c echo.Context) error {
	lat := c.QueryParam("lat")
	long := c.QueryParam("long")
	ctx := context.Background()
	stateCode := getStateCode(ctx, lat, long, *h.Redis)
	stateInfo := getStateCovidData(ctx, stateCode, *h.Redis, h.Col)
	return c.JSON(http.StatusOK, stateInfo)
}
