import logging
from typing import Type

from fastapi import Depends, FastAPI
from fastapi.responses import JSONResponse
from fastapi.security.http import HTTPAuthorizationCredentials
from fastapi_sqlalchemy import DBSessionMiddleware
from fastapi_sqlalchemy import db
from fastapi.middleware.cors import CORSMiddleware

from auth import (
    AccountSessionTokenBearer,
    AdminSessionTokenBearer,
    DeviceSessionTokenBearer,
)
from db import db_url, session_args
from firebase import firebase_dynamic_link
from schema import (
    AccountActivate,
    AccountCreate,
    AccountItem,
    AccountSession,
    BadRequest,
    DeviceCompleteItem,
    DeviceItem,
    DeviceItemMeasurementTime,
    DeviceTypeItem,
    Forbidden,
    HttpStatus,
    MeasurementsUploadFixed,
    MeasurementsUploadVariable,
    MeasurementsUploadResult,
    NotFound,
    Unauthorized,
)
from user import get_admin
import crud
from data.loader import csv_create_update

__version__ = '0.95'

logging.basicConfig(level=logging.DEBUG)

app = FastAPI(title='Twomes API', version=__version__)

app.add_middleware(DBSessionMiddleware, db_url=db_url,
                   session_args=session_args)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

admin_auth = AdminSessionTokenBearer()
account_auth = AccountSessionTokenBearer()
device_auth = DeviceSessionTokenBearer()


@app.on_event("startup")
# Fill DeviceType and Property values on API startup.
async def startup_event():
    csv_create_update()


def http_status(http_status_class: Type[HttpStatus], message: str) -> JSONResponse:
    return JSONResponse(status_code=http_status_class.code, content={'detail': message})


@app.post(
    '/account',
    response_model=AccountItem,
    responses={
        BadRequest.code: {'model': BadRequest},
        Unauthorized.code: {'model': Unauthorized},
    }
)
def account_create(account_input: AccountCreate,
                   authorization: HTTPAuthorizationCredentials = Depends(admin_auth)):
    pseudonym = account_input.pseudonym
    admin_session_token = authorization.credentials

    admin = get_admin(admin_session_token)
    if not admin:
        return http_status(Unauthorized, 'Invalid admin session token')

    if pseudonym:
        if crud.account_by_pseudonym(db.session, pseudonym):
            return http_status(BadRequest, 'Account pseudonym already in use')
    else:
        pseudonym = crud.generate_pseudonym(db.session)

    account = crud.account_create(db.session, pseudonym)

    location = account_input.location
    tz_name = account_input.tz_name
    crud.building_create(db.session, account, location, tz_name)

    url = firebase_dynamic_link(account.activation_token)

    return AccountItem(
        id=account.id,
        pseudonym=pseudonym,
        activation_token=account.activation_token,
        firebase_url=url
    )


@app.post(
    '/account/activate',
    response_model=AccountSession,
    responses={
        NotFound.code: {'model': NotFound}
    }
)
def account_activate(activation_token: AccountActivate):
    account = crud.account_by_token(
        db.session, activation_token.activation_token)
    if not account:
        return http_status(NotFound, 'No account found for provided activation token')

    session_token = crud.account_session_token(db.session, account)

    return AccountSession(session_token=session_token)


@app.post(
    '/account/device/provision',
    response_model=DeviceCompleteItem,
    responses={
        BadRequest.code: {'model': BadRequest},
        Unauthorized.code: {'model': Unauthorized},
        NotFound.code: {'model': NotFound}
    }
)
def account_device_provision(device_input: DeviceItem,
                             authorization: HTTPAuthorizationCredentials = Depends(account_auth)):
    device_name = device_input.name
    account_session_token = authorization.credentials

    account = crud.account_by_session(db.session, account_session_token)
    device_type = crud.device_type_by_name(
        db.session, device_input.device_type_name)

    if not account:
        return http_status(Unauthorized, 'Invalid account session token')

    if crud.device_by_name(db.session, device_name):
        return http_status(BadRequest, 'Device name already in use')

    if not crud.device_type_by_name(db.session, device_type.name):
        return http_status(BadRequest, f'Unknown device type "{device_type.name}"')

    created_device = crud.device_create(db.session, device_name, device_type)
    crud.device_activate(db.session, account, created_device)
    session_token = crud.device_session_token(db.session, created_device)

    complete_device: DeviceCompleteItem = {
        "id": created_device.id,
        "name": created_device.name,
        "device_type": created_device.device_type,
        "session_token": session_token,
        "created_on": created_device.created_on,
        "activated_on": created_device.activated_on,
    }
    return DeviceCompleteItem(**complete_device)


