import { ref } from 'vue'

export function useRecords() {
  const records = ref([])

  function connect() {
    const socket = new WebSocket("ws://localhost:16311/ws");
    socket.addEventListener("open", () => {
      // console.log("websocket connected")
      // socket.send("Hello Server!");
    });

    // Listen for messages
    socket.addEventListener("message", (event) => {
      // console.log("Message from server ", event.data);
      records.value.push(JSON.parse(event.data))
    });

    socket.addEventListener("close", () => {
      console.log("connect failed, reconnecting...")
      setTimeout(connect, 1000)
    });
  }
  connect()
  return { records }
}