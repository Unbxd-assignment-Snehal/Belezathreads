package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Product struct {
	ProductID      string      `json:"productid"`
	Title          string      `json:"title"`
	Price          float64     `json:"price"`
	Description    string      `json:"description"`
	CategoryID     int         `json:"categoryid"`
	Category       string      `json:"category"`
	ParentCategory sql.NullInt64 `json:"parentcategory"`
	Imagepath      string      `json:"imagepath"`
}

const GET_CAT1_PRODUCTS = `
SELECT P.productID, P.title, P.price, P.description, P.categoryID, C.category, C.parentcategory, I.imagePath
FROM PRODUCT P
JOIN CATEGORY C ON P.categoryid = C.categoryid
LEFT JOIN IMAGE I ON P.productID = I.productID
WHERE ((C.category = $1 AND C.parentcategory IS NULL) OR C.parentcategory = (
    SELECT C.categoryid FROM CATEGORY C WHERE C.category = $1
)) ;
`


func FilterCategoryHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		cat1 := params["cat1"]

		rows, err := db.Query(GET_CAT1_PRODUCTS, cat1)
		if err != nil {
			fmt.Println("Error querying database:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var products []Product

		for rows.Next() {
			var product Product
			err := rows.Scan(&product.ProductID, &product.Title, &product.Price, &product.Description, &product.CategoryID, &product.Category, &product.ParentCategory, &product.Imagepath)
			if err != nil {
				fmt.Println("Error scanning row:", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			products = append(products, product)
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
