package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

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

// invoice struct
type invoice struct {
	Id          int           `json:"id"`
	Series      string        `json:"series"`
	Number      int           `json:"number"`
	Date        string        `json:"date"`
	Customer    person        `json:"customer"`
	Issuer      person        `json:"issuer"`
	IsoCurrency string        `json:"iso_currency"`
	Items       []itemRow     `json:"items"`
	ItemsRender template.HTML `json:"items_render"`
	SubTotal    float32       `json:"sub_total"`
	Discount    float32       `json:"discount"`
	Total       float32       `json:"total"`
	Tax         float32       `json:"tax"`
	Final       float32       `json:"final"`
	Notes       string        `json:"notes"`
	PayMethod   string        `json:"pay_method"`
	Footer      string        `json:"footer"`
}

// Issued invoice struct (summary table)
type IssuedInvoice struct {
	Id       int    `json:"id"`
	Series   string `json:"series"`
	Number   int    `json:"number"`
	Date     string `json:"date"`
	Customer person `json:"customer"`
}

// Render invoice
func (inv invoice) render() *bytes.Buffer {
	// Render Items table
	tableBuf := bytes.NewBuffer([]byte{})
	tmplItem, err := template.New("itemTable").ParseFiles("templates/SimpleInvoice/ItemTable.html")
	if err != nil {
		log.Fatalf("Error parsing item template: %s", err)
	}
	for _, item := range inv.Items {
		err = tmplItem.Execute(tableBuf, item)
		if err != nil {
			log.Fatalf("Error executing item template: %s", err)
		}
	}
	inv.ItemsRender = template.HTML(tableBuf.String())

	// Render Invoice
	invoiceBuf := bytes.NewBuffer([]byte{})
	tmplInvoice, err := template.New("invoice").ParseFiles("templates/SimpleInvoice/SimpleInvoice.html")
	if err != nil {
		log.Fatalf("Error parsing invoice template: %s", err)
	}
	err = tmplInvoice.Execute(invoiceBuf, inv)
	if err != nil {
		log.Fatalf("Error executing invoice template: %s", err)
	}

	return invoiceBuf
}

// CLI for adding issuer data to the database
func createIssuer(db *sql.DB, state *int) {
	var usrInput string
	issuer := person{}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Name: ")
	scanner.Scan()
	issuer.Name = scanner.Text()
	fmt.Print("Tax identification number: ")
	scanner.Scan()
	issuer.TinNumber = scanner.Text()
	fmt.Print("Address: ")
	scanner.Scan()
	issuer.Address = scanner.Text()
	fmt.Print("City: ")
	scanner.Scan()
	issuer.City = scanner.Text()
	fmt.Print("Postal code: ")
	scanner.Scan()
	issuer.PostalCode = scanner.Text()
	fmt.Print("Country: ")
	scanner.Scan()
	issuer.Country = scanner.Text()
	fmt.Print("Phone: ")
	scanner.Scan()
	issuer.Phone = scanner.Text()
	fmt.Print("Email: ")
	scanner.Scan()
	issuer.Email = scanner.Text()

	// Summary:
	fmt.Printf(`
--- Summary ---
Name: %s
TIN: %s
Address: %s
City: %s
Postal code: %s
Country: %s
Phone: %s
Email: %s

Is this information correct? (y/N):
`, issuer.Name, issuer.TinNumber, issuer.Address, issuer.City, issuer.PostalCode, issuer.Country, issuer.Phone, issuer.Email)
	usrInput = "N"
	fmt.Scanln(&usrInput)
	if usrInput == "y" {
		// Save issuer data to database
		_, err := db.Exec("INSERT INTO User (Name, TinNumber, Address, City, PostalCode, Country, Phone, Email) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", issuer.Name, issuer.TinNumber, issuer.Address, issuer.City, issuer.PostalCode, issuer.Country, issuer.Phone, issuer.Email)
		if err != nil {
			log.Fatalf("Error saving issuer data to database: %s", err)
		}
		fmt.Println("Issuer data saved successfully")
	} else {
		fmt.Println("Issuer data not saved")
	}

	// Go back to main menu
	*state = 000
}

