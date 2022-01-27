package external

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Total struct {
	Confirmed int32 `json:"confirmed"`
}
type RawStateData struct {
	Total Total `json:"total"`
}
type StateData struct {
	StateCode  string
	ActiveCase int32
}

func GetCovidData() []StateData {
	toalCase := int32(0)
	resp, err := http.Get("https://data.covid19india.org/v4/min/data.min.json")
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	var data map[string]RawStateData
	json.Unmarshal(body, &data)
	var listResp []StateData
	for state, data := range data {
		toalCase = toalCase + data.Total.Confirmed
		listResp = append(
			listResp, StateData{StateCode: state, ActiveCase: data.Total.Confirmed})
	}
	listResp = append(
		listResp, StateData{StateCode: "total", ActiveCase: toalCase})
	return listResp
}
