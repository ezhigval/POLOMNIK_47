package postgres

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode"
)

// Goose splits SQL on semicolons unless a statement is wrapped in
// -- +goose StatementBegin / StatementEnd. Dollar-quoted article text
// with ";" would otherwise become invalid extra statements and fail CI migrate.
func TestGooseSQLMigrationsDoNotSplitInsideDollarQuotes(t *testing.T) {
	t.Parallel()

	dir := filepath.Join("..", "..", "..", "..", "migrations")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read migrations: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		name := entry.Name()
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			body, err := os.ReadFile(filepath.Join(dir, name))
			if err != nil {
				t.Fatalf("read %s: %v", name, err)
			}
			for i, stmt := range gooseUpStatements(string(body)) {
				keyword := leadingSQLKeyword(stmt)
				if !allowedGooseSQLKeyword[keyword] {
					t.Errorf("up statement %d starts with %q; goose likely split a string on ';'. Prefix:\n%s",
						i+1, keyword, truncateForTest(stmt, 180))
				}
			}
		})
	}
}

var allowedGooseSQLKeyword = map[string]bool{
	"alter": true, "begin": true, "comment": true, "commit": true,
	"create": true, "delete": true, "do": true, "drop": true,
	"grant": true, "insert": true, "revoke": true, "select": true,
	"set": true, "truncate": true, "update": true, "with": true,
}

func gooseUpStatements(sql string) []string {
	var (
		stmts   []string
		buf     strings.Builder
		ignore  bool
		started bool
	)
	scanner := bufio.NewScanner(strings.NewReader(sql))
	scanner.Buffer(make([]byte, 0, 1024), 4*1024*1024)
	flush := func() {
		stmt := strings.TrimSpace(buf.String())
		buf.Reset()
		if stmt != "" {
			stmts = append(stmts, stmt)
		}
	}
	for scanner.Scan() {
		line := scanner.Text()
		trim := strings.TrimSpace(line)
		if strings.Contains(trim, "+goose StatementBegin") {
			ignore = true
			continue
		}
		if strings.Contains(trim, "+goose StatementEnd") {
			ignore = false
			flush()
			continue
		}
		if strings.Contains(trim, "+goose Down") {
			break
		}
		if strings.Contains(trim, "+goose Up") {
			started = true
			continue
		}
		if !started {
			continue
		}
		buf.WriteString(line)
		buf.WriteByte('\n')
		if !ignore && gooseLineEndsStatement(line) {
			flush()
		}
	}
	flush()
	return stmts
}

func gooseLineEndsStatement(line string) bool {
	prev := ""
	for _, word := range strings.Fields(line) {
		if strings.HasPrefix(word, "--") {
			break
		}
		prev = word
	}
	return strings.HasSuffix(prev, ";")
}

func leadingSQLKeyword(stmt string) string {
	for _, line := range strings.Split(stmt, "\n") {
		trim := strings.TrimSpace(line)
		if trim == "" || strings.HasPrefix(trim, "--") {
			continue
		}
		var b strings.Builder
		for _, r := range trim {
			if unicode.IsLetter(r) {
				b.WriteRune(unicode.ToLower(r))
				continue
			}
			break
		}
		return b.String()
	}
	return ""
}

func truncateForTest(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
