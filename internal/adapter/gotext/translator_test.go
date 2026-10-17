package gotext

import (
	"bufio"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestReplaceVar(t *testing.T) {
	if got := replaceVar("Hello {name}, {name}!", "name", "Ada"); got != "Hello Ada, Ada!" {
		t.Fatalf("replaceVar() = %q", got)
	}
}

func TestTranslationKeysMatchAcrossLocales(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller() failed")
	}
	root := filepath.Join(filepath.Dir(file), "..", "..", "..", "locale")
	locales := []string{"en", "ru", "rs"}
	keys := make(map[string]map[string]struct{}, len(locales))
	for _, locale := range locales {
		keys[locale] = readPOKeys(t, filepath.Join(root, locale, "LC_MESSAGES", "bot.po"))
	}
	for _, locale := range locales[1:] {
		for key := range keys["en"] {
			if _, ok := keys[locale][key]; !ok {
				t.Errorf("locale %s is missing key %q", locale, key)
			}
		}
		for key := range keys[locale] {
			if _, ok := keys["en"][key]; !ok {
				t.Errorf("locale en is missing key %q", key)
			}
		}
	}
}

func readPOKeys(t *testing.T, path string) map[string]struct{} {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer file.Close()
	keys := map[string]struct{}{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "msgid ") {
			continue
		}
		value, err := strconv.Unquote(strings.TrimPrefix(line, "msgid "))
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}
		if value != "" {
			keys[value] = struct{}{}
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return keys
}
