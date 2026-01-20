const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000";

export async function createSession(): Promise<{ session_id: string }> {
  const res = await fetch(`${API_BASE}/session`, { method: "POST" });
  if (!res.ok) throw new Error("Failed to create session");
  return res.json();
}

export async function getRecommendation(sessionId: string): Promise<{
  name: string;
  description?: string;
}> {
  const res = await fetch(`${API_BASE}/recommendation?session_id=${sessionId}`);
  if (!res.ok) throw new Error("Failed to get recommendation");
  return res.json();
}

export async function sendSwipe(
  sessionId: string,
  foodName: string,
  action: "left" | "right" | "super"
): Promise<void> {
  const res = await fetch(`${API_BASE}/swipe`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      session_id: sessionId,
      food_name: foodName,
      action,
    }),
  });
  if (!res.ok) throw new Error("Failed to send swipe");
}

export async function getFoodInfo(foodName: string): Promise<{
  description: string;
  nutrients: string;
  similar_foods: string[];
}> {
  const res = await fetch(`${API_BASE}/food-info?food_name=${encodeURIComponent(foodName)}`);
  if (!res.ok) throw new Error("Failed to get food info");
  return res.json();
}
