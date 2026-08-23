package postformat

import (
	"testing"

	"palomnik/internal/ports"
)

func TestPostSkipsEmptyParts(t *testing.T) {
	got := Post(ports.PublishContent{Title: "  Заголовок  ", URL: "https://tikhvin-palomnik.ru/news/a"})
	want := "Заголовок\n\nhttps://tikhvin-palomnik.ru/news/a"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestEscapeHTML(t *testing.T) {
	got := EscapeHTML(`A & B <c>`)
	if got != "A &amp; B &lt;c&gt;" {
		t.Fatalf("got %q", got)
	}
}
