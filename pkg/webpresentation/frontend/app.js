import { createApp, fetchApps } from "/api.js";
import { currentRoute, navigate, startRouter } from "/router.js";
import { findApp, getState, setState, subscribe, upsertApp } from "/state.js";
import { renderCreateView } from "/views-create.js";
import { renderDetailView } from "/views-detail.js";
import { renderListView } from "/views-list.js";

const appRoot = document.getElementById("app");
const toast = document.getElementById("toast");

let activeRoute = currentRoute();
let toastTimer = null;

function showToast(message) {
  if (!message) {
    return;
  }
  clearTimeout(toastTimer);
  toast.textContent = message;
  toast.classList.add("visible");
  toastTimer = setTimeout(() => {
    toast.classList.remove("visible");
  }, 5000);
}

function normalizeError(error) {
  if (!error) {
    return "Unknown error.";
  }
  if (typeof error === "string") {
    return error;
  }
  if (error instanceof Error && error.message) {
    return error.message;
  }
  return "Request failed.";
}

async function loadApps() {
  setState({ loadingApps: true, errorMessage: "" });
  try {
    const apps = await fetchApps(0);
    setState({ apps: Array.isArray(apps) ? apps : [], loadingApps: false });
  } catch (error) {
    const message = normalizeError(error);
    setState({ loadingApps: false, errorMessage: message });
    showToast(message);
  }
}

async function handleGenerate(prompt) {
  setState({ creatingApp: true, errorMessage: "" });
  try {
    const app = await createApp(prompt);
    upsertApp(app);
    setState({ creatingApp: false, errorMessage: "" });
    navigate(`/app/${encodeURIComponent(app.slug)}`);
  } catch (error) {
    const message = normalizeError(error);
    setState({ creatingApp: false, errorMessage: message });
    showToast(message);
  }
}

function openSiteBySlug(slug) {
  const cleanSlug = String(slug || "").trim();
  if (!cleanSlug) {
    showToast("Slug is required.");
    return;
  }
  window.open(`/apps/${encodeURIComponent(cleanSlug)}/index.html`, "_blank", "noopener");
}

function render() {
  const state = getState();

  if (activeRoute.name === "create") {
    renderCreateView(appRoot, {
      creating: state.creatingApp,
      errorMessage: state.errorMessage,
      onBack: () => navigate("/"),
      onGenerate: handleGenerate,
      onOpenSite: openSiteBySlug,
    });
    return;
  }

  if (activeRoute.name === "detail") {
    renderDetailView(appRoot, {
      app: findApp(activeRoute.slug),
      onBack: () => navigate("/"),
      onCreate: () => navigate("/create"),
    });
    return;
  }

  renderListView(appRoot, {
    apps: state.apps,
    loading: state.loadingApps,
    errorMessage: state.errorMessage,
    onRefresh: loadApps,
    onOpenCreate: () => navigate("/create"),
    onOpenDetail: (slug) => navigate(`/app/${encodeURIComponent(slug)}`),
  });
}

async function onRouteChange(route) {
  activeRoute = route;
  render();

  if (route.name === "list") {
    await loadApps();
    return;
  }

  if (route.name === "detail" && !findApp(route.slug)) {
    await loadApps();
    render();
  }
}

subscribe(render);
startRouter(onRouteChange);
