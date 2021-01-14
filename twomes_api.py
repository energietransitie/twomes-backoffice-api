import ssl

from flask import Flask, jsonify, request
from flask_cors import CORS

from security.session_security import *
from api_helpers.api_helper_iot_device import *
from api_helpers.api_helper_mobile import *
from enelogic.new_enelogic_user import *

#this makes app a Flask application
app = Flask(__name__)
#this adds CORS to the application
CORS(app)

#this is for the sll certificate
context = ssl.SSLContext(ssl.PROTOCOL_TLSv1_2)
context.load_cert_chain('./security/cert.pem', './security/privkey.pem')

#this request is for the mobile applicaton
#this request is used to set a house in de database, this request expects a token and checks for token and key
@app.route('/set/house', methods=['POST'])
@check_for_token_and_key
def set_house():
    try:
        token = request.args.get('token')
        data = jwt.decode(token, get_secret_key_helper())
        result = set_house_helper(data, json.loads(request.data))
        return str(result)
    except:
        return jsonify({'message': 'Cannot create house, check sended data'}), 403

#this request is for a smart device
#this request is used to set a house central heating system in de database
@app.route('/set/house/centralHeatingTemperature', methods=['POST'])
def set_central_heating_temperature():
    try:
        result = set_central_heating_temperature_helper(json.loads(request.data))
        return str(result)
    except:
        return jsonify({'message': 'Cannot create central heating temperature, check sended data'}), 403

#this request is for a smart device
#this request is used to set a house smart meter in de database
@app.route('/set/house/smartMeter', methods=['POST'])
def set_smart_meter():
    try:
        result = set_smart_meter_helper(json.loads(request.data))
        return str(result)
    except:
        return jsonify({'message': 'Cannot create smart meter, check sended data'}), 403


#this request is for a smart device
#this request is used to set a house room temperature in de database
@app.route('/set/house/roomTemperature', methods=['POST'])
def set_room_temperature():
    try:
        result = set_room_temperature_helper(json.loads(request.data))
        return str(result)
    except:
        return jsonify({'message': 'Cannot create room temperature, check sended data'}), 403

#this request is for a smart device
#this request is used to set a house opentherm in de database
@app.route('/set/house/opentherm', methods=['POST'])
def set_opentherm():
    try:
        result = set_opentherm_helper(json.loads(request.data))
        return str(result)
    except:
        return jsonify({'message': 'Cannot create opentherm, check sended data'}), 403

#this request is for the mobile application
#this request is used to set house devices in de database before a house_id is generated
@app.route('/set/houseDevices', methods=['POST'])
def set_house_devices():
    try:
        result = set_house_devices_helper(json.loads(request.data))
        return str(result)
    except:
        return jsonify({'message': 'Cannot create opentherm, check sended data'}), 403

#this request is for the mobile application
#this request is used to set house devices in de database after a house_id is generated
@app.route('/set/devices', methods=['POST'])
@check_for_token_and_key
def set_devices():
    try:
        token = request.args.get('token')
        data = jwt.decode(token, get_secret_key_helper())
        result = set_devices_helper(data, json.loads(request.data))  
        return str(result)
    except:
        return jsonify({'message': 'Cannot create new devices, check sended data'}), 403 

#this request is for the mobile application
#this request is the first call to get the data which is sent to the API
@app.route('/set/house/setEnelogicUsers', methods=['POST'])
def set_enelogic():
    try:
        result = set_enelogic_user(json.loads(request.data))  
        return str(result)
    except:
        return jsonify({'message': 'Cannot create new devices, check sended data'}), 403  

#this request is for the mobile application
#this request is used to start a session, this needs to be done before the mobile application makes a request
#this session is used for the security flow (encryption, decryption, jwt) between the mobile application and this API
@app.route('/startSession', methods=['GET'])
def start_session():
    try:
        house_id = request.args.get('house_id')
        house_key = request.args.get('house_key')
        result = encrypt_secret_key(house_id, house_key)
        return str(result)
    except:
        return jsonify({'message': 'Cannot execute request, check sended data'}), 403

#this request is for the mobile application
#this request is used to get house devices from de database before a house_id is generated
@app.route('/get/houseDevices', methods=['GET'])
def get_house_devices():
    try:
        house_key = request.args.get('house_key')
        if check_for_valid_api_key(house_key):
            result = get_house_devices_helper(json.loads(request.data))
            return str(result)   
    except:
        return jsonify({'message': 'Cannot check registration, check sended data'}), 403

#this request is for the mobile application
#this request is used to get house registration from de database 
@app.route('/get/house/registration', methods=['GET'])
def get_registration():
    try:
        house_id = request.args.get('house_id')
        result = get_registration_helper(house_id)
        return str(result)   
    except:
        return jsonify({'message': 'Cannot check registration, check sended data'}), 403

#this request is for the mobile application
#this request is used to get house data from de database 
@app.route('/get/house/data', methods=['GET'])
@check_for_token_and_key
def get_house_data():
    try:
        token = request.args.get('token')
        data = jwt.decode(token, get_secret_key_helper())
        result = get_house_data_helper(data)
        return str(result)   
    except:
        return jsonify({'message': 'Cannot get house data, check sended data'}), 403

#this runs the flask application
app.run(host='0.0.0.0',port='443', ssl_context=context)
