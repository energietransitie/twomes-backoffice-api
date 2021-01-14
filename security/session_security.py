import jwt
import sys
import hashlib
import os
from cryptography.fernet import Fernet
from functools import wraps
from flask import jsonify, request
#this imports api_helpers as ahm
sys.path.append(os.path.join(os.path.dirname(__file__), '..', 'api_helpers'))
import api_helper_mobile as ahm

#this function returns encrypted sended key if api key is valid
def encrypt_secret_key(sended_key, api_key_client):
    try:
        if check_for_valid_api_key(api_key_client):
            f = Fernet(sended_key.replace(" ", "+"))
            token = f.encrypt(ahm.get_secret_key_helper().encode())
            return token.decode()
        else:
            return jsonify({'message': 'API key incorrect'}), 402
    except:
        return jsonify({'message': 'Cannot encrypt, check sended data'}), 403

#this function checks if api key is valid
def check_for_valid_api_key(api_key_client):
    api_key_client = hashlib.sha512(api_key_client.encode())
    api_key_client = api_key_client.hexdigest()
    if api_key_client == ahm.get_api_key_helper():
        return bool(True)
    else:
        return bool(False)

#this function checks if token and key are valid
def check_for_token_and_key(func):
    @wraps(func)
    def wrapped(*args, **kwargs):
        token = request.args.get('token')
        if not token:
            return jsonify({'message': 'Missing token'}, 403)
        try:
            data = jwt.decode(token, ahm.get_secret_key_helper())
            if data['APIkey'] and check_for_valid_api_key(data['APIkey']):
                return func(*args, **kwargs)
            else:
                abort(401)
        except:
            return jsonify({'message': 'Invalid token or key'}), 403
        return func(*args, **kwargs)
    return wrapped
