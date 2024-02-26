
package controller

import (
	"encoding/json"
	"net/http"
	"example.com/belezathreads/backend/src/model")

	func SearchUnbxdController(w http.ResponseWriter, r *http.Request) {
		// Extract query parameters from the request
		q := r.URL.Query().Get("q")
		pageno := r.URL.Query().Get("pageno")
		sort := r.URL.Query().Get("sort")
	
		// Hardcoded fields parameter
		fields := "title,price,description,imageUrl"
	
		// Call the model function to perform the Unbxd search
		unbxdResponse, err := model.SearchUnbxd(q, pageno, sort, fields)
		if err != nil {
			// Handle error
			http.Error(w, "Error contacting Unbxd API", http.StatusInternalServerError)
			return
		}
	
		// Convert the UnbxdResponse to JSON
		jsonResponse, err := json.Marshal(unbxdResponse)
		if err != nil {
			// Handle JSON encoding error
			http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
			return
		}
	
		// Set response headers
		w.Header().Set("Content-Type", "application/json")
	
		// Write the JSON response
		w.Write(jsonResponse)
	}