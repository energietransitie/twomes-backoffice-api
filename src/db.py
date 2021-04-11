import os

from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import Session

db_url_env = os.getenv('TWOMES_DB_URL')
assert db_url_env, 'Environment variable TWOMES_DB_URL not set. Format: user:pass@host/db'

db_url = f'mysql+pymysql://{db_url_env}?charset=utf8mb4'

engine = create_engine(db_url)

Base = declarative_base()


class AtomicSession(Session):

    def __init__(self):
        super().__init__(autocommit=False, autoflush=True, bind=engine)

    def __enter__(self):
        self.begin()
        return self

    def __exit__(self, exc_type, exc_value, traceback):
        # In case of an exception, no commit takes place, and an implicit
        # rollback will be done. If no exception, simply commit.
        if not exc_type:
            self.commit()

        self.close()
