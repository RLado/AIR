package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/c-bata/go-prompt"
	"github.com/rlado/air/setup"
)

// Config struct
type config struct {
	DbPath string `json:"db_path"`
}

// Person struct
type person struct {
	Id         int
	Name       string
	TinNumber  string
	Address    string
	City       string
	PostalCode string
	Country    string
	Phone      string
	Email      string
}

// Product/Service struct
type itemRow struct {
	IsoCurrency  string
	Concept      string
	ConceptNote  string
	UnitCost     float32
	UnitCostNote string
	SumCost      float32
	SumCostNote  string
	Discount     float32
	DiscountNote string
	Tax          float32
	TaxNote      string
	Total        float32
	TotalNote    string
}

// Invoice struct
type invoice struct {
	Id          int
	Series      string
	Number      int
	Date        string
	Customer    person
	Issuer      person
	IsoCurrency string
	Items       []itemRow
	SubTotal    float32
	Discount    float32
	Total       float32
	Tax         float32
	Final       float32
	Notes       string
	PayMethod   string
	Footer      string
}

// Simple invoice struct

// Prompt level 1 options
func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "users", Description: "Store the username and age"},
		{Text: "articles", Description: "Store the article text posted by user"},
		{Text: "comments", Description: "Store the text commented to articles"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func main() {
	// Flags
	config_path := flag.String("c", "./config/config.json", "Path to main config file. (Default: ./config/config.json)")
	init_flag := flag.Bool("i", false, "Initialize database. (Default: false)")
	flag.Parse()

	// Read config file
	var conf config
	file, err := os.Open(*config_path)
	if err != nil {
		log.Fatalf("Error opening config file: %s", err)
	}
	defer file.Close()

	jsonDecoder := json.NewDecoder(file)
	err = jsonDecoder.Decode(&conf)
	if err != nil {
		log.Fatalf("Error decoding config file: %s", err)
	}

	// Initialize database (if req.)
	if *init_flag {
		fmt.Print("Are you sure you want to initialize the database? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" {
			log.Print("Exiting...")
			return
		}

		err := setup.CreateDatabase(conf.DbPath)
		if err != nil {
			log.Fatalf("error creating a new database: %s", err)
		}
	}

	// Open database
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", conf.DbPath))
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	fmt.Println("Please select table.")
	t := prompt.Input("> ", completer)
	fmt.Println("You selected " + t)

}
