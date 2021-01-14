import mariadb
import requests
import sys
from flask.json import jsonify
import json
from datetime import datetime
from dateutil.relativedelta import *
import time

#database information
USER = '***'
PASSWORD = '***'
HOST = '***'
PORT = 0000
DATABASE = '***'

# this function creates a database connection based on the database information
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

# selects all data from table user and checks in which tier the user is
def select_enelogic_user_connection():
    try:
        conn = create_connection()
        cur = conn.cursor()
        result = cur.execute("SELECT * FROM enelogic_user")
        r = [dict((cur.description[i][0], value) \
            for i, value in enumerate(row)) for row in cur.fetchall()]
        conn.close()

        for item in r:
            client_id = item['client_id']
            tier = item['check_tier']
            access_token = item['access_token']
            if tier == 0:
                check_for_data(client_id, access_token)
            if tier == 1:
                first_call(access_token, client_id)
    except mariadb.Error as e: 
        print(f"Error: {e}")

# checks if data is available for the given access token
def check_for_data(client_id, access_token):
    try:
        token_url = "https://enelogic.com/api/measuringpoints/?access_token="+ str(access_token)
        result = requests.get(token_url)

        if result is not [] and result.status_code == 200:
            data = json.loads(result.text)
            print("result")
            print(result.json())
            if data and data[0]:
                print("UnitType exists")
                id1 = data[0]['id']
                id2 = data[1]['id']
                unitType1 = data[0]['unitType']
                unitType2 = data[1]['unitType']
                id_elec= ""
                id_gas = ""
            
                if(unitType1 == 0):
                    id_elec = id1
                    id_gas = id2
                if(unitType1 == 1):
                    id_elec = id2
                    id_gas = id1
                if not (id_elec is None and id_gas is None):
                    insert_id_to_enelogic_device(client_id, id_gas, id_elec)
                    update_tier(1, access_token)

    except:
        return jsonify({"message: ": "No data available"}), 403

#inserts properties(gas and electric) into the enelogic_device table
def insert_id_to_enelogic_device(client_id, id_gas, id_elec):
    try:
        conn = create_connection()
        cur = conn.cursor()
        result = cur.execute("INSERT INTO enelogic_device(client_id, enelogic_gas, enelogic_electric) VALUES (?, ?, ?)", (str(client_id), str(id_gas), str(id_elec),))
        conn.commit()
        conn.close()
        return result
    except mariadb.Error as e: 
        print(f"Error: {e}")

# updates the tier from a user
def update_tier(tier, access_token):      
    try:
        conn = create_connection()
        cur = conn.cursor()
        result = cur.execute("UPDATE enelogic_user SET check_tier = %s WHERE access_token = %s", (int(tier), str(access_token),))
        conn.commit()
        conn.close()
        return result
    except mariadb.Error as e: 
        print(f"Error: {e}")

# the first Enelogic requests for an user
def first_call(access_token, client_id):
    try:
        conn = create_connection()
        cur = conn.cursor()
        result = cur.execute("SELECT enelogic_gas, enelogic_electric FROM enelogic_device WHERE client_id = %s", (str(client_id),))
        r = [dict((cur.description[i][0], value) \
            for i, value in enumerate(row)) for row in cur.fetchall()]
        conn.close()
        var_gas = r[0]['enelogic_gas']
        var_elec = r[0]['enelogic_electric']
        request_13months(access_token, var_gas)
        request_13months(access_token, var_elec)
        request_40days(access_token, var_gas)
        request_40days(access_token, var_elec)
        update_tier(2, access_token)
    except mariadb.Error as e: 
        print(f"Error: {e}")

# Enelogic request for 13 months
def request_13months(access_token, property):
        current_date = datetime.today()
        date_13months = current_date - relativedelta(months=13)
        current_date_str = current_date.strftime('%Y-%m-%d')
        date_13months_str = date_13months.strftime('%Y-%m-%d')
       
        token_url = "https://enelogic.com/api/measuringpoints/"+property+"/datapoint/months/"+date_13months_str+"/"+current_date_str+"?access_token="+ str(access_token)
        result = requests.get(token_url)
        if result is not [] and result.status_code == 200:
            data = json.loads(result.text)
            insert_enelogic_measurements(data, property)

# Enelogic request for 40 days
def request_40days(access_token, property):
    current_date = datetime.today()
    date_40days = current_date - relativedelta(days=40)
    current_date_str = current_date.strftime('%Y-%m-%d')
    date_40days_str = date_40days.strftime('%Y-%m-%d')

    token_url = "https://enelogic.com/api/measuringpoints/" + property + "/datapoint/days/" + date_40days_str + "/" + current_date_str + "?access_token=" + str(access_token)
    result = requests.get(token_url)
    if result is not [] and result.status_code == 200:
        data = json.loads(result.text)
        insert_enelogic_measurements(data, property)

# inserts received JSON data into enelogic_measurements table
def insert_enelogic_measurements(data, property):
    for item in data:
        try:
            conn = create_connection()
            cur = conn.cursor()
            result = cur.execute("INSERT INTO enelogic_measurements(measurement_id, rate, property, value, datetime) VALUES (?, ?, ?, ?, ?)", (str(item['id']), str(item['rate']), str(property), str(item['quantity']), str(item['date']),))
            conn.commit()
            conn.close()
        except mariadb.Error as e: 
            print(f"Error: {e}")

# restart process after a day
def countdown(t):
    while t:
        mins, secs = divmod(t, 60)
        timer = '{:02d}:{:02d}'.format(mins, secs)
        time.sleep(1)
        t -= 1
        if t == 0:
            start = time.time()
            select_enelogic_user_connection()
            end = time.time()
            countdown(86400 - (int(end) - int(start)))

countdown(5)
