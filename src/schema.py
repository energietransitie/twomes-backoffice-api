from typing import Optional

from pydantic import BaseModel, condecimal, conint, constr

from model import Account, Building


class BadRequest(BaseModel):  # 400
    msg: str


class Unauthorized(BaseModel):  # 401
    msg: str


class NotFound(BaseModel):  # 404
    msg: str


class AccountLocation(BaseModel):
    longitude: condecimal(
        max_digits=Building.LOCATION_MAX_DIGITS,
        decimal_places=Building.LOCATION_DECIMAL_PLACES
    )
    latitude: condecimal(
        max_digits=Building.LOCATION_MAX_DIGITS,
        decimal_places=Building.LOCATION_DECIMAL_PLACES
    )


class AccountCreate(BaseModel):
    pseudonym: Optional[conint(ge=Account.PSEUDONYM_MIN, le=Account.PSEUDONYM_MAX)] = None
    location: Optional[AccountLocation] = None


class AccountItem(BaseModel):
    id: int
    pseudonym: int
    firebase_url: str


class AccountActivation(BaseModel):
    activation_token: str


class AccountSession(BaseModel):
    session_token: str


class DeviceCreate(BaseModel):
    device_type: str
    proof_of_presence_id: constr(strip_whitespace=True, min_length=8, max_length=1024)


class DeviceActivation(BaseModel):
    proof_of_presence_id: constr(strip_whitespace=True, min_length=8, max_length=1024)


class DeviceItem(BaseModel):
    id: int
    device_type: str


class DeviceTypeItem(BaseModel):
    device_type: str
    installation_manual_url: str

