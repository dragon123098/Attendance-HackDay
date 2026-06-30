package web

import (
	"image"
	_ "image/png"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/dragon123098/Attendance-HackDay.git/internal/view"
)

func TestSavedAvatarConfigDefaultsAndDropsUnavailableValues(t *testing.T) {
	resetAvatarTestState(t)
	app.OwnedShopItems["student1"] = []string{"hat_star"}
	app.AvatarConfigs["student1"] = &AvatarConfig{
		Base:      "missing_base",
		HairStyle: "hat_star",
		Clothing:  "cape_gold",
		Accessory: "missing_cosmetic",
		Effect:    "trail_rainbow",
	}

	cfg := savedAvatarConfig("student1")

	if cfg.Base != defaultAvatarBaseID {
		t.Fatalf("Base = %q, want %q", cfg.Base, defaultAvatarBaseID)
	}
	if cfg.HairStyle != "hat_star" {
		t.Fatalf("HairStyle = %q, want hat_star", cfg.HairStyle)
	}
	if cfg.Clothing != "" || cfg.Accessory != "" || cfg.Effect != "" {
		t.Fatalf("unavailable cosmetics were not dropped: %#v", cfg)
	}
}

func TestValidateAvatarConfigAcceptsOwnedCosmetics(t *testing.T) {
	resetAvatarTestState(t)
	app.OwnedShopItems["student1"] = []string{"hat_star", "cape_gold", "glasses_rocket", "trail_rainbow"}

	cfg, err := validateAvatarConfig("student1", &AvatarConfig{
		Base:      "mike",
		HairStyle: "hat_star",
		Clothing:  "cape_gold",
		Accessory: "glasses_rocket",
		Effect:    "trail_rainbow",
	})
	if err != nil {
		t.Fatalf("validateAvatarConfig returned error: %v", err)
	}

	if cfg.Base != "mike" || cfg.Effect != "trail_rainbow" {
		t.Fatalf("validated config did not preserve selections: %#v", cfg)
	}
}

func TestAvatarBaseCatalogUsesNormalizedStaticImages(t *testing.T) {
	if len(avatarBaseCatalog) != 6 {
		t.Fatalf("avatarBaseCatalog has %d options, want 6", len(avatarBaseCatalog))
	}

	for _, option := range avatarBaseCatalog {
		if !strings.HasPrefix(option.Image, "/static/images/avatars/") {
			t.Fatalf("%s image path = %q, want generated avatar path", option.ID, option.Image)
		}

		file, err := view.FS.Open(strings.TrimPrefix(option.Image, "/"))
		if err != nil {
			t.Fatalf("%s image does not exist: %v", option.ID, err)
		}

		cfg, _, err := image.DecodeConfig(file)
		if closeErr := file.Close(); closeErr != nil {
			t.Fatalf("close %s image: %v", option.ID, closeErr)
		}
		if err != nil {
			t.Fatalf("decode %s image config: %v", option.ID, err)
		}
		if cfg.Width != 512 || cfg.Height != 512 {
			t.Fatalf("%s image is %dx%d, want 512x512", option.ID, cfg.Width, cfg.Height)
		}
	}
}

