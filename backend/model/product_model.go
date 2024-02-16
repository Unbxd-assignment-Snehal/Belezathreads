package model

import (
	"database/sql"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)
const GET_PRODUCT = `
SELECT P.productID, P.title, P.price, P.description, P.categoryID, I.imagePath 
FROM PRODUCT P JOIN IMAGE I ON P.productid = I.productid 
WHERE P.productid = $1;` 
func GetProductHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestedProductID := vars["productID"]
		var (
			title       string
			price       float64
			description string
			categoryID  int
			imagePath   string
		)

		err := db.QueryRow(GET_PRODUCT, requestedProductID).Scan(&requestedProductID, &title, &price, &description, &categoryID, &imagePath)
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Error: Product not found"))
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %v", err)))
			return
		}

		fmt.Printf("Product ID: %s\nTitle: %s\nPrice: %.2f\nDescription: %s\nCategory ID: %d\nImage Path: %s\n",
			requestedProductID, title, price, description, categoryID, imagePath)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Product details displayed in terminal"))
	}
}
