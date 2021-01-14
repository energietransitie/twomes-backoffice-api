import mariadb
import sys
from datetime import datetime

#database information
USER = '***'
PASSWORD = '***'
HOST = '***'
PORT = 0000
DATABASE = '***'

#this function creates a database connection based on the database information
def create_connection():
    try:
        conn = mariadb.connect(
            user= USER,
            password= PASSWORD,
            host= HOST,
            port= PORT, 
            database= DATABASE)
        print("Success connecting to MariaDB Platform:")
        return conn
    except mariadb.Error as e:
        print(f"Error connecting to MariaDB Platform: {e}")
        sys.exit(1)
        

#this function inserts device measurements in the database based on data of the parameter (s)
def insert_device_measurements(device_mac_address, prop, value, datetime):
    try: 
        conn = create_connection()
        cur = conn.cursor()
        result = cur.execute("INSERT INTO device_measurements (device_mac_address, property, value, datetime) values (?, ?, ?, ?)", (device_mac_address, prop, value, datetime))
        conn.commit() 
        return (f"last inserted id: {cur.lastrowid}")
        conn.close()
    except mariadb.Error as e: 
        print(f"Error: {e}")
