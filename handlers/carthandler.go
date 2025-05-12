package handlers

import (
    "encoding/json"
    "net/http"
	"fmt"

	"strconv"
    "fashion/database"
)

type Product struct {
	Id              int    `json:"id"`
    ImageURL        string `json:"image_url"`
    Brand           string `json:"brand"`
    Para            string `json:"para"`
    Price           string `json:"price"`
    Rs              int    `json:"rs"`
    StrikedOffPrice string `json:"strikedoffprice"`
    Offer           string `json:"offer"`
    Atc             string `json:"atc"`
    Atw             string `json:"atw"`
    Category        string `json:"category"`
}

// POST /api/cart
func AddToCartHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // Handle preflight request
    if r.Method == http.MethodOptions {
        return
    }
    fmt.Println("AddToCartHandler called")

    var p Product
    if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
        fmt.Println("Error decoding:", err)
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    fmt.Printf("Received product: %+v\n", p)

    query := `INSERT INTO cart_items 
        (id, image_url, brand, para, price, rs, strikedoffprice, offer, atc, atw, category) 
        VALUES (?,?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

    _, err := database.Con.Exec(query, p.Id, p.ImageURL, p.Brand, p.Para, p.Price, p.Rs, p.StrikedOffPrice, p.Offer, p.Atc, p.Atw, p.Category)
    if err != nil {
        fmt.Println("Database error:", err)
        http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}



// DELETE /api/cart?id=1

func DeleteCartItemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	// Convert id string to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return
	}

	// Execute delete query
	_, err = database.Con.Exec("DELETE FROM cart_items WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item deleted successfully"})
}

// GET /api/cart
func DisplayCartItemsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	rows, err := database.Con.Query("SELECT id, image_url, brand, para, price, rs, strikedoffprice, offer, atc, atw, category FROM cart_items")
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var p Product
		err := rows.Scan( &p.Id, &p.ImageURL, &p.Brand, &p.Para, &p.Price, &p.Rs, &p.StrikedOffPrice, &p.Offer, &p.Atc, &p.Atw, &p.Category)
		if err != nil {
			http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
	}
}
