package replay

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestActionParsingAndValidation(t *testing.T) {
	cases := []struct{ request, input, kind, tile string }{
		{"discard", " b3 ", "discard", "B3"},
		{"discard", "b3 ting", "declare_waiting", "B3"},
		{"discard", "B3 TT", "declare_waiting", "B3"},
		{"discard", "", "discard_drawn", ""},
		{"discard", "ENDL", "cancel", ""},
		{"discard", "b3 junk TT", "invalid", ""},
		{"discard", "B99", "invalid", ""},
		{"confirm_pung", "yes", "accept", ""},
		{"confirm_kong", "no", "decline", ""},
		{"confirm_kong", "", "invalid", ""},
	}
	for _, tc := range cases {
		action := ParseAction(tc.request, tc.input)
		if action.Kind != tc.kind || string(action.Tile) != tc.tile {
			t.Fatalf("%q: %+v", tc.input, action)
		}
		if err := action.Validate(tc.request); err != nil {
			t.Fatal(err)
		}
	}
	for _, action := range []Action{{Kind: "accept", Tile: "B1"}, {Kind: "discard", Tile: "BOGUS"}, {Kind: "not_supported"}} {
		if err := action.Validate("discard"); err == nil {
			t.Fatalf("accepted malformed action %+v", action)
		}
	}
}

func TestLoadRejectsCorruptAndTruncatedFiles(t *testing.T) {
	for _, data := range []string{`{`, `{"format_version":999}`, `{"unknown_field":1}`, `{} {}`} {
		path := filepath.Join(t.TempDir(), "broken.json")
		if err := os.WriteFile(path, []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := Load(path); err == nil {
			t.Fatalf("loaded %s", data)
		}
	}
}

func TestDifferenceReportsField(t *testing.T) {
	expected := map[string]any{"state": map[string]any{"hand": []string{"B1"}}}
	actual := map[string]any{"state": map[string]any{"hand": []string{"B2"}}}
	if d := Difference(expected, actual); !strings.Contains(d, "state.hand") || !strings.Contains(d, "B1") || !strings.Contains(d, "B2") {
		t.Fatal(d)
	}
}
