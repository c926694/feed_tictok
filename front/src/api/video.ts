import { http } from "@/utils/http";
import { normalizeVideo, unwrapData } from "@/utils/normalize";
import type { RawVideo } from "@/types/backend";
import type { Video } from "@/types/domain";

interface FeedParams {
  limit?: number;
  lastScore?: string;
}

function pickVideoList(source: unknown): RawVideo[] {
  if (Array.isArray(source)) return source as RawVideo[];
  if (!source || typeof source !== "object") return [];
  const payload = source as Record<string, unknown>;
  const candidates = [payload.list, payload.items, payload.videos, payload.video_list, payload.feed_video_list];
  const found = candidates.find((entry) => Array.isArray(entry));
  return (found as RawVideo[] | undefined) ?? [];
}

export async function fetchFeedVideos(params: FeedParams = {}) {
  return fetchFeedByPath("/videos/feed", params);
}

export async function fetchHotVideos(limit = 5, lastScore?: string) {
  return fetchFeedByPath("/videos/feed/hot", { limit, lastScore });
}

export async function fetchFollowVideos(params: FeedParams = {}) {
  return fetchFeedByPath("/videos/feed/follow", params);
}

export async function fetchMyVideos(limit = 60) {
  const { data } = await http.get("/videos/me", {
    params: { limit }
  });
  const body = unwrapData<unknown>(data);
  return pickVideoList(body).map(normalizeVideo);
}

async function fetchFeedByPath(path: string, params: FeedParams = {}) {
  const { limit = 5, lastScore } = params;
  const { data } = await http.get(path, {
    params: {
      limit,
      _ts: Date.now(),
      ...(lastScore ? { last_score: lastScore } : {})
    }
  });
  const body = unwrapData<unknown>(data);
  const list = pickVideoList(body);
  return {
    videos: list.map(normalizeVideo),
    nextScore:
      typeof body === "object" && body
        ? String((body as Record<string, unknown>).last_score ?? (body as Record<string, unknown>).next_score ?? "")
        : ""
  };
}

export async function createVideo(payload: { title: string; description: string; cover: File; play: File }) {
  const formData = new FormData();
  formData.append("title", payload.title);
  formData.append("description", payload.description);
  formData.append("cover", payload.cover);
  formData.append("play", payload.play);
  await http.post("/videos/create", formData);
}

export async function deleteVideo(videoId: number) {
  await http.delete(`/videos/${videoId}`);
}

export function filterVideosByUser(videos: Video[], userId: number) {
  return videos.filter((video) => video.author.id === userId);
}
