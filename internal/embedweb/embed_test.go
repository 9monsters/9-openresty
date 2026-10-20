package embedweb

import (
	"io/fs"
	"strings"
	"testing"
)

func TestAllAssetsEmbedded(t *testing.T) {
	afs := AssetsFS()
	f, err := afs.Open("plugin-vue_export-helper-DlAUqK2U.js")
	if err != nil {
		t.Fatalf("AssetsFS cannot open plugin-vue_export-helper-DlAUqK2U.js: %v", err)
	}
	fi, err := f.Stat()
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}
	t.Logf("Found: %s (%d bytes)", fi.Name(), fi.Size())
	f.Close()

	// Verify no file in dist/assets starts with _ or . (Go embed excludes those).
	dir, err := afs.Open(".")
	if err != nil {
		t.Fatal(err)
	}
	entries, err := dir.(fs.ReadDirFile).ReadDir(-1)
	_ = dir.Close()
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		n := e.Name()
		if strings.HasPrefix(n, "_") || strings.HasPrefix(n, ".") {
			t.Errorf("file %q starts with _ or . → will be excluded by Go embed", n)
		}
	}
	t.Logf("Total assets: %d files", len(entries))
}