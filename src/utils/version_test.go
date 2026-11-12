package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveVersion(t *testing.T) {
	root := t.TempDir()
	oldPath := VersionFilePath
	VersionFilePath = filepath.Join(root, "VERSION")
	t.Cleanup(func() { VersionFilePath = oldPath })
	// A newer-looking changelog must never provide the application version.
	if err := os.Mkdir(filepath.Join(root, "config"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "config", "CHANGELOG.txt"), []byte("v999.0.0\n"), 0644); err != nil {
		t.Fatal(err)
	}
	oldLogPath := ChangeLogPath
	ChangeLogPath = filepath.Join(root, "config", "CHANGELOG.txt")
	t.Cleanup(func() { ChangeLogPath = oldLogPath })
	if got := resolveVersion(""); got != "dev" {
		t.Fatalf("missing VERSION: got %q", got)
	}
	for _, tc := range []struct{ name, content, injected, want string }{
		{"file version", "v9.1.0\n", "", "v9.1.0"},
		{"CRLF whitespace", " \tv9.1.0\r\n", "", "v9.1.0"},
		{"empty file", "", "", "dev"},
		{"blank file", " \r\n\t", "", "dev"},
		{"injected release", "v0.0.0\n", " v9.1.0 ", "v9.1.0"},
		{"injected test build", "v9.1.0\n", "dev", "dev"},
		{"blank injection", "v9.1.0\n", " \t", "v9.1.0"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := os.WriteFile(VersionFilePath, []byte(tc.content), 0644); err != nil {
				t.Fatal(err)
			}
			if got := resolveVersion(tc.injected); got != tc.want {
				t.Fatalf("got=%q want=%q", got, tc.want)
			}
		})
	}
	VersionFilePath = filepath.Join(root, "missing")
	if got := resolveVersion("v9.1.0"); got != "v9.1.0" {
		t.Fatalf("injection without metadata: got %q", got)
	}
}

func TestFindProjectRootVersionMarker(t *testing.T) {
	for _, marker := range []string{"VERSION", "go.mod"} {
		t.Run(marker, func(t *testing.T) {
			root := t.TempDir()
			nested := filepath.Join(root, "nested", "directory")
			if err := os.MkdirAll(nested, 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(root, marker), []byte("v9.1.0\n"), 0644); err != nil {
				t.Fatal(err)
			}
			t.Setenv("MAHJONG_ROOT", "")
			t.Chdir(nested)
			if got := findProjectRoot(); got != root {
				t.Fatalf("root=%q want=%q", got, root)
			}
			t.Setenv("MAHJONG_ROOT", nested)
			if got := findProjectRoot(); got != nested {
				t.Fatalf("override=%q want=%q", got, nested)
			}
		})
	}
}
