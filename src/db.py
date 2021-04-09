import os

from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker

db_url_env = os.getenv('TWOMES_DB_URL')
assert db_url_env, 'Environment variable TWOMES_DB_URL not set. Format: user:pass@host/db'

db_url = f'mysql+pymysql://{db_url_env}?charset=utf8mb4'

engine = create_engine(db_url)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

Base = declarative_base()
