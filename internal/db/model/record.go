package model

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Order struct {
	OrderID int64
	Owner   [32]byte
	Amount  int64
}

// 8 (OrderID) + 32 (Owner) + 8 (Amount) = 48 bytes.
const RecordSize = 48

// Serializes the Order struct into a byte slice.
func (o *Order) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, o.OrderID); err != nil {
		return nil, fmt.Errorf("failed to write OrderID: %w", err)
	}
	if err := binary.Write(buf, binary.LittleEndian, o.Owner); err != nil {
		return nil, fmt.Errorf("failed to write Owner: %w", err)
	}
	if err := binary.Write(buf, binary.LittleEndian, o.Amount); err != nil {
		return nil, fmt.Errorf("failed to write Amount: %w", err)
	}
	return buf.Bytes(), nil
}

// Deserializes a byte slice into an Order struct.
func (o *Order) UnmarshalBinary(data []byte) error {
	if len(data) != RecordSize {
		return fmt.Errorf("data size mismatch: expected %d bytes, got %d", RecordSize, len(data))
	}

	buf := bytes.NewReader(data)
	if err := binary.Read(buf, binary.LittleEndian, &o.OrderID); err != nil {
		return fmt.Errorf("failed to read OrderID: %w", err)
	}
	if err := binary.Read(buf, binary.LittleEndian, &o.Owner); err != nil {
		return fmt.Errorf("failed to read Owner: %w", err)
	}
	if err := binary.Read(buf, binary.LittleEndian, &o.Amount); err != nil {
		return fmt.Errorf("failed to read Amount: %w", err)
	}
	return nil
}
