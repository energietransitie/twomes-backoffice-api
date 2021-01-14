#CREATORS: Victor Woord, Boet Schrama, Gulsah Kurnaz, Ben van Ommen

import mariadb
import sys
import time
from datetime import datetime
from py_weernl import weerLive

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

#this function inserts knmi data in the database based on data of the parameter (s)
def insert_knmi_data(data):
    try: 
        conn = create_connection()
        cur = conn.cursor()
        result = cur.execute("INSERT INTO knmi_data (location, temp, feeling_temp, weather_condition, relative_humidity, wind_direction, wind_speed_ms, wind_force, wind_speed_knots, wind_speed_kmh, air_pressure, air_pressure_mm, dew_point, sight_km, datetime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", (str(data['plaats']), data['temp'], data['gtemp'], data['samenv'], data['lv'], data['windr'], data['windms'], data['winds'], data['windk'], data['windkmh'], data['luchtd'], data['ldmmhg'], data['dauwp'], data['zicht'], datetime.now()))
        conn.commit() 
        return (f"Last Inserted ID: {cur.lastrowid}")
        conn.close()
    except mariadb.Error as e: 
        print(f"Error: {e}")

#this function gets data using the API 
def get_knmi_data():
    place = "Assendorp"
    api_key = "a5f67c3bfb"
    w = weerLive(api=api_key)

    while True:
        data = w.getData(place)
        insert_knmi_data(data['liveweer'][0])
        time.sleep(300)

#this function will always run and repeat itself
def repeat():
    try:
        get_knmi_data()
    except:
        repeat()

repeat()