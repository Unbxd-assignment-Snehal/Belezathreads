package model

import (
	"database/sql"
	"fmt"
	"strconv"

	"example.com/belezathreads/backend/src/services"
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

func FilterCategoryModel2(db *sql.DB, cat1 string, cat2 string, pageNo string, sort string) ([]Product, error) {

	if pageNo == "" {
		pageNo = "1"
	}
	pageNoInt, err := strconv.Atoi(pageNo)
	if err != nil {
		return nil, err
	}

	pageSize := 10
	offset := (pageNoInt - 1) * pageSize

	query := fmt.Sprintf(GET_CAT1_CAT2_PRODUCTS_PAGINATED, sort)
	rows, err := services.QueryDB(db, query, cat1, cat2, pageSize, offset)
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
