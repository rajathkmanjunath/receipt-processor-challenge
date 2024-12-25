package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Receipt represents the structure of a receipt JSON.
type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
	ID           string `json:"id"`
	Points       int    `json:"points"`
}

// Item represents an item in the receipt.
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

func main() {
	router := mux.NewRouter()

	// Define API endpoints
	router.HandleFunc("/receipts/process", processReceipt).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", getPoints).Methods("GET")

	// Start the web server
	port := ":8080"
	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}

// Process a receipt and calculate points
func processReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt

	// Decode the JSON payload
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Calculate points
	receipt.Points = calculatePoints(&receipt)

	// Generate a unique ID for the receipt
	receipt.ID = uuid.New().String()

	// Store the receipt in memory (you can use a map or a slice)
	// For simplicity, we use a map where the key is the receipt ID
	// and the value is the Receipt struct.
	receipts[receipt.ID] = receipt

	// Send the ID in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": receipt.ID})
}

// Retrieve points for a receipt by ID
func getPoints(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	receiptID := params["id"]

	// Lookup the receipt by ID
	receipt, exists := receipts[receiptID]
	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	// Send the points in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"points": receipt.Points})
}

// Calculate points based on the receipt rules
func calculatePoints(receipt *Receipt) int {
	points := 0

	// Rule 1: One point for every alphanumeric character in the retailer name.
	for _, char := range receipt.Retailer {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			points++
		}
	}

	// Rule 2: 50 points if the total is a round dollar amount with no cents.
	totalFloat, err := strconv.ParseFloat(receipt.Total, 64)
	if err == nil && math.Mod(totalFloat, 1.0) == 0 {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25.
	if totalFloat > 0 && math.Mod(totalFloat, 0.25) == 0 {
		points += 25
	}

	// Rule 4: 5 points for every two items on the receipt.
	numItems := len(receipt.Items)
	points += (numItems / 2) * 5

	// Rule 5: If the trimmed length of the item description is a multiple of 3,
	// multiply the price by 0.2 and round up to the nearest integer.
	for _, item := range receipt.Items {
		trimmedLen := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLen%3 == 0 {
			priceFloat, err := strconv.ParseFloat(item.Price, 64)
			if err == nil {
				points += int(math.Round(priceFloat * 0.2))
			}
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd.
	purchaseDate, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err == nil && purchaseDate.Day()%2 != 0 {
		points += 6
	}

	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
	if err == nil && purchaseTime.After(time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC)) && purchaseTime.Before(time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)) {
		points += 10
	}

	return points
}

// Check if a rune is a letter or digit
func unicodeIsLetterOrDigit(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

// Store receipts in memory (replace this with a database in a production environment)
var receipts = make(map[string]Receipt)
