const API_URL = "http://localhost:8080";

export async function getHealth() {
  const response = await fetch(`${API_URL}/health`);

  if (!response.ok) {
    throw new Error("API request failed");
  }

  return response.json();
}