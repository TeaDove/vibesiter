const state = {
  apps: [],
  loadingApps: false,
  creatingApp: false,
  errorMessage: "",
};

const listeners = new Set();

export function getState() {
  return state;
}

export function setState(patch) {
  Object.assign(state, patch);
  listeners.forEach((listener) => listener(state));
}

export function subscribe(listener) {
  listeners.add(listener);
  return () => listeners.delete(listener);
}

export function upsertApp(nextApp) {
  const existingIndex = state.apps.findIndex((app) => app.slug === nextApp.slug);
  if (existingIndex === -1) {
    state.apps = [nextApp, ...state.apps];
  } else {
    state.apps[existingIndex] = nextApp;
  }
  listeners.forEach((listener) => listener(state));
}

export function findApp(slug) {
  return state.apps.find((app) => app.slug === slug) || null;
}
