package migrate

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"unicode"

	. "smatyx.com/shared/database"
)

type PgEntityCompareResult struct {
	NewEntityName string
	OldEntityName string

	AddFieldNames []string
	AddFieldTypes []int

	DropFieldNames []string

	UpdateFieldTypes map[string]int

	DropIndices []string
	DropUniques []string

	AddIndices map[string][]string
	AddUniques map[string][]string
}

type MigrateInstruction struct {
	UpdateSchema bool

	OldEntity *PgEntity
	NewEntity *PgEntity
}

var migrateInstructions []MigrateInstruction
func DoMigrate() string {
	for _, instruction := range migrateInstructions {
		if !instruction.UpdateSchema {
			query := NewWriteQuery(
				context.Background(), nil,
				MakeCreateEntitySQL(instruction.NewEntity))
			err := query.Submit()
			if err != nil {
				panic(err)
			}
		} else {
			compareResult := PgEntityCompare(
				instruction.NewEntity,
				instruction.OldEntity)
			query := NewWriteQuery(
				context.Background(), nil,
				compareResult.MakeUpdateEntitySQL())
			err := query.Submit()
			if err != nil {
				fmt.Printf("instruction.OldEntity: %v\n", instruction.OldEntity)
				panic(err)
			}
		}
	}

	sb := strings.Builder{}
	sb.WriteString(schemaFileHeader)
	stringBuilderAddDeclareList(
		&sb, "var", len(migrateInstructions),
		func(i int) string {
			return migrateInstructions[i].NewEntity.GoName
		},
		func(sb *strings.Builder, i int) string {
			entity := migrateInstructions[i].NewEntity
			return fmt.Sprintf(`database.NewPgEntity("%v", entities.%v{})`, entity.Name, entity.GoName)
		},
	)

	for _, instruction := range migrateInstructions {
		entity := instruction.NewEntity
		stringBuilderAddDeclareList(
			&sb, "const",
			len(entity.FieldNames),
			func(i int) string {
				return entity.GoName + "_" + entity.FieldGoNames[i] + "_Name"
			},
			func(sb *strings.Builder, i int) string {
				return fmt.Sprintf("`\"%v\".\"%v\"`", entity.Name, entity.FieldNames[i])
			},
		)
	}

	for _, instruction := range migrateInstructions {
		entity := instruction.NewEntity
		stringBuilderAddDeclareList(
			&sb, "const",
			len(entity.FieldNames),
			func(i int) string {
				return entity.GoName + "_" + entity.FieldGoNames[i] + "_Idx"
			},
			func(sb *strings.Builder, i int) string {
				return fmt.Sprintf("%d", i)
			},
		)
	}
	
	return sb.String()
}

func AddMigrate(oldEntityName string, entity *PgEntity) {
	oldEntity, _ := NewPgEntityFromDatabase(oldEntityName)
	// NOTE: assume that we don't have the oldEntity in
	// database
	if oldEntity == nil {
		migrateInstructions = append(migrateInstructions, MigrateInstruction{
			UpdateSchema: false,
			NewEntity:    entity,
		})
	} else {
		migrateInstructions = append(
			migrateInstructions,
			MigrateInstruction{
				UpdateSchema: true,
				OldEntity:    oldEntity,
				NewEntity:    entity,
			})
	}
}

func NewPgEntityFromDatabase(entityName string) (*PgEntity, error) {
	result := &PgEntity{
		Name:         entityName,
		FieldNames:   []string{},
		FieldPgTypes: []int{},
		Indices:      make(map[string][]string),
		Uniques:      make(map[string][]string),
	}

	processInformationSchema := func (rep *PgRep) error {
		var nm, tp string
		err := rep.Scan(&nm, &tp)
		if err != nil {
			return err
		}
		result.FieldNames = append(result.FieldNames, nm)
		pgType, ok := PgTypeReturns[strings.ToUpper(tp)]
		if !ok {
			return fmt.Errorf("Not supported pgType: '%v'\n", tp)
		}
		result.FieldPgTypes = append(result.FieldPgTypes, pgType)		
		return nil
	}

	query := NewReadQuery(
		context.Background(),
		processInformationSchema,
		`SELECT column_name, data_type FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = $1 ORDER BY ordinal_position ASC;`,
		entityName)

	err := query.Submit()
	if err != nil {
		return nil, err
	}

	if len(result.FieldNames) == 0 {
		return nil, errors.New("WTF")
	}

	query = NewReadQuery(
		context.Background(),
		func (rows *PgRep) error {
			var idx, def string

			err := rows.Scan(&idx, &def)
			if err != nil {
				return err
			}

			isUnique, fields, err := ParseCreateIndexFields(def)
			if err != nil {
				return err
			}

			if isUnique {
				// NOTE(duong): Ignore primary key
				if !(len(fields) == 1 && fields[0] == "id") {
					result.Uniques[idx] = fields
				}
			} else {
				result.Indices[idx] = fields
			}

			return nil
		},
		`SELECT indexname, indexdef FROM pg_indexes WHERE tablename = $1;`,
		entityName)
	err = query.Submit()

	if err != nil {
		log.Printf("err: %v\n", err)
		return nil, err
	}

	return result, nil
}

