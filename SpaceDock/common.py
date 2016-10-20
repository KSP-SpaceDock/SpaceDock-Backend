from flask import make_response, request
from flask_json import as_json, as_json_p
from flask_login import current_user
from functools import wraps
from SpaceDock.config import cfg
from SpaceDock.database import db
from SpaceDock.objects import Ability, Game
from wsgiref.handlers import format_date_time

import datetime
import json
import re
import time

def with_session(f):
    """
    Executes a function using a Database session
    """
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
    """
    Output transformer that returns a JSON formatted string
    """
    @wraps(f)
    def wrapper(*args, **kwargs):
        if request.args.get('callback'):
            return as_json_p(f)(*args, **kwargs)
        else:
            return as_json(f)(*args, **kwargs)
    return wrapper

# Return codes:
#   0: Everything is ok
#   1: Tried to edit a field that is listed in __lock__
#   2: Tried to patch a field that doesnt exist
#   3: Invalid type
def edit_object(object, patch):
    """
    Edits an object using a patch dictionary. Edits only fields that aren't listed in __lock__
    """
    patched_patch = {}
    for field in patch.keys():
        if field in dir(object):
            if '__lock__' in dir(object) and field in getattr(object, '__lock__') or field == '__lock__':
                return 1
        else:
            return 2
        if isinstance(getattr(object, field), (int, bool, str, float)):
            patched_patch[field] = patch[field]
        else:
            o = getattr(object, field)
            code = edit_object(o, patch[field])
            if not code == 0:
                return code
            patched_patch[field] = o
    for field2 in patched_patch:
        setattr(object, field, patched_patch[field])
    try:
        db.commit()
    except:
        db.rollback()
        return 3
    return 0

def user_has(ability, public=True, params=None):
    """
    Checks whether the user has the ability to view this site. Decorator function
    """
    def wrapper(func):
        @wraps(func)
        def inner(*args, **kwargs):
            # Check if the user is logged in
            if not current_user:
                return {'error': True, 'reasons': ['You need to be logged in to access this page'], 'codes': ['1035']}, 403
            if public and not current_user.public:
                return {'error': True, 'reasons': ['Only users with public profiles may access this page.'], 'codes': ['1000']}, 403

            # Get the specified ability
            desired_ability = Ability.query.filter(Ability.name == ability).first()
            user_abilities = []
            for role in current_user._roles:
                for ability_ in role.abilities:
                    user_abilities.append(ability_)
            user_params = {}
            for role in current_user._roles:
                user_params.update(json.loads(role.params))

            # Check whether the abilities match
            has = False
            if desired_ability in user_abilities:
                if params:
                    for p in params:
                        if re_in(get_param(ability, p, user_params), request.json.get(p)) or re_in(get_param(ability, p, user_params), kwargs.get(p)):
                            has = True
                else:
                    has = True
                if has:
                    return func(*args, **kwargs)
            return {'error': True, 'reasons': ['You don\'t have access to this page. You need to have the abilities: ' + ability], 'codes': ['1020']}, 403
        return inner

    # Make sure the ability exists
    desired_ability = Ability.query.filter(Ability.name == ability).first()
    if not desired_ability:
        desired_ability = Ability(ability)
        db.add(desired_ability)
        db.commit()

    return wrapper

def has_ability(ability, **params): # HAX
    """
    Checks whether the user has the ability to view this site.
    """
    def dummy():
        return None
    f = user_has(ability, **params)(dummy)
    return f() == None

def game_id(short):
    """
    Converts a game ID into a Gameshort
    """
    if not Game.query.filter(Game.short == short).first():
        return None
    return Game.query.filter(Game.short == short).first().id

def boolean(s):
    """
    Converts string to bool
    """
    if s == None:
        return False
    return str(s).lower() in ['true', 'yes', '1', 'y', 't']

def get_param(ability, param, p):
    """
    Gets the parameters for ability and param.
    """
    if ability in p.keys():
        if param in p[ability].keys():
            return p[ability][param]
    return None

def re_in(itr, value):
    """
    Check whether a value is in a list using regex
    """
    if itr == None:
        return False
    if value == None:
        return False
    for v in itr:
        if not re.match(str(v), value) == None:
            return True
    return False

def is_json(test):
    """
    Checks whether something is JSON formatted
    """
    try:
        json.loads(test)
        return True
    except ValueError as e:
        return False

def redirect(url, status=301):
    return None, status, {'Location': url}

def clamp_number(min_value, max_value, num):
	"""
	Clamps a number between a minimum and maximum value
	"""
	try:
		result = max(min(num, max_value), min_value)
		return result
	except TypeError as e:
		return -1

def cache(f):
    """
    Caches a response to reduce the load for the server
    """
    @wraps(f)
    def cache_func(*args, **kwargs): 
        now = datetime.datetime.now()
 
        response = make_response(f(*args, **kwargs))
        response.headers['Last-Modified'] = format_date_time(time.mktime(now.timetuple()))
            
        expires_time = now + datetime.timedelta(seconds=cfg.geti('cache-expires') * 60)

        response.headers['Cache-Control'] = 'public'
        response.headers['Expires'] = format_date_time(time.mktime(expires_time.timetuple()))
 
        return response
    return cache_func