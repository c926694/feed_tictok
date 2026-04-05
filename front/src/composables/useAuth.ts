import { computed, ref } from "vue";
import { clearToken, getToken, setToken } from "@/utils/storage";
import type { User } from "@/types/domain";

const token = ref(getToken());
const currentUser = ref<User | null>(null);

export function useAuth() {
  const isLoggedIn = computed(() => Boolean(token.value));

  const setAuthToken = (nextToken: string) => {
    token.value = nextToken;
    setToken(nextToken);
  };

  const clearAuth = () => {
    token.value = "";
    currentUser.value = null;
    clearToken();
  };

  return {
    token,
    currentUser,
    isLoggedIn,
    setAuthToken,
    clearAuth
  };
}
