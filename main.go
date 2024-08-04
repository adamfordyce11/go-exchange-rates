package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type ApiV4Response struct {
	Provider           string             `json:"provider"`
	Warnings           string             `json:"WARNING_UPGRADE_TO_V6"`
	Terms              string             `json:"terms`
	BaseCode           string             `json:"base"`
	Date               string             `json:"date"`
	TimeLastUpdateUnix int                `json:"time_last_updated"`
	ConversionRates    map[string]float64 `json:"rates"`
}

type ApiV6Response struct {
	Result             string             `json:"result"`
	Documentation      string             `json:"documentation"`
	TermsOfUse         string             `json:"terms_of_use"`
	TimeLastUpdateUnix int                `json:"time_last_update_unix"`
	TimeLastUpdateUtc  string             `json:"time_last_update_utc"`
	TimeNextUpdateUnix int                `json:"time_next_update_unix"`
	TimeNextUpdateUtc  string             `json:"time_next_update_utc"`
	BaseCode           string             `json:"base_code"`
	ConversionRates    map[string]float64 `json:"conversion_rates"`
}

type Currency struct {
	Source string
	Dest   string
}

type CurrencyRequest struct {
	Url     string
	Version string
	Request Currency
}

func (c *CurrencyRequest) readUrl() (*http.Response, error) {
	// Debug
	//fmt.Printf("Url: %s/latest/%s\n", c.Url, c.Request.Source)

	resp, err := http.Get(c.Url + "/latest/" + c.Request.Source)
	if err != nil {
		return &http.Response{}, fmt.Errorf("Error: %v ", err)
	}
	return resp, nil
}

func (c *CurrencyRequest) GetConversionRate() float64 {
	// Read the URL from the Currency API
	resp, err := c.readUrl()
	if err != nil {
		fmt.Errorf("Error reading url: %v ", err)
		return 99.9999
	}

	// Check the status code was acceptable, otherwise return nil
	if resp.StatusCode != 200 {
		fmt.Errorf("Error: %v ", resp.Status)
		return 99.9999
	}

	// Process the request for the version6 API
	if c.Version == "v6" {
		// Unmarshall the response body into the ApiResponse struct
		ApiResponse := &ApiV6Response{}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Errorf("Error reading response body: %v ", err)
			return 99.9996
		}
		if err := json.Unmarshal(body, ApiResponse); err != nil {
			fmt.Errorf("Error unmarshalling response body: %v ", err)
			return 99.9997
		}
		// Return the conversion rate
		return ApiResponse.ConversionRates[c.Request.Dest]
	} else {
		// Process the request for the version4 API
		// Unmarshall the response body into the ApiResponse struct
		ApiResponse := &ApiV4Response{}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Errorf("Error reading response body: %v ", err)
			return 99.9998
		}
		if err := json.Unmarshal(body, ApiResponse); err != nil {
			fmt.Errorf("Error unmarshalling response body: %v ", err)
			return 99.9999
		}
		// Return the conversion rate
		return ApiResponse.ConversionRates[c.Request.Dest]
	}
}

func NewCurrency(ApiKey, Version string) *CurrencyRequest {

	if Version == "" {
		Version = "v4"
	}

	Url := "https://api.exchangerate-api.com/v4"
	if Version == "v6" {
		Url = "https://v6.exchangerate-api.com/v6/" + ApiKey
	}

	Currency := Currency{
		Source: "GBP",
		Dest:   "USD",
	}

	return &CurrencyRequest{
		Url:     Url,
		Version: Version,
		Request: Currency,
	}
}

func getEnv(key, dflt string) string {
	if os.Getenv(key) != "" {
		return os.Getenv(key)
	}
	os.Setenv(key, dflt)
	return dflt
}

/*func main() {
ApiKey := getEnv("API_KEY", "API-KEY")
ApiVersion := getEnv("API_VERSION", "v4")
CurrencySource := getEnv("CURRENCY_SOURCE", "GBP")
CurrencyDest := getEnv("CURRENCY_DEST", "USD")
currencyRequest := NewCurrency(ApiKey, ApiVersion)
currencyRequest.Request.Source = CurrencySource
currencyRequest.Request.Dest = CurrencyDest
rate := currencyRequest.GetConversionRate()
fmt.Printf("Rate: %f\n", rate)

}*/
