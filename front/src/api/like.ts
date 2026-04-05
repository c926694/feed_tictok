import { http } from "@/utils/http";
import { unwrapData } from "@/utils/normalize";

function pickIsLiked(payload: unknown) {
  const body = unwrapData<unknown>(payload);
  if (!body || typeof body !== "object") return false;
  const data = body as Record<string, unknown>;
  return Boolean(data.is_liked ?? data.isLiked);
}

export async function switchVideoLike(videoId: number) {
  const { data } = await http.post(`/likes/video/switchLike/${videoId}`);
  return pickIsLiked(data);
}

export async function switchCommentLike(commentId: number) {
  const { data } = await http.post(`/likes/comment/switchLike/${commentId}`);
  return pickIsLiked(data);
}
