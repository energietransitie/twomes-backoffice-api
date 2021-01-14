import mariadb
import sys
from flask.json import jsonify
import json

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

#this function extracts the parameters form the json payload and store them into variables
def set_enelogic_user(data):
    try:
        house_id = data['house_id']
        client_id = data['client_id']
        client_secret = data['client_secret']
        client_access = data['access_token']
        client_refresh = data['refresh_token']
        result = insert_enelogic_user(house_id, client_id, client_secret, client_access, client_refresh)
        return str(result)
    except:
        return jsonify({'message': 'Check sended data'}), 403

#this function inserts the new Enelogic users into the table
def insert_enelogic_user(house_id, client_id, client_secret, client_access, client_refresh):
    try:
        conn = create_connection()
        cur = conn.cursor()
        result = cur.execute("INSERT INTO enelogic_user(house_id, client_id, client_secret, access_token, refresh_token, check_tier) VALUES (?, ?, ? , ?, ?, ?)", (int(house_id), str(client_id), str(client_secret), str(client_access), str(client_refresh), 0))
        conn.commit()
        return (f"Last Inserted ID: {cur.lastrowid}")
        conn.close()
    except mariadb.Error as e: 
        print(f"Error: {e}")