func (compareResult *PgEntityCompareResult) MakeUpdateEntitySQL() string {
	indexModifyBuilder := strings.Builder{}
	for _, unique := range compareResult.DropUniques {
		indexModifyBuilder.WriteString("\nDROP INDEX ")
		indexModifyBuilder.WriteString(WrapPostgresField(unique))
		indexModifyBuilder.WriteString(";")
	}

	for _, index := range compareResult.DropIndices {
		indexModifyBuilder.WriteString("\nDROP INDEX ")
		indexModifyBuilder.WriteString(WrapPostgresField(index))
		indexModifyBuilder.WriteString(";")
	}

	for unique, fields := range compareResult.AddUniques {
		indexModifyBuilder.WriteString("\nCREATE UNIQUE INDEX ")
		indexModifyBuilder.WriteString(WrapPostgresField(unique))
		indexModifyBuilder.WriteString(" ON ")
		indexModifyBuilder.WriteString(WrapPostgresField(compareResult.NewEntityName))
		indexModifyBuilder.WriteString(" ")

		indexModifyBuilder.WriteRune('(')
		for i, field := range fields {
			if i != 0 {
				indexModifyBuilder.WriteString(", ")
			}
			indexModifyBuilder.WriteString(WrapPostgresField(field))
		}
		indexModifyBuilder.WriteString(");")
	}

	for index, fields := range compareResult.AddIndices {
		indexModifyBuilder.WriteString("\nCREATE INDEX ")
		indexModifyBuilder.WriteString(WrapPostgresField(index))
		indexModifyBuilder.WriteString(" ON ")
		indexModifyBuilder.WriteString(WrapPostgresField(compareResult.NewEntityName))
		indexModifyBuilder.WriteString(" ")

		indexModifyBuilder.WriteRune('(')
		for i, field := range fields {
			if i != 0 {
				indexModifyBuilder.WriteString(", ")
			}
			indexModifyBuilder.WriteString(WrapPostgresField(field))
		}
		indexModifyBuilder.WriteString(");")
	}

	alterEntityBuilder := strings.Builder{}
	alterEntityBuilder.WriteString("ALTER TABLE ")
	alterEntityBuilder.WriteString(WrapPostgresField(compareResult.OldEntityName))

	needComma := false
	if compareResult.NewEntityName != compareResult.OldEntityName {
		needComma = true
		alterEntityBuilder.WriteString("\nRENAME TO ")
		alterEntityBuilder.WriteString(compareResult.NewEntityName)
	}

	for i, fieldName := range compareResult.AddFieldNames {
		fieldType := compareResult.AddFieldTypes[i]
		if needComma {
			alterEntityBuilder.WriteString(",\n")
		}
		alterEntityBuilder.WriteString("\nADD COLUMN ")
		alterEntityBuilder.WriteString(WrapPostgresField(fieldName))
		alterEntityBuilder.WriteString(" ")
		alterEntityBuilder.WriteString(PgTypeNames[fieldType])

		needComma = true
	}

	for _, fieldName := range compareResult.DropFieldNames {
		if needComma {
			alterEntityBuilder.WriteString(",\n")
		}
		alterEntityBuilder.WriteString("\nDROP COLUMN ")
		alterEntityBuilder.WriteString(WrapPostgresField(fieldName))

		needComma = true
	}

	for fieldName, fieldType := range compareResult.UpdateFieldTypes {
		if needComma {
			alterEntityBuilder.WriteString(",\n")
		}
		alterEntityBuilder.WriteString("\nALTER COLUMN ")
		alterEntityBuilder.WriteString(WrapPostgresField(fieldName))
		alterEntityBuilder.WriteString("\nTYPE ")
		alterEntityBuilder.WriteString(PgTypeNames[fieldType])

		needComma = true
	}
	alterEntityBuilder.WriteString(";")

	if needComma {
		return indexModifyBuilder.String() + "\n" + alterEntityBuilder.String()
	}

	return indexModifyBuilder.String()
}

