<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { useRecord } from '../hooks/useRecord'
import { Refresh } from '@element-plus/icons-vue'

const { records } = useRecord()

const columns = [
  {
    key: 'epc',
    dataKey: 'epc',
    title: 'EPC',
    width: 100,
    align: 'center'
  },
  // {
  //   key: 'antennaPort',
  //   dataKey: 'antennaPort',
  //   title: 'Antenna Port',
  //   width: 200,
  //   align: 'center'
  // },
  // {
  //   key: 'timestamp',
  //   dataKey: 'firstSeenTimeStamp',
  //   title: 'Timestamp',
  //   width: 200,
  //   align: 'center'
  // },
  {
    key: 'gaitSpeed',
    dataKey: 'gait_speed',
    title: 'Gait Speed (m/s)',
    width: 150,
    align: 'center'
  },
  {
    key: 'time',
    dataKey: 'time',
    title: 'Time',
    width: 100,
    align: 'center'
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
        <el-button @click="records = []" :icon="Refresh" size="small">Clear</el-button>
      </div>
    </template>
    <el-auto-resizer>
      <template #default="{ width }">
        <el-table-v2
          ref="tableRef"
          :columns="columns"
          :data="records"
          :width="width"
          :height="500"
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
</style>
