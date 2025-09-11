package controller

import (
	"bufio"
	"flatfile-order-records/internal/db"
	"flatfile-order-records/internal/model"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type CLI struct {
	database *db.Database
	reader   *bufio.Reader
}

func NewCLI(dbPath string) *CLI {
	return &CLI{
		database: db.NewDatabase(dbPath),
		reader:   bufio.NewReader(os.Stdin),
	}
}

func (c *CLI) Run() {
	fmt.Println("Order Records Flat-File Database System")
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
	orderID, ok := c.inputInt64("Enter OrderID: ")
	if !ok {
		return
	}

	ownerStr, ok := c.inputString("Enter Owner (max 32 chars): ")
	if !ok {
		return
	}
	var owner [32]byte
	copy(owner[:], ownerStr)

	amount, ok := c.inputInt64("Enter Amount: ")
	if !ok {
		return
	}

	order := model.Order{OrderID: orderID, Owner: owner, Amount: amount}
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
	c.displayRecords(records)
}

func (c *CLI) updateRecord() {
	orderID, ok := c.inputInt64("Enter OrderID to update: ")
	if !ok {
		return
	}

	ownerStr, ok := c.inputString("Enter new Owner (max 32 chars): ")
	if !ok {
		return
	}
	var owner [32]byte
	copy(owner[:], ownerStr)

	amount, ok := c.inputInt64("Enter new Amount: ")
	if !ok {
		return
	}

	newOrder := model.Order{OrderID: orderID, Owner: owner, Amount: amount}
	if err := c.database.UpdateById(orderID, newOrder); err != nil {
		log.Println("Error updating record:", err)
	} else {
		fmt.Println("Record updated successfully.")
	}
}

func (c *CLI) searchRecords() {
	orderID, ok := c.inputInt64("Enter OrderID to search for: ")
	if !ok {
		return
	}

	records, err := c.database.SearchById(orderID)
	if err != nil {
		log.Println("Error searching for records:", err)
		return
	}
	c.displayRecords(records)
}

func (c *CLI) inputInt64(prompt string) (int64, bool) {
	fmt.Print(prompt)
	str, err := c.reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading input:", err)
		return 0, false
	}
	val, err := strconv.ParseInt(strings.TrimSpace(str), 10, 64)
	if err != nil {
		fmt.Println("Invalid number format:", err)
		return 0, false
	}
	return val, true
}

func (c *CLI) inputString(prompt string) (string, bool) {
	fmt.Print(prompt)
	str, err := c.reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading input:", err)
		return "", false
	}
	return strings.TrimSpace(str), true
}

func (c *CLI) displayRecords(records []*model.Order) {
	if len(records) == 0 {
		fmt.Println("No records found.")
		return
	}
	fmt.Println("--------------- Records ---------------")
	for _, record := range records {
		fmt.Printf("OrderID: %d, Owner: %s, Amount: %d\n",
			record.OrderID,
			strings.Trim(string(record.Owner[:]), "\x00"),
			record.Amount)
	}
	fmt.Println("---------------------------------------")
}
