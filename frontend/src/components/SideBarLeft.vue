<script setup lang="ts">
import { ref } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { useTagHolder } from '../hooks/useTagHolder'
import { tags } from '../hooks/useStates'
const { registerTag, deleteTag } = useTagHolder()
const boxinput = ref('')

const handleOpen = (key: string, keyPath: string[]) => {
  console.log(key, keyPath)
}
const handleClose = (key: string, keyPath: string[]) => {
  console.log(key, keyPath)
}
</script>
<template>
  <el-menu default-active="0" class="side-bar" @open="handleOpen" @close="handleClose">
    <!-- <el-row align="middle" justify="space-around">
      <el-col :span="16">
        <el-input type="text" v-model="boxinput" placeholder="enter EPC here" />
      </el-col>
      <el-col :span="7" style="vertical-align: center">
        <el-button size="small" :icon="Plus" @click="registerTag(boxinput)"> </el-button>
      </el-col>
    </el-row> -->
    <br />
    <el-row>
      <div>
        <p>All registed tag status show below:</p>

        <ol v-for="(tag, index) in tags" v-bind:class="tag.led">
          {{
            tag.epc
          }}
          <p v-if="tag.led === 'GREY'">ACTIVE</p>
          <p v-if="tag.led === 'GREEN'">RECORDING</p>
          <p v-if="tag.led === 'RED'">HOLD</p> 
          <!-- <button @click="deleteTag(tags[index].epc24)" class="Button">Delete</button> -->
        </ol>
      </div>
    </el-row>
  </el-menu>
</template>

<style scoped>
.side-bar {
  height: 98vh;
}

ol {
  list-style-type: none;
  background: lightgray;
  margin: 10px auto;
  padding: 15px 15px 15px 15px;
  border-radius: 10px;
  display: flex;
  /* width: 410px; */
  /* width: 180%; */
  width:120%;
  font-size: medium;
  align-items: center;
  justify-content: space-between;
}

Button {
  width: 80px;
  height: 30px;
  font-size: medium;
  
}
ol.GREY {
  background: lightgray;
}
ol.RED {
  background: orange;
}
ol.GREEN {
  background: greenyellow;
}
.el-input {
  font-size: large;
}
</style>
