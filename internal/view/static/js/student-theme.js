(function () {
  const root = document.documentElement;
  const storageKey = "attendanceQuestTheme";
  const defaults = { mode: "light", background: "green" };
  const freeBackgrounds = new Set(["red", "blue", "green", "yellow", "orange", "pink", "purple"]);

  function availableBackgrounds() {
    const backgrounds = new Set(freeBackgrounds);
    document.querySelectorAll("[data-bg-value]").forEach((button) => {
      backgrounds.add(button.dataset.bgValue);
    });
    return backgrounds;
  }

  function readSettings() {
    try {
      return { ...defaults, ...JSON.parse(localStorage.getItem(storageKey) || "{}") };
    } catch {
      return defaults;
    }
  }

  function saveSettings(settings) {
    localStorage.setItem(storageKey, JSON.stringify(settings));
  }

  function applySettings(settings) {
    if (!availableBackgrounds().has(settings.background)) {
      settings.background = defaults.background;
    }

    root.dataset.studentMode = settings.mode;
    root.dataset.studentBg = settings.background;

    document.querySelectorAll("[data-theme-value]").forEach((button) => {
      button.classList.toggle("active", button.dataset.themeValue === settings.mode);
    });
    document.querySelectorAll("[data-bg-value]").forEach((button) => {
      button.classList.toggle("active", button.dataset.bgValue === settings.background);
    });
  }

  let settings = readSettings();
  applySettings(settings);

  document.addEventListener("DOMContentLoaded", () => {
    const panel = document.querySelector(".student-theme-panel");
    if (!panel) {
      return;
    }

    const toggle = panel.querySelector(".theme-panel-toggle");
    const menu = panel.querySelector(".theme-panel-menu");

    toggle.addEventListener("click", () => {
      const isOpen = menu.hasAttribute("hidden");
      menu.toggleAttribute("hidden", !isOpen);
      toggle.setAttribute("aria-expanded", String(isOpen));
    });

    panel.querySelectorAll("[data-theme-value]").forEach((button) => {
      button.addEventListener("click", () => {
        settings = { ...settings, mode: button.dataset.themeValue };
        saveSettings(settings);
        applySettings(settings);
      });
    });

    panel.querySelectorAll("[data-bg-value]").forEach((button) => {
      button.addEventListener("click", () => {
        settings = { ...settings, background: button.dataset.bgValue };
        saveSettings(settings);
        applySettings(settings);
      });
    });

    document.addEventListener("click", (event) => {
      if (panel.contains(event.target)) {
        return;
      }
      menu.hidden = true;
      toggle.setAttribute("aria-expanded", "false");
    });
  });
})();
