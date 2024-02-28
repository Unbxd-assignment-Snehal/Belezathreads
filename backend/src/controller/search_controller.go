package controller

import (
	"encoding/json"
	"net/http"
	"example.com/belezathreads/backend/src/model"
)

func SearchUnbxdController(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	pageno := r.URL.Query().Get("pageno")
	sort := r.URL.Query().Get("sort")
	fields := "title,price,description,imageUrl"

	unbxdResponse, err := model.SearchUnbxd(q, pageno, sort, fields)
	if err != nil {
		errorResponse := model.UnbxdResponse{
			Success: "",
			Response: model.UnbxdResponseData{
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

	successResponse := model.UnbxdResponse{
		Success: "ok",
		Response: model.UnbxdResponseData{
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
