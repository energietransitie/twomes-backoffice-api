import logging
from datetime import datetime, timezone
from typing import Optional
from secrets import token_urlsafe
import random

from sqlalchemy import desc, select
from sqlalchemy.orm import Session

from auth import session_token_generate, session_token_parts, session_token_verify
from model import (
    Account,
    Building,
    Device,
    DeviceType,
    Measurement,
)
from schema import AccountLocation


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

    pseudonym = random.choice(list(free))
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

    except Exception as e:
        logging.info(e)
        return None

    if not session_token_verify(session_token, account.session_token_hash):
        logging.info('Invalid session token')
        return None

    return account


def building_create(db: Session, account: Account, location: AccountLocation = None) -> Building:
    """
    Create a new Building
    """
    building = Building(account=account)
    if location:
        building.latitude = location.latitude
        building.longitude = location.longitude

    db.add(building)
    db.commit()
    db.refresh(building)

    return building


def device_type_by_name(db: Session, name: str) -> Optional[DeviceType]:
    """
    Get DeviceType instance by name
    """
    return db.query(DeviceType).filter(DeviceType.name == name).one_or_none()


def device_create(db: Session, device_type: DeviceType, proof_of_presence_id: str) -> Device:
    """
    Create a new Device
    """
    device = Device(
        device_type=device_type,
        proof_of_presence_id=proof_of_presence_id,
        created_on=datetime.now(timezone.utc),
    )

    db.add(device)
    db.commit()
    db.refresh(device)

    return device


def device_by_pop(db: Session, proof_of_presence_id: str) -> Optional[Device]:
    """
    Get Device instance by proof-of-presence identifier
    """
    query = db.query(Device).filter(Device.proof_of_presence_id == proof_of_presence_id)
    return query.one_or_none()


def device_activate(db: Session, account: Account, device: Device):
    """
    Active the device by assigning it to the building of an account.
    """
    device.building = account.building
    device.activated_on = datetime.now(timezone.utc)
    db.commit()


def device_by_account_and_id(db: Session, account: Account, device_id: int) -> Optional[Device]:
    """
    Get Device instance for an Account
    """
    query = select(Device).join(Device.building).filter(
        Device.id == device_id,
        Building.account_id == account.id
    )
    device = db.execute(query).scalars().one_or_none()

    return device


def device_latest_measurement_timestamp(db: Session, device_id: int) -> datetime:
    """
    Get the timestamp of the most recent Measurement for Device
    """
    measurements = select(Measurement).filter(Measurement.device_id == device_id)
    measurements = measurements.order_by(desc(Measurement.timestamp))
    measurement = db.execute(measurements).scalars().first()

    return measurement.timestamp if measurement else None
