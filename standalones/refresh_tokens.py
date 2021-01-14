#CREATORS: Victor Woord, Boet Schrama, Gulsah Kurnaz, Ben van Ommen

import mariadb
import requests
import sys
from flask.json import jsonify
import json
import time

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

#this function selects different parameters from Enelogic users
def select_enelogic_user_data():
    try:
        conn = create_connection()
        cur = conn.cursor()
        result = cur.execute("SELECT * FROM enelogic_user")
        r = [dict((cur.description[i][0], value) \
            for i, value in enumerate(row)) for row in cur.fetchall()]
        conn.close()

        for item in r:
            refresh_tokens(item['client_id'], item['client_secret'], item['refresh_token'])

        return (r if r else None) if False else r
    except mariadb.Error as e: 
        print(f"Error: {e}")

#this function does an API call to enelogic with the given parameters to refresh the tokens
def refresh_tokens(client_id, client_secret, refresh_token):
    try:
        auth = (client_id, client_secret)
        params = {
            "grant_type":"refresh_token",
            "refresh_token":str(refresh_token)
            }   
        token_url = "https://enelogic.com/oauth/v2/token"
        result = requests.post(token_url, auth=auth, data=params)
        data = json.loads(result.text)
        accessT = data['access_token']
        refreshT = data['refresh_token']
        insert_new_tokens(client_id, accessT, refreshT)
    except: 
        print(f"Error: Access token is not valid or missing")

#this function updates the new tokens int the table
def insert_new_tokens(client_id, access_token, refresh_token):
    try:
        conn = create_connection()
        cur = conn.cursor()
        result = cur.execute("UPDATE enelogic_user SET access_token = %s, refresh_token = %s WHERE client_id = %s",(str(access_token), str(refresh_token), str(client_id),))
        conn.commit()
        conn.close()
        return result
    except mariadb.Error as e: 
        print(f"Error: {e}")

#this countdown function counts the seconds until it's a day later and starts the refresh process again
def countdown(t):
    while t:
        mins, secs = divmod(t, 60)
        timer = '{:02d}:{:02d}'.format(mins, secs)
        time.sleep(1)
        t -= 1
        if t == 0:
            start = time.time()
            select_enelogic_user_data()
            end = time.time()
            countdown(259000 - (int(end)-int(start)))

countdown(5)