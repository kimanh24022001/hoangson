package types

import (
	"database/sql/driver"
	"errors"
	"math/big"
	"reflect"
	"time"

	"github.com/google/uuid"
)

type (
	Uuid     uuid.UUID
	String   string
	B1       bool
	I16      int16
	I32      int32
	I64      int64
	F32      float32
	F64      float64
	Serial16 int16
	Serial32 int32
	Serial64 int64
	Numeric  big.Int
	Time     time.Time
)

var (
	ReflectUUID     = reflect.TypeOf(Uuid{})
	ReflectString   = reflect.TypeOf(String(""))
	ReflectB1       = reflect.TypeOf(B1(false))
	ReflectI16      = reflect.TypeOf(I16(0))
	ReflectI32      = reflect.TypeOf(I32(0))
	ReflectI64      = reflect.TypeOf(I64(0))
	ReflectF32      = reflect.TypeOf(F32(0))
	ReflectF64      = reflect.TypeOf(F64(0))
	ReflectSerial16 = reflect.TypeOf(Serial16(0))
	ReflectSerial32 = reflect.TypeOf(Serial32(0))
	ReflectSerial64 = reflect.TypeOf(Serial64(0))
	ReflectNumeric  = reflect.TypeOf(Numeric{})
	ReflectTime     = reflect.TypeOf(Time{})
)

// TODO(duong): add comments??
func (self Uuid) String() string {
	return (uuid.UUID)(self).String()
}

func (self *Uuid) Scan(src any) error {
	return (*uuid.UUID)(self).Scan(src)
}

func (self Uuid) Value() (driver.Value, error) {
	return self.String(), nil
}

func (self Uuid) MarshalText() ([]byte, error) {
	return uuid.UUID(self).MarshalText()
}

func (self *Uuid) UnmarshalText(data []byte) error {
	return (*uuid.UUID)(self).UnmarshalText(data)
}

func (self Uuid) MarshalBinary() ([]byte, error) {
	return uuid.UUID(self).MarshalBinary()
}

func (self *Uuid) UnmarshalBinary(data []byte) error {
	return (*uuid.UUID)(self).UnmarshalBinary(data)
}

func (self Time) String() string {
	return time.Time(self).Format(time.RFC3339Nano)
}

func (self *Time) Scan(src any) error {
    switch t := src.(type) {
    case time.Time:
        *self = Time(t)
        return nil
    default:
        return errors.New("Column type is not types.Time")
	}
}

func (self Time) Value() (driver.Value, error) {
	return time.Time(self), nil
}

func (self Time) MarshalText() ([]byte, error) {
	return time.Time(self).MarshalText()
}

func (self *Time) UnmarshalText(data []byte) error {
	return (*time.Time)(self).UnmarshalText(data)
}

func (self Time) MarshalBinary() ([]byte, error) {
	return time.Time(self).MarshalBinary()
}

func (self *Time) UnmarshalBinary(data []byte) error {
	return (*time.Time)(self).UnmarshalBinary(data)
}
