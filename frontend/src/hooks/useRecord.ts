import { ref, watch } from 'vue'
import { useWebSocket } from '@vueuse/core'

export function useRecord(): any {
  const records = ref<any>([])
  const { data } = useWebSocket('ws://localhost:16311/ws', { autoReconnect: { delay: 2000 } })

  watch(data, () => {
    records.value.push(JSON.parse(data.value))
  })

  return { records }
}
