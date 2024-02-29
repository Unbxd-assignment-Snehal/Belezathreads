package services
import (

	"net/http"
	"encoding/json"
	// "net/url"
)
type UnbxdResponse struct {
	Success  string                 `json:"success"`
	Response UnbxdResponseData      `json:"response"`
	Debug    map[string]interface{} `json:"debug,omitempty"`
}

type UnbxdResponseData struct {
	NumberOfProducts int                    `json:"numberOfProducts"`
	Products         []map[string]interface{} `json:"products"`
}

func SearchUnbxd(unbxdURL string) (*UnbxdResponse, error) {

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