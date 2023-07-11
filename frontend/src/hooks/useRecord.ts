import { ref, watch } from 'vue'
import { useWebSocket } from '@vueuse/core'
import { useFetch } from '@vueuse/core'//

const tags = ref([])

export function useRecord(): any {
  const records = ref<any>([])
  const alltags = ref<any>([])
  // const { data } = useWebSocket('ws://localhost:16311/api/ui/ws', { autoReconnect: { delay: 2000 } })
  
  const getAllTags = async function () {
    const { data } = await useFetch("http://localhost:16311/api/ui/tag").json()
    tags.value = []
    for (const i in data.value.data) {
      tags.value.push(data.value.data[i])
    }

  }

  const registerTag = function (id: string) {
    useFetch("http://localhost:16311/api/ui/tag/"+id).post().text()//UI to server
    getAllTags() //click button, all tags from server
  }
  
  const deleteTag = function(id24:string){
    console.log("in deleteTag", id24)
    useFetch("http://localhost:16311/api/ui/tagdelete/"+id24).post().text()
    getAllTags()
  }
  
  
  const { data } = useWebSocket('ws://localhost:16311/api/ui/ws')
  watch(data, () => {//trigger from channel
    if(data.value.indexOf("gait_speed")>-1){
    console.log("in useRecord, data.value",data.value,data.value.length)

    records.value.push(JSON.parse(data.value))
    }
    else{
      // const dataobject =  JSON.parse(data.value)
      // console.log("!!!!",typeof dataobject,JSON.parse(data.value))
      getAllTags()
      console.log("#",tags)
      // for (const key in JSON.parse(data.value)){
      //   console.log("#",JSON.parse(data.value)[key])
      //   alltags.value.push(JSON.parse(data.value)[key])
      // }
      }


  }
  
  )

  return { records,tags,registerTag,getAllTags,deleteTag}
}
