<script setup lang="ts">
import { useTagHolder } from '../hooks/useTagHolder'

const { tags } = useTagHolder()
</script>

<template>
  <div>Current Tags: {{ tags }}</div>
  <div>Tags Status: </div>
  
</template>


<!-- UI wating data from sever,  -->