// CLI for listing issuer data from the database
func listIssuer(db *sql.DB, state *int) {
	rows, err := db.Query("SELECT * FROM User")
	if err != nil {
		log.Fatalf("Error querying issuer data: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var issuer person
		err := rows.Scan(&issuer.Id, &issuer.Name, &issuer.TinNumber, &issuer.Address, &issuer.City, &issuer.PostalCode, &issuer.Country, &issuer.Phone, &issuer.Email)
		if err != nil {
			log.Fatalf("Error scanning issuer data: %s", err)
		}
		fmt.Printf("%d - %s : %s : %s : %s : %s : %s : %s : %s", issuer.Id, issuer.Name, issuer.TinNumber, issuer.Address, issuer.City, issuer.PostalCode, issuer.Country, issuer.Phone, issuer.Email)
	}

	// Go back to main menu
	*state = 000
}

// CLI for deleting issuer data from the database
func deleteIssuer(db *sql.DB, state *int) {
	var usrInput string
	var issuer person

	fmt.Println("Enter the ID of the issuer you want to delete:")
	fmt.Scanln(&usrInput)

	row := db.QueryRow("SELECT * FROM User WHERE Id = ?", usrInput)
	err := row.Scan(&issuer.Id, &issuer.Name, &issuer.TinNumber, &issuer.Address, &issuer.City, &issuer.PostalCode, &issuer.Country, &issuer.Phone, &issuer.Email)
	if err != nil {
		fmt.Println("Issuer not found")
		*state = 000
		return
	}

	fmt.Printf("%d - %s : %s : %s : %s : %s : %s : %s : %s", issuer.Id, issuer.Name, issuer.TinNumber, issuer.Address, issuer.City, issuer.PostalCode, issuer.Country, issuer.Phone, issuer.Email)
	fmt.Println("Are you sure you want to delete the following issuer data? (y/N)")
	fmt.Scanln(&usrInput)
	if usrInput == "y" {
		_, err := db.Exec("DELETE FROM User WHERE Id = ?", issuer.Id)
		if err != nil {
			log.Fatalf("Error deleting issuer data: %s", err)
		}
		fmt.Println("Issuer data deleted successfully")
	} else {
		fmt.Println("Issuer data not deleted")
	}

	// Go back to main menu
	*state = 000
}

// CLI for adding customer data to the database
func createCustomer(db *sql.DB, state *int) {
	var usrInput string
	customer := person{}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Name: ")
	scanner.Scan()
	customer.Name = scanner.Text()
	fmt.Print("Tax identification number: ")
	scanner.Scan()
	customer.TinNumber = scanner.Text()
	fmt.Print("Address: ")
	scanner.Scan()
	customer.Address = scanner.Text()
	fmt.Print("City: ")
	scanner.Scan()
	customer.City = scanner.Text()
	fmt.Print("Postal code: ")
	scanner.Scan()
	customer.PostalCode = scanner.Text()
	fmt.Print("Country: ")
	scanner.Scan()
	customer.Country = scanner.Text()
	fmt.Print("Phone: ")
	scanner.Scan()
	customer.Phone = scanner.Text()
	fmt.Print("Email: ")
	scanner.Scan()
	customer.Email = scanner.Text()

	// Summary:
	fmt.Printf(`
--- Summary ---
Name: %s
TIN: %s
Address: %s
City: %s
Postal code: %s
Country: %s
Phone: %s
Email: %s

Is this information correct? (y/N):
`, customer.Name, customer.TinNumber, customer.Address, customer.City, customer.PostalCode, customer.Country, customer.Phone, customer.Email)
	usrInput = "N"
	fmt.Scanln(&usrInput)
	if usrInput == "y" {
		// Save customer data to database
		_, err := db.Exec("INSERT INTO Customers (Name, TinNumber, Address, City, PostalCode, Country, Phone, Email) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", customer.Name, customer.TinNumber, customer.Address, customer.City, customer.PostalCode, customer.Country, customer.Phone, customer.Email)
		if err != nil {
			log.Fatalf("Error saving customer data to database: %s", err)
		}
		fmt.Println("Customer data saved successfully")
	} else {
		fmt.Println("Customer data not saved")
	}

	// Go back to main menu
	*state = 000
}

// CLI for listing customer data from the database
func listCustomer(db *sql.DB, state *int) {
	rows, err := db.Query("SELECT * FROM Customers")
	if err != nil {
		log.Fatalf("Error querying customer data: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var customer person
		err := rows.Scan(&customer.Id, &customer.Name, &customer.TinNumber, &customer.Address, &customer.City, &customer.PostalCode, &customer.Country, &customer.Phone, &customer.Email)
		if err != nil {
			log.Fatalf("Error scanning customer data: %s", err)
		}
		fmt.Printf("%d - %s : %s : %s : %s : %s : %s : %s : %s", customer.Id, customer.Name, customer.TinNumber, customer.Address, customer.City, customer.PostalCode, customer.Country, customer.Phone, customer.Email)
	}

	// Go back to main menu
	*state = 000
}

// CLI for deleting customer data from the database
func deleteCustomer(db *sql.DB, state *int) {
	var usrInput string
	var customer person

	fmt.Println("Enter the ID of the customer you want to delete:")
	fmt.Scanln(&usrInput)

	row := db.QueryRow("SELECT * FROM User WHERE Id = ?", usrInput)
	err := row.Scan(&customer.Id, &customer.Name, &customer.TinNumber, &customer.Address, &customer.City, &customer.PostalCode, &customer.Country, &customer.Phone, &customer.Email)
	if err != nil {
		fmt.Println("Customer not found.")
		*state = 000
		return
	}

	fmt.Printf("%d - %s : %s : %s : %s : %s : %s : %s : %s", customer.Id, customer.Name, customer.TinNumber, customer.Address, customer.City, customer.PostalCode, customer.Country, customer.Phone, customer.Email)
	fmt.Println("Are you sure you want to delete the following customer data? (y/N)")
	fmt.Scanln(&usrInput)
	if usrInput == "y" {
		_, err := db.Exec("DELETE FROM User WHERE Id = ?", customer.Id)
		if err != nil {
			log.Fatalf("Error deleting customer data: %s", err)
		}
		fmt.Println("Customer data deleted successfully.")
	} else {
		fmt.Println("Customer data not deleted.")
	}

	// Go back to main menu
	*state = 000
}

// CLI for creating a new invoice
func createInvoice(db *sql.DB, state *int) {
	var usrInput string
	var inv invoice

	scanner := bufio.NewScanner(os.Stdin)

	// Request invoice data
	fmt.Print("Series: ")
	scanner.Scan()
	inv.Series = scanner.Text()

	// Get invoice number
	row := db.QueryRow("SELECT MAX(Number) FROM Invoices WHERE Series = ?", inv.Series)
	row.Scan(&inv.Number)
	inv.Number++

	// Get date in DD/MM/YYYY format
	inv.Date = time.Now().Format("02/01/2006")

	// Get customer data and query the database
	fmt.Print("Customer ID: ")
	fmt.Scanf("%d", &inv.Customer.Id)
	row = db.QueryRow("SELECT * FROM Customers WHERE Id = ?", inv.Customer.Id)
	err := row.Scan(&inv.Customer.Id, &inv.Customer.Name, &inv.Customer.TinNumber, &inv.Customer.Address, &inv.Customer.City, &inv.Customer.PostalCode, &inv.Customer.Country, &inv.Customer.Phone, &inv.Customer.Email)
	if err != nil {
		fmt.Printf("Error querying customer data. This customer ID might not extist: %s\n", err)
		*state = 000
		return
	}

	// Get issuer data and query the database
	fmt.Print("Issuer ID: ")
	fmt.Scanf("%d", &inv.Issuer.Id)
	row = db.QueryRow("SELECT * FROM User WHERE Id = ?", inv.Issuer.Id)
	err = row.Scan(&inv.Id, &inv.Issuer.Name, &inv.Issuer.TinNumber, &inv.Issuer.Address, &inv.Issuer.City, &inv.Issuer.PostalCode, &inv.Issuer.Country, &inv.Issuer.Phone, &inv.Issuer.Email)
	if err != nil {
		fmt.Printf("Error querying issuer data. This issuer ID might not exist: %s\n", err)
		*state = 000
		return
	}

	// Get ISO currency code
	fmt.Print("ISO Currency: ")
	fmt.Scanf("%s", &inv.IsoCurrency)

	// Get items
	fmt.Println("Items: ")
	for {
		fmt.Println("--- New Item ---")
		item := itemRow{}
		fmt.Print("Concept: ")
		scanner.Scan()
		item.Concept = scanner.Text()
		fmt.Print("Concept note: ")
		scanner.Scan()
		item.ConceptNote = scanner.Text()
		fmt.Print("Unit cost: ")
		fmt.Scanf("%f", &item.UnitCost)
		fmt.Print("Unit cost note: ")
		scanner.Scan()
		item.UnitCostNote = scanner.Text()
		fmt.Print("Number of units: ")
		var units float32
		fmt.Scanf("%f", &units)
		fmt.Print("Sum cost note: ")
		scanner.Scan()
		item.SumCostNote = scanner.Text()
		fmt.Print("Discount: ")
		fmt.Scanf("%f", &item.Discount)
		fmt.Print("Discount note: ")
		scanner.Scan()
		item.DiscountNote = scanner.Text()
		fmt.Print("Tax (%): ")
		fmt.Scanf("%f", &item.Tax)
		fmt.Print("Tax note: ")
		scanner.Scan()
		item.TaxNote = scanner.Text()
		fmt.Print("Total note: ")
		scanner.Scan()
		item.TotalNote = scanner.Text()

		// Calculate total
		item.SumCost = item.UnitCost * units
		item.Total = item.SumCost - item.Discount + ((item.SumCost - item.Discount) * item.Tax / 100)

		// Add currency code
		item.IsoCurrency = inv.IsoCurrency

		// Print item summary
		fmt.Printf(`
--- Summary ---
Concept: %s : %s
Unit cost: %f : %s
Sum cost: %f : %s
Discount: %f : %s
Tax: %f : %s
Total: %f : %s

`, item.Concept, item.ConceptNote, item.UnitCost, item.UnitCostNote, item.SumCost, item.SumCostNote, item.Discount, item.DiscountNote, item.Tax, item.TaxNote, item.Total, item.TotalNote)

		fmt.Print("Is this information correct? (y/N): ")
		usrInput = "N"
		fmt.Scanln(&usrInput)
		if usrInput == "y" {
			inv.Items = append(inv.Items, item)
		} else {
			fmt.Println("Item not added")
			if len(inv.Items) == 0 {
				continue
			}
		}

		// Ask if the user wants to add another item
		fmt.Print("Add another item? (y/N): ")
		usrInput = "N"
		fmt.Scanln(&usrInput)
		if usrInput != "y" {
			break
		}
	}

	// Get Notes, Payment method and Footer
	fmt.Print("Notes: ")
	stop_flag := false
	inputBuf := bytes.NewBuffer([]byte{})

	for {
		scanner.Scan()
		lastInput := scanner.Text()

		// Press enter on an empty line to finish
		if lastInput == "" && stop_flag {
			break
		} else if lastInput == "" && !stop_flag {
			stop_flag = true
		} else {
			stop_flag = false
		}

		inputBuf.WriteString(lastInput)
		inputBuf.WriteString("\n")
	}
	inv.Notes = inputBuf.String()

	fmt.Print("Payment method: ")
	stop_flag = false
	inputBuf = bytes.NewBuffer([]byte{})

	for {
		scanner.Scan()
		lastInput := scanner.Text()

		// Press enter on an empty line to finish
		if lastInput == "" && stop_flag {
			break
		} else if lastInput == "" && !stop_flag {
			stop_flag = true
		} else {
			stop_flag = false
		}

		inputBuf.WriteString(lastInput)
		inputBuf.WriteString("\n")
	}
	inv.PayMethod = inputBuf.String()

	fmt.Print("Footer: ")
	scanner.Scan()
	inv.Footer = scanner.Text()

	// Summary
	fmt.Printf(`
--- Summary ---
Invoice number: %s%d
Date: %s
Customer: %d - %s : %s : %s : %s : %s : %s : %s : %s
Issuer: %d - %s : %s : %s : %s : %s : %s : %s : %s
ISO Currency: %s
Items: %v
Subtotal: %f
Discount: %f
Tax: %f
Total: %f
Final: %f
Notes: %s
Payment method: %s
Footer: %s

`, inv.Series, inv.Number, inv.Date, inv.Customer.Id, inv.Customer.Name, inv.Customer.TinNumber, inv.Customer.Address, inv.Customer.City, inv.Customer.PostalCode, inv.Customer.Country, inv.Customer.Phone, inv.Customer.Email, inv.Issuer.Id, inv.Issuer.Name, inv.Issuer.TinNumber, inv.Issuer.Address, inv.Issuer.City, inv.Issuer.PostalCode, inv.Issuer.Country, inv.Issuer.Phone, inv.Issuer.Email, inv.IsoCurrency, inv.Items, inv.SubTotal, inv.Discount, inv.Tax, inv.Total, inv.Final, inv.Notes, inv.PayMethod, inv.Footer)
	fmt.Print("Is this information correct? (y/N): ")
	usrInput = "N"
	fmt.Scanln(&usrInput)
	if usrInput == "y" {
		// Save invoice data to database
		invJson, err := json.Marshal(inv)
		if err != nil {
			log.Fatalf("Error marshalling invoice data: %s", err)
		}
		_, err = db.Exec("INSERT INTO IssuedInvoices (Series, Number, CustomerId, Date, Data) VALUES (?, ?, ?, ?, ?)", inv.Series, inv.Number, inv.Customer.Id, inv.Date, string(invJson))
		if err != nil {
			log.Fatalf("Error saving invoice data to database: %s", err)
		}
		fmt.Println("Invoice data saved successfully")
	} else {
		fmt.Println("Invoice data not saved")
	}

	// Go back to main menu
	*state = 000
}

// CLI for listing invoices from the database
func listInvoice(db *sql.DB, state *int) {
	rows, err := db.Query("SELECT Id, Series, Number, CustomerId, Date FROM IssuedInvoices")
	if err != nil {
		log.Fatalf("Error querying invoice data: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var inv_data IssuedInvoice
		err = rows.Scan(&inv_data.Id, &inv_data.Series, &inv_data.Number, &inv_data.Customer.Id, &inv_data.Date)
		if err != nil {
			log.Fatalf("Error scanning invoice data: %s", err)
		}
		// Get customer data
		row := db.QueryRow("SELECT Name FROM Customers WHERE Id = ?", inv_data.Customer.Id)
		err = row.Scan(&inv_data.Customer.Name)
		if err != nil {
			log.Fatalf("Error querying customer data: %s", err)
		}
		fmt.Printf("%d - %s%06d : %s : %s\n", inv_data.Id, inv_data.Series, inv_data.Number, inv_data.Date, inv_data.Customer.Name)
	}

	// Go back to main menu
	*state = 000
}

// Command line interface
func cli(db *sql.DB) {
	var state int

	// Main CLI loop
	for {
		switch state {

		// Main menu
		case 000:
			fmt.Println(`

------ Main Menu ------
Please select an option:
  001 - Set issuer data
  002 - List issuer data
  003 - Delete issuer data
  
  101 - Create new customer
  102 - List existing customer
  103 - Delete existing customer
  
  201 - Create new invoice
  202 - List existing invoice
  
  999 - Quit`)
			fmt.Print("> ")
			fmt.Scanln(&state)

		// Set issuer info 0xx
		case 001: // Set issuer data
			createIssuer(db, &state)
		case 002: // List issuer data
			listIssuer(db, &state)
		case 003: // Delete issuer data
			deleteIssuer(db, &state)

		// Customer info 1xx
		case 101: // Create new customer
			createCustomer(db, &state)
		case 102: // List existing customer
			listCustomer(db, &state)
		case 103: // Delete existing customer
			deleteCustomer(db, &state)

		// Invoice info 2xx
		case 201: // Create new invoice
			createInvoice(db, &state)
		case 202: // List existing invoice
			listInvoice(db, &state)

		// Quit
		case 999: // Quit
			return
		}
	}
}

// Web server to show rendered invoices
func webServer(db *sql.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, "static/index.html")
	})

	http.HandleFunc("/render/{page}", func(w http.ResponseWriter, r *http.Request) {
		page := r.PathValue("page")
		log.Printf("Request received: %s", page)

		// Look for the requested invoice id
		var inv invoice
		var inv_data string
		row := db.QueryRow("SELECT Data FROM IssuedInvoices WHERE Id = ?", page)
		err := row.Scan(&inv_data)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Invoice not found"))
			return
		}

		// Unmarshal invoice data
		err = json.Unmarshal([]byte(inv_data), &inv)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error unmarshalling invoice data"))
			return
		}

		// Render invoice
		w.Header().Set("Content-Type", "text/html")
		w.Write(inv.render().Bytes())
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	// Flags
	config_path := flag.String("c", "../config/config.json", "Path to main config file. (Default: ../config/config.json)")
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

	// Run the web server
	go webServer(db)

	// Run the command line interface
	fmt.Println(`
       d8888 8888888 8888888b.  
      d88888   888   888   Y88b 
     d88P888   888   888    888 
    d88P 888   888   888   d88P 
   d88P  888   888   8888888P"  
  d88P   888   888   888 T88b   
 d8888888888   888   888  T88b  
d88P     888 8888888 888   T88b 


Welcome to Air!`)
	cli(db)

}
