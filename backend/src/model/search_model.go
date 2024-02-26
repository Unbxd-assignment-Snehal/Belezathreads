// model/filter_model.go

package model

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const UnbxdAPIEndpoint = "https://search.unbxd.io/fb853e3332f2645fac9d71dc63e09ec1/demo-unbxd700181503576558/search"

func SearchUnbxd(q, pageno, sort, fields string) ([]map[string]interface{}, error) {
	// Construct the URL for the Unbxd API with query parameters
	unbxdURL, err := buildUnbxdURL(q, pageno, sort, fields)
	if err != nil {
		return nil, err
	}

	// Make a GET request to the Unbxd API
	response, err := http.Get(unbxdURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Decode the response body into a map
	var responseBody map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	// Extract the 'products' field from the response body
	products, ok := responseBody["response"].(map[string]interface{})["products"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format from Unbxd API")
	}

	// Convert the list of products to a slice of maps
	var productList []map[string]interface{}
	for _, product := range products {
		if productMap, ok := product.(map[string]interface{}); ok {
			productList = append(productList, productMap)
		}
	}

	return productList, nil
}

func buildUnbxdURL(q, pageno, sort, fields string) (string, error) {
	// Parse the Unbxd API endpoint
	baseURL, err := url.Parse(UnbxdAPIEndpoint)
	if err != nil {
		return "", err
	}

	// Prepare the query parameters
	params := url.Values{}
	params.Set("q", q)
	params.Set("rows", "10")
	params.Set("start", pageno)
	params.Set("sorting", sort)
	params.Set("fields", fields) // Add the 'fields' parameter

	// Add the query parameters to the URL
	baseURL.RawQuery = params.Encode()

	// Construct the final URL as a string
	finalURL := baseURL.String()
	fmt.Printf("%v", finalURL)
	return finalURL, nil
}
