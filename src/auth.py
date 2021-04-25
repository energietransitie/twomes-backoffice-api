from typing import Tuple
import secrets

from fastapi.security.http import HTTPBearer
from jwt.utils import base64url_decode, base64url_encode
from passlib.context import CryptContext

pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")


def session_token_generate(identifier: int) -> Tuple[str, str]:
    """
    Generate a session token consisting of the following:
        <base64 encode identifier>.<random token>

    Returns both the token and the hash of the token
    """
    id_base64 = base64url_encode(str(identifier).encode()).decode('ascii')
    token = secrets.token_urlsafe(32)

    session_token = f'{id_base64}.{token}'
    session_token_hash = pwd_context.hash(session_token)

    return session_token, session_token_hash


def session_token_parts(session_token) -> Tuple[int, str]:
    """
    Split the session token, in an identifier and secret token part
    """
    parts = session_token.split('.')
    if len(parts) != 2:
        raise ValueError('Illegal session token')

    identifier = int(base64url_decode(parts[0]))
    token = parts[1]

    return identifier, token


def session_token_verify(session_token, session_token_hash) -> bool:
    return pwd_context.verify(session_token, session_token_hash)


class AccountSessionTokenBearer(HTTPBearer):
    pass


class AdminSessionTokenBearer(HTTPBearer):
    pass


class DeviceSessionTokenBearer(HTTPBearer):
    pass

