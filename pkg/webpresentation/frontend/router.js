function parseHash(rawHash) {
  const hash = rawHash || "#/";

  if (hash === "#" || hash === "#/" || hash === "") {
    return { name: "list" };
  }

  if (hash === "#/create") {
    return { name: "create" };
  }

  const appMatch = hash.match(/^#\/app\/([^/]+)$/);
  if (appMatch) {
    return { name: "detail", slug: decodeURIComponent(appMatch[1]) };
  }

  return { name: "list" };
}

export function currentRoute() {
  return parseHash(window.location.hash);
}

export function startRouter(onRouteChange) {
  const emit = () => onRouteChange(currentRoute());
  window.addEventListener("hashchange", emit);
  emit();
}

export function navigate(path) {
  const next = path.startsWith("#") ? path : `#${path}`;
  if (window.location.hash === next) {
    window.dispatchEvent(new HashChangeEvent("hashchange"));
    return;
  }
  window.location.hash = next;
}
