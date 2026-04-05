import axios, { type AxiosError } from "axios";
import { getToken } from "@/utils/storage";
import { useToast } from "@/composables/useToast";

const { showToast } = useToast();

export const http = axios.create({
  baseURL: "/api",
  timeout: 12000
});

http.interceptors.request.use((config) => {
  const token = getToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

http.interceptors.response.use(
  (response) => response,
  (error: AxiosError<{ msg?: string; message?: string }>) => {
    const status = error.response?.status;
    const msg = error.response?.data?.msg ?? error.response?.data?.message ?? error.message;

    if (status === 401) {
      showToast("登录已失效，请重新登录");
      if (window.location.pathname !== "/login") {
        window.location.href = "/login";
      }
    } else if (msg) {
      showToast(msg);
    }
    return Promise.reject(error);
  }
);
