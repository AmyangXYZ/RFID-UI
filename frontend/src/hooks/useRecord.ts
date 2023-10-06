import { watch } from 'vue'
import { useWebSocket } from '@vueuse/core'
import { useFetch } from '@vueuse/core'
import { tags, records } from './useStates'

export function useRecord(): any {
  const getAllTags = async function () {
    const { data } = await useFetch('/api/ui/tag').json()
    // const { data } = await useFetch("http://192.168.0.100:16311/api/ui/tag").json()

    tags.value = []
    for (const i in data.value.data) {
      tags.value.push(data.value.data[i])
    }
  }
  const getData = function () {
    const { data } = useWebSocket('ws://' + window.location.host + '/api/ui/ws')
    // const { data } = useWebSocket('ws://192.168.0.100:16311/api/ui/ws')


    watch(data, () => {
      if (data.value.indexOf('gait_speed') > -1) {
    records.value.push(JSON.parse(data.value))
      } else {
      getAllTags()
      }
    })
  }
  return { records, tags, getAllTags, getData }
  }
  