(function () {
  const root = document.documentElement;
  const themeStorageKey = "attendanceQuestTheme";
  const sidebarStorageKey = "attendanceQuestSidebar";
  const modes = new Set(["light", "dark"]);
  const backgrounds = new Set([
    "red", "blue", "green", "yellow", "orange", "pink", "purple",
    "beach", "forest", "sky", "meadow", "mountain", "sunset",
  ]);

  let settings = {};
  try {
    settings = JSON.parse(localStorage.getItem(themeStorageKey) || "{}");
  } catch {
    settings = {};
  }

  root.dataset.studentMode = modes.has(settings.mode) ? settings.mode : "dark";
  root.dataset.studentBg = backgrounds.has(settings.background) ? settings.background : "green";
  root.dataset.sidebar = localStorage.getItem(sidebarStorageKey) === "collapsed" ? "collapsed" : "expanded";
})();
