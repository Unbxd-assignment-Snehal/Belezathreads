package controller

import (
	"encoding/json"
	"net/http"

	"example.com/belezathreads/backend/src/model"
	"example.com/belezathreads/backend/src/services"
)

func SearchUnbxdController(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	pageno := r.URL.Query().Get("pageno")
	sort := r.URL.Query().Get("sort")
	fields := "title,price,description,imageUrl"

	unbxdResponse, err := model.BuildUnbxdURL(q, pageno, sort, fields)
	if err != nil {
		errorResponse := services.UnbxdResponse{
			Success: "",
			Response: services.UnbxdResponseData{
				NumberOfProducts: 0,
				Products:         nil,
			},
			Debug: map[string]interface{}{
				"reason": err.Error(),
			},
		}

		jsonResponse, _ := json.Marshal(errorResponse)
		http.Error(w, string(jsonResponse), http.StatusInternalServerError)
		return
	}

	successResponse := services.UnbxdResponse{
		Success: "ok",
		Response: services.UnbxdResponseData{
			NumberOfProducts: len(unbxdResponse.Response.Products),
			Products:         unbxdResponse.Response.Products,
		},
	}

	jsonResponse, err := json.Marshal(successResponse)
	if err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
