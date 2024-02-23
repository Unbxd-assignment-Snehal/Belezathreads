package main
import (
	"database/sql"
	"fmt"
	"net/http"
	"example.com/belezathreads/backend/src/model"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "apparels"
)

func main() {
	psqlconn := fmt.Sprintf("host = %s port = %d user = %s password = %s dbname = %s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Println("Error opening database connection:", err)
		return
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Error pinging database:", err)
		return
	}
	defer db.Close()
	fmt.Println("Database connection successful!")
	
	
	
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


	router.HandleFunc("/product/{productID}", model.GetProductHandler(db)).Methods("GET")
	router.HandleFunc("/products/{cat1}", model.FilterCategoryHandler(db)).Methods("GET")
	router.HandleFunc("/products/{cat1}/{cat2}", model.FilterCategoryHandler2(db)).Methods("GET")

	port := ":8080"
	fmt.Printf("Server is running on http://localhost%s\n", port)
	http.ListenAndServe(port, router)
}



