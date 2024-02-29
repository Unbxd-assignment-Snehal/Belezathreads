package model

import (
	// "encoding/json"
	// "net/http"
	"net/url"

	"example.com/belezathreads/backend/src/services"
)

const UnbxdAPIEndpoint = "https://search.unbxd.io/fb853e3332f2645fac9d71dc63e09ec1/demo-unbxd700181503576558/search"





func BuildUnbxdURL(q, pageno, sort, fields string) (*services.UnbxdResponse, error) {
	baseURL, err := url.Parse(UnbxdAPIEndpoint)
	if err != nil {
		return nil, err
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
	resp, err := services.SearchUnbxd(finalURL)
	if err!=nil {
		return nil, err
	}
	return resp, nil
}
