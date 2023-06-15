<script setup lang="ts">
import { useTagHolder } from '../hooks/useTagHolder'

const { tags } = useTagHolder()
</script>

<template>
  <div>Current Tags: {{ tags }}</div>
</template>
