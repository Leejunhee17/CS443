from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base

Base = declarative_base()

from .block import Block
from .transaction import Transaction


def create_database(connection_string: str):
    engine = create_engine(connection_string)
    Base.metadata.create_all(engine)
