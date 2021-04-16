from typing import Optional

from pydantic import BaseModel, condecimal, conint, constr

from model import Account, Building


class BadRequest(BaseModel):
    message: str


class AccountLocation(BaseModel):
    longitude: condecimal(
        max_digits=Building.LOCATION_SCALE,
        decimal_places=Building.LOCATION_PRECISION
    )
    latitude: condecimal(
        max_digits=Building.LOCATION_SCALE,
        decimal_places=Building.LOCATION_PRECISION
    )


class AccountCreate(BaseModel):
    pseudonym: Optional[conint(ge=Account.PSEUDONYM_MIN, le=Account.PSEUDONYM_MAX)] = None
    location: Optional[AccountLocation] = None


class AccountInfo(BaseModel):
    pseudonym: int
    firebase_url: str


class DeviceCreate(BaseModel):
    device_type: str
    proof_of_presence_id: constr(strip_whitespace=True, min_length=8, max_length=1024)


class Device(BaseModel):
    id: int
    device_type: str

    class Config:
        orm_mode = True
