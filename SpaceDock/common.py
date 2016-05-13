from flask import session, request, Response, abort
from flask_json import as_json_p
from flask.ext.login import current_user
from werkzeug.utils import secure_filename
from functools import wraps
import SpaceDock.database
from SpaceDock.objects import User, Permission

db = SpaceDock.database.db
Base = SpaceDock.database.Base

import json
import urllib
import requests
import xml.etree.ElementTree as ET
import re

def firstparagraph(text):
    try:
        para = text.index("\n\n")
        return text[:para + 2]
    except:
        try:
            para = text.index("\r\n\r\n")
            return text[:para + 4]
        except:
            return text

def remainingparagraphs(text):
    try:
        para = text.index("\n\n")
        return text[para + 2:]
    except:
        try:
            para = text.index("\r\n\r\n")
            return text[para + 4:]
        except:
            return ""

def dumb_object(model):
    if type(model) is list:
        return [dumb_object(x) for x in model]

    result = {}

    for col in model._sa_class_manager.mapper.mapped_table.columns:
        a = getattr(model, col.name)
        if not isinstance(a, Base):
            result[col.name] = a

    return result

def wrap_mod(mod):
    details = dict()
    details['mod'] = mod
    if len(mod.versions) > 0:
        details['latest_version'] = mod.versions[0]
        details['safe_name'] = secure_filename(mod.name)[:64]
        details['details'] = '/mod/' + str(mod.id) + '/' + secure_filename(mod.name)[:64]
        details['dl_link'] = '/mod/' + str(mod.id) + '/' + secure_filename(mod.name)[:64] + '/download/' + mod.versions[0].friendly_version
    else:
        return None
    return details

def with_session(f):
    @wraps(f)
    def go(*args, **kw):
        try:
            ret = f(*args, **kw)
            db.commit()
            return ret
        except:
            db.rollback()
            db.close()
            raise
    return go

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

#TODO: This approach has problems.
#def accessrequired(f):
#    @wraps(f)
#    def wrapper(*args, **kwargs):
#        if not current_user or current_user.confirmation or not has_access(current_user, request.path):
#            return jsonify({'error': True, 'accessErrors': 'You don\'t have the permission to access this page.'}), 401
#        else:
#            return f(*args, **kwargs)
#    return wrapper

def json_output(f):
    @wraps(f)
    def wrapper(*args, **kwargs):
        def jsonify_wrap(obj):
            jsonification = json.dumps(obj)
            return Response(jsonification, mimetype='application/json')
        
        if request.args.get('callback'):
            result = as_json_p(f)(*args, **kwargs)
        else:
            result = f(*args, **kwargs)
            
        if isinstance(result, tuple):
            return jsonify_wrap(result[0]), result[1]
        if isinstance(result, dict):
            return jsonify_wrap(result)
        if isinstance(result, list):
            return jsonify_wrap(result)

        # This is a fully fleshed out response, return it immediately
        return result

    return wrapper

def cors(f):
    @wraps(f)
    def wrapper(*args, **kwargs):
        res = f(*args, **kwargs)
        if request.headers.get('x-cors-status', False):
            if isinstance(res, tuple):
                json_text = res[0].data
                code = res[1]
            else:
                json_text = res.data
                code = 200

            o = json.loads(json_text)
            o['x-status'] = code

            return jsonify(o)

        return res

    return wrapper

def has_access(user, rule):
    if not user or not rule:
        return False
    
    # Get matching permission
    perms = Permission.query.filter(Permission.user_id == user.id).all() # I love you too sqlalchemy
    for p in perms:
        if check_permission(p.rule, p.params, rule): return True
    return False

def check_permission(rule, params, url):
    perm = Permission.query.filter(Permission.params == params).filter(Permission.rule == rule).first()
    iURL = perm.rule
    params_ = perm.get_params()
    for key in params_:
        iURL = iURL.replace('<' + key + '>', str(params_[key]).replace(', ', '|').replace('\'', '').replace('[', '(?:').replace(']', ')'))
    print(iURL)
    return not re.search('^' + iURL + '$', url) == None

def edit_object(object, patch):
    for field in patch:
        if field in dir(object):
            setattr(object, field, patch[field])

