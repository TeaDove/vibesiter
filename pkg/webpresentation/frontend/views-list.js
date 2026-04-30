function formatDate(input) {
  const date = new Date(input);
  if (Number.isNaN(date.getTime())) {
    return "Unknown";
  }
  return new Intl.DateTimeFormat("en-US", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(date);
}

function statusClass(status) {
  if (status === "ACTIVE") {
    return "ok";
  }
  if (status === "CREATING") {
    return "warn";
  }
  return "muted";
}

function createPreview(slug) {
  const wrap = document.createElement("div");
  wrap.className = "preview-wrap";

  const fallback = document.createElement("div");
  fallback.className = "preview-fallback";
  fallback.textContent = "Preview is loading...";

  const frame = document.createElement("iframe");
  frame.className = "preview-frame";
  frame.loading = "lazy";
  frame.src = `/apps/${encodeURIComponent(slug)}/index.html`;
  frame.title = `Preview ${slug}`;
  frame.referrerPolicy = "same-origin";
  frame.addEventListener("load", () => fallback.classList.add("hidden"));
  frame.addEventListener("error", () => {
    fallback.classList.remove("hidden");
    fallback.textContent = "Preview is unavailable.";
  });

  wrap.append(frame, fallback);
  return wrap;
}

function createCard(app, handlers) {
  const card = document.createElement("article");
  card.className = "site-card";

  const preview = createPreview(app.slug);
  preview.style.cursor = "pointer";
  preview.addEventListener("click", () => handlers.onOpenDetail(app.slug));

  const content = document.createElement("div");
  content.className = "site-content";

  const title = document.createElement("h3");
  title.className = "site-title";
  title.textContent = app.title || app.slug;

  const description = document.createElement("p");
  description.className = "site-description";
  description.textContent = app.description || "No description yet.";

  const status = document.createElement("span");
  status.className = `status-pill ${statusClass(app.status)}`;
  status.textContent = app.status || "UNKNOWN";

  const created = document.createElement("div");
  created.className = "notice";
  created.textContent = `Created: ${formatDate(app.createdAt)}`;

  const actions = document.createElement("div");
  actions.className = "inline-actions";

  const briefButton = document.createElement("button");
  briefButton.className = "btn btn-secondary";
  briefButton.type = "button";
  briefButton.textContent = "Open brief";
  briefButton.addEventListener("click", () => handlers.onOpenDetail(app.slug));

  const openSiteLink = document.createElement("a");
  openSiteLink.className = "btn btn-primary";
  openSiteLink.href = `/apps/${encodeURIComponent(app.slug)}/index.html`;
  openSiteLink.target = "_blank";
  openSiteLink.rel = "noopener noreferrer";
  openSiteLink.textContent = "Open site";

  actions.append(briefButton, openSiteLink);
  content.append(title, description, status, created, actions);
  card.append(preview, content);
  return card;
}

export function renderListView(root, input) {
  root.innerHTML = "";

  const card = document.createElement("section");
  card.className = "page-card";

  const header = document.createElement("header");
  header.className = "page-header";

  const headingWrap = document.createElement("div");
  const heading = document.createElement("h1");
  heading.className = "page-title";
  heading.textContent = "Generated Sites";
  const subtitle = document.createElement("p");
  subtitle.className = "page-subtitle";
  subtitle.textContent = "Browse your generated projects and launch any site instantly.";
  headingWrap.append(heading, subtitle);

  const actions = document.createElement("div");
  actions.className = "inline-actions";

  const refreshButton = document.createElement("button");
  refreshButton.className = "btn btn-secondary";
  refreshButton.type = "button";
  refreshButton.textContent = input.loading ? "Refreshing..." : "Refresh";
  refreshButton.disabled = input.loading;
  refreshButton.addEventListener("click", input.onRefresh);

  const createButton = document.createElement("button");
  createButton.className = "btn btn-primary";
  createButton.type = "button";
  createButton.textContent = "Create site";
  createButton.addEventListener("click", input.onOpenCreate);

  actions.append(refreshButton, createButton);
  header.append(headingWrap, actions);
  card.append(header);

  if (input.errorMessage) {
    const error = document.createElement("p");
    error.className = "error-box";
    error.textContent = input.errorMessage;
    card.append(error);
  }

  if (input.loading && input.apps.length === 0) {
    const loading = document.createElement("p");
    loading.className = "notice";
    loading.textContent = "Loading sites...";
    card.append(loading);
    root.append(card);
    return;
  }

  if (input.apps.length === 0) {
    const empty = document.createElement("div");
    empty.className = "empty";
    empty.textContent = "No sites yet. Create your first one.";
    card.append(empty);
    root.append(card);
    return;
  }

  const grid = document.createElement("div");
  grid.className = "cards-grid";
  input.apps.forEach((app) => {
    grid.append(
      createCard(app, {
        onOpenDetail: input.onOpenDetail,
      }),
    );
  });
  card.append(grid);
  root.append(card);
}
