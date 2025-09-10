package main

import (
	"flatfile-orders-record/internal/controller"
	"log"
	"os"
)

const dbFile = "orders.db"

func main() {
	log.SetOutput(os.Stdout)
	log.Println("Starting the flat-file database orders record system...")

	cli := controller.NewCLI(dbFile)
	cli.Run()
}
