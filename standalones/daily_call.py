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
            user=USER,
            password=PASSWORD,
            host=HOST,
            port=PORT,
            database=DATABASE)
        print("Success connecting to MariaDB Platform:")
        return conn
    except mariadb.Error as e:
        print(f"Error connecting to MariaDB Platform: {e}")
        sys.exit(1)

# selects all data from table user
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
            if tier == 2:
                daily_call_check(access_token, client_id)
    except mariadb.Error as e:
        print(f"Error: {e}")

# selects gas and electric property from user in tier 2
def daily_call_check(access_token, client_id):
    try:
        conn = create_connection()
        cur = conn.cursor()
        result = cur.execute("SELECT enelogic_gas, enelogic_electric FROM enelogic_device WHERE client_id = %s",
                             (str(client_id),))
        r = [dict((cur.description[i][0], value) \
                  for i, value in enumerate(row)) for row in cur.fetchall()]
        conn.close()
        var_gas = r[0]['enelogic_gas']
        var_elec = r[0]['enelogic_electric']
        daily_call(access_token, var_gas)
        daily_call(access_token, var_elec)
    except mariadb.Error as e:
        print(f"Error: {e}")

# Enelogic request for given property
def daily_call(access_token, property):
        current_date = datetime.today()
        date_4_days = current_date - relativedelta(days=4)
        date_4_days_str = date_4_days.strftime('%Y-%m-%d')
        current_date_str = current_date.strftime('%Y-%m-%d')
        token_url = "https://enelogic.com/api/measuringpoints/" + property + "/datapoints/" + date_4_days_str + "/" + date_4_days_str + "?access_token=" + str(access_token)
        result = requests.get(token_url)
        if result is not [] and result.status_code == 200:
            data = json.loads(result.text)
            insert_enelogic_measurements(data, property)

# inserts received json data into eneloic_measurements table
def insert_enelogic_measurements(data, property):
    for item in data:
        try:
            conn = create_connection()
            cur = conn.cursor()
            result = cur.execute("INSERT INTO enelogic_measurements(rate, property, value, datetime) VALUES (?, ?, ?, ?)", (str(item['rate']), str(property), str(item['quantity']), str(item['timezone']),))
            conn.commit()
            conn.close()
        except:
            print("Error")

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
