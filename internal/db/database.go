package db

import (
	"flatfile-order-records/internal/model"
	"fmt"
	"io"
	"os"
	"sync"
)

type Database struct {
	mu       sync.RWMutex
	filePath string
}

func NewDatabase(filePath string) *Database {
	return &Database{
		filePath: filePath,
	}
}

func (d *Database) ensureFileExists() (*os.File, error) {
	return os.OpenFile(d.filePath, os.O_RDWR|os.O_CREATE, 0644)
}

func (d *Database) Insert(order model.Order) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	file, err := d.ensureFileExists()
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		return fmt.Errorf("failed to seek to end of file: %w", err)
	}

	data, err := order.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write record: %w", err)
	}

	return nil
}

func (d *Database) ReadAll() ([]*model.Order, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.scan(nil)
}

func (d *Database) UpdateById(orderID int64, newOrder model.Order) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	file, err := d.ensureFileExists()
	if err != nil {
		return err
	}
	defer file.Close()

	offset := int64(0)
	for {
		data := make([]byte, model.RecordSize)
		n, err := file.Read(data)
		if err == io.EOF {
			return fmt.Errorf("record with OrderID %d not found", orderID)
		}
		if err != nil {
			return fmt.Errorf("failed to read record: %w", err)
		}
		if n < model.RecordSize {
			return fmt.Errorf("corrupt file or incomplete record")
		}

		var record model.Order
		if err := record.UnmarshalBinary(data); err != nil {
			return err
		}

		if record.OrderID == orderID {
			if _, err := file.Seek(offset, io.SeekStart); err != nil {
				return err
			}

			newData, err := newOrder.MarshalBinary()
			if err != nil {
				return err
			}

			_, err = file.Write(newData)
			if err != nil {
				return fmt.Errorf("failed to write updated record: %w", err)
			}
			return nil
		}
		offset += model.RecordSize
	}
}

func (d *Database) SearchById(orderID int64) ([]*model.Order, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	filter := func(order *model.Order) bool {
		return order.OrderID == orderID
	}

	return d.scan(filter)
}

// Iterate through the database file and returns records that match the filter
// If filter is nil then return all records
func (d *Database) scan(filter func(*model.Order) bool) ([]*model.Order, error) {
	file, err := d.ensureFileExists()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []*model.Order
	data := make([]byte, model.RecordSize)

	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read from file: %w", err)
		}
		if n < model.RecordSize {
			break
		}

		var order model.Order
		if err := order.UnmarshalBinary(data); err != nil {
			continue
		}

		if filter == nil || filter(&order) {
			records = append(records, &order)
		}
	}
	return records, nil
}
