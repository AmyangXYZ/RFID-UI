import {ref} from 'vue'

export function useRecords() {
  const socket = new WebSocket("ws://localhost:16311/ws");
  const records = ref([])
  // Connection opened
  socket.addEventListener("open", (event) => {
    // socket.send("Hello Server!");
  });

  // Listen for messages
  socket.addEventListener("message", (event) => {
    // console.log("Message from server ", event.data);
    records.value.push(JSON.parse(event.data))
  });

  return { records }
}