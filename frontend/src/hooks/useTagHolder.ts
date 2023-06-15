import { ref } from 'vue'

export function useTagHolder() {
  const tags = ref([{ id: 8951231, name: 'xx' }])
  return { tags }
}
