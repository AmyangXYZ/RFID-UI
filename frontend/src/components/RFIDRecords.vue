<script setup>
import { useRecords } from '../hooks/useRecords.js'
const { records } = useRecords()
</script>

<template>
  <vs-card>
    <template #header>
      <vs-row>
        <vs-col vs-w="4">
          <h3>RFID Records</h3>
        </vs-col>
        <vs-col vs-offset="7.5" vs-w=".5">
          <vs-button @click="records = []" color="primary" type="filled">Clear</vs-button>
        </vs-col>
      </vs-row>
      
      
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
