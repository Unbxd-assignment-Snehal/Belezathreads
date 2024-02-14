package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	_ "github.com/lib/pq"
)

type Category struct {
	CatLevel1Name string `json:"catlevel1Name"`
	CatLevel2Name string `json:"catlevel2Name"`
}

type Product struct {
	UniqueID          string  `json:"uniqueId"`
	Title             string  `json:"title"`
	Price             float64 `json:"price"`
	ProductDescription string  `json:"productDescription"`
}

type Image struct{	
	ImagePath		  string `json:"productImage"`
	UniqueID          string  `json:"uniqueId"`

}


func ingestData(db *sql.DB) {
	jsonData, err := os.ReadFile("sample.json")
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
			CatLevel2Name: item["catlevel2Name"].(string),
		}
		_, err := db.Exec("INSERT INTO CATEGORY (category, parentCategory) VALUES ($1, $2) ON CONFLICT DO NOTHING", category.CatLevel2Name, category.CatLevel1Name)
		if err != nil {
			fmt.Println("Error inserting data into CATEGORY table:", err)
			return
		}
	}

	for _, item := range data {
		product := Product{
			UniqueID:          item["uniqueId"].(string),
			Title:             item["title"].(string),
			Price:             item["price"].(float64),
			ProductDescription: item["productDescription"].(string),			
		}
		_, err := db.Exec("INSERT INTO PRODUCT (productID, title, price, description) VALUES ($1, $2, $3, $4)", product.UniqueID, product.Title, product.Price, product.ProductDescription)
		if err != nil {
			fmt.Println("Error inserting data into PRODUCT table:", err)
			return
		}
	}

	for _, item := range data {
		image := Image{
			ImagePath: item["productImage"].(string),
			UniqueID: item["uniqueId"].(string),
		}
		_, err := db.Exec("INSERT INTO IMAGE (imagepath, productid) VALUES ($1, $2)", image.ImagePath, image.UniqueID)
		if err != nil {
			fmt.Println("Error inserting data into IMAGE table:", err)
			return
		}
	}

	fmt.Println("Data insertion successful!")
}