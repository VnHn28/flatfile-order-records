package controller

import (
	"bufio"
	"flatfile-orders-record/internal/db"
	"flatfile-orders-record/internal/db/model"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const dbFile = "orders.db"

type CLI struct {
	database *db.Database
	reader   *bufio.Reader
}

func NewCLI() *CLI {
	return &CLI{
		database: db.NewDatabase(dbFile),
		reader:   bufio.NewReader(os.Stdin),
	}
}

func (c *CLI) Run() {
	fmt.Println("Orders Record Flat-File Database System")
	fmt.Println("---------------------------------------")

	for {
		c.printMenu()
		input, err := c.reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading input:", err)
			return
		}
		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			c.insertRecord()
		case "2":
			c.updateRecord()
		case "3":
			c.readAllRecords()
		case "4":
			c.searchRecords()
		case "5":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func (c *CLI) printMenu() {
	fmt.Println("\nSelect an action:")
	fmt.Println("1. Insert new record")
	fmt.Println("2. Update record by OrderID")
	fmt.Println("3. Read all records")
	fmt.Println("4. Search records by OrderID")
	fmt.Println("5. Exit")
	fmt.Print("> ")
}

func (c *CLI) insertRecord() {
	fmt.Print("Enter OrderID: ")
	orderIDStr, err := c.reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading OrderID:", err)
		return
	}
	orderID, err := strconv.ParseInt(strings.TrimSpace(orderIDStr), 10, 64)
	if err != nil {
		fmt.Println("Invalid OrderID:", err)
		return
	}

	fmt.Print("Enter Owner (max 32 chars): ")
	ownerStr, err := c.reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading Owner:", err)
		return
	}
	ownerStr = strings.TrimSpace(ownerStr)
	var owner [32]byte
	copy(owner[:], ownerStr)

	fmt.Print("Enter Amount: ")
	amountStr, err := c.reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading Amount:", err)
		return
	}
	amount, err := strconv.ParseInt(strings.TrimSpace(amountStr), 10, 64)
	if err != nil {
		fmt.Println("Invalid Amount:", err)
		return
	}

	order := model.Order{
		OrderID: orderID,
		Owner:   owner,
		Amount:  amount,
	}

	if err := c.database.Insert(order); err != nil {
		log.Println("Error inserting record:", err)
	} else {
		fmt.Println("Record inserted successfully.")
	}
}

func (c *CLI) readAllRecords() {
	records, err := c.database.ReadAll()
	if err != nil {
		log.Println("Error reading records:", err)
		return
	}
	if len(records) == 0 {
		fmt.Println("No records found.")
		return
	}
	for _, record := range records {
		fmt.Printf("OrderID: %d, Owner: %s, Amount: %d\n",
			record.OrderID,
			strings.Trim(string(record.Owner[:]), "\x00"),
			record.Amount)
	}
}

func (c *CLI) updateRecord() {
	fmt.Print("Enter OrderID to update: ")
	orderIDStr, err := c.reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading OrderID:", err)
		return
	}
	orderID, err := strconv.ParseInt(strings.TrimSpace(orderIDStr), 10, 64)
	if err != nil {
		fmt.Println("Invalid OrderID:", err)
		return
	}

	fmt.Print("Enter new Owner (max 32 chars): ")
	ownerStr, err := c.reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading Owner:", err)
		return
	}
	ownerStr = strings.TrimSpace(ownerStr)
	var owner [32]byte
	copy(owner[:], ownerStr)

	fmt.Print("Enter new Amount: ")
	amountStr, err := c.reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading Amount:", err)
		return
	}
	amount, err := strconv.ParseInt(strings.TrimSpace(amountStr), 10, 64)
	if err != nil {
		fmt.Println("Invalid Amount:", err)
		return
	}

	newOrder := model.Order{
		OrderID: orderID,
		Owner:   owner,
		Amount:  amount,
	}

	if err := c.database.UpdateById(orderID, newOrder); err != nil {
		log.Println("Error updating record:", err)
	} else {
		fmt.Println("Record updated successfully.")
	}
}

func (c *CLI) searchRecords() {
	fmt.Print("Enter OrderID to search for: ")
	orderIDStr, err := c.reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading OrderID:", err)
		return
	}
	orderID, err := strconv.ParseInt(strings.TrimSpace(orderIDStr), 10, 64)
	if err != nil {
		fmt.Println("Invalid OrderID:", err)
		return
	}

	records, err := c.database.SearchById(orderID)
	if err != nil {
		log.Println("Error searching for records:", err)
		return
	}
	if len(records) == 0 {
		fmt.Println("No records found with that OrderID.")
		return
	}
	for _, record := range records {
		fmt.Printf("OrderID: %d, Owner: %s, Amount: %d\n",
			record.OrderID,
			strings.Trim(string(record.Owner[:]), "\x00"),
			record.Amount)
	}
}
