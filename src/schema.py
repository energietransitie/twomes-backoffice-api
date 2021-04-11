from typing import Optional

from pydantic import BaseModel, condecimal, conint

from model import Account, Building


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