@app.get(
    '/device_type/{device_name}',
    response_model=DeviceTypeItem,
    responses={
        BadRequest.code: {'model': BadRequest},
        Unauthorized.code: {'model': Unauthorized},
        NotFound.code: {'model': NotFound}
    }
)
def device_type(device_name: str,
                authorization: HTTPAuthorizationCredentials = Depends(account_auth)):

    account_session_token = authorization.credentials

    account = crud.account_by_session(db.session, account_session_token)
    if not account:
        return http_status(Unauthorized, 'Invalid account session token')

    device = crud.device_by_name(db.session, device_name)
    if not device:
        return http_status(NotFound, f'Device {device_name} not found')

    device_type = device.device_type

    return device_type


@app.get(
    '/device/{device_name}',
    response_model=DeviceItemMeasurementTime,
    responses={
        BadRequest.code: {'model': BadRequest},
        Unauthorized.code: {'model': Unauthorized},
        NotFound.code: {'model': NotFound}
    }
)
def device_read(device_name: str,
                authorization: HTTPAuthorizationCredentials = Depends(account_auth)):

    account_session_token = authorization.credentials

    account = crud.account_by_session(db.session, account_session_token)
    if not account:
        return http_status(Unauthorized, 'Invalid account session token')

    device = crud.device_by_name(db.session, device_name)
    if not device or not device.building or device.building.account_id != account.id:
        return http_status(NotFound, f'Device {device_name} not found')

    timestamp = crud.device_latest_measurement_timestamp(db.session, device.id)
    device.latest_measurement_timestamp = timestamp

    return device


@app.get(
    '/device',
    response_model=DeviceCompleteItem,
    responses={
        BadRequest.code: {'model': BadRequest},
        Unauthorized.code: {'model': Unauthorized},
        NotFound.code: {'model': NotFound}
    }
)
def device_read_self(authorization: HTTPAuthorizationCredentials = Depends(device_auth)):

    device_session_token = authorization.credentials

    device = crud.device_by_session(db.session, device_session_token)
    if not device:
        return http_status(Unauthorized, 'Invalid device session token')
    if not device.building_id:
        return http_status(Forbidden, 'Device not attached to account')

    return device


@app.post(
    '/device/measurements/fixed-interval',
    response_model=MeasurementsUploadResult,
    responses={
        BadRequest.code: {'model': BadRequest},
        Unauthorized.code: {'model': Unauthorized},
        Forbidden.code: {'model': Forbidden},
        NotFound.code: {'model': NotFound}
    }
)
def device_upload_fixed(measurements_upload: MeasurementsUploadFixed,
                        authorization: HTTPAuthorizationCredentials = Depends(device_auth)):

    device_session_token = authorization.credentials

    device = crud.device_by_session(db.session, device_session_token)
    if not device:
        return http_status(Unauthorized, 'Invalid device session token')
    if not device.building_id:
        return http_status(Forbidden, 'Device not attached to account')

    properties = device.device_type.properties
    valid_property_names = {p.name for p in properties}
    property_names = {
        item.property_name for item in measurements_upload.property_measurements}

    invalid_property_names = property_names - valid_property_names
    if invalid_property_names:
        return http_status(BadRequest, f'Invalid property name(s): {invalid_property_names}')

    data = crud.upload_fixed_to_variable(measurements_upload)
    upload = crud.device_upload_variable(db.session, device, properties, data)

    return MeasurementsUploadResult(size=upload.size, server_time=upload.server_time)


@app.post(
    '/device/measurements/variable-interval',
    response_model=MeasurementsUploadResult,
    responses={
        BadRequest.code: {'model': BadRequest},
        Unauthorized.code: {'model': Unauthorized},
        Forbidden.code: {'model': Forbidden},
        NotFound.code: {'model': NotFound}
    }
)
def device_upload_variable(measurements_upload: MeasurementsUploadVariable,
                           authorization: HTTPAuthorizationCredentials = Depends(device_auth)):

    device_session_token = authorization.credentials

    device = crud.device_by_session(db.session, device_session_token)
    if not device:
        return http_status(Unauthorized, 'Invalid device session token')
    if not device.building_id:
        return http_status(Forbidden, 'Device not attached to account')

    properties = device.device_type.properties
    valid_property_names = {p.name for p in properties}
    property_names = {
        item.property_name for item in measurements_upload.property_measurements}

    invalid_property_names = property_names - valid_property_names
    if invalid_property_names:
        return http_status(BadRequest, f'Invalid property name(s): {invalid_property_names}')

    upload = crud.device_upload_variable(
        db.session, device, properties, measurements_upload)

    return MeasurementsUploadResult(size=upload.size, server_time=upload.server_time)
