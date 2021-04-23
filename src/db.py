import logging
import os
import threading

from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import Query, Session

logger = logging.getLogger(__name__)

db_url_env = os.getenv('TWOMES_DB_URL')
assert db_url_env, 'Environment variable TWOMES_DB_URL not set. Format: user:pass@host/db'

db_url = f'mysql+pymysql://{db_url_env}?charset=utf8mb4'

session_args = {
    'autocommit': False,
    'autoflush': True,
}

Base = declarative_base()


# Wrapper functions for SQLAlchemy Session, to follow session creation
session_init = Session.__init__
session_close = Session.close


def session_init_wrap(self, *args, **kwargs):
    logger.debug(f'{threading.get_ident()}: session create')
    session_init(self, *args, **kwargs)


def session_close_wrap(self):
    logger.debug(f'{threading.get_ident()}: session close')
    session_close(self)


# Patch the SQLAlchemy Session __init__ function
Session.__init__  = session_init_wrap
Session.close = session_close_wrap
