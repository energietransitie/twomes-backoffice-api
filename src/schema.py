from datetime import timedelta
from decimal import Decimal
from enum import Enum
from typing import ClassVar, List, Optional

from pydantic import BaseModel, condecimal, conint, constr

from field import Datetime, Timezone
from model import Account, Building


class HttpStatus(BaseModel):
    code: ClassVar[int]
    detail: str


class BadRequest(HttpStatus):
    code = 400


class Unauthorized(HttpStatus):
    code = 401


class Forbidden(HttpStatus):
    code = 403


class NotFound(HttpStatus):
    code = 404


class BuildingLocation(BaseModel):
    longitude: condecimal(
        max_digits=Building.LOCATION_MAX_DIGITS,
        decimal_places=Building.LOCATION_DECIMAL_PLACES
    )
    latitude: condecimal(
        max_digits=Building.LOCATION_MAX_DIGITS,
        decimal_places=Building.LOCATION_DECIMAL_PLACES
    )


class BuildingDefaults:
    """
    Building defaults for the 50 Tinten Groen Assendorp project. Using an approximate
    location (center of Assendorperplein, Zwolle), to avoid pinpointing individual homes.
    """
    location = BuildingLocation(
        latitude=Decimal('52.50655'),
        longitude=Decimal('6.09961')
    )
    tz_name = 'Europe/Amsterdam'


class AccountCreate(BaseModel):
    pseudonym: Optional[conint(
        ge=Account.PSEUDONYM_MIN, le=Account.PSEUDONYM_MAX)] = None
    location: Optional[BuildingLocation] = BuildingDefaults.location
    tz_name: Optional[Timezone] = BuildingDefaults.tz_name


class AccountItem(BaseModel):
    id: int
    pseudonym: int
    activation_token: str
    firebase_url: str


class AccountActivate(BaseModel):
    activation_token: str


class SessionToken(BaseModel):
    session_token: str


class AccountSession(SessionToken):
    pass


class DeviceSession(SessionToken):
    pass


class DeviceCreate(BaseModel):
    name: constr(regex=r'TWOMES-([0-9A-F]){6}/')
    device_type: str


class DeviceTypeItem(BaseModel):
    name: str
    display_name: str
    installation_manual_url: str

    class Config:
        orm_mode = True


class PropertyCompleteItem(BaseModel):
    id: int
    name: str
    unit: Optional[str]

    class Config:
        orm_mode = True


class DeviceTypeCompleteItem(BaseModel):
    id: int
    name: str
    installation_manual_url: str
    properties: List[PropertyCompleteItem]

    class Config:
        orm_mode = True


class DeviceItem(BaseModel):
    name: str
    device_type_name: str
    activation_token: str

    class Config:
        orm_mode = True


class DeviceCompleteItem(BaseModel):
    id: int
    name: str
    device_type: DeviceTypeItem
    activation_token: str
    created_on: Datetime
    activated_on: Optional[Datetime]

    class Config:
        orm_mode = True

class DeviceVerify(BaseModel):
    activation_token: constr(strip_whitespace=True, min_length=8, max_length=1024)

class DeviceItemMeasurementTime(DeviceItem):
    latest_measurement_timestamp: Optional[Datetime]

    class Config:
        orm_mode = True


class MeasurementValue(BaseModel):
    timestamp: Datetime
    value: str


class TimestampType(str, Enum):
    start = 'start'
    end = 'end'


class PropertyMeasurementsFixed(BaseModel):
    property_name: str
    timestamp: Datetime
    timestamp_type: TimestampType
    interval: timedelta
    measurements: List[str]


class PropertyMeasurementsVariable(BaseModel):
    property_name: str
    measurements: List[MeasurementValue]


class MeasurementsUploadFixed(BaseModel):
    upload_time: Datetime
    property_measurements: List[PropertyMeasurementsFixed]


class MeasurementsUploadVariable(BaseModel):
    upload_time: Datetime
    property_measurements: List[PropertyMeasurementsVariable]


class MeasurementsUploadResult(BaseModel):
    server_time: Datetime
    size: int
