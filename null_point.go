package spatial

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
)

type NullPoint struct {
	Point Point
	Valid bool
}

func (np *NullPoint) Scan(val interface{}) error {
	if val == nil {
		np.Point, np.Valid = Point{}, false
		return nil
	}

	point := &Point{}
	err := point.Scan(val)
	if err != nil {
		np.Point, np.Valid = Point{}, false
		return nil
	}
	np.Point = Point{
		Lat: point.Lat,
		Lng: point.Lng,
	}
	np.Valid = true

	return nil
}

func (np NullPoint) Value() (driver.Value, error) {
	if !np.Valid {
		return nil, nil
	}
	return np.Point, nil
}

// NullPointFrom creates a new NullPoint that will never be blank.
func NullPointFrom(s Point) NullPoint {
	return NewNullPoint(s, true)
}

// NullPointFromPtr creates a new NullPoint that be null if s is nil.
func NullPointFromPtr(s *Point) NullPoint {
	if s == nil {
		return NewNullPoint(Point{}, false)
	}
	return NewNullPoint(*s, true)
}

// NewNullPoint creates a new NullPoint
func NewNullPoint(s Point, valid bool) NullPoint {
	return NullPoint{
		Point: s,
		Valid: valid,
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *NullPoint) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		s.Point = Point{}
		s.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &s.Point); err != nil {
		return err
	}

	s.Valid = true
	return nil
}

// NullBytes is a global byte slice of JSON null
var NullBytes = []byte("null")

// MarshalJSON implements json.Marshaler.
func (s NullPoint) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return NullBytes, nil
	}
	return json.Marshal(s.Point)
}

// MarshalText implements encoding.TextMarshaler.
func (s NullPoint) MarshalText() ([]byte, error) {
	if !s.Valid {
		return []byte{}, nil
	}
	return s.Point.MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *NullPoint) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		s.Valid = false
		return nil
	}

	s.Point = Point{}
	err := s.Point.UnmarshalText(text)
	s.Valid = err == nil
	return nil
}

// SetValid changes this NullPoint's value and also sets it to be non-null.
func (s *NullPoint) SetValid(v Point) {
	s.Point = v
	s.Valid = true
}

// Ptr returns a pointer to this NullPoint's value, or a nil pointer if this NullPoint is null.
func (s NullPoint) Ptr() *Point {
	if !s.Valid {
		return nil
	}
	return &s.Point
}

// IsZero returns true for null Points, for potential future omitempty support.
func (s NullPoint) IsZero() bool {
	return !s.Valid
}

// Randomize for sqlboiler
func (s *NullPoint) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	if shouldBeNull {
		s.Point = Point{}
		s.Valid = false
	} else {
		s.Point = Point{
			Lng: float64(nextInt()%10)/10.0 + float64(nextInt()%10),
			Lat: float64(nextInt()%10)/10.0 + float64(nextInt()%10),
		}
		s.Valid = true
	}
}
