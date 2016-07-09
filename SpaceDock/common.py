from flask import session, request, Response, abort
from flask_json import as_json_p, as_json
from flask_login import current_user
from werkzeug.utils import secure_filename
from functools import wraps
from SpaceDock.database import db, Base
from SpaceDock.objects import Role, Ability

import urllib
import requests
import xml.etree.ElementTree as ET
import re

def game_id(short):
    return Game.query.filter(Game.short == short).first().id

def boolean(s):
    return s.lower() in ['true', 'yes', '1', 'y', 't']

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

def user_has(*ability):
    def wrapper(func):
        @wraps(func)
        def inner(*args, **kwargs):
            desired_ability = Ability.query.filter_by(Ability.name == ability).first()
            user_abilities = []
            if not current_user:
                return {'error': True, 'reasons': ['You need to be logged in to access this page']}, 400
            for role in current_user._roles:
                user_abilities += role.abilities
            has = True
            for a in ability:
                if not a in user_abilities:
                    has = False
            if has:
                return func(*args, **kwargs)
            else:
                return {'error': True, 'reasons': ['You don\'t have access to this page. You need to have the abilities: ' + ','.join(ability)]}, 400
        return inner
    return wrapper

def user_is(*role):
    def wrapper(func):
        @wraps(func)
        def inner(*args, **kwargs):
            if not current_user:
                return {'error': True, 'reasons': ['You need to be logged in to access this page']}, 400
            has = True
            for r in role:
                if not r in current_user.roles:
                    has = False
            if has:
                return func(*args, **kwargs)
            return {'error': True, 'reasons': ['You don\'t have access to this page You need to have the roles: ' + ','.join(role)]}, 400
        return inner
    return wrapper
