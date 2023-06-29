<script setup lang="ts">
import { useTagHolder } from '../hooks/useTagHolder'

const { tags } = useTagHolder()
</script>

<template>
  <div> 
    <ul>
      <ol v-for="tag in tags">
        {{ tag.epc }}--{{ tag.led }}
      </ol>
    </ul>
  </div>
  <!-- <div>Current Tags: {{ tags }}</div> -->
  
</template>


<!-- UI wating data from sever,  -->
