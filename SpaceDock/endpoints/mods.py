from sqlalchemy import desc
from flask_login import current_user
from werkzeug.utils import secure_filename
from SpaceDock.common import *
from SpaceDock.config import cfg
from SpaceDock.database import db
from SpaceDock.formatting import mod_info
from SpaceDock.objects import *
from SpaceDock.routing import route

import json
import os
import time


@route('/api/mods')
def mod_list():
    """
    Returns a list of all mods
    """
    results = list()
    for mod in Mod.query.order_by(desc(Mod.id)).filter(Mod.published):
        results.append(mod_info(mod))
    return {'error': False, 'count': len(results), 'data': results}

@route('/api/mods/<modid>')
def mod_info(modid):
    """
    Returns information for one mod
    """
    if not modid.isdigit() or not Mod.query.filter(Mod.id == int(modid)).first():
        return {'error': True, 'reasons': ['The modid is invalid']}, 400
    # Get the mod
    mod = Mod.query.filter(Mod.id == int(modid)).first()
    return {'error': False, 'count': 1, 'data': mod_info(mod)}

@route('/api/mods/<modid>/edit', methods=['POST'])
@user_has('mods-edit', params=['modid'])
@with_session
def mod_edit(modid):
    """
    Edits a mod, based on the request parameters. Required fields: data
    """
    errors = list()
    if not modid.isdigit() or not Mod.query.filter(Mod.id == int(modid)).first():
        errors.append('The Mod ID is invalid.')
    if not request.form.get('data') or not is_json(request.form.get('data')):
        errors.append('The patch data is invalid.')
    if any(errors):
        return {'error': True, 'reasons': errors}, 400

    # Get variables
    parameters = json.loads(request.form.get('data'))

    # Get the matching mod and edit it
    mod = Mod.query.filter(Mod.id == int(modid)).first()
    edit_object(mod, parameters)
    return {'error': False}

@route('/api/mods/add', methods=['POST'])
@user_has('mods-add', params=['gameshort'])
@with_session
def add_mod():
    """
    Adds a mod, based on the request parameters. Required fields: name, gameshort, license
    """
    # Get variables
    name = request.form.get('name')
    short = request.form.get('gameshort')
    license = request.form.get('license')

    # Check the vars
    errors = list()
    if not name:
        errors.append('Invalid mod name.')
    if Mod.query.filter(Mod.name == name).first():
        errors.append('A mod with this name does already exist.')
    if not short:
        errors.append('Invalid gameshort.')
    if not license:
        errors.append('Invalid License.')
    if any(errors):
        return {'error': True, 'reasons': errors}, 400

    # Add new mod
    mod = Mod(name, current_user.id, game_id(short), license)
    db.add(mod)
    role = Role.query.filter(Role.name == current_user.username).first()
    role.add_abilities('mods-edit', 'mods-remove')
    role.add_param('mods-edit', 'modid', str(mod.id))
    role.add_param('mods-remove', 'name', name)
    return {'error': False}

@route('/api/mods/remove', methods=['POST'])
@user_has('mods-remove', params=['gameshort', 'name']) # We might want to allow deletion of own mods. Gameshort is here to allow per-game moderators.
@with_session
def remove_mod():
    """
    Removes a mod, based on the request parameters. Required fields: name, gameshort
    """
    # Get variables
    name = request.form.get('name')
    short = request.form.get('gameshort')

    # Check the vars
    errors = list()
    if not name:
        errors.append('Invalid mod name.')
    if not short:
        errors.append('Invalid gameshort.')
    if name and short and not Mod.query.filter(Mod.name == name).filter(Mod.game_id == game_id(short)).first():
        errors.append('A mod with theese parameters does not exist.')
    if any(errors):
        return {'error': True, 'reasons': errors}, 400

    # Add new mod
    mod = Mod.query.filter(Mod.name == name).filter(Mod.game_id == game_id(short)).first()
    db.remove(mod)
    role = Role.query.filter(Role.name == current_user.username).first()
    if not any(current_user.mods):
        role.remove_abilities('mods-edit', 'mods-remove')
    role.remove_param('mods-edit', 'modid', str(mod.id))
    role.remove_param('mods-remove', 'name', name)
    return {'error': False}

@route('/api/mods/<modid>/edit', methods=['POST'])
@user_has('mods-edit', params=['modid'])
@with_session
def mod_updateBG( modid):
    """
    Updates a mod background. Required fields: image
    """
    errors = list()
    if not modid.isdigit() or not Mod.query.filter(Mod.id == int(modid)).first():
        errors.append('The Mod ID is invalid.')
    if not request.files.get('image'):
        errors.append('The background is invalid.')
    if any(errors):
        return {'error': True, 'reasons': errors}, 400

    # Find the mod
    mod = Mod.query.filter(Mod.id == int(modid)).first()

    # Get the file and save it to disk
    f = request.files['image']
    filetype = os.path.splitext(os.path.basename(f.filename))[1]
    if not filetype in ['.png', '.jpg']:
        return {'error': True, 'reasons': ['This file type is not acceptable.']}, 400
    filename = secure_filename(mod.name) + '-' + str(time.time()) + filetype
    base_path = os.path.join(secure_filename(mod.user.username) + '_' + str(mod.user.id), secure_filename(mod.name))
    full_path = os.path.join(cfg['storage'], base_path)
    if not os.path.exists(full_path):
        os.makedirs(full_path)
    path = os.path.join(full_path, filename)
    try:
        os.remove(os.path.join(cfg['storage'], mod.background))
    except:
        pass # who cares
    f.save(path)
    mod.background = os.path.join(base_path, filename)
    return {'error': False}