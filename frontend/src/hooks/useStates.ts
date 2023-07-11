// import { ref, watch } from 'vue'
// import { useWebSocket } from '@vueuse/core'
// import { useFetch } from '@vueuse/core'


// const tags = ref([])

// export function useStates(){
//     const alltags = ref<any>([])
//     // const { data } = useWebSocket('ws://localhost:16311/api/ui/ws', { autoReconnect: { delay: 2000 } })
//     const { data } = useWebSocket('ws://localhost:16311/api/ui/ws')

    

//     watch(data, () => {
//       if(data.value.indexOf("gait_speed")<=-1){
      
//     //old method
//         // // const dataobject =  JSON.parse(data.value)
//         // // console.log("!!!!",typeof dataobject,JSON.parse(data.value))
//         // for (const key in JSON.parse(data.value)){
//         //   console.log("#",JSON.parse(data.value)[key])
//         //   alltags.value.push(JSON.parse(data.value)[key])
//         // }

//         getAllTags()
        
        
//         }
  
  
//     }
    
//     )
//     const getAllTags = async function () {
//         const { datatags } = await useFetch("http://localhost:16311/api/ui/tag").json()
//         tags.value = []
//         console.log("debug",datatags)
//         // for (const i in datatags.value.data) {
//         //   tags.value.push(datatags.value.data[i])
//         // }
    
//       }
  
//     return { tags,getAllTags}
//   }