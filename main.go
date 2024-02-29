package main
import (
	"fmt"
	"net/http"
	"example.com/belezathreads/backend/src/controller"
	"example.com/belezathreads/backend/src/services"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	dbConfig := services.DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   "apparels",
	}

	db, err := services.NewDBConnection(dbConfig)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer services.CloseDB(db)

	
	router := mux.NewRouter()
	router.HandleFunc("/ingestion", func(w http.ResponseWriter, r *http.Request) {
		err := createTables(db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error creating tables: %v", err)))
			return
		}
		ingestData(db)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Data ingestion successful"))
	}).Methods("POST")


	router.HandleFunc("/product/{productID}", controller.GetProductController(db)).Methods("GET")
	router.HandleFunc("/products/{cat1}", controller.FilterCategoryController(db)).Methods("GET")
	router.HandleFunc("/products/{cat1}/{cat2}", controller.FilterCategoryController2(db)).Methods("GET")
	router.HandleFunc("/search", controller.SearchUnbxdController).Methods("GET")
	port := ":8080"
	fmt.Printf("Server is running on http://localhost%s\n", port)
	http.ListenAndServe(port, router)
}



