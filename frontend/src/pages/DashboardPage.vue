<script setup>
import Records from '../components/RFIDRecords.vue'
import SideBarLeft from '../components/SideBarLeft.vue'
import SideBarRight from '../components/SideBarRight.vue'
</script>

<template>
  <vs-row vs-w="12" vs-type="flex" vs-justify="space-between">
    <vs-col vs-w="0.8">
      <SideBarLeft />
    </vs-col>
    <vs-col vs-w="10">
      <Records />
    </vs-col>
    <vs-col vs-w="0.8">
      <SideBarRight />
    </vs-col>
  </vs-row>
</template>
