import { ref } from "vue";

const visible = ref(false);
const message = ref("");
let timer: number | undefined;

export function useToast() {
  const showToast = (content: string, duration = 1800) => {
    message.value = content;
    visible.value = true;
    window.clearTimeout(timer);
    timer = window.setTimeout(() => {
      visible.value = false;
    }, duration);
  };

  return {
    visible,
    message,
    showToast
  };
}
