package gotext

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"tgbot/config"
)

func TestTranslatorLoadsLocalesAndFallsBack(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(file), "..", "..", "..")
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(old); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	}()
	tr := NewTranslator(&config.Config{Languages: "en"}, []string{"en", "ru"})
	if got := tr.T("language_btn", "en"); got == "language_btn" {
		t.Fatal("expected loaded translation")
	}
	if got := tr.T("language_btn", "missing"); got == "language_btn" {
		t.Fatal("expected default language fallback")
	}
	if got := tr.T("missing_key", "en"); got != "missing_key" {
		t.Errorf("missing key = %q", got)
	}
	if got := tr.T("missing_key", "none", map[string]any{"x": "y"}); got != "missing_key" {
		t.Errorf("missing locale = %q", got)
	}
	if got := tr.T("language_btn", "ru", map[string]any{"x": "y"}); got == "" {
		t.Fatal("translation unexpectedly empty")
	}
}

func TestTranslatorWithoutLocaleReturnsKey(t *testing.T) {
	tr := NewTranslator(&config.Config{Languages: "en"}, nil)
	if got := tr.T("key", "en"); got != "key" {
		t.Errorf("T() = %q", got)
	}
}
