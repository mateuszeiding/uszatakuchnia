package db

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
)

func TestPgxConfigDisablesStatementCache(t *testing.T) {
	for _, dsn := range []string{
		"postgres://u:p@example-pooler.test:5432/db?sslmode=disable",
		"host=127.0.0.1 port=5433 user=u password=p dbname=db sslmode=disable",
		// an explicit mode in the DSN does not bring the cache back
		"postgres://u:p@127.0.0.1:5433/db?sslmode=disable&default_query_exec_mode=cache_statement",
	} {
		c, err := pgxConfig(dsn)
		if err != nil {
			t.Fatalf("pgxConfig(%q): %v", dsn, err)
		}
		if c.DefaultQueryExecMode != pgx.QueryExecModeExec {
			t.Errorf("pgxConfig(%q): mode %v, want %v", dsn, c.DefaultQueryExecMode, pgx.QueryExecModeExec)
		}
	}
}

func TestPgxConfigErrorHidesPassword(t *testing.T) {
	_, err := pgxConfig("postgres://u:s3cret-pass@host:notaport/db")
	if err == nil {
		t.Fatal("want an error for a bad port")
	}
	if strings.Contains(err.Error(), "s3cret-pass") {
		t.Errorf("error leaks the password: %v", err)
	}
}

// Every connection goes through db.Open; a direct gorm.Open(postgres.Open(...)) would
// bring back pgx's prepared statement cache.
func TestNoDirectPostgresOpen(t *testing.T) {
	err := filepath.WalkDir("..", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && (d.Name() == "node_modules" || d.Name() == "web" || strings.HasPrefix(d.Name(), ".") && p != "..") {
			return filepath.SkipDir
		}
		if d.IsDir() || !strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go") {
			return nil
		}
		f, err := parser.ParseFile(token.NewFileSet(), p, nil, 0)
		if err != nil {
			return err
		}
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			pkg, ok := sel.X.(*ast.Ident)
			if ok && pkg.Name == "postgres" && (sel.Sel.Name == "Open" || sel.Sel.Name == "New") &&
				filepath.Clean(p) != filepath.Join("..", "db", "open.go") {
				t.Errorf("%s: postgres.%s — open the database with db.Open", p, sel.Sel.Name)
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
