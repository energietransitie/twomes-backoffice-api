from datetime import datetime, timezone
from typing import Optional
from secrets import token_urlsafe
import random

from db import Session
from model import (
    Account,
    Building,
    Device,
    DeviceType,
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


def building_create(db: Session, account: Account, location: AccountLocation = None) -> Building:
    """
    Create a new Building
    """
    building = Building(account=account)
    if location:
        building.latitude = location.latitude
        building.longitude = location.longitude
    db.add(building)

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
    )
    db.add(device)

    return device