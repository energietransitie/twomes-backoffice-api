import logging
import secrets
from datetime import (
    datetime,
    timezone,
)
from typing import (
    List,
    Optional,
)
from secrets import token_urlsafe

from sqlalchemy import (
    desc,
    select,
)
from sqlalchemy.orm import Session

from auth import (
    session_token_generate,
    session_token_parts,
    session_token_verify,
)
from model import (
    Account,
    Building,
    Device,
    DeviceType,
    Measurement,
    Property,
    Upload,
)
from schema import (
    BuildingLocation,
    MeasurementsUploadFixed,
    MeasurementsUploadVariable,
    MeasurementValue,
    PropertyMeasurementsFixed,
    PropertyMeasurementsVariable,
    TimestampType,
)


def generate_pseudonym(db: Session) -> int:
    """
    Generate an Account pseudonym not already in use.
    """
    max_values = Account.PSEUDONYM_MAX - Account.PSEUDONYM_MIN

    # We hold the full range of possible values in memory - put a cap
    # on the maximum number of pseudonyms, to not run out of memory.
    assert max_values < 1000000, 'Pseudonym range too large'

    taken = {r[0] for r in db.query(Account.pseudonym).all()}

    free = set(range(Account.PSEUDONYM_MIN, Account.PSEUDONYM_MAX + 1)) - taken
    assert free, 'No more pseudonyms left'

    pseudonym = secrets.choice(list(free))
    return pseudonym


def account_create(db: Session, pseudonym: int) -> Account:
    """
    Create a new Account
    """
    created_on = datetime.now(timezone.utc)
    activation_token = token_urlsafe(32)

    account = Account(
        pseudonym=pseudonym,
        created_on=created_on,
        activation_token=activation_token
    )

    db.add(account)
    db.commit()
    db.refresh(account)

    return account


def account_by_token(db: Session, activation_token: str) -> Optional[Account]:
    """
    Get Account by activation token
    """
    return db.query(Account).filter(Account.activation_token == activation_token).one_or_none()


def account_by_pseudonym(db: Session, pseudonym: int) -> Optional[Account]:
    """
    Get Account by pseudonym
    """
    return db.query(Account).filter(Account.pseudonym == pseudonym).one_or_none()


def account_session_token(db: Session, account: Account) -> str:
    """
    Get session token for an account.

    An account session token consists of two parts:
        <base64_encoded account id>.<random token>

    The hash of the session token is stored in the database.
    """
    session_token, session_token_hash = session_token_generate(account.id)

    account.session_token_hash = session_token_hash
    account.activated_on = datetime.now(timezone.utc)
    db.commit()

    return session_token


def account_by_session(db: Session, session_token: str) -> Optional[Account]:
    """
    Get Account by an account session token
    """
    try:
        account_id, _ = session_token_parts(session_token)
        account: Account = db.get(Account, account_id)

        # account may be None, in case of invalid account_id

    except Exception as e:
        logging.info(e)
        return None

    if not account or not session_token_verify(session_token, account.session_token_hash):
        logging.info('Invalid session token')
        return None

    if account.activation_token:
        # Unset the activation token, at first usage of the session token
        account.activation_token = None
        db.commit()

    return account


def building_create(db: Session,
                    account: Account,
                    location: BuildingLocation,
                    tz_name: str) -> Building:
    """
    Create a new Building
    """
    building = Building(account=account)
    building.latitude = location.latitude
    building.longitude = location.longitude
    building.tz_name = tz_name

    db.add(building)
    db.commit()
    db.refresh(building)

    return building


def device_type_by_name(db: Session, name: str) -> Optional[DeviceType]:
    """
    Get DeviceType instance by name
    """
    return db.query(DeviceType).filter(DeviceType.name == name).one_or_none()


def device_create(db: Session,
                  name: str,
                  device_type: DeviceType) -> Device:
    """
    Create a new Device
    """
    device = Device(
        name=name,
        device_type=device_type,
        created_on=datetime.now(timezone.utc),
    )

    db.add(device)
    db.commit()
    db.refresh(device)

    return device


def device_activate(db: Session, account: Account, device: Device):
    """
    Active the device by assigning it to the building of an account.
    """
    device.building = account.building
    device.activated_on = datetime.now(timezone.utc)
    db.commit()


def device_by_name(db: Session, name: str) -> Optional[Device]:
    """
    Get Device instance by name
    """
    query = db.query(Device).filter(Device.name == name)
    return query.one_or_none()


def device_latest_measurement_timestamp(db: Session, device_id: int) -> datetime:
    """
    Get the timestamp of the most recent Measurement for Device
    """
    measurements = select(Measurement).filter(
        Measurement.device_id == device_id)
    measurements = measurements.order_by(desc(Measurement.timestamp))
    measurement = db.execute(measurements).scalars().first()

    return measurement.timestamp if measurement else None


def device_session_token(db: Session, device: Device) -> str:
    """
    Get session token for a device.

    A device session token consists of two parts:
        <base64_encoded device id>.<random token>

    The hash of the session token is stored in the database.
    """
    session_token, session_token_hash = session_token_generate(device.id)

    device.session_token_hash = session_token_hash
    db.commit()

    return session_token


def device_by_session(db: Session, session_token: str) -> Optional[Device]:
    """
    Get Device by a device session token
    """
    try:
        device_id, _ = session_token_parts(session_token)
        device: Device = db.get(Device, device_id)

        # device may be None, in case of invalid device_id

    except Exception as e:
        logging.info(e)
        return None

    if not device or not session_token_verify(session_token, device.session_token_hash):
        logging.info('Invalid session token')
        return None

    return device


def measurements_fixed_to_variable(data: PropertyMeasurementsFixed) -> PropertyMeasurementsVariable:
    size = len(data.measurements)
    delta = data.interval

    if data.timestamp_type == TimestampType.start:
        start = data.timestamp
    elif data.timestamp_type == TimestampType.end:
        start = data.timestamp - (size - 1) * delta
    else:
        raise ValueError(f'Illegal timestamp type "{data.timestamp_type}"')

    timestamps = [start + n * delta for n in range(size)]

    result = PropertyMeasurementsVariable(
        property_name=data.property_name,
        measurements=[
            MeasurementValue(timestamp=i[0], value=i[1])
            for i in zip(timestamps, data.measurements)
        ]
    )
    return result


def upload_fixed_to_variable(data: MeasurementsUploadFixed) -> MeasurementsUploadVariable:
    result = MeasurementsUploadVariable(
        upload_time=data.upload_time,
        property_measurements=[
            measurements_fixed_to_variable(pm)
            for pm in data.property_measurements
        ]
    )
    return result


def device_upload_variable(db: Session,
                           device: Device,
                           properties: List[Property],
                           data: MeasurementsUploadVariable) -> Upload:
    """
    Save measurements data from a single upload of a device,  with variable
    measurement timestamps. Property names in the upload data must be valid.
    """
    server_time = datetime.now(timezone.utc)
    size = sum([len(l.measurements) for l in data.property_measurements])

    upload = Upload(
        device=device,
        server_time=server_time,
        device_time=data.upload_time,
        size=size,
    )

    db.add(upload)
    db.commit()
    db.refresh(upload)

    property_ids = {p.name: p.id for p in properties}

    measurements = [
        {
            'device_id': device.id,
            'property_id': property_ids[item.property_name],
            'upload_id': upload.id,
            'timestamp': measurement.timestamp,
            'value': measurement.value,
        }
        for item in data.property_measurements for measurement in item.measurements
    ]

    db.bulk_insert_mappings(Measurement, measurements)
    db.commit()

    return upload
