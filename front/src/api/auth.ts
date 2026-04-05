import { http } from "@/utils/http";
import { normalizeUser, unwrapData } from "@/utils/normalize";
import type { ApiEnvelope, RawUser } from "@/types/backend";
import type { User } from "@/types/domain";

interface LoginPayload {
  username: string;
  password: string;
}

export async function login(payload: LoginPayload) {
  const { data } = await http.post("/users/login", payload);
  const envelope = data as ApiEnvelope<unknown>;
  const rawData = envelope?.data;
  let token = "";
  if (typeof rawData === "string") {
    token = rawData;
  } else if (rawData && typeof rawData === "object") {
    const body = rawData as Record<string, unknown>;
    token = String(body.token ?? body.access_token ?? body.jwt ?? "");
  }
  return {
    token,
    raw: data
  };
}

export async function fetchMe() {
  const { data } = await http.get("/users/me");
  const body = unwrapData<RawUser | { user?: RawUser }>(data);
  if ("user" in (body as { user?: RawUser })) {
    return normalizeUser((body as { user?: RawUser }).user);
  }
  return normalizeUser(body as RawUser);
}

export async function updateMyProfile(payload: { nickname?: string; avatar?: File | null }) {
  const formData = new FormData();
  if (payload.nickname && payload.nickname.trim()) {
    formData.append("nickname", payload.nickname.trim());
  }
  if (payload.avatar) {
    formData.append("avatar", payload.avatar);
  }
  const { data } = await http.post("/users/me", formData);
  const body = unwrapData<RawUser | { user?: RawUser }>(data);
  if ("user" in (body as { user?: RawUser })) {
    return normalizeUser((body as { user?: RawUser }).user);
  }
  return normalizeUser(body as RawUser);
}

export async function logout() {
  await http.delete("/users/logout");
}

export async function registerUser(payload: { username: string; password: string; re_password: string }) {
  const { data } = await http.post("/users/register", payload);
  const body = unwrapData<unknown>(data);
  return body;
}

export function pickUserFromRaw(raw: Record<string, unknown>): User | null {
  const maybeUser = raw.user as RawUser | undefined;
  if (maybeUser) return normalizeUser(maybeUser);
  return null;
}
