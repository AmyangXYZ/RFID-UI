<script setup lang="ts">
import { ref } from 'vue'
import Logo from '@/assets/uconn-health-logo.png'

const activeIndex = ref('1')
const handleSelect = (key: string, keyPath: string[]) => {
  console.log(key, keyPath)
}
</script>

<template>
  <el-menu
    :default-active="activeIndex"
    class="el-menu-demo"
    mode="horizontal"
    :ellipsis="false"
    @select="handleSelect"
  >
    <el-menu-item index="0"><img height="32" :src="Logo" /></el-menu-item>
    <div class="flex-grow" />
    <el-sub-menu index="2">
      <template #title>Parameters</template>
      <p>Distance: 4 meters</p>
      <p>Reader Name: SpeedwayR-14-D3-31</p>
      <p>Rx Sensitivity: -50</p>
      <p>Tx Power: 15</p>
    </el-sub-menu>
  </el-menu>
</template>

<style scoped>
.flex-grow {
  flex-grow: 1;
}
p {
  padding: 10px;
  font-size: 18px;
}
</style>
