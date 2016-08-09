from flask import request
from flask_login import current_user
from SpaceDock.common import *
from SpaceDock.database import db
from SpaceDock.formatting import pack_info
from SpaceDock.routing import route
from SpaceDock.objects import ModList

@route('/api/packs/')
def packs_list():
    """
    Outputs a list of modpacks
    """
    result = list()
    for pack in ModList.query.all():
        result.append(pack_info(pack))
    return {'error': False, 'count': len(result), 'data': result}

@route('/api/packs/<gameshort>/<packid>')
def packs_info(gameshort, packid):
    """
    Returns info for a specific modpack
    """
    if not packid.isdigit() or not ModList.query.filter(ModList.id == int(packid)).first():
        return {'error': True, 'reasons': ['The pack ID is invalid']}, 400
    if not ModList.query.filter(ModList.id == int(packid)).filter(ModList.game_id == game_id(gameshort)).first():
        return {'error': True, 'reasons': ['The gameshort is invalid.']}, 400

    # Get the pack
    pack = ModList.query.filter(ModList.id == int(packid)).first()
    return {'error': False, 'count': 1, 'data': pack_info(pack)}

@route('/api/packs/add', methods=['POST'])
@user_has('packs-add', params=['gameshort'])
@with_session
def packs_add():
    """
    Creates a new modlist. Required fields: name, gameshort
    """
    name = request.form.get('name')
    gameshort = request.form.get('gameshort')

    # Check the vars
    errors = list()
    if not name:
        errors.append('Invalid mod name.')
    if ModList.query.filter(ModList.name == name).first():
        errors.append('A modlist with this name does already exist.')
    if not short or not game_id(short):
        errors.append('Invalid gameshort.')
    if any(errors):
        return {'error': True, 'reasons': errors}, 400

    # Make the new list
    list = ModList()
    list.name = name
    list.game = Game.query.filter(Game.short == gameshort).first()
    list.user = current_user
    db.add(list)
    current_user.add_roles(name)    
    role = Role.query.filter(Role.name == name).first()
    role.add_abilities('packs-edit', 'mods-remove')
    role.add_param('packs-edit', 'packid', str(list.id))
    role.add_param('packs-remove', 'name', name)    
    db.add(role)
    db.commit()
    return {'error': False, 'count': 1, 'data': pack_info(list)}

@route('/api/packs/<gameshort>/<packid>/edit', methods=['POST'])
@user_has('packs-edit', params=['gameshort', 'packid'])
@with_session
def packs_edit(gameshort, packid):
    """
    Edits a modlist based on patch data. Required fields: data
    """
    if not packid.isdigit() or not ModList.query.filter(ModList.id == int(packid)).first():
        return {'error': True, 'reasons': ['The pack ID is invalid']}, 400
    if not ModList.query.filter(ModList.id == int(packid)).filter(ModList.game_id == game_id(gameshort)).first():
        return {'error': True, 'reasons': ['The gameshort is invalid.']}, 400
    if not request.form.get('data') or not is_json(request.form.get('data')):
        return {'error': True, 'reasons': ['The patch data is invalid.']}, 400

    # Get variables
    parameters = json.loads(request.form.get('data'))

    # Get the list
    list = ModList.query.filter(ModList.id == int(packid)).first()
    code = edit_object(list, parameters)

    # Error check
    if code == 2:
        return {'error': True, 'reasons': ['You tried to edit a value that doesn\'t exist.']}, 400
    elif code == 1:
        return {'error': True, 'reasons': ['You tried to edit a value that is marked as read-only.']}, 400
    else:
        return {'error': False, 'count': 1, 'data': pack_info(list)}


    