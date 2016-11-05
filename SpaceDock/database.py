from flask_sqlalchemy import SQLAlchemy
from sqlalchemy import Column, String
from SpaceDock.config import cfg

import json

connection = SQLAlchemy()

# Base class for database queries        
Base = connection.Model

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

db = connection.session

def init_db(app):
    
    # Configure database
    app.config['SQLALCHEMY_DATABASE_URI'] = cfg['connection-string']
    app.config['SQLALCHEMY_POOL_SIZE'] = cfg.geti('db-pool-size')
    app.config['SQLALCHEMY_POOL_RECYCLE'] = cfg.geti('db-pool-recycle')
    app.config['SQLALCHEMY_COMMIT_ON_TEARDOWN'] = False  # We continue to use @with_session, because that doesn't waste resources on routes that dont edit the db
    connection.init_app(app)
    connection.app = app

    import SpaceDock.objects
    connection.create_all()