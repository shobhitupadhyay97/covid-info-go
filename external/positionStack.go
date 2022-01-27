package external

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Loacation struct {
	StateCode string `json:"region_code"`
}

func GetLoacationInfo(lat string, long string) Loacation {
	url := "http://api.positionstack.com/v1/reverse?access_key=" + os.Getenv("POSITION_STACK_ACCESS_KEY") + "&query=" + lat + "," + long
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Can't get reverse geocode %v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error parsing resp %v", err)
	}
	defer resp.Body.Close()
	var data map[string][]Loacation
	json.Unmarshal(body, &data)
	locationRespList, exist := data["data"]
	if exist {
		if len(locationRespList) != 0 {
			return locationRespList[0]
		}
	}
	return Loacation{StateCode: ""}
}