func TestAvatarViewShowsAvailableAndLockedOptions(t *testing.T) {
	resetAvatarTestState(t)
	app.OwnedShopItems["student1"] = []string{"hat_star"}
	app.AvatarConfigs["student1"] = &AvatarConfig{
		Base:      "d_money",
		HairStyle: "hat_star",
	}

	req := authedAvatarRequest(t, http.MethodGet, "/avatar", nil)
	rec := httptest.NewRecorder()

	avatarView(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	for _, want := range []string{"BrainRot", "D-Money", "Gerald", "Mike", "Milk Man", "Salary Man", "Star Hat", "Rocket Glasses", "Locked"} {
		if !strings.Contains(body, want) {
			t.Fatalf("avatar page did not contain %q\n%s", want, body)
		}
	}
}

func TestAvatarPreviewDoesNotPersistSelection(t *testing.T) {
	resetAvatarTestState(t)
	app.OwnedShopItems["student1"] = []string{"hat_star"}

	form := url.Values{
		"base":       {"milkman"},
		"hair_style": {"hat_star"},
	}
	req := authedAvatarRequest(t, http.MethodPost, "/avatar/preview", form)
	rec := httptest.NewRecorder()

	avatarPreviewView(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if _, ok := app.AvatarConfigs["student1"]; ok {
		t.Fatal("preview persisted avatar config")
	}
	if !strings.Contains(rec.Body.String(), "Previewing unsaved avatar changes.") {
		t.Fatal("preview message was not rendered")
	}
	if !strings.Contains(rec.Body.String(), "Milk Man") {
		t.Fatal("preview selection was not rendered")
	}
}

func TestAvatarSavePersistsSelection(t *testing.T) {
	resetAvatarTestState(t)
	withTempWorkingDir(t)
	app.OwnedShopItems["student1"] = []string{"hat_star"}

	form := url.Values{
		"base":       {"brainrot"},
		"hair_style": {"hat_star"},
	}
	req := authedAvatarRequest(t, http.MethodPost, "/avatar", form)
	rec := httptest.NewRecorder()

	avatarSaveView(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusSeeOther)
	}
	cfg := app.AvatarConfigs["student1"]
	if cfg == nil {
		t.Fatal("avatar config was not saved")
	}
	if cfg.Base != "brainrot" || cfg.HairStyle != "hat_star" {
		t.Fatalf("saved config = %#v", cfg)
	}
	if _, err := os.Stat("data/data.json"); err != nil {
		t.Fatalf("expected persisted data file: %v", err)
	}
}

func TestAvatarSaveRejectsLockedAndUnknownCosmetics(t *testing.T) {
	cases := []struct {
		name    string
		form    url.Values
		message string
	}{
		{
			name: "locked",
			form: url.Values{
				"base":       {"gerald"},
				"hair_style": {"hat_star"},
			},
			message: "You can only equip cosmetics you own.",
		},
		{
			name: "unknown",
			form: url.Values{
				"base":       {"gerald"},
				"hair_style": {"not_real"},
			},
			message: "Choose a valid avatar option.",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resetAvatarTestState(t)

			req := authedAvatarRequest(t, http.MethodPost, "/avatar", tc.form)
			rec := httptest.NewRecorder()

			avatarSaveView(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
			}
			if _, ok := app.AvatarConfigs["student1"]; ok {
				t.Fatal("invalid avatar config was saved")
			}
			if !strings.Contains(rec.Body.String(), tc.message) {
				t.Fatalf("response did not contain %q\n%s", tc.message, rec.Body.String())
			}
		})
	}
}

func resetAvatarTestState(t *testing.T) {
	t.Helper()

	previousApp := app
	sessionMu.Lock()
	previousSessions := sessionStore
	sessionStore = map[string]sessionRecord{}
	sessionMu.Unlock()

	app = AppState{
		Users: map[string]*User{
			"student1": &User{
				Name:        "Test Student",
				Role:        "student",
				Email:       "student@example.com",
				ClassroomID: "classroom1",
				UserID:      "student1",
			},
		},
		ShopItems: map[string]*ShopItem{
			"hat_star":       &ShopItem{ID: "hat_star", Name: "Star Hat", Price: 5},
			"cape_gold":      &ShopItem{ID: "cape_gold", Name: "Golden Cape", Price: 12},
			"glasses_rocket": &ShopItem{ID: "glasses_rocket", Name: "Rocket Glasses", Price: 10},
			"trail_rainbow":  &ShopItem{ID: "trail_rainbow", Name: "Rainbow Trail", Price: 8},
		},
		OwnedShopItems: map[string][]string{},
		AvatarConfigs:  map[string]*AvatarConfig{},
	}

	t.Cleanup(func() {
		app = previousApp
		sessionMu.Lock()
		sessionStore = previousSessions
		sessionMu.Unlock()
	})
}

func authedAvatarRequest(t *testing.T, method, target string, form url.Values) *http.Request {
	t.Helper()

	sessionRecorder := httptest.NewRecorder()
	if err := createSession(sessionRecorder, "student1"); err != nil {
		t.Fatalf("createSession: %v", err)
	}

	var body *strings.Reader
	if form == nil {
		body = strings.NewReader("")
	} else {
		body = strings.NewReader(form.Encode())
	}

	req := httptest.NewRequest(method, target, body)
	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.AddCookie(sessionRecorder.Result().Cookies()[0])

	return req
}

func withTempWorkingDir(t *testing.T) {
	t.Helper()

	previousDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("change working directory: %v", err)
	}

	t.Cleanup(func() {
		if err := os.Chdir(previousDir); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	})
}
