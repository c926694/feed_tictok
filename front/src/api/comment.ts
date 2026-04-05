import { http } from "@/utils/http";
import { normalizeComment, unwrapData } from "@/utils/normalize";
import type { RawComment } from "@/types/backend";

function pickCommentList(source: unknown): RawComment[] {
  if (Array.isArray(source)) return source as RawComment[];
  if (!source || typeof source !== "object") return [];
  const payload = source as Record<string, unknown>;
  const candidates = [payload.list, payload.items, payload.comments, payload.comment_list];
  const found = candidates.find((entry) => Array.isArray(entry));
  return (found as RawComment[] | undefined) ?? [];
}

export async function fetchCommentList(videoId: number) {
  const { data } = await http.get(`/comments/list/${videoId}`);
  const body = unwrapData<unknown>(data);
  return pickCommentList(body).map(normalizeComment);
}

export async function createComment(payload: { video_id: number; content: string }) {
  await http.post("/comments", payload);
}

export async function deleteComment(commentId: number) {
  await http.delete(`/comments/${commentId}`);
}
