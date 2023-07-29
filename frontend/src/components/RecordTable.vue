<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { useRecord } from '../hooks/useRecord'
import { Refresh } from '@element-plus/icons-vue'
import {records  } from '../hooks/useStates'

// const { records} = useRecord()

// import { records } from '../hooks/useStates'
const {getData} = useRecord()
getData()
//for register input box
const boxinput = ref('')

const handleOpen = (key: string, keyPath: string[]) => {
  console.log(key, keyPath)
}
const handleClose = (key: string, keyPath: string[]) => {
  console.log(key, keyPath)
}

//for record 
const columns = [
  {
    key: 'epc',
    dataKey: 'epc',
    title: 'EPC',
    width: 250,
    align: 'center'
  },

  {
    key: 'gaitSpeed',
    dataKey: 'gait_speed',
    title: 'Gait Speed (m/s)',
    width: 150,
    align: 'left'
  },
  {
    key: 'time',
    dataKey: 'time',
    title: 'Time',
    width: 150,
    align: 'left'
  }
]

const tableRef = ref()
watch(
  records,
  () => {
    if (records.value.length > 0) {
      nextTick(() => {
        tableRef.value?.scrollToRow(records.value.length)
      })
    }
  },

  { deep: true }
)
</script>






<template>
  <el-card style="width: 100%">
    <template #header>
      <div class="card-header">
        <span>Records</span>
        <el-button @click="records = []" :icon="Refresh">Clear</el-button>
      </div>
    </template>
    <el-auto-resizer>
      <template #default="{width }">
        <el-table-v2
          ref="tableRef"
          :columns="columns"
          :data="records"
          :width="width"
          :height="600"
          fixed
        />
      </template>
    </el-auto-resizer>
  </el-card>


</template>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.el-table-v2{
  font-size: large;
}

</style>
