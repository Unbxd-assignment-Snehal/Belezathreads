
package model

import (
	"database/sql"
	"fmt"
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

func GetProductModel(db *sql.DB, requestedProductID string) (ProductResponse, error) {
	var response ProductResponse

	err := db.QueryRow(GET_PRODUCT, requestedProductID).Scan(
		&response.ProductID, &response.Title, &response.Price,
		&response.Description, &response.CategoryID, &response.ImagePath,
	)

	if err == sql.ErrNoRows {
		return response, fmt.Errorf("Error: Product not found")
	} else if err != nil {
		return response, fmt.Errorf("Error: %v", err)
	}

	return response, nil
}
