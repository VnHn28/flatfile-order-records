package controller

import (
	"flatfile-order-records/internal/db"
	"flatfile-order-records/internal/model"
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Fyne-based GUI controller
type GUI struct {
	database *db.Database
	app      fyne.App
	window   fyne.Window
}

func NewGUI(dbPath string) *GUI {
	a := app.New()
	win := a.NewWindow("Order Records Flat-File Database")

	gui := &GUI{
		database: db.NewDatabase(dbPath),
		app:      a,
		window:   win,
	}
	return gui
}

// Starts the Fyne application
func (g *GUI) Run() {
	g.window.SetContent(g.createContent())
	g.window.Resize(fyne.NewSize(800, 600))
	g.window.ShowAndRun()
}

// Builds main UI layout
func (g *GUI) createContent() fyne.CanvasObject {
	// --- Display Area ---
	display := widget.NewMultiLineEntry()
	display.Disable()
	display.Wrapping = fyne.TextWrapWord

	// --- Insert Form ---
	insertOrderIDEntry := widget.NewEntry()
	insertOrderIDEntry.SetPlaceHolder("Order ID (e.g., 101)")
	insertOwnerEntry := widget.NewEntry()
	insertOwnerEntry.SetPlaceHolder("Owner Name")
	insertAmountEntry := widget.NewEntry()
	insertAmountEntry.SetPlaceHolder("Amount (e.g., 500)")

	var insertForm *widget.Form
	insertForm = &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Order ID", Widget: insertOrderIDEntry},
			{Text: "Owner", Widget: insertOwnerEntry},
			{Text: "Amount", Widget: insertAmountEntry},
		},
		OnSubmit: func() {
			order, ok := g.parseOrderForm(insertOrderIDEntry, insertOwnerEntry, insertAmountEntry)
			if !ok {
				return
			}

			if err := g.database.Insert(order); err != nil {
				dialog.ShowError(err, g.window)
			} else {
				dialog.ShowInformation("Success", "Record inserted successfully.", g.window)
				insertForm.Items[0].Widget.(*widget.Entry).SetText("")
				insertForm.Items[1].Widget.(*widget.Entry).SetText("")
				insertForm.Items[2].Widget.(*widget.Entry).SetText("")
				insertForm.Refresh()
			}
		},
	}

	// --- Update Form ---
	updateOrderIDEntry := widget.NewEntry()
	updateOrderIDEntry.SetPlaceHolder("Order ID to update")
	updateOwnerEntry := widget.NewEntry()
	updateOwnerEntry.SetPlaceHolder("New Owner Name")
	updateAmountEntry := widget.NewEntry()
	updateAmountEntry.SetPlaceHolder("New Amount")

	var updateForm *widget.Form
	updateForm = &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Order ID", Widget: updateOrderIDEntry},
			{Text: "New Owner", Widget: updateOwnerEntry},
			{Text: "New Amount", Widget: updateAmountEntry},
		},
		OnSubmit: func() {
			order, ok := g.parseOrderForm(updateOrderIDEntry, updateOwnerEntry, updateAmountEntry)
			if !ok {
				return
			}

			if err := g.database.UpdateById(order.OrderID, order); err != nil {
				dialog.ShowError(err, g.window)
			} else {
				dialog.ShowInformation("Success", "Record updated successfully.", g.window)
				updateForm.Items[0].Widget.(*widget.Entry).SetText("")
				updateForm.Items[1].Widget.(*widget.Entry).SetText("")
				updateForm.Items[2].Widget.(*widget.Entry).SetText("")
				updateForm.Refresh()
			}
		},
	}

	// --- Search ---
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search by Order ID...")
	searchButton := widget.NewButton("Search", func() {
		if searchEntry.Text == "" {
			dialog.ShowInformation("Info", "Search field is empty.", g.window)
			return
		}
		orderID, err := strconv.ParseInt(searchEntry.Text, 10, 64)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid Order ID format"), g.window)
			return
		}
		records, err := g.database.SearchById(orderID)
		if err != nil {
			dialog.ShowError(err, g.window)
			return
		}
		g.updateDisplay(display, records)
	})

	// --- Display All Button ---
	displayAllButton := widget.NewButton("Display All Records", func() {
		records, err := g.database.ReadAll()
		if err != nil {
			dialog.ShowError(err, g.window)
			return
		}
		g.updateDisplay(display, records)
	})

	// --- Layout ---
	forms := container.NewAppTabs(
		container.NewTabItem("Insert Record", insertForm),
		container.NewTabItem("Update Record", updateForm),
	)

	actions := container.NewVBox(
		container.NewGridWithColumns(2, searchEntry, searchButton),
		displayAllButton,
	)

	leftPanel := container.NewVBox(forms, actions)
	return container.NewBorder(nil, nil, leftPanel, nil, display)
}

// Formats and shows records in the display area
func (g *GUI) updateDisplay(display *widget.Entry, records []*model.Order) {
	if len(records) == 0 {
		display.SetText("No records found.")
		return
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Found %d record(s):\n\n", len(records)))
	for _, r := range records {
		// Trim null bytes for clean display
		owner := strings.Trim(string(r.Owner[:]), "\x00")
		line := fmt.Sprintf("OrderID: %d\nOwner: %s\nAmount: %d\n-----------------\n", r.OrderID, owner, r.Amount)
		builder.WriteString(line)
	}
	display.SetText(builder.String())
}

// Reads and validates data from form entries
func (g *GUI) parseOrderForm(idEntry, ownerEntry, amountEntry *widget.Entry) (model.Order, bool) {
	orderID, err := strconv.ParseInt(idEntry.Text, 10, 64)
	if err != nil {
		dialog.ShowError(fmt.Errorf("invalid Order ID: must be a number"), g.window)
		return model.Order{}, false
	}

	amount, err := strconv.ParseInt(amountEntry.Text, 10, 64)
	if err != nil {
		dialog.ShowError(fmt.Errorf("invalid Amount: must be a number"), g.window)
		return model.Order{}, false
	}

	var owner [32]byte
	copy(owner[:], ownerEntry.Text)

	return model.Order{
		OrderID: orderID,
		Owner:   owner,
		Amount:  amount,
	}, true
}
