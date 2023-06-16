import { ref } from 'vue'
import { useFetch } from '@vueuse/core'

// global variable
const tags = ref([])

export function useTagHolder() {//function get tag data from server
  const getAllTags = async function () {
    const { data } = await useFetch("http://localhost:16311/api/ui/tag").json()
    tags.value = []
    for (const i in data.value.data) {
      tags.value.push(data.value.data[i]) //push to UI
    }

  }
  const registerTag = function (id: string) {//UI to server
    useFetch("http://localhost:16311/api/ui/tag/"+id).post().text()
    getAllTags()
  }
  return { tags, getAllTags, registerTag }
}
