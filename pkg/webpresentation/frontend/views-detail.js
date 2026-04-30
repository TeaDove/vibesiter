function formatDate(input) {
  const date = new Date(input);
  if (Number.isNaN(date.getTime())) {
    return "Unknown";
  }
  return new Intl.DateTimeFormat("en-US", {
    dateStyle: "full",
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

export function renderDetailView(root, input) {
  root.innerHTML = "";

  const section = document.createElement("section");
  section.className = "page-card";

  if (!input.app) {
    const title = document.createElement("h1");
    title.className = "page-title";
    title.textContent = "Site not found";

    const subtitle = document.createElement("p");
    subtitle.className = "page-subtitle";
    subtitle.textContent = "The selected site is unavailable in the current list.";

    const actions = document.createElement("div");
    actions.className = "inline-actions";

    const backButton = document.createElement("button");
    backButton.type = "button";
    backButton.className = "btn btn-secondary";
    backButton.textContent = "Back to list";
    backButton.addEventListener("click", input.onBack);

    actions.append(backButton);
    section.append(title, subtitle, actions);
    root.append(section);
    return;
  }

  const heading = document.createElement("h1");
  heading.className = "page-title";
  heading.textContent = input.app.title || input.app.slug;

  const subtitle = document.createElement("p");
  subtitle.className = "page-subtitle";
  subtitle.textContent = input.app.description || "No description available.";

  const metaList = document.createElement("div");
  metaList.className = "meta-list";
  metaList.innerHTML = `
    <div class="meta-row"><span class="meta-key">Slug</span><span class="mono">${input.app.slug}</span></div>
    <div class="meta-row"><span class="meta-key">Created</span><span>${formatDate(input.app.createdAt)}</span></div>
    <div class="meta-row"><span class="meta-key">Status</span><span><span class="status-pill ${statusClass(input.app.status)}">${input.app.status || "UNKNOWN"}</span></span></div>
  `;

  const actions = document.createElement("div");
  actions.className = "inline-actions";

  const openButton = document.createElement("a");
  openButton.className = "btn btn-primary";
  openButton.href = `/apps/${encodeURIComponent(input.app.slug)}/index.html`;
  openButton.target = "_blank";
  openButton.rel = "noopener noreferrer";
  openButton.textContent = "Open site";

  const backButton = document.createElement("button");
  backButton.className = "btn btn-secondary";
  backButton.type = "button";
  backButton.textContent = "Back to list";
  backButton.addEventListener("click", input.onBack);

  const createButton = document.createElement("button");
  createButton.className = "btn btn-secondary";
  createButton.type = "button";
  createButton.textContent = "Create another";
  createButton.addEventListener("click", input.onCreate);

  actions.append(openButton, backButton, createButton);

  const preview = document.createElement("div");
  preview.className = "hero-preview";
  preview.innerHTML = `<iframe title="Site preview ${input.app.slug}" src="/apps/${encodeURIComponent(input.app.slug)}/index.html" loading="lazy" referrerpolicy="same-origin"></iframe>`;

  section.append(heading, subtitle, metaList, actions, preview);
  root.append(section);
}
