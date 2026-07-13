package web

const shopItemSlotTheme = "background"

type themeBackgroundOption struct {
	ID          string
	ShopItemID  string
	Label       string
	Description string
}

var specialThemeBackgroundCatalog = []themeBackgroundOption{
	{
		ID:          "beach",
		ShopItemID:  "background_beach",
		Label:       "Beach",
		Description: "A soft shoreline with sand, sea, and gentle sky colors.",
	},
	{
		ID:          "forest",
		ShopItemID:  "background_forest",
		Label:       "Forest",
		Description: "A calm woodland palette with mossy greens and warm light.",
	},
	{
		ID:          "sky",
		ShopItemID:  "background_sky",
		Label:       "Sky",
		Description: "Airy blues and cloud-soft highlights for a clear day.",
	},
	{
		ID:          "meadow",
		ShopItemID:  "background_meadow",
		Label:       "Meadow",
		Description: "Pastel grass, tiny blooms, and a quiet afternoon feel.",
	},
	{
		ID:          "mountain",
		ShopItemID:  "background_mountain",
		Label:       "Mountain",
		Description: "Cool ridge colors with misty lavender shadows.",
	},
	{
		ID:          "sunset",
		ShopItemID:  "background_sunset",
		Label:       "Sunset",
		Description: "Peach, rose, and gold tones for a mellow evening glow.",
	},
}

func themeBackgroundByShopItemID(itemID string) (themeBackgroundOption, bool) {
	for _, background := range specialThemeBackgroundCatalog {
		if background.ShopItemID == itemID {
			return background, true
		}
	}
	return themeBackgroundOption{}, false
}

func ownedThemeBackgroundOptionViews(ownedIDs []string) []ThemeBackgroundOptionView {
	views := make([]ThemeBackgroundOptionView, 0, len(specialThemeBackgroundCatalog))
	for _, background := range specialThemeBackgroundCatalog {
		if !ownsShopItem(ownedIDs, background.ShopItemID) {
			continue
		}
		views = append(views, ThemeBackgroundOptionView{
			ID:    background.ID,
			Label: background.Label,
		})
	}
	return views
}
