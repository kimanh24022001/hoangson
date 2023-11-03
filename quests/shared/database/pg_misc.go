package database

import (
	"errors"
	"strconv"
	"strings"
)

var PresetVariables = initPresetVariables(10_000)

func initPresetVariables(capacity int) []string {
	result := make([]string, capacity)

	for i := 0; i < capacity; i++ {
		result[i] = "$" + strconv.FormatInt(int64(i+1), 10)
	}

	return result
}

func StringVariables(lo, hi int) string {
	return strings.Join(PresetVariables[lo:hi], ",")
}

// NOTE(duong): Build the string with content:
//
// INSERT INTO %v(%v) VALUES (%v),(%v),(%v) ...
func StringPgInsert(
	tableName, columns string,
	startIdx, columnCount, rowCount int) string {

	sb := strings.Builder{}
	// TODO: pre-allocate

	// NOTE(duong): table name might not contain any ", so this might
	// be redundant.
	sb.WriteString(`INSERT INTO "`)
	sb.WriteString(strings.ReplaceAll(tableName, `"`, `""`))
	sb.WriteString(`" (`)
	sb.WriteString(columns)
	sb.WriteString(`) VALUES `)

	for i := 0; i < rowCount; i++ {
		if i != 0 {
			sb.WriteRune(',')
		}

		endIdx := startIdx + columnCount
		variables := strings.Join(PresetVariables[startIdx:endIdx], ",")
		startIdx = endIdx

		sb.WriteRune('(')
		sb.WriteString(variables)
		sb.WriteRune(')')
	}

	return sb.String()
}

type SearchCondition struct {
	Operator int
	Field    string
	Value    string
}

func StringPgUpdate(
	tableName string, columns []string,
	startIdx int, searchConditions string) string {

	sb := strings.Builder{}

	sb.WriteString(`UPDATE "`)
	sb.WriteString(strings.ReplaceAll(tableName, `"`, `""`))
	sb.WriteString(`" SET`)

	for i, column := range columns {
		if i != 0 {
			sb.WriteRune(',')
		}

		sb.WriteRune('"')
		sb.WriteString(strings.ReplaceAll(column, `"`, `""`))
		sb.WriteString(`"=`)
		sb.WriteString(PresetVariables[startIdx])

		startIdx++
	}

	return sb.String()
}

func ValidateAffectedCountEqual(expecetd int) ProcessRepFunc {
	result := func(rep *PgRep) error {
		affected := rep.RowsAffected()

		if affected != expecetd {
			return errors.New("Process validate error")
		}
		return nil
	}

	return result
}
