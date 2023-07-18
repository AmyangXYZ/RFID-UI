import { ref } from 'vue'
import { useFetch } from '@vueuse/core'

// global variable
//original tage holder, use to show tags from go as a list
const tags = ref([])

export function useTagHolder() {//function get tag data from server
  const getAllTags = async function () {
    const { data } = await useFetch("http://192.168.1.49:16311/api/ui/tag").json()
    console.log("getAllTags function useTagholder",typeof data)
    tags.value = []
    for (const i in data.value.data) {
      tags.value.push(data.value.data[i])
    }

  }
  const registerTag = function (id: string) {
    useFetch("http://192.168.1.49:16311/api/ui/tag/"+id).post().text()//UI to server
    getAllTags() //click button, all tags from server
    
  }

  const deleteTag = function(id24:string){
    console.log("in deleteTag", id24)
    useFetch("http://192.168.1.49:16311/api/ui/tagdelete/"+id24).post().text()
    getAllTags()
    
  }
  return { tags, getAllTags, registerTag,deleteTag }
}
