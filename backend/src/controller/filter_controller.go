package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"example.com/belezathreads/backend/src/model"
)

func FilterCategoryController(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		cat1 := params["cat1"]

		pageNo := r.URL.Query().Get("pageno")
		sort := r.URL.Query().Get("sort")

		products, err := model.FilterCategoryModel(db, cat1, pageNo, sort)
		if err != nil {
			fmt.Println("Error in model:", err)
			if err == sql.ErrNoRows {
				errorResponse := map[string]string{"error": "Product not found"}
				response, err := json.Marshal(errorResponse)
				if err != nil {
					fmt.Println("Error marshaling JSON:", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				w.Write(response)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(products) == 0 {
			errorResponse := map[string]string{"error": "Product not found"}
			response, err := json.Marshal(errorResponse)
			if err != nil {
				fmt.Println("Error marshaling JSON:", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write(response)
			return
		}

		response, err := json.Marshal(products)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
