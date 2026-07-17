(function () {
  const dialog = document.querySelector("[data-avatar-preview-dialog]");
  if (!dialog) {
    return;
  }

  const closeButton = dialog.querySelector("[data-avatar-preview-close]");
  closeButton?.addEventListener("click", () => dialog.close());
  dialog.addEventListener("close", () => window.location.replace("/avatar"));

  if (dialog.hasAttribute("data-open-on-load")) {
    dialog.showModal();
  }
})();
