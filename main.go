package main

import (
	"flag"
	"flatfile-order-records/internal/controller"
	"log"
	"os"
)

const dbFile = "orders.db"

func main() {
	log.SetOutput(os.Stdout)

	guiFlag := flag.Bool("gui", false, "Run the graphical user interface")
	flag.Parse()

	if *guiFlag {
		log.Println("Starting the GUI application...")
		gui := controller.NewGUI(dbFile)
		gui.Run()
	} else {
		log.Println("Starting the CLI application...")
		cli := controller.NewCLI(dbFile)
		cli.Run()
	}
}
