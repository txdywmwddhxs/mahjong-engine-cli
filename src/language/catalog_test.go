package language

import (
	"testing"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

func TestCatalogCoverage_AllKeysPresent(t *testing.T) {
	t.Helper()

	if _, ok := catalogs[utils.Chinese]; !ok {
		t.Fatalf("missing catalogs entry for %q", utils.Chinese)
	}
	if _, ok := catalogs[utils.English]; !ok {
		t.Fatalf("missing catalogs entry for %q", utils.English)
	}

	for _, key := range AllKeys {
		zh, ok := zhCN[key]
		if !ok {
			t.Fatalf("missing zhCN translation for key %q", key)
		}
		if zh == "" {
			t.Fatalf("empty zhCN translation for key %q", key)
		}

		en, ok := enUS[key]
		if !ok {
			t.Fatalf("missing enUS translation for key %q", key)
		}
		if en == "" {
			t.Fatalf("empty enUS translation for key %q", key)
		}
	}
}
