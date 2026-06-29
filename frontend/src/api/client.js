const BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080/api/v1";

class ApiError extends Error {
  constructor(message, status) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

async function request(path, options = {}) {
  const res = await fetch(`${BASE_URL}${path}`, {
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
    ...options,
  });

  if (res.status === 204) {
    return null;
  }

  const isJSON = res.headers.get("content-type")?.includes("application/json");
  const body = isJSON ? await res.json().catch(() => null) : null;

  if (!res.ok) {
    const message = body?.message || body?.error || `request failed (${res.status})`;
    throw new ApiError(message, res.status);
  }

  return body;
}

export const api = {
  register: (payload) =>
    request("/auth/register", { method: "POST", body: JSON.stringify(payload) }),
  login: (payload) =>
    request("/auth/login", { method: "POST", body: JSON.stringify(payload) }),
  logout: () => request("/auth/logout", { method: "POST" }),
  me: () => request("/auth/me"),

  listTasks: (params = {}) => {
    const query = new URLSearchParams();
    if (params.completed !== undefined && params.completed !== null) {
      query.set("completed", String(params.completed));
    }
    if (params.archived !== undefined && params.archived !== null) {
      query.set("archived", String(params.archived));
    }
    if (params.priority) query.set("priority", params.priority);
    if (params.tag) query.set("tag", params.tag);
    if (params.folderId) query.set("folder_id", params.folderId);
    if (params.limit) query.set("limit", String(params.limit));
    if (params.offset) query.set("offset", String(params.offset));
    const qs = query.toString();
    return request(`/tasks${qs ? `?${qs}` : ""}`);
  },
  createTask: (payload) =>
    request("/tasks", { method: "POST", body: JSON.stringify(payload) }),
  patchTask: (id, payload) =>
    request(`/tasks/${id}`, { method: "PATCH", body: JSON.stringify(payload) }),
  archiveTask: (id) => request(`/tasks/${id}/archive`, { method: "POST" }),
  unarchiveTask: (id) => request(`/tasks/${id}/unarchive`, { method: "POST" }),
  permanentlyDeleteTask: (id) => request(`/tasks/${id}/permanent`, { method: "DELETE" }),
  completeTask: (id) => request(`/tasks/${id}/complete`, { method: "POST" }),
  uncompleteTask: (id) => request(`/tasks/${id}/uncomplete`, { method: "POST" }),

  listTags: () => request("/tags"),

  listFolders: () => request("/folders"),
  createFolder: (payload) =>
    request("/folders", { method: "POST", body: JSON.stringify(payload) }),
  deleteFolder: (id) => request(`/folders/${id}`, { method: "DELETE" }),
};

export { ApiError };
