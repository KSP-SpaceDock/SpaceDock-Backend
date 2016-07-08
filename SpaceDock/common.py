from flask import session, request, Response, abort
from flask_json import as_json_p, as_json
from flask.ext.login import current_user
from werkzeug.utils import secure_filename
from functools import wraps
from SpaceDock.database import db, Base
from SpaceDock.objects import User, Permission, Game

import urllib
import requests
import xml.etree.ElementTree as ET
import re

def game_id(short):
    return Game.query.filter(Game.short == short).first().id

def with_session(f):
    @wraps(f)
    def wrapper(*args, **kw):
        try:
            ret = f(*args, **kw)
            db.commit()
            return ret
        except:
            db.rollback()
            db.close()
            raise
    return wrapper

def json(f):
    @wraps(f)
    def wrapper(*fargs, **fkwargs):
        if request.args.get('callback'):
            return as_json_p(f)(*fargs, **fkwargs)
        else:
            return as_json(f)(*fargs, **fkwargs)
    return wrapper

def loginrequired(f):
    @wraps(f)
    def wrapper(*args, **kwargs):
        if not current_user or current_user.confirmation:
            return {'error': True, 'accessErrors': 'You need to be logged in to access this page.'}, 401
        else:
            return f(*args, **kwargs)
    return wrapper

def adminrequired(f):
    @wraps(f)
    def wrapper(*args, **kwargs):
        if not current_user or current_user.confirmation or not current_user.admin:
            return {'error': True, 'accessErrors': 'You don\'t have the permission to access this page.'}, 401
        else:
            return f(*args, **kwargs)
    return wrapper

def edit_object(object, patch):
    for field in patch:
        if field in dir(object):
            setattr(object, field, patch[field])

