const backendURL: string = "http://localhost:8080"

export async function apiFetch<T>(url: string, init?: RequestInit): Promise<T> {
  const response = await fetch(backendURL + url, init);
  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.error || "Invalid request.");
  }
  return await response.json() as T;
}