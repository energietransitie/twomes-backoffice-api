import logging

from fastapi import Depends, FastAPI
from fastapi.responses import JSONResponse
from fastapi.security.http import HTTPBearer, HTTPAuthorizationCredentials
from fastapi_sqlalchemy import DBSessionMiddleware
from fastapi_sqlalchemy import db

from db import db_url, session_args
from firebase import firebase_dynamic_link
from schema import (
    AccountActivate, AccountCreate, AccountItem, AccountSession,
    BadRequest, DeviceActivate, DeviceCreate, DeviceItem,
    DeviceItemMeasurementTime, NotFound, Unauthorized
)
import crud

logging.basicConfig(level=logging.DEBUG)

app = FastAPI()

app.add_middleware(DBSessionMiddleware, db_url=db_url, session_args=session_args)

auth = HTTPBearer()


@app.post(
    '/account',
    response_model=AccountItem,
    responses={
        400: {'model': BadRequest}
    }
)
def account_create(account_input: AccountCreate):
    pseudonym = account_input.pseudonym

    if pseudonym:
        if crud.account_by_pseudonym(db.session, pseudonym):
            return JSONResponse(
                status_code=400,
                content={'msg': 'Account pseudonym already in use'}
            )
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
        404: {'model': NotFound}
    }
)
def account_activate(activation_token: AccountActivate):
    account = crud.account_by_token(db.session, activation_token.activation_token)
    if not account:
        return JSONResponse(
            status_code=404,
            content={'msg': 'No account found for provided activation token'}
        )

    session_token = crud.account_session_token(db.session, account)

    return AccountSession(session_token=session_token)


@app.post(
    '/device',
    response_model=DeviceItem,
    responses={
        400: {'model': BadRequest}
    }
)
def device_create(device_input: DeviceCreate):
    device_type_name = device_input.device_type
    proof_of_presence_id = device_input.proof_of_presence_id

    device_type = crud.device_type_by_name(db.session, device_type_name)
    if not device_type:
        return JSONResponse(
            status_code=400,
            content={'msg': f'Unknown device type "{device_type_name}"'}
        )

    if crud.device_by_pop(db.session, proof_of_presence_id):
        return JSONResponse(
            status_code=400,
            content={'msg': f'Proof-of-presence identifier already in use'}
        )

    device = crud.device_create(db.session, device_type, proof_of_presence_id)
    return device


@app.post(
    '/device/activate',
    response_model=DeviceItem,
    responses={
        400: {'model': BadRequest},
        401: {'model': Unauthorized},
        404: {'model': NotFound}
    }
)
def device_activate(device_activation: DeviceActivate,
                    authorization: HTTPAuthorizationCredentials = Depends(auth)):

    proof_of_presence_id = device_activation.proof_of_presence_id
    account_session_token = authorization.credentials

    device = crud.device_by_pop(db.session, proof_of_presence_id)
    if not device:
        return JSONResponse(
            status_code=404,
            content={'msg': 'No device found for provided proof-of-presence id'}
        )
    if device.activated_on:
        return JSONResponse(
            status_code=400,
            content={'msg': 'Device already activated'}
        )

    account = crud.account_by_session(db.session, account_session_token)
    if not account:
        return JSONResponse(
            status_code=401,
            content={'msg': 'Invalid account session token'}
        )

    crud.device_activate(db.session, account, device)

    return device


@app.get(
    '/device/{device_id}',
    response_model=DeviceItemMeasurementTime,
    responses={
        400: {'model': BadRequest},
        401: {'model': Unauthorized},
        404: {'model': NotFound}
    }
)
def device_read(device_id: int,
                authorization: HTTPAuthorizationCredentials = Depends(auth)):

    account_session_token = authorization.credentials

    account = crud.account_by_session(db.session, account_session_token)
    if not account:
        return JSONResponse(
            status_code=401,
            content={'msg': 'Invalid account session token'}
        )

    device = crud.device_by_account_and_id(db.session, account, device_id)
    if not device:
        return JSONResponse(
            status_code=404,
            content={'msg': f'Device {device_id} not found'}
        )
    print(device, device.id)

    timestamp = crud.device_latest_measurement_timestamp(db.session, device_id)
    device.latest_measurement_timestamp = timestamp

    return device
