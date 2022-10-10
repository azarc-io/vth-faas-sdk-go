package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	filepath.Walk(".", func(filePath string, info fs.FileInfo, err error) error {
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}
		fset := token.NewFileSet()
		//fullPath := path.Join(filePath, info.Name())
		node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
		if err != nil {
			fmt.Printf("Error parsing file: %s; error: %s", filePath, err.Error())
			os.Exit(255)
		}

		// traverse all tokens
		ast.Inspect(node, func(n ast.Node) bool {
			switch t := n.(type) {
			// find variable declarations
			case *ast.TypeSpec:
				// which are public
				if t.Name.IsExported() {
					switch t.Type.(type) {
					// and are interfaces
					case *ast.InterfaceType:
						if t.Name.Name == "Job" {
							fmt.Printf("### JOB ###: %#v \n", t)
							pos := fset.Position(t.Pos())
							fmt.Printf("%s:%d:%d\n", pos.Filename, pos.Line, pos.Column)

						}
					}
				}
			}
			return true
		})
		return nil
	})
}

// gopls implementation pkg/api/job.go:12:2
