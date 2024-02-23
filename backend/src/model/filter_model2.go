package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

const GET_CAT1_CAT2_PRODUCTS_PAGINATED = `
SELECT P.productID, P.title, P.price, P.description, P.categoryID, C.category, C.parentcategory, I.imagePath
FROM PRODUCT P
JOIN CATEGORY C ON P.categoryid = C.categoryid
LEFT JOIN IMAGE I ON P.productID = I.productID
WHERE (C.category = $2 AND C.parentcategory = (
    SELECT categoryid FROM CATEGORY WHERE category = $1
))
ORDER BY P.price %s
LIMIT $3 OFFSET $4;
`

func FilterCategoryHandler2(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		cat1 := vars["cat1"]
		cat2 := vars["cat2"]

		// Get query parameters
		pageNo := r.URL.Query().Get("pageno")
		sort := r.URL.Query().Get("sort")

		// Set default values
		if pageNo == "" {
			pageNo = "1"
		}
		if sort != "asc" && sort != "desc" {
			sort = "asc"
		}

		// Convert pageNo to integer
		pageNoInt, err := strconv.Atoi(pageNo)
		if err != nil {
			fmt.Println("Error converting pageno to integer:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pageSize := 10
		offset := (pageNoInt - 1) * pageSize

		// Construct the SQL query with pagination and sorting
		query := fmt.Sprintf(GET_CAT1_CAT2_PRODUCTS_PAGINATED, sort)
		rows, err := db.Query(query, cat1, cat2, pageSize, offset)
		if err != nil {
			fmt.Println("Error querying database:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		defer rows.Close()

		var products []Product

		for rows.Next() {
			var product Product
			err := rows.Scan(&product.ProductID, &product.Title, &product.Price, &product.Description, &product.CategoryID, &product.Category, &product.ParentCategory, &product.Imagepath)
			if err != nil {
				fmt.Println(http.StatusInternalServerError)
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
