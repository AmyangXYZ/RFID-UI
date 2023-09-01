import { useFetch } from '@vueuse/core'
import { tags } from './useStates'

export function useTagHolder() {
  //function get tag data from server
  const getAllTags = async function () {
    const { data } = await useFetch('/api/ui/tag').json()
    console.log('getAllTags function useTagholder', typeof data, data)
    tags.value = []
    for (const i in data.value.data) {
      tags.value.push(data.value.data[i])
    }
  }

  const registerTag = function (id: string) {
    console.log('in registerTag', id)

    useFetch('/api/ui/tag/' + id)
      .post()
      .text()
    setTimeout(getAllTags, 250)
  }

  const deleteTag = function (id24: string) {
    console.log('in deleteTag', id24)
    useFetch('/api/ui/tagdelete/' + id24)
      .post()
      .text()
    setTimeout(getAllTags, 250)
  }
  return { getAllTags, registerTag, deleteTag }
}
