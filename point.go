// Based on https://github.com/jinzhu/gorm/issues/142
package spatial

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type Point struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}

func (p *Point) Scan(val interface{}) error {
	b, err := hex.DecodeString(string(val.([]uint8)))
	if err != nil {
		return err
	}
	r := bytes.NewReader(b)
	var wkbByteOrder uint8
	if err := binary.Read(r, binary.LittleEndian, &wkbByteOrder); err != nil {
		return err
	}

	var byteOrder binary.ByteOrder
	switch wkbByteOrder {
	case 0:
		byteOrder = binary.BigEndian
	case 1:
		byteOrder = binary.LittleEndian
	default:
		return fmt.Errorf("invalid byte order %d", wkbByteOrder)
	}

	var wkbGeometryType uint64
	if err := binary.Read(r, byteOrder, &wkbGeometryType); err != nil {
		return err
	}

	if err := binary.Read(r, byteOrder, p); err != nil {
		return err
	}

	return nil
}

func (p Point) Value() (driver.Value, error) {
	return p, nil
}

// MarshalText implements encoding.TextMarshaler.
func (p Point) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("SRID=4326;POINT(%v %v)", p.Lng, p.Lat)), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *Point) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		return nil
	}

	*s = Point{}
	err := s.Scan(text)
	return err
}