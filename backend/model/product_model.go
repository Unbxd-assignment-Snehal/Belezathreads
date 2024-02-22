package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

const GET_PRODUCT = `
SELECT P.productID, P.title, P.price, P.description, P.categoryID, I.imagePath 
FROM PRODUCT P JOIN IMAGE I ON P.productid = I.productid 
WHERE P.productid = $1 ;`

type ProductResponse struct {
	ProductID    string  `json:"productID"`
	Title        string  `json:"title"`
	Price        float64 `json:"price"`
	Description  string  `json:"description"`
	CategoryID   int     `json:"categoryID"`
	ImagePath    string  `json:"imagePath"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func GetProductHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestedProductID := vars["productID"]
		var response ProductResponse

		err := db.QueryRow(GET_PRODUCT, requestedProductID).Scan(
			&response.ProductID, &response.Title, &response.Price,
			&response.Description, &response.CategoryID, &response.ImagePath,
		)

		if err == sql.ErrNoRows {
			errorResponse := ErrorResponse{Message: "Error: Product not found"}

			jsonError, _ := json.Marshal(errorResponse)

			w.Header().Set("Content-Type", "application/json")

			w.WriteHeader(http.StatusNotFound)
			w.Write(jsonError)
			return
		} else if err != nil {
			errorResponse := ErrorResponse{Message: fmt.Sprintf("Error: %v", err)}

			jsonError, _ := json.Marshal(errorResponse)

			w.Header().Set("Content-Type", "application/json")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(jsonError)
			return
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			errorResponse := ErrorResponse{Message: fmt.Sprintf("Error: %v", err)}

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
