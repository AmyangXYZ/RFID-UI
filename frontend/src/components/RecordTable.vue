<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { useRecord } from '../hooks/useRecord'
import { Refresh } from '@element-plus/icons-vue'
// const { tags } = useTagHolder()


const { records,tags,deleteTag } = useRecord()

const columns = [
  {
    key: 'epc',
    dataKey: 'epc',
    title: 'EPC',
    width: 200,
    align: 'center'
  },

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


  <div> 
    <ul>
      
      <ol v-for="(tag,index) in tags" v-bind:class="tag.led ">

        {{tag.epc}} 
        <p v-if="tag.led ==='GREY'">ACTIVE</p>
        <p v-if="tag.led ==='GREEN'">(ACTIVE) Passing First Antenna</p>
        <p v-if="tag.led ==='RED'">(HOLD) Passing Second Antenna </p>
        <button @click="deleteTag(tags[index].epc24)">Delet</button>
      </ol>
    </ul>
  </div>
</template>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
ol{
            list-style-type: none;
            background:lightgray;
            margin: 10px auto;
            padding: 10px 20px;
            border-radius: 10px;
            display: flex;
            align-items: center;
            justify-content: space-between;

        }

ol.GREY{
    background: lightgray;
}
ol.RED{
  background: orange ;
}
ol.GREEN{
  background: greenyellow ;
}

</style>
