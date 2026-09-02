export function emailAssetBaseUrl() {
  // React Email preview serves files from `--dir`/static at /static.
  // Ignore an exported EMAIL_ASSETS_BASE_URL so images do not point at :8080.
  if (process.env.EMAIL_PREVIEW === "1" || process.env.npm_lifecycle_event === "dev") {
    return "";
  }
  if (process.env.EMAIL_ASSETS_BASE_URL) {
    return process.env.EMAIL_ASSETS_BASE_URL.replace(/\/$/, "");
  }

  return "";
}

export function goBoilerplateAsset(path: string) {
  const base = emailAssetBaseUrl();
  return `${base}/static/go_boilerplate/${path}`;
}
