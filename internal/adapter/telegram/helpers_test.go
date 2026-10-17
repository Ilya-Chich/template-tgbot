package telegram

import (
	"errors"
	"fmt"
	"testing"

	"tgbot/internal/core/domain"
	"tgbot/internal/core/port"

	"gopkg.in/telebot.v3"
)

type testContext struct {
	telebot.Context
	sender    *telebot.User
	message   *telebot.Message
	callback  *telebot.Callback
	text      string
	sent      []any
	sendOpts  [][]any
	responses []*telebot.CallbackResponse
	sendErr   error
}

func (c *testContext) Sender() *telebot.User       { return c.sender }
func (c *testContext) Message() *telebot.Message   { return c.message }
func (c *testContext) Callback() *telebot.Callback { return c.callback }
func (c *testContext) Text() string                { return c.text }
func (c *testContext) Send(v any, opts ...any) error {
	c.sent = append(c.sent, v)
	c.sendOpts = append(c.sendOpts, opts)
	return c.sendErr
}

func (c *testContext) Respond(r ...*telebot.CallbackResponse) error {
	c.responses = append(c.responses, r...)
	return nil
}

type fakeTranslator struct{ values map[string]string }

func (t fakeTranslator) T(key, lang string, _ ...map[string]any) string {
	if v, ok := t.values[key+":"+lang]; ok {
		return v
	}
	return key + ":" + lang
}

type fakeRepo struct {
	lang   string
	err    error
	person *domain.Person
}

func (r *fakeRepo) AddNewPerson(int64, string, string, string, *string) error { return nil }
func (r *fakeRepo) GetPerson(int64) (*domain.Person, error)                   { return r.person, r.err }
func (r *fakeRepo) GetPersonLang(int64) (string, error)                       { return r.lang, r.err }
func (r *fakeRepo) UpdateLang(int64, string) error                            { return nil }

type fakeLanguages struct {
	result  *port.LanguageResult
	err     error
	gotID   int64
	gotLang string
}

func (f *fakeLanguages) ChooseLanguage(id int64, lang string) (*port.LanguageResult, error) {
	f.gotID = id
	f.gotLang = lang
	return f.result, f.err
}

type fakeStarter struct {
	result *port.StartResult
	err    error
	args   []any
}

func (f *fakeStarter) Start(id int64, username, fullname, lang string) (*port.StartResult, error) {
	f.args = []any{id, username, fullname, lang}
	return f.result, f.err
}

type fakeLogger struct{ calls int }

func (f *fakeLogger) Error(string, error, string) { f.calls++ }
func (f *fakeLogger) Fatal(string, error, string) {}

func TestKeyboardBuilders(t *testing.T) {
	tr := fakeTranslator{values: map[string]string{"language_btn:en": "Language"}}
	menu := UserMenu(tr, "en")
	if len(menu.ReplyKeyboard) != 1 || menu.ReplyKeyboard[0][0].Text != "Language" {
		t.Fatalf("unexpected menu: %#v", menu.ReplyKeyboard)
	}
	kb := ChoosingLangInline([]string{"en", "xx"})
	if len(kb.InlineKeyboard) != 2 || kb.InlineKeyboard[0][0].Text != languageNames["en"] ||
		kb.InlineKeyboard[1][0].Text != "xx" {
		t.Fatalf("unexpected inline keyboard: %#v", kb.InlineKeyboard)
	}
}

func TestLanguageAndTextHandlers(t *testing.T) {
	tr := fakeTranslator{values: map[string]string{"language_btn:en": "Choose"}}
	c := &testContext{sender: &telebot.User{ID: 7}, text: "Choose"}
	repo := &fakeRepo{lang: "ru"}
	if err := handleOnText(c, tr, repo, []string{"en"}, "en"); err != nil || len(c.sent) != 1 {
		t.Fatalf("handleOnText: %v %#v", err, c.sent)
	}
	c.text = "other"
	c.sent = nil
	if err := handleOnText(c, tr, repo, []string{"en"}, "en"); err != nil || len(c.sent) != 0 {
		t.Fatal("unexpected response for unrelated text")
	}
	if got := currentLanguage(c, &fakeRepo{err: errors.New("db")}, "en"); got != "en" {
		t.Errorf("currentLanguage fallback = %q", got)
	}
}

