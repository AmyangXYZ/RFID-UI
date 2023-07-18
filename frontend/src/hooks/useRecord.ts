import { ref, watch } from 'vue'
import { useWebSocket } from '@vueuse/core'
import { useFetch } from '@vueuse/core'


const tags = ref([])
const records = ref<any>([])


export function useRecord(): any {
  // const { data } = useWebSocket('ws://localhost:16311/api/ui/ws', { autoReconnect: { delay: 2000 } })
  
  const getAllTags = async function () {
    const { data } = await useFetch("http://192.168.1.49:16311/api/ui/tag").json()
    console.log("getAllTags function useRecord",data)
    tags.value = []
    for (const i in data.value.data) {
      tags.value.push(data.value.data[i])
    }

  }

  
  
  
  const { data } = useWebSocket('ws://192.168.1.49:16311/api/ui/ws')

  watch(data, () => {//trigger from channel
    if(data.value.indexOf("gait_speed")>-1){
    console.log("in useRecord, data.value",data.value,data.value.length)

    records.value.push(JSON.parse(data.value))
    console.log("in useRecord, data.value, why not show, show record",records)

    }
    else{
      // const dataobject =  JSON.parse(data.value)
      // console.log("!!!!",typeof dataobject,JSON.parse(data.value))
      getAllTags()
      // console.log("#",tags)
      // for (const key in JSON.parse(data.value)){
      //   console.log("#",JSON.parse(data.value)[key])
      //   alltags.value.push(JSON.parse(data.value)[key])
      // }
      }


  }
  
  )

  return { records,tags,getAllTags}
}
