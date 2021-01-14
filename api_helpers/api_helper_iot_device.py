#CREATORS: Victor Woord, Boet Schrama, Gulsah Kurnaz, Ben van Ommen

import pytz
from flask import jsonify
from datetime import datetime
from db_handlers.db_handler_iot_device import *

#this function extracts data from the payload (json format). 
#after that, times per measurement are calculated based on intervals that are included in the payload. 
#finally, function 'insert_central_heating_temperature' is called to put the data in the database.
def set_central_heating_temperature_helper(data):
    try:
        #extract data from the payload and convert it into variables
        device_mac_address = data['id']
        unix_time = data['dataSpec']['lastTime']
        time_interval = data['dataSpec']['interval']
        total_interval = data['dataSpec']['total']  - 1

        #this while loop takes care of the order of retrieving the pipe temperatures. 
        #the corresponding time is calculated by subtracting the interval(in seconds) from the valid time.
        while total_interval >= 0 :  
            converted_time = convert_unix_to_utc_time(unix_time)    # Convert to valid datetime
            pipe_temp1 = data['data']['pipeTemp1']
            pipe_temp2 = data['data']['pipeTemp2']

            result = insert_device_measurements(device_mac_address, 'pipeTemp1', pipe_temp1[total_interval], converted_time)
            result = insert_device_measurements(device_mac_address, 'pipeTemp2', pipe_temp2[total_interval], converted_time)

            unix_time = int(unix_time) - int(time_interval)
            total_interval = total_interval - 1 
            
        return str(result)
    except:
        #if an error occurs, return this as a message
        return jsonify({'message': 'Check sended data'}), 403

#this function extracts data from the payload (json format). 
#after that, times per measurement are calculated based on intervals that are included in the payload. 
#finally, function 'insert_smart_meter' is called to put the data in the database.
def set_smart_meter_helper(data):
    try:
        
        #extract data from the payload and convert it into variables
        unix_time = data['dataSpec']['lastTime']
        interval = data['dataSpec']['interval']  
        total_interval = data['dataSpec']['total'] - 1
        device_mac_address = data['id'] 
        dsmr = data['data']['dsmr']
        electricity_delivered_to_t1 = data['data']['evt1']
        electricity_delivered_to_t2 = data['data']['evt2']
        electricity_delivered_by_t1 = data['data']['egt1']
        electricity_delivered_by_t2 = data['data']['egt2']
        tariff_indicator = data['data']['ht']
        electricity_received = data['data']['ehv']
        electricity_delivered = data['data']['ehl']
        gas = data['data']['gas']
        property_values_list = []
        
        for item in data['data']:
            property_values_list.append(item)
        #this while loop takes care of the order of retrieving the pipe temperatures. 
        #the corresponding time is calculated by subtracting the interval(in seconds) from the valid time.

        while total_interval >= 0:  
            counter = 0
            converted_time = convert_unix_to_utc_time(unix_time)
            time_gas = data['data']['tgas'][total_interval]     
            time_gas_converted = convert_unix_to_utc_time(time_gas)
            property_values = [dsmr[total_interval], electricity_delivered_to_t1[total_interval], electricity_delivered_to_t2[total_interval], electricity_delivered_by_t1[total_interval], 
                electricity_delivered_by_t2[total_interval], tariff_indicator[total_interval], electricity_received[total_interval], electricity_delivered[total_interval], gas[total_interval], 
                time_gas_converted]

            for item in property_values:
                result = insert_device_measurements(device_mac_address, property_values_list[counter], item, converted_time)
                counter = counter + 1

            unix_time = int(unix_time) - int(interval)
            total_interval = total_interval - 1
            
        return str(result)
    except:
        return jsonify({'message': 'Check sended data'}), 403

#this function extracts data from the payload (json format). 
#after that, times per measurement are calculated based on intervals that are included in the payload. 
#finally, function 'insert_room_temperature' is called to put the data in the database.
def set_room_temperature_helper(data):
    device_mac_address = data['id']
    unix_time = data['dataSpec']['lastTime']
    time_interval = data['dataSpec']['interval']  
    total_interval = data['dataSpec']['total'] - 1
    room_temp = data['data']['roomTemp']  

    #this while loop takes care of the order of retrieving the pipe temperatures. 
    #the corresponding time is calculated by subtracting the interval(in seconds) from the converted time.
    while total_interval >= 0:
        converted_time = convert_unix_to_utc_time(unix_time)

        result = insert_device_measurements(device_mac_address, 'roomTemp', room_temp[total_interval], converted_time)
        unix_time = int(unix_time) - int(time_interval)
        total_interval = total_interval - 1
    return str(result)

#this function organizes the data that has been sent and passes it on to the function insert_device_measurements
def set_opentherm_helper(data):
    device_mac_address = data['deviceMac']
    unix_time = data['time']
    converted_time = convert_unix_to_utc_time(unix_time)

    for obj in data['measurements']:
        result = insert_device_measurements(device_mac_address, obj['property'], obj['value'], converted_time)

    return str(result)

#this function converts unix to utc 
def convert_unix_to_utc_time(time):
    tz = pytz.timezone('Europe/Amsterdam')
    converted_time = datetime.fromtimestamp(int(time), tz).isoformat()
    return converted_time