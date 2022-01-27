package external

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Total struct {
	Confirmed int64 `json:"confirmed"`
}
type covidRespData struct {
	AN struct {
		Total Total `json:"total"`
	} `json:"AN"`
	AP struct {
		Total Total `json:"total"`
	} `json:"AP"`
	AR struct {
		Total Total `json:"total"`
	} `json:"AR"`
	AS struct {
		Total Total `json:"total"`
	} `json:"AS"`
	BR struct {
		Total Total `json:"total"`
	} `json:"BR"`
	CH struct {
		Total Total `json:"total"`
	} `json:"CH"`
	CT struct {
		Total Total `json:"total"`
	} `json:"CT"`
	DL struct {
		Total Total `json:"total"`
	} `json:"DL"`
	DN struct {
		Total Total `json:"total"`
	} `json:"DN"`
	GA struct {
		Total Total `json:"total"`
	} `json:"GA"`
	GJ struct {
		Total Total `json:"total"`
	} `json:"GJ"`
	HP struct {
		Total Total `json:"total"`
	} `json:"HP"`
	HR struct {
		Total Total `json:"total"`
	} `json:"HR"`
	JH struct {
		Total Total `json:"total"`
	} `json:"JH"`
	JK struct {
		Total Total `json:"total"`
	} `json:"JK"`
	KA struct {
		Total Total `json:"total"`
	} `json:"KA"`
	KL struct {
		Total Total `json:"total"`
	} `json:"KL"`
	LA struct {
		Total Total `json:"total"`
	} `json:"LA"`
	LD struct {
		Total Total `json:"total"`
	} `json:"LD"`
	MH struct {
		Total Total `json:"total"`
	} `json:"MH"`
	ML struct {
		Total Total `json:"total"`
	} `json:"ML"`
	MN struct {
		Total Total `json:"total"`
	} `json:"MN"`
	MP struct {
		Total Total `json:"total"`
	} `json:"MP"`
	MZ struct {
		Total Total `json:"total"`
	} `json:"MZ"`
	NL struct {
		Total Total `json:"total"`
	} `json:"NL"`
	OR struct {
		Total Total `json:"total"`
	} `json:"OR"`
	PB struct {
		Total Total `json:"total"`
	} `json:"PB"`
	PY struct {
		Total Total `json:"total"`
	} `json:"PY"`
	RJ struct {
		Total Total `json:"total"`
	} `json:"RJ"`
	SK struct {
		Total Total `json:"total"`
	} `json:"SK"`
	TG struct {
		Total Total `json:"total"`
	} `json:"TG"`
	TN struct {
		Total Total `json:"total"`
	} `json:"TN"`
	TR struct {
		Total Total `json:"total"`
	} `json:"TR"`
	TT struct {
		Total Total `json:"total"`
	} `json:"TT"`
	UP struct {
		Total Total `json:"total"`
	} `json:"UP"`
	UT struct {
		Total Total `json:"total"`
	} `json:"UT"`
	WB struct {
		Total Total `json:"total"`
	} `json:"WB"`
}

func GetCovidData() covidRespData {
	url := "https://data.covid19india.org/v4/min/data.min.json"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	var data covidRespData
	json.Unmarshal(body, &data)
	return data
}
