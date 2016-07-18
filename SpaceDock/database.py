from sqlalchemy import create_engine
from sqlalchemy.orm import scoped_session, sessionmaker
from sqlalchemy.ext.declarative import declarative_base
from SpaceDock.app import cfg

engine = create_engine(cfg['connection-string'])
db = scoped_session(sessionmaker(autocommit=False, autoflush=False, bind=engine))

Base = declarative_base()
Base.query = db.query_property()

def init_db():
    import SpaceDock.objects
    Base.metadata.create_all(bind=engine)