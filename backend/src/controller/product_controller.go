
package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"example.com/belezathreads/backend/src/model"
)

func GetProductController(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestedProductID := vars["productID"]

		response, err := model.GetProductModel(db, requestedProductID)
		if err != nil {
			errorResponse := model.ErrorResponse{Message: fmt.Sprintf("Error: %v", err)}

			jsonError, _ := json.Marshal(errorResponse)

			w.Header().Set("Content-Type", "application/json")

			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

			w.Write(jsonError)
			return
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			errorResponse := model.ErrorResponse{Message: fmt.Sprintf("Error: %v", err)}

			jsonError, _ := json.Marshal(errorResponse)

			w.Header().Set("Content-Type", "application/json")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(jsonError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}
