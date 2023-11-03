package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"

	"smatyx.com/shared/cast"
)

type Entity struct {
	StructName string
	OldName    string
	NewName    string
}

var Entities = make([]Entity, 0, 5_000)
var CallsInterfaces = make([]string, 0, 5_000)

func ValidName(name string) bool {
	if len(name) < 2 {
		return false
	}

	for _, r := range name {
		if !(unicode.IsDigit(r) || unicode.IsLetter(r) || r == '_') {
			return false
		}
	}
	return true
}

func ProcessStruct(pkgName string, decl *ast.GenDecl) {
	entityName := ""
	for _, spec := range decl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}
		entityName = typeSpec.Name.Name
	}

	if decl.Doc != nil {
		for _, comment := range decl.Doc.List {
			commentContent := strings.TrimSpace(comment.Text[2:])
			if !strings.HasPrefix(commentContent, "entity:table") {
				continue
			}

			nameContent := strings.TrimSpace(commentContent[len("entity:table"):])
			if len(nameContent) == 0 {
				name := cast.StringPlural(cast.StringLowerSnakeCase(entityName))
				Entities = append(Entities, Entity{
					StructName: entityName,
					OldName: name,
					NewName: name,
				})
			} else if strings.HasPrefix(nameContent, "(") && strings.HasSuffix(nameContent, ")") {
				nameContent = nameContent[1 : len(nameContent)-1]
				if strings.Contains(nameContent, "->") {
					parts := strings.Split(nameContent, "->")
					firstName := strings.TrimSpace(parts[0])
					secondName := strings.TrimSpace(parts[1])

					if len(secondName) == 0 {
						secondName = cast.StringLowerSnakeCase(entityName)
					}

					if !ValidName(firstName) || !ValidName(secondName) {
						panic(
							fmt.Sprintf("Error while reading %s.%s: Invalid table declaration format. " +
								"Should be \"{firstName} -> {secondName}\"",
							pkgName, entityName))
					}
					Entities = append(Entities, Entity{
						StructName: entityName,
						OldName: firstName,
						NewName: firstName,
					})
				} else {
					if !ValidName(nameContent) {
						panic(
							fmt.Sprintf("Error while reading %s.%s: Invalid table declaration format. " +
								"The string should only contain digits, letters, and '_'.",
								pkgName, entityName))
					}
					Entities = append(Entities, Entity{
						StructName: entityName,
						OldName: nameContent,
						NewName: nameContent,
					})
				}
			} else {
				panic(fmt.Sprintf("Error while reading %s.%s. %s", pkgName, entityName, nameContent))
			}
		}
	}
}

func ParseEntities(directory string) {
	fset := token.NewFileSet()
	packages, err := parser.ParseDir(fset, directory, nil, parser.ParseComments)

	if err != nil {
		fmt.Println(err)
		return
	}

	for pkgName, pkg := range packages {
		for _, file := range pkg.Files {
			for _, node := range file.Decls {
				decl, ok := node.(*ast.GenDecl)
				if !ok {
					continue
				}
				ProcessStruct(pkgName, decl)
			}
		}
	}
}

func ProcessInterface(pkgName string, decl *ast.GenDecl) {
	interfaceName := ""
	for _, spec := range decl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}
		interfaceName = typeSpec.Name.Name
	}

	if decl.Doc != nil {
		for _, comment := range decl.Doc.List {
			commentContent := strings.TrimSpace(comment.Text[2:])
			if commentContent !=  "entity:calls" {
				continue
			}
			CallsInterfaces = append(CallsInterfaces, interfaceName)
		}
	}
}

var basePath = func () string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(b)
}()
func realPath(pathName string) string{
	if len(pathName) == 0 && pathName[0] == '/' {
		return pathName
	}

	return basePath + "/" + pathName
}

func main() {
	ParseEntities(realPath("../../shared/entities/"))
	WriteMigrateCode(realPath("../do_migrate/main.go"))
}

func WriteMigrateCode(fileName string) {
	sb := strings.Builder{}
	sb.WriteString(
		`// NOTE(auto): This file is auto-generated. Please don't modify.
package main

import (
	"os"
	"path/filepath"
	"runtime"

	"smatyx.com/shared/meta"
	"smatyx.com/shared/database"
	"smatyx.com/shared/entities"
	"smatyx.com/shared/migrate"
)

var basePath = func () string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(b)
}()
func realPath(pathName string) string{
	if len(pathName) == 0 && pathName[0] == '/' {
		return pathName
	}

	return basePath + "/" + pathName
}

func main() {
	callsFile := realPath("../../shared/meta/db_calls.go")
	entitiesFile := realPath("../../shared/meta/db_entities.go")

	os.WriteFile(callsFile, []byte{}, 0666)
	os.WriteFile(entitiesFile, []byte{}, 0666)

	callsFd, err := os.OpenFile(callsFile, os.O_WRONLY | os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer callsFd.Close()

	entitiesFd, err := os.OpenFile(entitiesFile, os.O_WRONLY | os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer entitiesFd.Close()

	callsFd.WriteString(meta.DbCallsHeader)
	entitiesFd.WriteString(meta.DbEntitiesHeader)
`)

	for _, entity := range Entities {
		content := fmt.Sprintf(
			`	// NOTE(auto): Migrate entities.%v
	{
		pgEntity := database.NewPgEntity("%v", entities.%v{})
		migrate.AddMigrate("%v", pgEntity)
		callsFd.WriteString(meta.MakeDbCalls(pgEntity))
		entitiesFd.WriteString(meta.MakeDbEntities(pgEntity))
	}

`, entity.StructName, entity.NewName, entity.StructName, entity.OldName)

		sb.WriteString(content)
	}

	sb.WriteString(`	migrate.DoMigrate()
`)
	sb.WriteRune('}')

	fd, err := os.OpenFile(fileName, os.O_TRUNC | os.O_CREATE | os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer fd.Close()
	fd.WriteString(sb.String())
	

	fmt.Print(`File "cmd/do_migrate/main.go" successfully created! 🚀
To initiate the database migration and forge entity-related functions, execute the following:

	$ go run cmd/do_migrate/main.go

`)
}
