// filter_model2.go
package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
)


const GET_CAT1_CAT2_PRODUCTS = "SELECT P.productID, P.title, P.price, P.description, P.categoryID, C.category, C.parentcategory FROM PRODUCT P JOIN CATEGORY C ON P.categoryid = C.categoryid WHERE C.categoryid = $1 AND C.parentcategory = $2"

// FilterCategoryHandler2 handles category filtering based on two parameters (cat1 and cat2)
func FilterCategoryHandler2(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		cat1 := params["cat1"]
		cat2 := params["cat2"]

		categoryID, err := strconv.Atoi(cat1)
		if err != nil {
			fmt.Println("Error converting cat1 to integer:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		parentCategoryID, err := strconv.Atoi(cat2)
		if err != nil {
			fmt.Println("Error converting cat2 to integer:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		rows, err := db.Query(GET_CAT1_CAT2_PRODUCTS, categoryID, parentCategoryID)
		if err != nil {
			fmt.Println("Error querying database:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var products []Product

		for rows.Next() {
			var product Product
			err := rows.Scan(&product.ProductID, &product.Title, &product.Price, &product.Description, &product.CategoryID, &product.Category, &product.ParentCategory)
			if err != nil {
				fmt.Println("Error scanning row:", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			products = append(products, product)
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
