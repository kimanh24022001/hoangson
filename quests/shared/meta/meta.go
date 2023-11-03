package meta

import (
	"fmt"
	"strings"

	"smatyx.com/shared/database"

	_ "embed"
)

//go:embed templates/func_insert
var funcInsertTemplate string
//go:embed templates/func_batch_insert
var funcBatchInsertTemplate string
//go:embed templates/func_update
var funcUpdateTemplate string
//go:embed templates/func_batch_update
var funcBatchUpdateTemplate string

const DbCallsHeader = `// NOTE(auto): This file is auto-generated. Please don't modify.
package meta

import (
	"smatyx.com/shared/database"
	"smatyx.com/shared/entities"
	"smatyx.com/shared/server"
)


`

const DbEntitiesHeader = `// NOTE(auto): This file is auto-generated. Please don't modify.
package meta

import (
	"smatyx.com/shared/database"
	"smatyx.com/shared/entities"
)


`

func MakeDbCalls(entity *database.PgEntity) string {
	sb := strings.Builder{}

	// NOTE(duong): insert stuff
	{
		qb := database.NewInsertBuilder(
			entity.Name,
			entity.AllWrappedFields(),
			make([]any, len(entity.FieldNames)))

		queryText := qb.String()

		argsTextBuilder := strings.Builder{}
		for i, field := range entity.FieldGoNames {
			if i != 0 {
				argsTextBuilder.WriteString(", ")
			}
			argsTextBuilder.WriteString("entity.")
			argsTextBuilder.WriteString(field)
		}

		// NOTE(duong): insert function
		{
			s := fmt.Sprintf(
				funcInsertTemplate,
				entity.GoName,
				entity.GoName,
				queryText,
				argsTextBuilder.String())
			sb.WriteString(s)
			sb.Write([]byte{'\n'})
		}

		// NOTE(duong): insert batch function
		{
			s := fmt.Sprintf(
				funcBatchInsertTemplate,
				entity.GoName,
				entity.GoName,
				queryText,
				argsTextBuilder.String())
			sb.WriteString(s)
			sb.Write([]byte{'\n'})
		}
	}

	// NOTE(duong): update stuff
	{
		qb := database.NewUpdateBuilder(
			entity.Name,
			entity.FieldNames,
			make([]any, len(entity.FieldNames)))
		qb.WhereVar(entity.Name, "id", database.Op_Eq, nil)

		queryText := qb.String()

		argsTextBuilder := strings.Builder{}
		for i, field := range entity.FieldGoNames {
			if i != 0 {
				argsTextBuilder.WriteString(", ")
			}
			argsTextBuilder.WriteString("entity.")
			argsTextBuilder.WriteString(field)
		}
		argsTextBuilder.WriteString(", ")
		argsTextBuilder.WriteString("entity.")
		argsTextBuilder.WriteString("Id")

		// NOTE(duong): update function
		{
			s := fmt.Sprintf(
				funcUpdateTemplate,
				entity.GoName,
				entity.GoName,
				queryText,
				argsTextBuilder.String())
			sb.WriteString(s)
			sb.Write([]byte{'\n'})
		}

		// NOTE(duong): update batch function
		{
			s := fmt.Sprintf(
				funcBatchUpdateTemplate,
				entity.GoName,
				entity.GoName,
				queryText,
				argsTextBuilder.String())
			sb.WriteString(s)
			sb.Write([]byte{'\n'})
		}
	}

	return sb.String()
}

func MakeDbEntities(entity *database.PgEntity) string {
	sb := strings.Builder{}


	sb.WriteString(
		fmt.Sprintf(`var %v = database.NewPgEntity("%v", entities.%v{})
`,
			entity.GoName,
			entity.Name,
			entity.GoName))

	stringBuilderAddDeclareList(
		&sb, "const",
		len(entity.FieldNames),
		func(i int) string {
			return entity.GoName + "_" + entity.FieldGoNames[i] // + "_Name"
		},
		func(sb *strings.Builder, i int) string {
			return fmt.Sprintf("`\"%v\".\"%v\"`", entity.Name, entity.FieldNames[i])
		},
	)

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

	return sb.String()
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
