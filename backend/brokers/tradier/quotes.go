package tradier

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/tidwall/gjson"
)

type Quote struct {
	Type             string  `json:"type"`
	Symbol           string  `json:"symbol"`
	Size             int     `json:"size"`
	Last             float64 `json:"last"`
	Open             float64 `json:"open"`
	High             float64 `json:"high"`
	Low              float64 `json:"low"`
	Bid              float64 `json:"bid"`
	Ask              float64 `json:"ask"`
	Close            float64 `json:"close"`
	PrevClose        float64 `json:"prevclose"`
	Change           float64 `json:"change"`
	ChangePercentage float64 `json:"change_percentage"`
	Volume           int     `json:"volume"`
	AverageVolume    int     `json:"average_volume"`
	LastVolume       int     `json:"last_volume"`
	Description      string  `json:"description"`
}

//
// GetQuotes - Get a quote.
//
func (t *Api) GetQuotes(symbols []string) ([]types.Quote, error) {
	// No symbols, no quotes.
	if len(symbols) == 0 {
		return nil, nil
	}

	var quotes []Quote

	// Setup http client
	client := &http.Client{}

	// Get url to api
	apiUrl := apiBaseUrl

	if t.Sandbox {
		apiUrl = sandBaseUrl
	}

	// Setup api request
	req, _ := http.NewRequest("GET", apiUrl+"/markets/quotes", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprint("Bearer ", t.ApiKey))

	// Build query string
	q := req.URL.Query()
	q.Add("symbols", strings.Join(symbols, ","))
	req.URL.RawQuery = q.Encode()

	// Do API request
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	// Close Body
	defer res.Body.Close()

	//fmt.Println(res.Header.Get("X-Ratelimit-Allowed"))
	//fmt.Println(res.Header.Get("X-Ratelimit-Used"))
	//fmt.Println(res.Header.Get("X-Ratelimit-Available"))
	//fmt.Println(res.Header.Get("X-Ratelimit-Expiry"))

	// Make sure the api responded with a 200
	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprint("GetQuotes API did not return 200, It returned ", res.StatusCode))
	}

	// Read the data we got.
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	// Did we get one quote or many?
	if strings.Contains(string(body), "[") {

		vo := gjson.Get(string(body), "quotes.quote")

		err := json.Unmarshal([]byte(vo.String()), &quotes)

		if err != nil {
			return nil, fmt.Errorf("%s: %s", err, body)
		}

	} else {

		var res map[string]map[string]Quote

		err := json.Unmarshal(body, &res)

		if err != nil {
			return nil, fmt.Errorf("%s: %s", err, body)
		}

		quote, ok := res["quotes"]["quote"]

		if !ok {
			return nil, nil
		}

		quotes = []Quote{quote}

	}

	// Normalize the data
	realQuotes := []types.Quote{}

	for _, row := range quotes {

		realQuotes = append(realQuotes, types.Quote{
			Type:             row.Type,
			Symbol:           row.Symbol,
			Size:             row.Size,
			Last:             row.Last,
			Open:             row.Open,
			High:             row.High,
			Low:              row.Low,
			Bid:              row.Bid,
			Ask:              row.Ask,
			Close:            row.Close,
			PrevClose:        row.PrevClose,
			Change:           row.Change,
			ChangePercentage: row.ChangePercentage,
			Volume:           row.Volume,
			AverageVolume:    row.AverageVolume,
			LastVolume:       row.LastVolume,
			Description:      row.Description,
		})

	}

	// Return happy
	return realQuotes, nil
}

/* End File */
