import matplotlib.pyplot as plt 
# Reload and process the file without using pandas

# File path
file_path = "rfiddata/data4.txt"

# Initialize lists to store extracted data
timestamps = []
rssi_values = []

# Read the file and extract relevant data
with open(file_path, "r") as file:
    for line in file:
        # Process only lines that start with "--- {"
        if line.startswith("--- {"):
            parts = line.strip().strip("--- ").strip("{}").split()
            if len(parts) >= 4:  # Ensure the line has at least 4 elements
                timestamps.append(int(parts[2]))  # Timestamp (third item)
                rssi_values.append(int(parts[3]))  # RSSI (fourth item)

# Convert timestamps to seconds (assuming timestamps are in microseconds)
timestamps_sec = [(t - timestamps[0]) / 1e6 for t in timestamps]  # Normalize to start from 0 sec

# Plot the extracted RSSI values over time in seconds
plt.figure(figsize=(12, 6))
plt.plot(timestamps_sec, rssi_values, marker="o", linestyle="-", label="RSSI Values")
plt.xlabel("Time (Seconds)")
plt.ylabel("RSSI")
plt.title("RSSI Variation Over Time")
plt.legend()
plt.grid()

# Display the plot
plt.show()
