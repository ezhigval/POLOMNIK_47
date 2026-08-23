package tests

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestDomainDoesNotImportOuterLayers(t *testing.T) {
	forbiddenImports := []string{
		"palomnik/internal/application",
		"palomnik/internal/ports",
		"palomnik/internal/adapters",
		"palomnik/internal/config",
		"palomnik/internal/logger",
		"palomnik/internal/validation",
		"github.com/gofiber/fiber",
		"database/sql",
		"github.com/redis",
		"github.com/go-redis",
		"bitrix",
		"onec",
	}

	files, err := filepath.Glob("../internal/domain/*.go")
	if err != nil {
		t.Fatalf("glob domain files: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("expected domain files")
	}

	fset := token.NewFileSet()
	for _, file := range files {
		parsed, err := parser.ParseFile(fset, file, nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse %s: %v", file, err)
		}

		ast.Inspect(parsed, func(node ast.Node) bool {
			importSpec, ok := node.(*ast.ImportSpec)
			if !ok {
				return true
			}

			path, err := strconv.Unquote(importSpec.Path.Value)
			if err != nil {
				t.Fatalf("unquote import in %s: %v", file, err)
			}

			for _, forbidden := range forbiddenImports {
				if strings.HasPrefix(path, forbidden) || strings.Contains(path, forbidden) {
					t.Fatalf("domain file %s imports forbidden outer dependency %q", file, path)
				}
			}
			return true
		})
	}
}
