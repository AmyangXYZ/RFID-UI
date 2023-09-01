import { ref, watch } from 'vue'
import type {Tag,Record} from  './useDef'
export const tags  = ref<Tag[]>([])
export const records = ref<Record[]>([])

// export const tags  = ref([])
// export const records = ref([])