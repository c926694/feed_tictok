import { http } from "@/utils/http";
import { unwrapData } from "@/utils/normalize";

export async function switchFollow(userId: number) {
  const { data } = await http.post(`/follows/switchFollow/${userId}`);
  const body = unwrapData<unknown>(data);
  if (!body || typeof body !== "object") return false;
  const payload = body as Record<string, unknown>;
  return Boolean(payload.is_follow ?? payload.isFollow);
}
