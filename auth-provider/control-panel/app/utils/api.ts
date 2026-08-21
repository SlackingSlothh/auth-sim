export async function apiFetch(path: string, init: RequestInit = {}) {
  const base = (import.meta.env.VITE_BACKEND_URL as string) || "";
  const url = base + path;
  console.log(import.meta.env.VITE_BACKEND_URL);
  console.log(new URL("/users", import.meta.env.VITE_BACKEND_URL).toString());

  const res = await fetch(url, {
    credentials: "include",
    ...init,
    headers: {
      ...(init.headers || {}),
    },
  });

  return res;
}

export async function apiJson(path: string, init: RequestInit = {}) {
  const res = await apiFetch(path, init);
  const text = await res.text();
  try {
    return { res, data: text ? JSON.parse(text) : null };
  } catch {
    return { res, data: text };
  }
}
