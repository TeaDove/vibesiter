const REQUEST_TIMEOUT_MS = 30000;

async function requestJson(url, options = {}) {
  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), REQUEST_TIMEOUT_MS);

  try {
    const response = await fetch(url, {
      ...options,
      headers: {
        "Content-Type": "application/json",
        ...(options.headers || {}),
      },
      signal: controller.signal,
    });

    if (!response.ok) {
      let message = `${response.status} ${response.statusText}`;
      try {
        const payload = await response.json();
        if (payload && typeof payload.error === "string") {
          message = payload.error;
        }
      } catch (_) {
        const rawText = await response.text();
        if (rawText) {
          message = rawText;
        }
      }
      throw new Error(message);
    }

    return response.json();
  } catch (error) {
    if (error.name === "AbortError") {
      throw new Error("Request timed out.");
    }
    throw error;
  } finally {
    clearTimeout(timeout);
  }
}

export async function fetchApps(page = 0) {
  return requestJson(`/apps?page=${encodeURIComponent(page)}`);
}

export async function createApp(userPrompt) {
  return requestJson("/apps/", {
    method: "POST",
    body: JSON.stringify({ userPrompt }),
  });
}
