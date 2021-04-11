from datetime import datetime, timezone
from random import randint
from secrets import token_urlsafe
from typing import Set

from fastapi import FastAPI

from db import AtomicSession
from model import Account, Building
from schema import AccountCreate, AccountInfo

app = FastAPI()


def generate_pseudonym(taken: Set[int] = None) -> int:
    """
    Generate an Account pseudonym (int), not in taken.
    """
    max_values = Account.PSEUDONYM_MAX - Account.PSEUDONYM_MIN
    assert len(taken) < 0.5 * max_values, 'Too few free pseudonyms'

    while True:
        pseudonym = randint(Account.PSEUDONYM_MIN, Account.PSEUDONYM_MAX)
        if pseudonym not in taken:
            break

    return pseudonym


@app.post('/account')
async def account_create(account_input: AccountCreate) -> AccountInfo:
    pseudonym = account_input.pseudonym
    activation_token = token_urlsafe(32)

    with AtomicSession() as session:
        if not pseudonym:
            existing_pseudonyms = {r[0] for r in session.query(Account.pseudonym).all()}
            pseudonym = generate_pseudonym(existing_pseudonyms)

        account = Account(
            pseudonym= pseudonym,
            created_on=datetime.now(timezone.utc),
            activation_token=activation_token
        )
        session.add(account)

        building = Building(
            account=account
        )
        if account_input.location:
            building.latitude = account_input.location.latitude
            building.longitude = account_input.location.longitude
        session.add(building)

    url = f'http://bla/firebase/{activation_token}'

    return AccountInfo(pseudonym=pseudonym, firebase_url=url)
