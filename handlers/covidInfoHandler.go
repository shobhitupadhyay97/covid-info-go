package handlers

import (
	"context"
	"covid-info-go/database"
	"covid-info-go/external"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
)

type CovidStateInfo struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"state_code" bson:"state_code"`
	ActiveCase int                `json:"active_case" bson:"active_case"`
	LastSync   time.Time          `json:"sync_time" bson:"sync_time"`
}

type CovidInfoHandler struct {
	Col   database.CollectionAPI
	Redis database.RedisAPI
}

func (h *CovidInfoHandler) InsertCovidInfo(c echo.Context) error {
	covidInfo := external.GetCovidData()
	today := time.Now()
	for _, stateData := range covidInfo {
		filter := bson.D{{Key: "state_code", Value: stateData.StateCode}}
		update := bson.D{{Key: "$set", Value: bson.D{
			{Key: "active_case", Value: stateData.ActiveCase},
			{Key: "sync_time", Value: today}}}}
		opts := options.Update().SetUpsert(true)
		_, err := h.Col.UpdateOne(context.TODO(), filter, update, opts)
		if err != nil {
			fmt.Println(err)
		}
	}
	return c.JSON(http.StatusCreated, covidInfo)
}

func (h *CovidInfoHandler) InsertCovidInfoV2(c echo.Context) error {
	covidInfo := external.GetCovidData()
	today := time.Now()
	g, _ := errgroup.WithContext(context.Background())
	for _, stateData := range covidInfo {
		filter := bson.D{{Key: "state_code", Value: stateData.StateCode}}
		update := bson.D{{Key: "$set", Value: bson.D{
			{Key: "active_case", Value: stateData.ActiveCase},
			{Key: "sync_time", Value: today}}}}
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

func (h *CovidInfoHandler) GetCovidInfo(c echo.Context) error {
	// "26.45918", "80.329288"
	lat := c.QueryParam("lat")
	long := c.QueryParam("long")
	locationInfo := external.GetLoacationInfo(lat, long)
	var stateInfo CovidStateInfo
	filter := bson.D{{Key: "state_code", Value: locationInfo.StateCode}}
	if err := h.Col.FindOne(context.Background(), filter).Decode(&stateInfo); err != nil {
		log.Fatal(err)
	}
	return c.JSON(http.StatusOK, stateInfo)
}
