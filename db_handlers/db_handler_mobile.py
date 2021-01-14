#CREATORS: Victor Woord, Boet Schrama, Gulsah Kurnaz, Ben van Ommen
import mariadb

import sys
import json
from datetime import datetime
from json import JSONEncoder

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

#this is used in select function
class DateTimeEncoder(JSONEncoder):
        #override the default method
        def default(self, obj):
            if isinstance(obj, (type(datetime.date), type(datetime))):
                return obj.isoformat()

#this is used in select function
def query_db(query, house_id, one=False):
    conn = create_connection()
    cur = conn.cursor()
    cur.execute(query, house_id)
    r = [dict((cur.description[i][0], value) \
               for i, value in enumerate(row)) for row in cur.fetchall()]
    cur.connection.close    ()
    return (r[0] if r else None) if one else r

#this function inserts a house in the database based on data of the parameter (s)
def insert_house(house_id, postal_code, house_number, house_number_addition):
    try: 
        conn = create_connection()
        cur = conn.cursor()
        result = cur.execute("INSERT INTO house(house_id, postal_code, house_number, house_number_addition) VALUES (?, ?, ?, ?)", (house_id, postal_code, house_number, house_number_addition))
        conn.commit() 
        return (f"Last Inserted ID: {cur.lastrowid}")
        conn.close()
    except mariadb.Error as e: 
        print(f"Error: {e}")

#this function inserts a device in the database based on data of the parameter (s)
def insert_device(house_id, device_mac_address):
    try:
        conn = create_connection()
        cur = conn.cursor()
        result = cur.execute("INSERT INTO device(house_id, device_mac_address) VALUES (?, ? )", (house_id, device_mac_address))
        conn.commit() 
        return (f"Last Inserted ID: {cur.lastrowid}")
        conn.close()
    except mariadb.Error as e: 
        print(f"Error: {e}")

#this function inserts house devices in the database based on data of the parameter (s)
def insert_house_devices(postal_code, house_number, house_number_addition, device):
    try:
        conn = create_connection()
        cur = conn.cursor()
        result = cur.execute("INSERT INTO house_devices(postal_code, house_number, house_number_addition, device) VALUES (?, ?, ?, ? )", (postal_code, house_number, house_number_addition, device))
        conn.commit() 
        return (f"Last Inserted ID: {cur.lastrowid}")
        conn.close()
    except mariadb.Error as e: 
        print(f"Error: {e}")

#this function selects house devices from the database based on data of the parameter (s)
def select_house_devices(postal_code, house_number, house_number_addition):
    try: 
        devices_query = query_db("SELECT device FROM house_devices WHERE postal_code = %s AND house_number = %s AND house_number_addition = %s", (str(postal_code), house_number, str(house_number_addition)))
        devices_query_json_output = json.dumps(devices_query,  cls=DateTimeEncoder)
        return devices_query_json_output
    except mariadb.Error as e: 
        print(f"Error: {e}")

#this function selects house data from the database based on data of the parameter (s)
def select_house_data(house_id):
    device_measurements_query = query_db("SELECT device_measurements.device_mac_address, device_measurements.property, device_measurements.value, device_measurements.datetime FROM device_measurements LEFT JOIN device ON device_measurements.device_mac_address = device.device_mac_address LEFT JOIN house ON device.house_id = house.house_id WHERE house.house_id =%s AND device_measurements.datetime >= DATE(NOW()) - INTERVAL 7 DAY ORDER BY device_measurements.datetime ASC", (str(house_id),))
    device_measurements_query_json_output = json.dumps(device_measurements_query,  cls=DateTimeEncoder)
    return device_measurements_query_json_output

#this function selects registered from the database based on data of the parameter (s)
def select_registered(house_id):
    temp_house_id = house_id
    conn = create_connection()
    cur = conn.cursor()
    cur.execute("SELECT registered FROM house WHERE house_id = %s", (str(house_id),))
    r = [dict((cur.description[i][0], value) \
        for i, value in enumerate(row)) for row in cur.fetchall()]
    registered = r[0]["registered"]
    cur.connection.close()
    if registered == 1:
        return 1
    else:
        registered = 1
        update_registered(temp_house_id, registered)
        return 0 

#this function selects api key from the database based on data of the parameter (s)
def select_api_key():
    conn = create_connection()
    cur = conn.cursor()
    cur.execute("SELECT api_key FROM security LIMIT 1")
    r = [dict((cur.description[i][0], value) \
        for i, value in enumerate(row)) for row in cur.fetchall()]
    cur.connection.close    ()
    return (r[0] if r else None) if False else r[0]["api_key"]

#this function selects secret key from the database based on data of the parameter (s)
def select_secret_key():
    conn = create_connection()
    cur = conn.cursor()
    cur.execute("SELECT secret_key FROM security LIMIT 1")
    r = [dict((cur.description[i][0], value) \
        for i, value in enumerate(row)) for row in cur.fetchall()]
    cur.connection.close    ()
    return (r[0] if r else None) if False else r[0]["secret_key"]

#this function updates registered in the database based on data of the parameter (s)
def update_registered(house_id, registered):
    try:
        conn = create_connection()
        cur = conn.cursor()
        result = cur.execute("UPDATE house SET registered = %s WHERE house_id = %s",(int(registered), str(house_id),))
        conn.commit()
        conn.close()
        return result
    except mariadb.Error as e: 
        print(f"Error: {e}")