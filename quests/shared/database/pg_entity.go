package database

import (
	"fmt"
	"reflect"
	"strings"

	"smatyx.com/shared/cast"
)

type PgEntity struct {
	Name   string
	GoName string

	FieldNames   []string
	FieldGoNames []string
	FieldPgTypes []int
	FieldGoTypes []reflect.Type

	FieldFullNames []string

	Indices map[string][]string
	Uniques map[string][]string
}

// Create new PgEntity from struct
func NewPgEntity(entityName string, s any) *PgEntity {
	tp := reflect.TypeOf(s)
	tpFieldCount := tp.NumField()

	result := &PgEntity{
		Name:           entityName,
		GoName:         tp.Name(),
		FieldNames:     make([]string, 0, tpFieldCount),
		FieldFullNames: make([]string, 0, tpFieldCount),
		FieldPgTypes:   make([]int, 0, tpFieldCount),
		FieldGoNames:   make([]string, 0, tpFieldCount),
		FieldGoTypes:   make([]reflect.Type, 0, tpFieldCount),
		Indices:        make(map[string][]string),
		Uniques:        make(map[string][]string),
	}

	for i := 0; i < tpFieldCount; i++ {
		field := tp.Field(i)
		fieldName := field.Name
		fieldType := field.Type
		fieldTag := field.Tag

		fieldUnique := fieldTag.Get("unique")
		fieldIndex := fieldTag.Get("index")

		// fieldType = fieldType.Elem()

		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		if fieldType.Kind() == reflect.Array {
			// TODO
		}

		pgType, ok := PgReflectTypes[fieldType]
		if ok {
			name := cast.StringLowerSnakeCase(fieldName)
			result.FieldNames = append(result.FieldNames, name)
			result.FieldFullNames = append(result.FieldFullNames, makeFullFieldName(entityName, name))
			result.FieldPgTypes = append(result.FieldPgTypes, pgType)
			result.FieldGoNames = append(result.FieldGoNames, fieldName)
			result.FieldGoTypes = append(result.FieldGoTypes, fieldType)

			if len(fieldUnique) != 0 {
				if result.Uniques[fieldUnique] == nil {
					result.Uniques[fieldUnique] = []string{}
				}
				result.Uniques[fieldUnique] = append(result.Uniques[fieldUnique], name)
			}

			if len(fieldIndex) != 0 {
				if result.Indices[fieldIndex] == nil {
					result.Indices[fieldIndex] = []string{}
				}
				result.Indices[fieldIndex] = append(result.Uniques[fieldUnique], name)
			}
		}
	}

	return result
}

func (entity *PgEntity) MergeFields(list []int) string {
	sb := strings.Builder{}

	allocLen := 0
	for _, idx := range list {
		allocLen += len(entity.FieldNames[idx]) + 3
	}

	sb.Grow(allocLen)

	sb.WriteRune('"')
	sb.WriteString(strings.ReplaceAll(entity.FieldNames[list[0]], `"`, `""`))
	sb.WriteRune('"')

	for i := 1; i < len(list); i++ {
		idx := list[i]
		sb.WriteRune(',')

		sb.WriteRune('"')
		sb.WriteString(strings.ReplaceAll(entity.FieldNames[idx], `"`, `""`))
		sb.WriteRune('"')
	}

	return sb.String()
}

func (entity *PgEntity) MergeFullFields(list []int) string {
	sb := strings.Builder{}

	allocLen := 0
	for _, idx := range list {
		allocLen += len(entity.FieldFullNames[idx]) + 1
	}

	sb.Grow(allocLen)

	sb.WriteString(entity.FieldFullNames[list[0]])
	for i := 1; i < len(list); i++ {
		idx := list[i]
		sb.WriteRune(',')
		sb.WriteString(entity.FieldFullNames[idx])
	}

	return sb.String()
}

func (entity *PgEntity) AllWrappedFields() []string {
	result := make([]string, 0, len(entity.FieldNames))
	for _, field := range entity.FieldNames {
		result = append(result, `"` + strings.ReplaceAll(field, `"`, `""`) + `"`)
	}

	return result
}

func (entity *PgEntity) MergeAllFields() string {
	list := entity.FieldNames

	sb := strings.Builder{}

	allocLen := 0
	for _, field := range list {
		allocLen += len(field) + 3
	}

	sb.Grow(allocLen)

	sb.WriteRune('"')
	sb.WriteString(strings.ReplaceAll(list[0], `"`, `""`))
	sb.WriteRune('"')

	for i := 1; i < len(list); i++ {
		sb.WriteRune(',')

		sb.WriteRune('"')
		sb.WriteString(strings.ReplaceAll(list[i], `"`, `""`))
		sb.WriteRune('"')
	}

	return sb.String()
}

func (entity *PgEntity) MergeAllFullFields() string {
	return strings.Join(entity.FieldFullNames, ",")
}

func (entity *PgEntity) String() string {
	namesBuilder := strings.Builder{}
	typesBuilder := strings.Builder{}
	indicesBuilder := strings.Builder{}
	uniquesBuilder := strings.Builder{}

	for i, name := range entity.FieldNames {
		pgType := entity.FieldPgTypes[i]
		if i != 0 {
			namesBuilder.WriteString(", ")
			typesBuilder.WriteString(", ")
		}
		namesBuilder.WriteRune('"')
		namesBuilder.WriteString(name)
		namesBuilder.WriteRune('"')

		typesBuilder.WriteString(PgTypeReflectNames[pgType])
	}

	i := 0
	for key, index := range entity.Indices {
		if i != 0 {
			indicesBuilder.WriteString("\n")
			indicesBuilder.WriteString("              > \"")
		} else {
			indicesBuilder.WriteString("     > \"")
		}
		i++

		indicesBuilder.WriteString(key)
		indicesBuilder.WriteString("\"")
		indicesBuilder.WriteString(": {")
		for j, name := range index {
			if j != 0 {
				indicesBuilder.WriteString(", ")
			}
			indicesBuilder.WriteString("\"")
			indicesBuilder.WriteString(name)
			indicesBuilder.WriteString("\"")
		}
		indicesBuilder.WriteString("}")
	}

	i = 0
	for key, unique := range entity.Uniques {
		if i != 0 {
			uniquesBuilder.WriteString("\n")
			uniquesBuilder.WriteString("              > \"")
		} else {
			uniquesBuilder.WriteString("     > \"")
		}
		i++

		uniquesBuilder.WriteString(key)
		uniquesBuilder.WriteString("\"")
		uniquesBuilder.WriteString(": {")
		for j, name := range unique {
			if j != 0 {
				uniquesBuilder.WriteString(", ")
			}
			uniquesBuilder.WriteString("\"")
			uniquesBuilder.WriteString(name)
			uniquesBuilder.WriteString("\"")
		}
		uniquesBuilder.WriteString("}")
	}

	return fmt.Sprintf(`====PgEntity:
Name:         "%s"
FieldNames:   {%s}
FieldPgTypes: {%s}
Indices: %s
Uniques: %s
=========
`,
		entity.Name,
		namesBuilder.String(),
		typesBuilder.String(),
		indicesBuilder.String(),
		uniquesBuilder.String())
}

func makeFullFieldName(tableName, columnName string) string {
	return fmt.Sprintf(`"%v"."%v"`,
		strings.ReplaceAll(tableName, `"`, `""`),
		strings.ReplaceAll(columnName, `"`, `""`))
}
