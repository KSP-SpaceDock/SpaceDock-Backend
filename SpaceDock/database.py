from sqlalchemy import create_engine, Column, String
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker, scoped_session
from SpaceDock.config import cfg

import json

engine = create_engine(cfg['connection-string'], pool_size=20, max_overflow=100)
db = scoped_session(sessionmaker(autocommit=False, autoflush=False, bind=engine))

# Base class for database objects, to support metadata (plugins)
class MetaObject():
    meta = Column(String(512), server_default='{}')

    # Adds a new metadata, or updates an existing one
    def __getitem__(self, key):
        """
        Returns a metadata value
        """
        if not test_is_json(self.meta):
            data = '{}'
        data = json.loads(self.meta)
        if key in data:
            return data[key]
        return None

    def __setitem__(self, key, value):
        """
        Sets a metadata value
        """
        if not test_is_json(self.meta):
            data = '{}'
        data = json.loads(self.meta)
        data[key] = value
        if value == None:
            del data[key]
        self.meta = json.dumps(data)

def test_is_json(test):
    """
    Checks whether something is JSON formatted
    """
    try:
        json.loads(test)
        return True
    except ValueError as e:
        return False
        
Base = declarative_base()
Base.query = db.query_property()

def init_db():
    import SpaceDock.objects
    Base.metadata.create_all(bind=engine)