func TestChooseLanguageBranches(t *testing.T) {
	tr := fakeTranslator{}
	base := &testContext{sender: &telebot.User{ID: 9}, callback: &telebot.Callback{Data: "bad"}}
	if err := onChooseLang(base, &fakeLanguages{}, tr, &fakeRepo{}, "en"); err != nil || len(base.responses) != 1 {
		t.Fatal("invalid callback branch")
	}
	for _, tc := range []struct {
		name string
		err  error
		key  string
	}{
		{"unsupported", domain.ErrInvalidLanguage, "callback_unsupported_language"}, {"persist", domain.ErrLanguagePersistence, "callback_update_language_failed"},
	} {
		c := &testContext{sender: &telebot.User{ID: 9}, callback: &telebot.Callback{Data: "choose:" + tc.name}}
		if err := onChooseLang(c, &fakeLanguages{err: tc.err}, tr, &fakeRepo{}, "en"); err != nil ||
			len(c.responses) != 1 {
			t.Errorf("%s: %v", tc.name, err)
		}
	}
	c := &testContext{sender: &telebot.User{ID: 9}, callback: &telebot.Callback{Data: "choose:ru"}}
	l := &fakeLanguages{result: &port.LanguageResult{Language: "ru", MessageKey: "welcome"}}
	if err := onChooseLang(c, l, tr, &fakeRepo{}, "en"); err != nil || l.gotLang != "ru" || len(c.sent) != 1 ||
		len(c.responses) != 1 {
		t.Fatalf("success branch: %v %#v", err, c)
	}
	c.sendErr = fmt.Errorf("send")
	if err := onChooseLang(c, l, tr, &fakeRepo{}, "en"); err == nil {
		t.Fatal("expected send error")
	}
}

func TestStartAndMiddleware(t *testing.T) {
	tr := fakeTranslator{}
	logger := &fakeLogger{}
	c := &testContext{sender: &telebot.User{ID: 1, Username: "bob", FirstName: "Bob", LanguageCode: "en"}}
	if err := onStart(c, &fakeStarter{result: &port.StartResult{Language: "en", MessageKey: "hi"}}, tr, logger); err != nil ||
		len(c.sent) != 1 {
		t.Fatal(err)
	}
	c.sender.IsBot = true
	c.sent = nil
	if err := onStart(c, &fakeStarter{}, tr, logger); err != nil || len(c.sent) != 0 {
		t.Fatal("bot should be ignored")
	}
	c.sender.IsBot = false
	c.sent = nil
	if err := onStart(c, &fakeStarter{err: errors.New("db")}, tr, logger); err != nil || logger.calls != 1 {
		t.Fatal("start error branch")
	}
	c.message = &telebot.Message{Sender: c.sender, Chat: &telebot.Chat{ID: 2}, Text: "hello"}
	c.callback = nil
	called := false
	if err := LoggingMiddleware(func(telebot.Context) error { called = true; return nil })(c); err != nil || !called {
		t.Fatal("logging middleware")
	}
	if !IsBotBlockedError(errors.New(BotBlockedErrorMsg)) || !IsBotBlockedError(errors.New(BotDeactivatedErrorMsg)) ||
		IsBotBlockedError(nil) {
		t.Fatal("blocked error matching")
	}
}

func TestHTMLParseMiddleware(t *testing.T) {
	c := &testContext{}
	if err := HTMLParseMiddleware(func(ctx telebot.Context) error { return ctx.Send("hello", &telebot.ReplyMarkup{}) })(c); err != nil {
		t.Fatal(err)
	}
	if len(c.sendOpts) != 1 {
		t.Fatal("send was not captured")
	}
	if _, ok := c.sendOpts[0][0].(*telebot.SendOptions); !ok {
		t.Fatalf("expected SendOptions, got %#v", c.sendOpts[0])
	}
	c.sent = nil
	c.sendOpts = nil
	photo := &telebot.Photo{Caption: "<b>x</b>"}
	if err := HTMLParseMiddleware(func(ctx telebot.Context) error { return ctx.Send(photo) })(c); err != nil {
		t.Fatal(err)
	}
	if opts, ok := c.sendOpts[0][0].(*telebot.SendOptions); !ok || opts.ParseMode != telebot.ModeHTML {
		t.Fatalf("photo opts = %#v", c.sendOpts)
	}
}

func TestLogErrorOptionsAndCallbackLogging(t *testing.T) {
	LogError(
		"oops",
		errors.New("bad"),
		WithUserID(4),
		WithServerName("srv"),
		WithSessionID("sess"),
		WithComponent("comp"),
	)
	c := &testContext{
		sender:   &telebot.User{ID: 3, Username: "u"},
		callback: &telebot.Callback{Sender: &telebot.User{ID: 3, Username: "u"}, Data: "x"},
	}
	if err := LoggingMiddleware(func(telebot.Context) error { return nil })(c); err != nil {
		t.Fatal(err)
	}
}
