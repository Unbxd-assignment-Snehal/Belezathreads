package model

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const UnbxdAPIEndpoint = "https://search.unbxd.io/fb853e3332f2645fac9d71dc63e09ec1/demo-unbxd700181503576558/search"

type UnbxdResponse struct {
	Success  string                 `json:"success"`
	Response UnbxdResponseData      `json:"response"`
	Debug    map[string]interface{} `json:"debug,omitempty"`
}

type UnbxdResponseData struct {
	NumberOfProducts int                    `json:"numberOfProducts"`
	Products         []map[string]interface{} `json:"products"`
}

func SearchUnbxd(q, pageno, sort, fields string) (*UnbxdResponse, error) {
	unbxdURL, err := buildUnbxdURL(q, pageno, sort, fields)
	if err != nil {
		return nil, err
	}

	response, err := http.Get(unbxdURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var responseBody UnbxdResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
	
}

func buildUnbxdURL(q, pageno, sort, fields string) (string, error) {
	baseURL, err := url.Parse(UnbxdAPIEndpoint)
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Set("q", q)
	params.Set("rows", "10")
	params.Set("start", pageno)

	if sort == "asc" {
		sort = "price asc"
	}
	if sort == "desc" {
		sort = "price desc"
	}
	params.Set("sort", sort)
	params.Set("fields", fields)

	baseURL.RawQuery = params.Encode()

	finalURL := baseURL.String()
	fmt.Printf("%v", finalURL)
	return finalURL, nil
}
