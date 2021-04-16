from fastapi import FastAPI
from fastapi.responses import JSONResponse

from db import AtomicSession
from firebase import firebase_dynamic_link
from schema import (
    AccountCreate, AccountInfo, BadRequest,
    DeviceCreate, Device
)
import crud

app = FastAPI()



@app.post(
    '/account',
    response_model=AccountInfo,
    responses={
        400: {'model': BadRequest}
    }
)
async def account_create(account_input: AccountCreate):
    pseudonym = account_input.pseudonym

    with AtomicSession() as db:
        if pseudonym:
            if crud.account_by_pseudonym(db, pseudonym):
                return JSONResponse(
                    status_code=400,
                    content={'message': 'Account pseudonym already in use'}
                )
        else:
            pseudonym = crud.generate_pseudonym(db)

        account = crud.account_create(db, pseudonym)

        location = account_input.location if account_input.location else None
        crud.building_create(db, account, location)

    url = firebase_dynamic_link(account.activation_token)

    return AccountInfo(pseudonym=pseudonym, firebase_url=url)


@app.post(
    '/device',
    response_model=Device,
    responses={
        400: {'model': BadRequest}
    }
)
async def device_create(device_input: DeviceCreate):
    device_type_name = device_input.device_type
    proof_of_presence_id = device_input.proof_of_presence_id

    with AtomicSession() as db:
        device_type = crud.device_type_by_name(db, device_type_name)
        if not device_type:
            return JSONResponse(
                status_code=400,
                content={'message': f'Unknown device type "{device_type_name}"'}
            )

        device = crud.device_create(db, device_type, proof_of_presence_id)
        if not device:
            return JSONResponse(
                status_code=400,
                content={'message': f'Proof-of-presence identifier already in use'}
            )

        return device
