from flask import session, request, Response, abort
from flask_json import as_json_p, as_json
from flask_login import current_user
from werkzeug.utils import secure_filename
from functools import wraps
from SpaceDock.database import db, Base
from SpaceDock.objects import Role, Ability
from sqlalchemy import Column

import urllib
import requests
import xml.etree.ElementTree as ET
import re
import json

def game_id(short):
    return Game.query.filter(Game.short == short).first().id

def boolean(s):
    return s.lower() in ['true', 'yes', '1', 'y', 't']

def get_param(ability, param, p):
    if ability in p.keys():
        if param in p[ability].keys():
            return p[ability][param]
    return None

def re_in(itr, value):
    if itr == None:
        return False
    for v in itr:
        if not re.match(str(v), value) == None:
            return True
    return False

def is_json(test):
    try:
        s = json.loads(test)
        return True
    except ValueError as e:
        return False

def has_ability(ability, **params): # HAX
    def dummy():
        return None
    f = user_has(ability, **params)(dummy)
    return f() == None

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

def json_output(f):
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
            if '__lock__' in dir(object):
                if field in getattr(object, '__lock__') or field == '__lock__': # We might want a function to report theese guys
                    continue
            if not type(getattr(object, field)) == Column:
                continue
            if isinstance(getattr(object, field), (int, bool, str, float)):
                setattr(object, field, patch[field])
            else:
                setattr(object, field, edit_object(getattr(object, field), patch[field]))
    return object

def user_has(ability, **params):
    def wrapper(func):
        @wraps(func)
        def inner(*args, **kwargs):
            desired_ability = Ability.query.filter(Ability.name == ability).first()
            user_abilities = []
            user_params = {}
            if not current_user:
                return {'error': True, 'reasons': ['You need to be logged in to access this page']}, 400
            for role in current_user._roles:
                user_abilities += role.abilities
                user_params += json.loads(role.params)
            has = True
            if desired_ability in user_abilities:
                if 'params' in params:
                    for p in params['params']:
                        if not re_in(get_param(ability, p, user_params), kwargs[p]) and not re_in(get_param(ability, p, user_params), request.form.get(p)):
                            has = False
                if has:
                    return func(*args, **kwargs)
                else:
                    return {'error': True, 'reasons': ['You don\'t have access to this page. You need to have the abilities: ' + ability]}, 400
            else:
                return {'error': True, 'reasons': ['You don\'t have access to this page. You need to have the abilities: ' + ability]}, 400
        return inner
    return wrapper

# This is not supported atm
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
