package http

import (
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Każda funkcja Vercela w api/ musi zacząć od `if resp.Begin(w, r) { return }` —
// inaczej endpoint zapisujący byłby otwarty dla każdego.
func TestEveryHandlerStartsWithBegin(t *testing.T) {
	var files []string
	err := filepath.WalkDir("../api", func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(p, ".go") && !strings.HasSuffix(p, "_test.go") {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no handlers found in api/")
	}

	handlers := 0
	for _, p := range files {
		f, err := parser.ParseFile(token.NewFileSet(), p, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || !fn.Name.IsExported() || !isHandlerSig(fn.Type) {
				continue
			}
			handlers++
			if !startsWithBegin(fn) {
				t.Errorf("%s: %s must start with `if resp.Begin(w, r) { return }`", p, fn.Name.Name)
			}
		}
	}
	if handlers < len(files) {
		t.Errorf("found %d handlers in %d files — each api file should export a handler", handlers, len(files))
	}
}

func isHandlerSig(ft *ast.FuncType) bool {
	if ft.Params == nil || len(ft.Params.List) != 2 {
		return false
	}
	sel, ok := ft.Params.List[0].Type.(*ast.SelectorExpr)
	return ok && sel.Sel.Name == "ResponseWriter"
}

func startsWithBegin(fn *ast.FuncDecl) bool {
	if fn.Body == nil || len(fn.Body.List) == 0 {
		return false
	}
	ifs, ok := fn.Body.List[0].(*ast.IfStmt)
	if !ok || ifs.Init != nil || len(ifs.Body.List) != 1 {
		return false
	}
	if _, ok := ifs.Body.List[0].(*ast.ReturnStmt); !ok {
		return false
	}
	call, ok := ifs.Cond.(*ast.CallExpr)
	if !ok {
		return false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Begin" {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	return ok && pkg.Name == "resp"
}

const key = "0123456789abcdef0123456789abcdef-begin"

func run(method, origin, authz string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, "/api/recipes", nil)
	if origin != "" {
		r.Header.Set("Origin", origin)
	}
	if authz != "" {
		r.Header.Set("Authorization", authz)
	}
	w := httptest.NewRecorder()
	if !Begin(w, r) {
		w.WriteHeader(http.StatusOK) // handler poszedłby dalej
	}
	return w
}

func TestBeginGuardsWrites(t *testing.T) {
	t.Setenv("API_ADMIN_KEY", key)
	for _, m := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
		if c := run(m, "", "").Code; c != http.StatusUnauthorized {
			t.Errorf("%s without token: %d, want 401", m, c)
		}
		if c := run(m, "", "Bearer wrong").Code; c != http.StatusUnauthorized {
			t.Errorf("%s with wrong key: %d, want 401", m, c)
		}
		if c := run(m, "", "Bearer "+key).Code; c != http.StatusOK {
			t.Errorf("%s with key: %d, want pass", m, c)
		}
	}
	if c := run(http.MethodGet, "", "").Code; c != http.StatusOK {
		t.Errorf("GET without token: %d, want pass", c)
	}
}

func TestBeginWithoutKeyConfiguredFailsClosed(t *testing.T) {
	t.Setenv("API_ADMIN_KEY", "")
	if c := run(http.MethodPost, "", "Bearer ").Code; c != http.StatusUnauthorized {
		t.Fatalf("got %d, want 401", c)
	}
}

func TestCORS(t *testing.T) {
	t.Setenv("CORS_ORIGINS", "https://extra.example")
	cases := map[string]bool{
		"http://localhost:5173": true,
		"https://extra.example": true,
		"https://evil.example":  false,
		"":                      false,
	}
	for origin, allowed := range cases {
		w := run(http.MethodOptions, origin, "")
		if w.Code != http.StatusNoContent {
			t.Errorf("preflight %q: %d", origin, w.Code)
		}
		got := w.Header().Get("Access-Control-Allow-Origin")
		if allowed && got != origin {
			t.Errorf("%q should be allowed, ACAO=%q", origin, got)
		}
		if !allowed && got != "" {
			t.Errorf("%q should be blocked, ACAO=%q", origin, got)
		}
	}
}
