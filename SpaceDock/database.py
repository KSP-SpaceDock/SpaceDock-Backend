from sqlalchemy import create_engine
from sqlalchemy.orm import scoped_session, sessionmaker
from sqlalchemy.ext.declarative import declarative_base

#Database used in wrapper functions in common.py
db = None
#Python doesn't support passing types into modules so this shit has to be singleton in the module
Base = None

class Database:
    def __init__(self, cfg):
        self.cfg = cfg
        self.engine = create_engine(self.cfg['connection-string'])
        self.db = scoped_session(sessionmaker(autocommit=False, autoflush=False, bind=self.engine))
        global db
        if not db:
            db = self.db
        global Base
        if not Base:
            Base = declarative_base()
            Base.query = self.db.query_property()

    def init_db(self):
        Base.metadata.create_all(bind=engine)
        
    def get_database(self):
        return self.db