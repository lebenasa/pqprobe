package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	_ "github.com/lib/pq"

	"github.com/lebenasa/pqprobe"
	"github.com/pkg/errors"
)

type (
	// CodeGen data
	CodeGen struct {
		PackageName string
		Table       pqprobe.Table
	}
)

func usage() {
	fmt.Println("Generate code to access specific database.")
	fmt.Println("See template/sample.go.template for sample template")
	fmt.Println("Usage: codegen [database connection string] [table name] [template file]")
	fmt.Println("Example: codegen postgres://user:pass@host/database musics example/codegen/templates/sample.go")
	flag.PrintDefaults()
}

func main() {
	// Flags
	flag.Usage = usage
	flagPackageName := flag.String("package-name", "", "package name for generated code, defaults to table name")
	flagOutput := flag.String("output", "", "write generated code to file instead of stdout")
	flagOutputShort := flag.String("o", "", "write generated code to file instead of stdout")
	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		return
	}

	connectionString := flag.Arg(0)
	tableName := flag.Arg(1)
	templateFile := flag.Arg(2)
	packageName := tableName
	if flagPackageName != nil && *flagPackageName != "" {
		packageName = *flagPackageName
	}

	output := os.Stdout
	var err error
	if flagOutput != nil && *flagOutput != "" {
		output, err = os.Create(*flagOutput)
		if err != nil {
			log.Fatalf("Opening file for writing: %v", err)
		}
	}
	if output == os.Stdout && flagOutputShort != nil && *flagOutputShort != "" {
		output, err = os.Create(*flagOutputShort)
		if err != nil {
			log.Fatalf("Opening file for writing: %v", err)
		}
	}

	// Apps logic
	prober, err := pqprobe.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Connecting to database %v: %v", connectionString, err)
	}

	table, err := prober.QueryTable(tableName)
	if err != nil {
		log.Fatalf("Querying table fields: %v", errors.Cause(err))
	}

	code, err := ioutil.ReadFile(templateFile)
	if err != nil {
		log.Fatalf("Read file: %v", err)
	}

	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"notAtEnd": func(i, len int) bool {
			return i < len-1
		},
	}

	codeTemplate := template.Must(template.New("code").Funcs(funcMap).Parse(string(code)))
	err = codeTemplate.Execute(output, CodeGen{PackageName: packageName, Table: table})
	if err != nil {
		log.Fatalf("Executing template: %v", err)
	}

	os.Exit(0)
}
