import { useFetch } from '@vueuse/core'
import { tags } from './useStates'



export function useTagHolder() {
  const distdata = {
    data: "4.0",
  };
  //function get tag data from server
  const getAllTags = async function () {
    const { data } = await useFetch('/api/ui/tag').json()
    // const { data } = await useFetch("http://192.168.0.111:16311/api/ui/tag").json()
    console.log('getAllTags function useTagholder', typeof data, data)
    tags.value = []
    for (const i in data.value.data) {
      tags.value.push(data.value.data[i]) 
    }
  }

  const registerTag = function (id: string) {
    console.log('in registerTag', id)

    useFetch('/api/ui/tag/' + id)
    // useFetch("http://192.168.0.111:16311/api/ui/tag/"+id).post().text()//UI to server
      .post()
      .text()
    setTimeout(getAllTags, 250)
  }

  const deleteTag = function (id24: string) {
    console.log('in deleteTag', id24)
    useFetch('/api/ui/tagdelete/' + id24)
    // useFetch("http://192.168.0.111:16311/api/ui/tagdelete/"+id24).post().text()
      .post()
      .text()
    setTimeout(getAllTags, 250)
  }

  const enterDistance = function(UIdist: string){
    console.log('in enterDistance', UIdist)
    distdata.data = UIdist
    // useFetch("http://192.168.0.111:16311/api/ui/distance/"+UIdist).post().text()//UI to server
    useFetch('/api/ui/distance/' + UIdist) 
    .post()
      .text()
    setTimeout(getAllTags, 250)
  }


  return { getAllTags, registerTag, deleteTag, enterDistance, distdata }
}
