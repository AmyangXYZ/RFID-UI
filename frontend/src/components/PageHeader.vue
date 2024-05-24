<script setup lang="ts">
import { ref } from 'vue'
import Logo from '@/assets/uconn-health-logo.png'
import { useTagHolder } from '../hooks/useTagHolder'


// const { enterDistance, distdata } = useTagHolder()
const boxinput = ref('')
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
      <!-- <ol v-for="(tag, index) in tags" v-bind:class="tag.led">
        <div v-if = "index == tags.length - 1">
          <p>Current Distance:{{ tag.dist }} </p>>
        </div>
        
      </ol> -->
      <!-- <p>Current Distance: {{distdata.data}} meters </p> -->
      <p>Current Distance: 4.0 meters </p>
      <!-- <p>Distance:   
        <input type="text" v-model="boxinput"   />
        meters
        <button @click="enterDistance(boxinput)" >Update</button>
      </p> -->
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
input{
  width:60px;
  padding:4px
}
button{
  padding:4px;
  margin:5px;
  background:rgb(163, 185, 202);
}

</style> 
