<script setup lang="ts">
import PageHeader from './components/PageHeader.vue'
import RecordTable from './components/RecordTable.vue'
import TagHolder from './components/TagHolder.vue'
import SideBarLeft from './components/SideBarLeft.vue'
</script>

<template>
  <el-container>
    <el-header><PageHeader /></el-header>
    <el-main>
      <el-row>
        <el-col :span="4">
          <SideBarLeft />
        </el-col>
        <el-col :offset="2" :span="12">
          <RecordTable />
          <TagHolder style="margin-top: 30px;"/>
        </el-col>
      </el-row>
    </el-main>
  
  </el-container>
</template>

<style scoped></style>
