<script setup>
import { useRecords } from '../hooks/useRecords.js'
const { records } = useRecords()
</script>

<template>
  <vs-card>
    <template #header>
      <h3>RFID Records</h3>
    </template>
    <vs-table :key="records.length" :data="records" stripe>
      <template #thead>
        <vs-th>
          <p>EPC</p>
        </vs-th>
        <vs-th>
          <p>Speed</p>
        </vs-th>
        <vs-th>
          <p>Timestamp</p>
        </vs-th>
      </template>

      <template v-slot="{ data }">
        <vs-tr :data="tr" :key="i" v-for="(tr, i) in data">
          <vs-td>
            <p>{{ tr.epc }}</p>
          </vs-td>
          <vs-td>
            <p>{{ tr.speed }}</p>
          </vs-td>
          <vs-td>
            <p>{{ tr.timestamp }}</p>
          </vs-td>
        </vs-tr>
      </template>
    </vs-table>
  </vs-card>
</template>

<style scoped></style>
