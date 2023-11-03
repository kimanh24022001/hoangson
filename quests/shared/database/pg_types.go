package database

import (
	"reflect"

	"smatyx.com/shared/types"
)

const (
	PgTypeUUID = iota
	PgTypeText
	PgTypeBoolean
	PgTypeSmallInt
	PgTypeInteger
	PgTypeBigInt
	PgTypeReal
	PgTypeDoublePrecision
	PgTypeSmallSerial
	PgTypeSerial
	PgTypeBigSerial
	PgTypeTime
	PgTypeNumeric
	PgTypeCount
)

var PgTypeReflectNames = []string{
	PgTypeUUID:            "PgTypeUUID",
	PgTypeText:            "PgTypeText",
	PgTypeBoolean:         "PgTypeBoolean",
	PgTypeSmallInt:        "PgTypeSmallInt",
	PgTypeInteger:         "PgTypeInteger",
	PgTypeBigInt:          "PgTypeBigInt",
	PgTypeReal:            "PgTypeReal",
	PgTypeDoublePrecision: "PgTypeDoublePrecision",
	PgTypeSmallSerial:     "PgTypeSmallSerial",
	PgTypeSerial:          "PgTypeSerial",
	PgTypeBigSerial:       "PgTypeBigSerial",
	PgTypeTime:            "PgTypeTime",
	PgTypeNumeric:         "PgTypeNumeric",
}

var PgTypeNames = []string{
	PgTypeUUID:            "UUID",
	PgTypeText:            "TEXT",
	PgTypeBoolean:         "BOOLEAN",
	PgTypeSmallInt:        "SMALLINT",
	PgTypeInteger:         "INTERGER",
	PgTypeBigInt:          "BIGINT",
	PgTypeReal:            "REAL",
	PgTypeDoublePrecision: "DOUBLE PRECISION",
	PgTypeSmallSerial:     "SMALLSERIAL",
	PgTypeSerial:          "SERIAL",
	PgTypeBigSerial:       "BIGSERIAL",
	PgTypeTime:            "TIMESTAMP WITHOUT TIME ZONE",
	PgTypeNumeric:         "NUMERIC(32)",
}

var PgTypeReturns = func() map[string]int {
	result := make(map[string]int)

	for i, name := range PgTypeNames {
		result[name] = i
	}

	return result
}()

var PgReflectTypes = map[reflect.Type]int{
	types.ReflectUUID:     PgTypeUUID,
	types.ReflectString:   PgTypeText,
	types.ReflectB1:       PgTypeBoolean,
	types.ReflectI16:      PgTypeSmallInt,
	types.ReflectI32:      PgTypeInteger,
	types.ReflectI64:      PgTypeBigInt,
	types.ReflectF32:      PgTypeReal,
	types.ReflectF64:      PgTypeDoublePrecision,
	types.ReflectSerial16: PgTypeSmallSerial,
	types.ReflectSerial32: PgTypeSerial,
	types.ReflectSerial64: PgTypeBigSerial,
	types.ReflectTime:     PgTypeTime,
	types.ReflectNumeric:  PgTypeNumeric,
}

type EntityCallType int

const (
	EntityCallGet EntityCallType = iota
	EntityCallCount
	EntityCallSearch
	EntityCallCreate
	EntityCallUpdate
	EntityCallDelete
	EntityCallExist
)
// type PgError pq.Error
