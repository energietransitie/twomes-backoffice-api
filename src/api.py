import logging
from typing import Type

from fastapi import Depends, FastAPI
from fastapi.responses import JSONResponse
from fastapi.security.http import (
    HTTPBearer,
    HTTPAuthorizationCredentials,
)
from fastapi_sqlalchemy import DBSessionMiddleware
from fastapi_sqlalchemy import db

from db import db_url, session_args
from firebase import firebase_dynamic_link
from schema import (
    AccountActivate,
    AccountCreate,
    AccountItem,
    AccountSession,
    BadRequest,
    DeviceVerify,
    DeviceCompleteItem,
    DeviceCreate,
    DeviceItem,
    DeviceItemMeasurementTime,
    DeviceSession,
    Forbidden,
    HttpStatus,
    MeasurementsUpload,
    MeasurementsUploadResult,
    NotFound,
    Unauthorized,
)
import crud

__version__ = '0.8'

logging.basicConfig(level=logging.DEBUG)

app = FastAPI(title='Twomes API', version=__version__)

app.add_middleware(DBSessionMiddleware, db_url=db_url, session_args=session_args)

auth = HTTPBearer()


def http_status(http_status_class: Type[HttpStatus], message: str) -> JSONResponse:
    return JSONResponse(status_code=http_status_class.code, content={'detail': message})


@app.post(
    '/account',
    response_model=AccountItem,
    responses={
        BadRequest.code: {'model': BadRequest}
    }
)
def account_create(account_input: AccountCreate):
    pseudonym = account_input.pseudonym

    if pseudonym:
        if crud.account_by_pseudonym(db.session, pseudonym):
            return http_status(BadRequest, 'Account pseudonym already in use')
    else:
        pseudonym = crud.generate_pseudonym(db.session)

    account = crud.account_create(db.session, pseudonym)

    location = account_input.location if account_input.location else None
    crud.building_create(db.session, account, location)

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
    account = crud.account_by_token(db.session, activation_token.activation_token)
    if not account:
        return http_status(NotFound, 'No account found for provided activation token')

    session_token = crud.account_session_token(db.session, account)

    return AccountSession(session_token=session_token)


@app.post(
    '/device',
    response_model=DeviceItem,
    responses={
        BadRequest.code: {'model': BadRequest}
    }
)
def device_create(device_input: DeviceCreate):
    device_type_name = device_input.device_type
    proof_of_presence_id = device_input.proof_of_presence_id

    device_type = crud.device_type_by_name(db.session, device_type_name)
    if not device_type:
        return http_status(BadRequest, f'Unknown device type "{device_type_name}"')

    if crud.device_by_pop(db.session, proof_of_presence_id):
        return http_status(BadRequest, 'Proof-of-presence identifier already in use')

    device = crud.device_create(db.session, device_type, proof_of_presence_id)
    return device


@app.post(
    '/device/activate',
    response_model=DeviceItem,
    responses={
        BadRequest.code: {'model': BadRequest},
        Unauthorized.code: {'model': Unauthorized},
        NotFound.code: {'model': NotFound}
    }
)
def device_activate(device_verify: DeviceVerify,
                    authorization: HTTPAuthorizationCredentials = Depends(auth)):

    proof_of_presence_id = device_verify.proof_of_presence_id
    account_session_token = authorization.credentials

    device = crud.device_by_pop(db.session, proof_of_presence_id)
    if not device:
        return http_status(NotFound, 'No device found for provided proof-of-presence id')
    if device.activated_on:
        return http_status(BadRequest, 'Device already activated')

    account = crud.account_by_session(db.session, account_session_token)
    if not account:
        return http_status(Unauthorized, 'Invalid account session token')

    crud.device_activate(db.session, account, device)

    return device


@app.get(
    '/device/{device_id}',
    response_model=DeviceItemMeasurementTime,
    responses={
        BadRequest.code: {'model': BadRequest},
        Unauthorized.code: {'model': Unauthorized},
        NotFound.code: {'model': NotFound}
    }
)
def device_read(device_id: int,
                authorization: HTTPAuthorizationCredentials = Depends(auth)):

    account_session_token = authorization.credentials

    account = crud.account_by_session(db.session, account_session_token)
    if not account:
        return http_status(Unauthorized, 'Invalid account session token')

    device = crud.device_by_account_and_id(db.session, account, device_id)
    if not device:
        return http_status(NotFound, f'Device {device_id} not found')

    timestamp = crud.device_latest_measurement_timestamp(db.session, device_id)
    device.latest_measurement_timestamp = timestamp

    return device


@app.post(
    '/device/session',
    response_model=DeviceSession,
    responses={
        Forbidden.code: {'model': Forbidden},
        NotFound.code: {'model': NotFound}
    }
)
def device_session(device_verify: DeviceVerify):

    proof_of_presence_id = device_verify.proof_of_presence_id

    device = crud.device_by_pop(db.session, proof_of_presence_id)
    if not device:
        return http_status(NotFound, 'No device found for provided proof-of-presence id')
    if not device.building_id:
        return http_status(Forbidden, 'Device not attached to account')

    session_token = crud.device_session_token(db.session, device)

    return DeviceSession(session_token=session_token)


@app.get(
    '/device',
    response_model=DeviceCompleteItem,
    responses={
        BadRequest.code: {'model': BadRequest},
        Unauthorized.code: {'model': Unauthorized},
        NotFound.code: {'model': NotFound}
    }
)
def device_read_self(authorization: HTTPAuthorizationCredentials = Depends(auth)):

    device_session_token = authorization.credentials

    device = crud.device_by_session(db.session, device_session_token)
    if not device:
        return http_status(Unauthorized, 'Invalid device session token')
    if not device.building_id:
        return http_status(Forbidden, 'Device not attached to account')

    return device


@app.post(
    '/device/measurements',
    response_model=MeasurementsUploadResult,
    responses={
        BadRequest.code: {'model': BadRequest},
        Unauthorized.code: {'model': Unauthorized},
        Forbidden.code: {'model': Forbidden},
        NotFound.code: {'model': NotFound}
    }
)
def device_upload(measurements_upload: MeasurementsUpload,
                  authorization: HTTPAuthorizationCredentials = Depends(auth)):

    device_session_token = authorization.credentials

    device = crud.device_by_session(db.session, device_session_token)
    if not device:
        return http_status(Unauthorized, 'Invalid device session token')
    if not device.building_id:
        return http_status(Forbidden, 'Device not attached to account')

    valid_property_ids = {p.id for p in device.device_type.properties}
    property_ids = {item.property_id for item in measurements_upload.items}

    invalid_property_ids = property_ids - valid_property_ids
    if invalid_property_ids:
        return http_status(BadRequest, f'Invalid property identifier(s): {invalid_property_ids}')

    upload = crud.device_upload(db.session, device, measurements_upload)

    return MeasurementsUploadResult(size=upload.size, server_time=upload.server_time)
