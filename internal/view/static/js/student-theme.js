(function () {
  const root = document.documentElement;
  const storageKey = "attendanceQuestTheme";
  const sidebarStorageKey = "attendanceQuestSidebar";
  const defaults = { mode: "dark", background: "green" };
  const freeBackgrounds = new Set(["red", "blue", "green", "yellow", "orange", "pink", "purple"]);
  const savedSidebarState = localStorage.getItem(sidebarStorageKey) === "collapsed" ? "collapsed" : "expanded";

  root.dataset.sidebar = savedSidebarState;

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
    const sidebarToggle = document.querySelector("[data-sidebar-toggle]");
    const sidebarToggleIcon = sidebarToggle ? sidebarToggle.querySelector(".sidebar-toggle-icon") : null;

    function applySidebarState(state) {
      const isCollapsed = state === "collapsed";
      root.dataset.sidebar = state;

      if (!sidebarToggle) {
        return;
      }

      sidebarToggle.setAttribute("aria-expanded", String(!isCollapsed));
      sidebarToggle.setAttribute("aria-label", isCollapsed ? "Expand sidebar" : "Minimize sidebar");

      if (sidebarToggleIcon) {
        sidebarToggleIcon.textContent = isCollapsed ? "\u203a" : "\u2039";
      }
    }

    applySidebarState(savedSidebarState);

    if (sidebarToggle) {
      sidebarToggle.addEventListener("click", () => {
        const nextState = root.dataset.sidebar === "collapsed" ? "expanded" : "collapsed";
        localStorage.setItem(sidebarStorageKey, nextState);
        applySidebarState(nextState);
      });
    }

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
