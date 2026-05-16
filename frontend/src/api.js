const API_URL = import.meta.env.VITE_API_URL || "/api";
let authToken = localStorage.getItem("token") || "";

export function setAuthToken(token) {
  authToken = token || "";
  if (authToken) {
    localStorage.setItem("token", authToken);
  } else {
    localStorage.removeItem("token");
  }
}

export async function api(path, options = {}) {
  const headers = { "Content-Type": "application/json", ...(options.headers || {}) };
  if (authToken) headers.Authorization = `Bearer ${authToken}`;
  const res = await fetch(`${API_URL}${path}`, {
    headers,
    ...options,
  });
  if (!res.ok) {
    let message = res.statusText;
    try {
      const data = await res.json();
      message = data.error || message;
    } catch {
      // ignore non-json errors
    }
    throw new Error(message);
  }
  return res.json();
}
