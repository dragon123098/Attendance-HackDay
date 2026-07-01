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
	if len(avatarBaseCatalog) != 9 {
		t.Fatalf("avatarBaseCatalog has %d options, want 9", len(avatarBaseCatalog))
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

func TestAvatarCosmeticCatalogUsesVisualOverlayImages(t *testing.T) {
	if len(avatarCosmeticCatalog) != 12 {
		t.Fatalf("avatarCosmeticCatalog has %d options, want 12", len(avatarCosmeticCatalog))
	}

	for _, option := range avatarCosmeticCatalog {
		if !strings.HasPrefix(option.Image, "/static/images/cosmetics/") {
			t.Fatalf("%s image path = %q, want cosmetic overlay path", option.ID, option.Image)
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

func TestAvatarPreviewIncludesSelectedCosmeticLayers(t *testing.T) {
	preview := buildAvatarPreview(&AvatarConfig{
		Base:      "gerald",
		HairStyle: "hat_star",
		Clothing:  "cape_gold",
		Accessory: "glasses_rocket",
		Effect:    "trail_rainbow",
	})

	if len(preview.Layers) != 4 {
		t.Fatalf("preview has %d layers, want 4", len(preview.Layers))
	}

	for _, layer := range preview.Layers {
		if layer.Image == "" {
			t.Fatalf("layer %s has empty image", layer.ID)
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
	for _, want := range []string{"BrainRot", "D-Money", "Funk Rapper", "Gerald", "Gopher", "Mike", "Milk Man", "Peter", "Salary Man", "Star Hat", "Wizard Hat", "Rocket Glasses", "Sparkle Aura", "Locked"} {
		if !strings.Contains(body, want) {
			t.Fatalf("avatar page did not contain %q\n%s", want, body)
		}
	}
}

func TestSeedShopItemsAddsVisualCosmeticsAndSpecialBackgrounds(t *testing.T) {
	resetAvatarTestState(t)
	withTempWorkingDir(t)
	app.ShopItems = map[string]*ShopItem{
		"hat_star": &ShopItem{ID: "hat_star", Name: "Custom Star Hat", Price: 99},
	}

	seedShopItems()

	wantCount := len(avatarCosmeticCatalog) + len(specialThemeBackgroundCatalog)
	if len(app.ShopItems) != wantCount {
		t.Fatalf("shop item count = %d, want %d", len(app.ShopItems), wantCount)
	}
	if app.ShopItems["hat_star"].Name != "Custom Star Hat" {
		t.Fatal("seedShopItems overwrote an existing shop item")
	}
	for _, cosmetic := range avatarCosmeticCatalog {
		if _, ok := app.ShopItems[cosmetic.ID]; !ok {
			t.Fatalf("missing seeded shop item %q", cosmetic.ID)
		}
	}
	for _, background := range specialThemeBackgroundCatalog {
		item, ok := app.ShopItems[background.ShopItemID]
		if !ok {
			t.Fatalf("missing seeded background item %q", background.ShopItemID)
		}
		if item.Price != themeBackgroundPrice {
			t.Fatalf("%s price = %d, want %d", background.ShopItemID, item.Price, themeBackgroundPrice)
		}
	}
}

func TestGetShopItemViewsIncludesVisualPreviews(t *testing.T) {
	resetAvatarTestState(t)
	app.OwnedShopItems["student1"] = []string{"hat_star", "glasses_rocket", "background_beach"}

	avatarItems, backgroundItems, owned := getShopItemViews("student1")

	if len(avatarItems) != len(avatarCosmeticCatalog) {
		t.Fatalf("avatar item count = %d, want %d", len(avatarItems), len(avatarCosmeticCatalog))
	}
	if len(backgroundItems) != len(specialThemeBackgroundCatalog) {
		t.Fatalf("background item count = %d, want %d", len(backgroundItems), len(specialThemeBackgroundCatalog))
	}

	if len(owned) != 3 {
		t.Fatalf("owned item count = %d, want 3", len(owned))
	}

	for _, view := range avatarItems {
		if view.Image == "" {
			t.Fatalf("avatar shop item %q has empty image", view.ID)
		}
		if view.Slot == "" {
			t.Fatalf("avatar shop item %q has empty slot", view.ID)
		}
	}
	for _, view := range backgroundItems {
		if view.Slot != shopItemSlotTheme {
			t.Fatalf("background shop item %q slot = %q, want %q", view.ID, view.Slot, shopItemSlotTheme)
		}
		if view.ThemeBackgroundID == "" {
			t.Fatalf("background shop item %q has empty preview ID", view.ID)
		}
	}
	for _, view := range owned {
		if view.Image == "" && view.ThemeBackgroundID == "" {
			t.Fatalf("owned item %q has empty visual preview", view.ID)
		}
	}
}

func TestOwnedThemeBackgroundOptionsOnlyIncludesPurchasedBackgrounds(t *testing.T) {
	resetAvatarTestState(t)
	app.OwnedShopItems["student1"] = []string{"background_beach", "hat_star"}

	options := ownedThemeBackgroundOptionViews("student1")

	if len(options) != 1 {
		t.Fatalf("owned background count = %d, want 1", len(options))
	}
	if options[0].ID != "beach" || options[0].Label != "Beach" {
		t.Fatalf("owned background option = %#v", options[0])
	}
}

func TestShopViewShowsOwnedSpecialBackgroundSwatches(t *testing.T) {
	resetAvatarTestState(t)
	app.OwnedShopItems["student1"] = []string{"background_beach"}

	req := authedAvatarRequest(t, http.MethodGet, "/shop", nil)
	rec := httptest.NewRecorder()

	shopView(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	for _, want := range []string{`data-bg-value="red"`, `data-bg-value="beach"`, `data-theme-preview="beach"`} {
		if !strings.Contains(body, want) {
			t.Fatalf("shop page did not contain %q\n%s", want, body)
		}
	}
	for _, want := range []string{"Avatar Items", "Backgrounds"} {
		if !strings.Contains(body, want) {
			t.Fatalf("shop page did not contain section header %q\n%s", want, body)
		}
	}
	if strings.Contains(body, `data-bg-value="forest"`) {
		t.Fatalf("shop page showed unowned forest theme swatch\n%s", body)
	}
}

func TestGetCoinBalanceIncludesManualAdjustments(t *testing.T) {
	resetAvatarTestState(t)
	app.ManualCoinAdjustments["student1"] = 25
	app.Transactions = []CoinTransaction{
		{UserID: "student1", Amount: 2, Description: "Attendance reward"},
		{UserID: "student1", Amount: -5, Description: "Purchased Star Hat"},
	}

	if balance := getCoinBalance("student1"); balance != 32 {
		t.Fatalf("coin balance = %d, want 32", balance)
	}
}

func TestShopBuyCanPurchaseSpecialBackground(t *testing.T) {
	resetAvatarTestState(t)
	withTempWorkingDir(t)
	app.Transactions = []CoinTransaction{
		{UserID: "student1", Amount: 10, Description: "Test bonus"},
	}

	form := url.Values{"item_id": {"background_beach"}}
	req := authedAvatarRequest(t, http.MethodPost, "/shop/buy", form)
	rec := httptest.NewRecorder()

	shopBuyView(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusSeeOther)
	}
	if !userOwnsShopItem("student1", "background_beach") {
		t.Fatal("special background was not added to owned shop items")
	}
	if balance := getCoinBalance("student1"); balance != 5 {
		t.Fatalf("coin balance = %d, want 5", balance)
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
		ShopItems:             seededShopItemMap(),
		OwnedShopItems:        map[string][]string{},
		AvatarConfigs:         map[string]*AvatarConfig{},
		ManualCoinAdjustments: map[string]int{},
	}

	t.Cleanup(func() {
		app = previousApp
		sessionMu.Lock()
		sessionStore = previousSessions
		sessionMu.Unlock()
	})
}

func seededShopItemMap() map[string]*ShopItem {
	items := map[string]*ShopItem{}
	for _, item := range seededShopItems() {
		copied := *item
		items[item.ID] = &copied
	}
	return items
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
