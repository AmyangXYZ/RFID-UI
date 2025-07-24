import requests
import json
import re

def extract_rfid_data(filepath):
    rfid_data_list = []
    pattern = re.compile(r'--- (\d+)')
    data_pattern = re.compile(r'--- \{(E[0-9A-F]+)\s+(\d+)\s+(\d+)\s+(-?\d+)\}')
    
    with open(filepath, 'r') as file:
        lines = [line.strip() for line in file if line.startswith('---')]
    
    i = 0
    while i < len(lines):
        match = pattern.match(lines[i])
        if match:
            count = int(match.group(1))
            temp_list = []
            i += 1
            for _ in range(count):
                if i < len(lines):
                    data_match = data_pattern.match(lines[i])
                    if data_match:
                        epc, antenna_port, timestamp, peak_rssi = data_match.groups()
                        temp_list.append({
                            "epc": epc,
                            "antennaPort": int(antenna_port),
                            "firstSeenTimeStamp": int(timestamp),
                            "peakRssi": int(peak_rssi)
                        })
                    i += 1
            if temp_list:
                rfid_data_list.extend(temp_list)
        else:
            i += 1
    return rfid_data_list

def send_to_server(rfid_data_list, server_url):
    payload = {"tag_reads": rfid_data_list}
    headers = {'Content-Type': 'application/json'}
    response = requests.post(server_url, data=json.dumps(payload), headers=headers)
    return response.status_code, response.text

if __name__ == "__main__":
    # file_path = "rfiddata/geriatrics/02042025/data14.txt"
    file_path = "rfiddata/southington/07212025/data1.txt"  # Update with the correct path if needed

    server_endpoint = "http://localhost:16311/api/reader/connect"  # Adjust if the server runs on a different host
    
    rfid_data = extract_rfid_data(file_path)
    print(len(rfid_data))
    if rfid_data:
        status, response_text = send_to_server(rfid_data, server_endpoint)
        print(f"Server Response ({status}): {response_text}")
    else:
        print("No RFID data extracted from the file.")
