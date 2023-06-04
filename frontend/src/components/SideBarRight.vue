<script setup>
import { ref } from 'vue'

const active = ref(false)
</script>

<template>
  <div ref="parentSidebar">
    <vs-button @click="active = !active" color="primary" type="filled">Open Sidebar</vs-button>
    <vs-sidebar
      default-index="1"
      position-right
      color="primary"
      spacer
      v-model="active"
    >
      <vs-sidebar-group title="Item 1">
        <vs-sidebar-item index="1" icon="question_answer"> Value 1 </vs-sidebar-item>
      </vs-sidebar-group>
      <vs-sidebar-group title="Item 2">
        <vs-sidebar-item index="2" icon="question_answer"> Value 2 </vs-sidebar-item>
      </vs-sidebar-group>
      <vs-sidebar-group title="Item 3">
        <vs-sidebar-item index="3" icon="question_answer"> Value 3 </vs-sidebar-item>
      </vs-sidebar-group>
      <vs-sidebar-group title="Item 4">
        <vs-sidebar-item index="4" icon="question_answer"> Value 4 </vs-sidebar-item>
      </vs-sidebar-group>
    </vs-sidebar>
  </div>
</template>

<style scoped>
.vs-sidebar-group h4 {
  padding: 15px;
  font-size: 4rem;
}

.vs-sidebar-group.vs-sidebar-group-open h4 {
  font-size: 1.5rem;
}
</style>
