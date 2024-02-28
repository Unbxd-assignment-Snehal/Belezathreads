
package model

import (
	"database/sql"
	"fmt"
	"example.com/belezathreads/backend/src/services"
	"strconv"
	
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

const GET_CAT1_PRODUCTS_PAGINATED = `
SELECT P.productID, P.title, P.price, P.description, P.categoryID, C.category, C.parentcategory, I.imagePath
FROM PRODUCT P
JOIN CATEGORY C ON P.categoryid = C.categoryid
LEFT JOIN IMAGE I ON P.productID = I.productID
WHERE ((C.category = $1 AND C.parentcategory IS NULL) OR C.parentcategory = (
    SELECT C.categoryid FROM CATEGORY C WHERE C.category = $1
))
ORDER BY P.price %s
LIMIT $2 OFFSET $3;
`

func FilterCategoryModel(db *sql.DB, cat1 string, pageNo string, sort string) ([]Product, error) {
	if pageNo == "" {
		pageNo = "1"
	}
	pageNoInt, err := strconv.Atoi(pageNo)
	if err != nil {
		return nil, err
	}

	pageSize := 10
	offset := (pageNoInt - 1) * pageSize

	query := fmt.Sprintf(GET_CAT1_PRODUCTS_PAGINATED, sort)
	rows, err := services.QueryDB(db, query, cat1, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ProductID, &product.Title, &product.Price, &product.Description, &product.CategoryID, &product.Category, &product.ParentCategory, &product.Imagepath)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}
