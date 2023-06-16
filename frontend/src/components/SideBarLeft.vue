<script setup lang="ts">
import { ref } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { useTagHolder } from '../hooks/useTagHolder'

const { registerTag } = useTagHolder()

const selectableTagIDs = ref([
  {
    id: '00000945'
  },
  {
    id: '18145536'
  }
])
const selectedTagID = ref('')
const input = ref('')   //here

const handleOpen = (key: string, keyPath: string[]) => {
  console.log(key, keyPath)
}
const handleClose = (key: string, keyPath: string[]) => {
  console.log(key, keyPath)
}
</script>



<template>
  <el-menu default-active="0" class="side-bar" @open="handleOpen" @close="handleClose">
    <el-row align="middle" justify="space-around"> 
      <el-col :span="16">
        <el-select v-model="selectedTagID" class="m-2" placeholder="Select">
          <el-option
            v-for="item in selectableTagIDs"
            :key="item.id"
            :label="item.id"
            :value="item.id"
          />
        </el-select>
        
      </el-col>
      <el-col :span="6" style="vertical-align: center">
        <el-button size="small" :icon="Plus" @click="registerTag(selectedTagID)"> </el-button>
      </el-col> 
    </el-row>

    <el-row align="middle" justify="space-around"> 
      <el-col :span="16">
        <el-input>
          <el-input v-model="input" placeholder="Please input" />
        </el-input>
      </el-col>
      <el-col :span="6" style="vertical-align: center">
        <el-button size="small" :icon="Plus" @click="registerTag(input)"> </el-button>
      </el-col>
    </el-row>

  </el-menu>


</template>

<style scoped>
.side-bar {
  height: 15vh;
}
</style>
