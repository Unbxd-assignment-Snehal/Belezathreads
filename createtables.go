package main

import (
	"database/sql"
	"fmt"

	"example.com/belezathreads/backend/src/services"
)

func createTables(db *sql.DB) error {
	_, err := services.ExecDB(db, CREATE_CATEGORY_TABLE)
	if err != nil {
		return fmt.Errorf("error creating category table: %v", err)
	}

	_, err = services.ExecDB(db, CREATE_PRODUCT_TABLE)
	if err != nil {
		return fmt.Errorf("error creating product table: %v", err)
	}

	_, err = services.ExecDB(db, CREATE_IMAGE_TABLE)
	if err != nil {
		return fmt.Errorf("error creating image table: %v", err)
	}

	fmt.Println("Table creation successful")
	return nil
}
