package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"example.com/belezathreads/backend/src/services"
	_ "github.com/lib/pq"
)

type Category struct {
	CatLevel1Name string `json:"catlevel1Name"`
	CatLevel2Name string `json:"catlevel2Name"`
}

type Product struct {
	UniqueID           string  `json:"uniqueId"`
	Title              string  `json:"title"`
	Price              float64 `json:"price"`
	ProductDescription string  `json:"productDescription"`
	CategoryID         int
}

type Image struct {
	ImagePath string `json:"productImage"`
	UniqueID  string `json:"uniqueId"`
}

func ingestData(db *sql.DB) {
	jsonData, err := os.ReadFile("data/out.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	var data []map[string]interface{}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON data:", err)
		return
	}

	for _, item := range data {
		category := Category{
			CatLevel1Name: item["catlevel1Name"].(string),
		}

		if catLevel2Name, ok := item["catlevel2Name"].(string); ok {
			category.CatLevel2Name = catLevel2Name
		} else {
			category.CatLevel2Name = " "
		}

		var parentCategoryID int
		err := services.QueryRowDB(db, "SELECT categoryid FROM CATEGORY WHERE category = $1", category.CatLevel1Name).Scan(&parentCategoryID)

		if err == sql.ErrNoRows {
			_, err := services.ExecDB(db, "INSERT INTO CATEGORY (category, parentCategory) VALUES ($1, $2) ON CONFLICT DO NOTHING", category.CatLevel1Name, nil)
			if err != nil {
				fmt.Println("Error inserting parent category:", err)
				return
			}

			err = services.QueryRowDB(db, "SELECT categoryid FROM CATEGORY WHERE category = $1", category.CatLevel1Name).Scan(&parentCategoryID)
			if err != nil {
				fmt.Println("Error retrieving categoryid of the parent category:", err)
				return
			}
		} else if err != nil {
			fmt.Println("Error checking if parent category exists:", err)
			return
		}
		if category.CatLevel2Name != " " {

			_, err = services.ExecDB(db, "INSERT INTO CATEGORY (category, parentCategory) VALUES ($1, $2) ON CONFLICT DO NOTHING", category.CatLevel2Name, parentCategoryID)
			if err != nil {
				fmt.Println("Error inserting child category:", err)
				return
			}

		}
	}

	for _, item := range data {
		product := Product{
			UniqueID: item["uniqueId"].(string),
			Title:    item["title"].(string),
			Price:    item["price"].(float64),
		}
		if productDescription, ok := item["productDescription"].(string); ok {
			product.ProductDescription = productDescription
		} else {
			product.ProductDescription = ""
		}

		var parentCategoryID int
		err := services.QueryRowDB(db, "SELECT categoryid FROM CATEGORY WHERE category = $1", item["catlevel1Name"]).Scan(&parentCategoryID)
		if err != nil {
			fmt.Println("Error retrieving categoryid of catlevel1Name:", err)
			return
		}
		var categoryID int
		if catLevel2Name, ok := item["catlevel2Name"].(string); ok {
			err := services.QueryRowDB(db, "SELECT categoryid FROM CATEGORY WHERE category = $1 AND parentCategory = $2", catLevel2Name, parentCategoryID).Scan(&categoryID)
			if err != nil {
				fmt.Println("Error retrieving categoryID from CATEGORY table:", err)
				return
			}
		} else {
			err := services.QueryRowDB(db, "SELECT categoryid FROM CATEGORY WHERE  category = $1", item["catlevel1Name"]).Scan(&categoryID)
			if err == sql.ErrNoRows {
				fmt.Println("Category not found in CATEGORY table", err.Error())
				return
			} else if err != nil {
				fmt.Println("Error retrieving categoryID from CATEGORY table:", err)
				return
			}
		}
		_, err = services.ExecDB(db, "INSERT INTO PRODUCT (productID, title, price, description, categoryID) VALUES ($1, $2, $3, $4, $5)", product.UniqueID, product.Title, product.Price, product.ProductDescription, categoryID)
		if err != nil {
			fmt.Println("Error inserting data into PRODUCT table:", err)
			return
		}
	}

	for _, item := range data {
		image := Image{
			ImagePath: item["productImage"].(string),
			UniqueID:  item["uniqueId"].(string),
		}
		_, err := services.ExecDB(db, "INSERT INTO IMAGE (imagepath, productid) VALUES ($1, $2)", image.ImagePath, image.UniqueID)
		if err != nil {
			fmt.Println("Error inserting data into IMAGE table:", err)
			return
		}
	}

	fmt.Println("Data insertion successful!")
}