func PgEntityCompare(newEntity, oldEntity *PgEntity) PgEntityCompareResult {
	result := PgEntityCompareResult{
		NewEntityName:    newEntity.Name,
		OldEntityName:    oldEntity.Name,
		AddFieldNames:    []string{},
		AddFieldTypes:    []int{},
		DropFieldNames:   []string{},
		UpdateFieldTypes: map[string]int{},
		DropIndices:      []string{},
		DropUniques:      []string{},
		AddIndices:       map[string][]string{},
		AddUniques:       map[string][]string{},
	}

	for i, newFieldName := range newEntity.FieldNames {
		check := false
		newFieldType := newEntity.FieldPgTypes[i]

		for j, oldFieldName := range oldEntity.FieldNames {
			oldFieldType := oldEntity.FieldPgTypes[j]

			if newFieldName == oldFieldName {
				if newFieldType != oldFieldType {
					result.UpdateFieldTypes[newFieldName] = newFieldType
				}
				check = true
				break
			}
		}

		if !check {
			result.AddFieldNames = append(result.AddFieldNames, newFieldName)
			result.AddFieldTypes = append(result.AddFieldTypes, newFieldType)
		}
	}

	for _, oldFieldName := range oldEntity.FieldNames {
		check := false

		for _, newFieldName := range newEntity.FieldNames {
			if newFieldName == oldFieldName {
				check = true
				break
			}
		}

		if !check {
			result.DropFieldNames = append(result.AddFieldNames, oldFieldName)
		}
	}

	for unique, newFields := range newEntity.Uniques {
		oldFields, ok := oldEntity.Uniques[unique]
		if !ok {
			result.AddUniques[unique] = newFields
			continue
		}

		if !reflect.DeepEqual(oldFields, newFields) {
			result.AddUniques[unique] = newFields
			result.DropUniques = append(result.DropUniques, unique)
		}
	}

	for unique := range oldEntity.Uniques {
		_, ok := newEntity.Uniques[unique]
		if !ok {
			result.DropUniques = append(result.DropUniques, unique)
		}
	}

	for index, newFields := range newEntity.Indices {
		oldFields, ok := oldEntity.Indices[index]
		if !ok {
			result.AddIndices[index] = newFields
			continue
		}

		if !reflect.DeepEqual(oldFields, newFields) {
			result.AddIndices[index] = newFields
			result.DropIndices = append(result.DropIndices, index)
		}
	}

	for index := range oldEntity.Indices {
		_, ok := newEntity.Indices[index]
		if !ok {
			result.DropIndices = append(result.DropIndices, index)
		}
	}

	return result
}

