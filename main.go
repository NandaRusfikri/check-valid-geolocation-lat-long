package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Results struct {
	Components Component `json:"components"`
	Confidence int    `json:"confidence"`
	Formatted  string `json:"formatted"`
	Geometry   struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"geometry"`
}
type Component struct {
	ISO31661Alpha2 string   `json:"ISO_3166-1_alpha-2"`
	ISO31661Alpha3 string   `json:"ISO_3166-1_alpha-3"`
	ISO31662       []string `json:"ISO_3166-2"`
	Category       string   `json:"_category"`
	Type           string   `json:"_type"`
	Attraction     string   `json:"attraction"`
	Continent      string   `json:"continent"`
	Country        string   `json:"country"`
	CountryCode    string   `json:"country_code"`
	Postcode       string   `json:"postcode"`
	Road           string   `json:"road"`
	State          string   `json:"state"`
	StateCode      string   `json:"state_code"`
	Subdistrict    string   `json:"subdistrict"`
	Village        string   `json:"village"`
}

type ResponseAPI struct {
	ID string `json:"id"`
	ATMID string `json:"atm_id"`
	Latitude string `json:"latitude"`
	Longitude string `json:"longitude"`
	Documentation string `json:"documentation"`
	Licenses      []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"licenses"`
	Results []Results `json:"results"`
	Status struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"status"`
	StayInformed struct {
		Blog    string `json:"blog"`
		Twitter string `json:"twitter"`
	} `json:"stay_informed"`
	Thanks    string `json:"thanks"`
	Timestamp struct {
		CreatedHttp string `json:"created_http"`
		CreatedUnix int    `json:"created_unix"`
	} `json:"timestamp"`
	TotalResults int `json:"total_results"`
}
type SchemaDatasource struct {
	Id string
	ATMId string
	Latitude string
	Longitude string
}

func ExtrackData(data [][]string) []SchemaDatasource {
	var shoppingList []SchemaDatasource
	for i, line := range data {
		if i > 0 { // omit header line
			var rec SchemaDatasource
			for j, field := range line {
				if j == 0 {
					rec.Id = field
				} else if j == 1 {
					rec.ATMId = field
				}else if j ==2 {
					rec.Latitude = field
				} else if j ==3 {
					rec.Longitude = field

				}
			}
			shoppingList = append(shoppingList, rec)
		}
	}
	return shoppingList
}
type CallAPIDto struct {
	Method       string
	Url          string
	ContentType  string
	Headers      map[string]interface{}
	BodyRequest  string
	BodyResponse string
	HttpCode     int
}

func (d *CallAPIDto) Validate() error {
	if d.Method == "" {
		return errors.New("method required")
	}
	if d.Url == "" {
		return errors.New("url required")
	}

	return nil
}
func CallAPI(data *CallAPIDto) (err error) {
	//Timenow := time.Now()

	if err = data.Validate(); err != nil {
		return err
	}

	client := &http.Client{}
	var request *http.Request
	if data.BodyRequest != "" {
		request, err = http.NewRequest(data.Method, data.Url, bytes.NewBuffer([]byte(data.BodyRequest)))
		if err != nil {
			return err
		}
	} else {
		request, err = http.NewRequest(data.Method, data.Url, nil)
		if err != nil {
			return err
		}
	}

	request.Header.Set("Content-Type", data.ContentType)

	if data.Headers != nil && len(data.Headers) > 0 {
		for key, header := range data.Headers {
			request.Header.Set(key, header.(string))
		}
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			fmt.Println("error on close response body: ", err)
		}
	}()

	bodyResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	data.BodyResponse = string(bodyResponse)
	data.HttpCode = response.StatusCode


	return nil
}


func GetAPI (lat string,long string) (ResponseAPI){
	url := fmt.Sprintf("https://api.opencagedata.com/geocode/v1/json?q=%v+%v&key=03c48dae07364cabb7f121d8c1519492&no_annotations=1&language=en",lat,long)

	headers := make(map[string]interface{})

	request := CallAPIDto{
		Method:      "GET",
		Url:         url,
		ContentType: "application/json",
		Headers:     headers,
	}

	err := CallAPI(&request)
	if err != nil {
		fmt.Println("err CallAPI",err)
	}
	//fmt.Println("request.BodyResponse",request.BodyResponse)

	//var rawResponse map[string]interface{}
	rawResponse := ResponseAPI{}
	err = json.Unmarshal([]byte(request.BodyResponse), &rawResponse)


	for _, result := range rawResponse.Results {
		fmt.Printf("Country: %+v State:%v PosCode: %v subdistrict:%v Village:%v\n", result.Components.Country,
			result.Components.State,result.Components.Postcode,result.Components.Subdistrict,result.Components.Village)

	}


	return rawResponse

}
func main()  {
	f, err := os.Open("datasource.csv")
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// convert records to array of structs
	ListExtrackData := ExtrackData(data)

	var ListData []ResponseAPI

	fmt.Println("Total Row Datasource : ",len(ListExtrackData))
	fmt.Println("Wait Generate Data")
	for _, item := range ListExtrackData {
		Hasil:= GetAPI(item.Latitude,item.Longitude)
		Hasil.ATMID = item.ATMId
		Hasil.ID = item.Id
		Hasil.Latitude = item.Latitude
		Hasil.Longitude = item.Longitude
		ListData = append(ListData,Hasil)

	}

	timeNow := time.Now()
	path := fmt.Sprintf("output")
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		fmt.Println("err ",err)
	}
	fileName := fmt.Sprintf("output-%d%02d%02d-%v.csv", timeNow.Year(), timeNow.Month(), timeNow.Day(),timeNow.Unix())

	fullpath := filepath.Join(path, fileName)
	err = os.Remove(fullpath)
	if err != nil {
		fmt.Println("err ",err)
		//return
	}

	csvFile, err := os.Create(fullpath)
	if err != nil {
		log.Println(err)
	}

	if err != nil {

		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(csvFile)

	var header []string
	header = append(header, "NO")
	header = append(header, "ATM ID")
	header = append(header, "Latitude")
	header = append(header, "Longitude")
	header = append(header, "Country")
	header = append(header, "Postcode")
	header = append(header, "State")
	header = append(header, "Subdistrict")
	header = append(header, "Village")
	header = append(header, "Road")
	header = append(header, "Formated")
	header = append(header, "URL Google Map")

	if err := w.Write(header); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	for _, item := range ListData {

		var row []string
		row = append(row, item.ID)
		row = append(row, item.ATMID)
		row = append(row, item.Latitude)
		row = append(row, item.Longitude)
		for _, result := range item.Results {
			row = append(row, result.Components.Country)
			row = append(row, result.Components.Postcode)
			row = append(row, result.Components.State)
			row = append(row, result.Components.Subdistrict)
			row = append(row, result.Components.Village)
			row = append(row, result.Components.Road)
			row = append(row, result.Formatted)
		}
		row = append(row, fmt.Sprintf("https://www.google.co.id/maps/@%v,%v,19z",item.Latitude,item.Longitude))  // ref id


		if err := w.Write(row); err != nil {
			log.Fatalln("error writing record to file", err)
		}

	}
	if err := w.Write([]string{""}); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	w.Flush()
	csvFile.Close()


	fmt.Println("Done Generate Data")

}

