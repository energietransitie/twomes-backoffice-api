from sqlalchemy import create_engine
from sqlalchemy.orm import scoped_session, sessionmaker

from db import db_url

engine = create_engine(db_url)
session = scoped_session(sessionmaker(bind=engine))
