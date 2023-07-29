

import os 
import subprocess
import time

if __name__ == '__main__':


    time.sleep(3)
    print("inpython printing")
    os.chdir("C:\\Users\\sh311mini_1\\Desktop\\RFID-UI\\frontend")
    command = "npm run dev -- --host"
    os.system(command )

