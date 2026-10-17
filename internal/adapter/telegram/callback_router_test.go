package telegram

import "testing"

func TestCleanCallbackData(t *testing.T) {
	if got := cleanCallbackData("\fchoose_lang:en"); got != "choose_lang:en" {
		t.Fatalf("cleanCallbackData() = %q", got)
	}
	if got := cleanCallbackData("choose_lang:en"); got != "choose_lang:en" {
		t.Fatalf("cleanCallbackData() changed clean data to %q", got)
	}
}
