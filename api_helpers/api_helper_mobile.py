from flask import jsonify
from datetime import datetime
from cryptography.fernet import Fernet
from db_handlers.db_handler_mobile import *

#this function organizes the data that has been sent and passes it on to the function insert_house
def set_house_helper(house_data, data):
    try:
        house_id = house_data['house_id']
        postal_code = data['postal_code']
        house_number = data['house_number']
        house_number_addition = data['house_number_addition']
        result = insert_house(house_id, postal_code, house_number, house_number_addition)
        return str(result)
    except:
        return jsonify({'message': 'Check sended data'}), 403

#this function organizes the data that has been sent and passes it on to the function insert_device
def set_devices_helper(house_data, data):
    try:
        house_id = house_data['house_id']
        smart_meter_id = data['smart_meter_id']
        temperature_sensor_id = data['temperature_sensor_id']
        system_id= data['system_id']

        result = insert_device(house_id, smart_meter_id)
        result = insert_device(house_id, temperature_sensor_id)
        result = insert_device(house_id, system_id)
        return str(result)
    except:
        return jsonify({'message': 'Check sended data'}), 403

#this function organizes the data that has been sent and passes it on to the function insert_house_devices
def set_house_devices_helper(data):
    try:
        postal_code = data['postal_code']
        house_number = data['house_number']
        house_number_addition = data['house_number_addition']
        device = data['device']

        result = insert_house_devices(postal_code, house_number, house_number_addition, device)
        return str(result)
    except:
        return jsonify({'message': 'Check sended data'}), 403

#this function organizes the data that has been sent and passes it on to the function select_house_data
def get_house_data_helper(data):
    try:
        house_id = data['house_id']
        result = select_house_data(house_id)
        return str(result)
    except:
        return jsonify({'message': 'Check sended data'}), 403

#this function organizes the data that has been sent and passes it on to the function select_registered
def get_registration_helper(house_id):
    try:
        result = select_registered(house_id)
        return str(result)
    except:
        return jsonify({'message': 'Check sended data'}), 403

#this function organizes the data that has been sent and passes it on to the function select_house_devices
def get_house_devices_helper(data):
    try:
        postal_code = data['postal_code']
        house_number = data['house_number']
        house_number_addition = data['house_number_addition']

        result = select_house_devices(postal_code, house_number, house_number_addition)
        return str(result)
    except:
        return jsonify({'message': 'Check sended data'}), 403

#this function organizes the data that has been sent and passes it on to the function select_api_key
def get_api_key_helper():
    try:
        result = select_api_key()
        return str(result)
    except:
        return jsonify({'message': 'Something went wrong'}), 403

#this function organizes the data that has been sent and passes it on to the function select_secret_key
def get_secret_key_helper():
    try:
        secret_key = select_secret_key()
        key = 'PsQGv2aw3iHI0nPT/jsWwqObvOyp4dkRtfbUcU+qd7M='
        f = Fernet(key)
        plain_key = f.decrypt(secret_key.encode())
        return plain_key.decode()
    except:
        return jsonify({'message': 'Something went wrong'}), 403

