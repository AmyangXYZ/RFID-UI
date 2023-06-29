import { ref, watch } from 'vue'
import { useWebSocket } from '@vueuse/core'

export function useRecord(): any {
  const records = ref<any>([])
  // const { data } = useWebSocket('ws://localhost:16311/api/ui/ws', { autoReconnect: { delay: 2000 } })
  const { data } = useWebSocket('ws://localhost:16311/api/ui/ws')
  watch(data, () => {
    console.log("in useRecord, data.value",data.value,data.value.length)
    //how to identify data from dirrerent cannel send from server?

    records.value.push(JSON.parse(data.value))
  })

  return { records }
}
