package main
import (
	"database/sql"
	"fmt"
	"net/http"
	"example.com/belezathreads/backend/model"
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

	_, err = db.Exec(CREATE_CATEGORY_TABLE)
	if err != nil {
		fmt.Println("Error creating category table:", err)
	}
	_, err= db.Exec(CREATE_PRODUCT_TABLE)
	if err != nil {
		fmt.Println("Error creating  product table:", err)
	}
	_, err = db.Exec(CREATE_IMAGE_TABLE)

	if err != nil {
		fmt.Println("Error creating image table:", err)
		return
	}
	fmt.Println("Table creation successful or already exists.")

	
	// insertSample1 := `insert into CATEGORY ("category", "parentcategory") values('New Arrivalss', 1)`
	// _, err = db.Exec(insertSample1)
	// CheckError(err)

	// insertSample2 := `insert into PRODUCT ("productid", "title", "price", "description", "categoryid") values($1, $2, $3, $4, $5)`
	// _, err = db.Exec(insertSample2, "123", "it is", 56.34, "come buy", 1)
	// CheckError(err)

	// insertSample3 := `insert into IMAGE ("imageid", "imagepath", "productid") values($1, $2, $3)`
	// _, err = db.Exec(insertSample3, "123", "https://images.express.com/is/image/expressfashion/0020_01705319_0001?cache=on&wid=361&fmt=jpeg&qlt=75,1&resmode=sharp2&op_usm=1,1,5,0&defaultImage=Photo-Coming-Soon", "123")
	// CheckError(err)

	router := mux.NewRouter()

	router.HandleFunc("/ingestion", func(w http.ResponseWriter, r *http.Request) {
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



