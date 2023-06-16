import { ref, watch } from 'vue'
import { useWebSocket } from '@vueuse/core'

export function useRecord(): any {
  const records = ref<any>([])
  // const { data } = useWebSocket('ws://localhost:16311/api/ui/ws', { autoReconnect: { delay: 2000 } })
  const { data } = useWebSocket('ws://localhost:16311/api/ui/ws')

  watch(data, () => {
    console.log("in useRecord, data print", data)
    records.value.push(JSON.parse(data.value))
  })

  return { records }
}