func MakeCreateEntitySQL(entity *PgEntity) string {
	fieldsBuilder := strings.Builder{}
	for i := 0; i < len(entity.FieldNames); i++ {
		if i != 0 {
			fieldsBuilder.WriteString(",\n")
		} else {
			fieldsBuilder.WriteString("\n")
		}

		fieldsBuilder.WriteRune('"')
		fieldsBuilder.WriteString(entity.FieldNames[i])
		fieldsBuilder.WriteRune('"')

		fieldsBuilder.WriteRune(' ')
		fieldsBuilder.WriteString(PgTypeNames[entity.FieldPgTypes[i]])

		if entity.FieldNames[i] == "id" {
			fieldsBuilder.WriteString(" PRIMARY KEY")
		}
	}

	uniqueBuilder := strings.Builder{}
	for unique, fields := range entity.Uniques {
		uniqueBuilder.WriteString(
			fmt.Sprintf(`CREATE UNIQUE INDEX IF NOT EXISTS %s ON %s `,
				WrapPostgresField(unique),
				WrapPostgresField(entity.Name)))

		uniqueBuilder.WriteRune('(')
		for i, field := range fields {
			if i != 0 {
				uniqueBuilder.WriteString(", ")
			}
			uniqueBuilder.WriteString(WrapPostgresField(field))
		}
		uniqueBuilder.WriteString(");\n")
	}

	indexBuilder := strings.Builder{}
	for index, fields := range entity.Indices {
		indexBuilder.WriteString(
			fmt.Sprintf(`CREATE INDEX IF NOT EXISTS %s ON %s `,
				WrapPostgresField(index),
				WrapPostgresField(entity.Name)))

		indexBuilder.WriteRune('(')
		for i, field := range fields {
			if i != 0 {
				indexBuilder.WriteString(", ")
			}
			indexBuilder.WriteString(WrapPostgresField(field))
		}
		indexBuilder.WriteString(");\n")
	}

	return fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS "%s" (%s);
%s
%s`,
		entity.Name,
		fieldsBuilder.String(),
		uniqueBuilder.String(),
		indexBuilder.String())
}

// Pattern:
// CREATE (might have UNIQUE here) INDEX {index_name} ON {entity_name} USING {algo} ({fields})
func ParseCreateIndexFields(query string) (bool, []string, error) {
	result := []string{}

	queryRunes := []rune(query)

	fieldsCloseIndex := len(queryRunes) - 1
	if queryRunes[fieldsCloseIndex] != ')' {
		return false, nil, errors.New("Invalid CREATE INDEX")
	}

	inText := false
	inQuote := false
	happyClose := false
	endFieldIndex := fieldsCloseIndex - 1

	for i := fieldsCloseIndex - 1; i >= 0; i-- {
		r := queryRunes[i]

		if !inQuote {
			if r == '"' {
				endFieldIndex = i
				inQuote = true
			}

			if !inText {
				if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
					endFieldIndex = i + 1
					inText = true
				}
			} else {
				if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_') {
					field := UnwrapPostgresRunesField(queryRunes[i+1 : endFieldIndex])
					result = append(result, field)
					inText = false
				}
			}

			if r == '(' {
				happyClose = true
				break
			}
		} else {
			if r == '"' {
				if queryRunes[i-1] == '"' {
					i--
				} else {
					field := UnwrapPostgresRunesField(queryRunes[i : endFieldIndex+2])
					result = append(result, field)
					inText = false
					inQuote = false
				}
			}
		}
	}

	if !happyClose {
		return false, nil, errors.New(fmt.Sprintf("Can not find the corrensponding '(' of \n=> %s", query))
	}

	// NOTE: reverse the result
	{
		i := 0
		j := len(result) - 1
		for i < j {
			result[i], result[j] = result[j], result[i]
			i++
			j--
		}
	}

	return strings.HasPrefix(query, "CREATE UNIQUE"), result, nil
}

func UnwrapPostgresRunesField(field []rune) string {
	if field[0] == '"' {
		field = field[1 : len(field)-2]
	}

	sb := strings.Builder{}
	skip := false
	for _, c := range field {
		if skip {
			skip = false
			continue
		}
		if c == '"' {
			skip = true
		}
		sb.WriteRune(rune(c))
	}
	return sb.String()
}

func WrapPostgresField(field string) string {
	return "\"" + strings.ReplaceAll(field, `"`, `""`) + "\""
}

func stringBuilderAddDeclareList(
	stringBuilder *strings.Builder,
	declareKeyword string,
	listLen int, 
	nameFn func(i int) string,
	valueFn func(sb *strings.Builder, i int) string) {

	stringBuilder.WriteString(declareKeyword)
	stringBuilder.WriteString(" (\n")

	maxNameLen := 0
	for i := 0; i < listLen; i++ {
		name := nameFn(i)
		maxNameLen = max(len(name), maxNameLen)
	}

	for i := 0; i < listLen; i++ {
		name := nameFn(i)
		value := valueFn(stringBuilder, i)
		line := fmt.Sprintf("	%v %*v %v\n", name, maxNameLen - len(name) + 1, "=", value)
		stringBuilder.WriteString(line)
	}

	stringBuilder.WriteString(")\n\n")
}
