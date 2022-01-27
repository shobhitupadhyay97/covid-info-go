package handlers

import (
	"covid-info-go/external"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CovidInfoHandler struct {
}

func InsertCovidInfo(c echo.Context) error {
	covidInfo := external.GetCovidData()
	return c.JSON(http.StatusCreated, covidInfo)
}
