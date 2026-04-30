export function renderCreateView(root, input) {
  root.innerHTML = "";

  const section = document.createElement("section");
  section.className = "page-card";

  const heading = document.createElement("h1");
  heading.className = "page-title";
  heading.textContent = "Create a New Site";

  const subtitle = document.createElement("p");
  subtitle.className = "page-subtitle";
  subtitle.textContent = "Describe your idea, generate the site, then open it immediately.";

  const form = document.createElement("form");
  form.className = "form";

  const promptLabel = document.createElement("label");
  promptLabel.textContent = "Site prompt";

  const promptInput = document.createElement("textarea");
  promptInput.className = "textarea";
  promptInput.name = "userPrompt";
  promptInput.placeholder = "A portfolio page with projects, social links, and contact section.";
  promptInput.required = true;
  promptInput.disabled = input.creating;

  const actions = document.createElement("div");
  actions.className = "inline-actions";

  const generateButton = document.createElement("button");
  generateButton.className = "btn btn-primary";
  generateButton.type = "submit";
  generateButton.textContent = input.creating ? "Generating..." : "Generate";
  generateButton.disabled = input.creating;

  const backButton = document.createElement("button");
  backButton.className = "btn btn-secondary";
  backButton.type = "button";
  backButton.textContent = "Back to list";
  backButton.disabled = input.creating;
  backButton.addEventListener("click", input.onBack);

  actions.append(generateButton, backButton);
  promptLabel.append(promptInput);
  form.append(promptLabel, actions);

  form.addEventListener("submit", async (event) => {
    event.preventDefault();
    const prompt = promptInput.value.trim();
    if (!prompt) {
      return;
    }
    await input.onGenerate(prompt);
  });

  if (input.errorMessage) {
    const error = document.createElement("p");
    error.className = "error-box";
    error.textContent = input.errorMessage;
    form.prepend(error);
  }

  const quickOpenCard = document.createElement("section");
  quickOpenCard.className = "page-card";

  const quickOpenTitle = document.createElement("h2");
  quickOpenTitle.className = "page-title";
  quickOpenTitle.textContent = "Quick open by slug";

  const quickOpenForm = document.createElement("form");
  quickOpenForm.className = "form";

  const slugLabel = document.createElement("label");
  slugLabel.textContent = "Slug";

  const slugInput = document.createElement("input");
  slugInput.className = "input";
  slugInput.placeholder = "my-generated-site";
  slugInput.required = true;
  slugLabel.append(slugInput);

  const quickOpenButton = document.createElement("button");
  quickOpenButton.className = "btn btn-secondary";
  quickOpenButton.type = "submit";
  quickOpenButton.textContent = "Open site";

  quickOpenForm.append(slugLabel, quickOpenButton);
  quickOpenForm.addEventListener("submit", (event) => {
    event.preventDefault();
    const slug = slugInput.value.trim();
    if (!slug) {
      return;
    }
    input.onOpenSite(slug);
  });

  quickOpenCard.append(quickOpenTitle, quickOpenForm);
  section.append(heading, subtitle, form);
  root.append(section, quickOpenCard);
}
