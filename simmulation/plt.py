import matplotlib.pyplot as plt

# File path
file_path = "rfiddata/data12.txt"

timestamps, rssi_values, colors = [], [], []

with open(file_path, "r") as file:
    for line in file:
        if line.startswith("--- {"):
            parts = line.strip().strip("--- ").strip("{}").split()
            if len(parts) >= 4:  
                timestamp = int(parts[2])  # Timestamp 
                rssi = int(parts[3])       # RSSI 
                antenna = int(parts[1])    # Antenna port

                timestamps.append(timestamp)
                rssi_values.append(rssi)

                if antenna == 17:
                    colors.append("blue")
                elif antenna == 9:
                    colors.append("red")

timestamps_sec = [(t - timestamps[0]) / 1e6 for t in timestamps]

plt.figure(figsize=(12, 6))

for i in range(1, len(timestamps_sec)):
    plt.plot(timestamps_sec[i-1:i+1], rssi_values[i-1:i+1], linestyle="-", color=colors[i])

plt.xlabel("Time (Seconds)")
plt.ylabel("RSSI")
plt.title("RSSI Variation Over Time")
plt.grid()

plt.show()
