from flask import request
from flask_login import current_user
from SpaceDock.common import edit_object, game_id, user_has, with_session
from SpaceDock.database import db
from SpaceDock.formatting import pack_info
from SpaceDock.objects import ModList, Game, Role, ModListItem, Mod
from SpaceDock.routing import route

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
        return {'error': True, 'reasons': ['The pack ID is invalid'], 'codes': ['2135']}, 400
    if not ModList.query.filter(ModList.id == int(packid)).filter(ModList.game_id == game_id(gameshort)).first():
        return {'error': True, 'reasons': ['The gameshort is invalid.'], 'codes': ['2125']}, 400

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
    name = request.json.get('name')
    gameshort = request.json.get('gameshort')

    # Check the vars
    errors = ()
    codes = ()
    if not name:
        errors.append('Invalid modlist name.')
        codes.append('2137')
    if ModList.query.filter(ModList.name == name).first():
        errors.append('A modlist with this name already exists.')
        codes.append('2025')
    if not gameshort or not game_id(gameshort):
        errors.append('Invalid gameshort.')
        codes.append('2125')
    if any(errors):
        return {'error': True, 'reasons': errors, 'codes': codes}, 400

    # Make the new list
    pack = ModList(name, Game.query.filter(Game.short == gameshort).first(), current_user)
    db.add(pack)
    current_user.add_roles(name)    
    role = Role.query.filter(Role.name == name).first()
    role.add_abilities('packs-edit', 'packs-remove')
    role.add_param('packs-edit', 'packid', str(pack.id))
    role.add_param('packs-remove', 'name', name)    
    db.add(role)
    db.flush()
    return {'error': False, 'count': 1, 'data': pack_info(pack)}

@route('/api/packs/<gameshort>/<packid>/edit', methods=['POST'])
@user_has('packs-edit', params=['gameshort', 'packid'])
@with_session
def packs_edit(gameshort, packid):
    """
    Edits a modlist based on patch data. Required fields: data
    """
    if not packid.isdigit() or not ModList.query.filter(ModList.id == int(packid)).first():
        return {'error': True, 'reasons': ['The pack ID is invalid'], 'codes': ['2135']}, 400
    if not ModList.query.filter(ModList.id == int(packid)).filter(ModList.game_id == game_id(gameshort)).first():
        return {'error': True, 'reasons': ['The gameshort is invalid.'], 'codes': ['2125']}, 400

    # Get the list
    pack = ModList.query.filter(ModList.id == int(packid)).first()
    code = edit_object(pack, request.json)

    # Error check
    if code == 3:
        return {'error': True, 'reasons': ['The value you submitted is invalid'], 'codes': ['2180']}, 400
    elif code == 2:
        return {'error': True, 'reasons': ['You tried to edit a value that doesn\'t exist.'], 'codes': ['3090']}, 400
    elif code == 1:
        return {'error': True, 'reasons': ['You tried to edit a value that is marked as read-only.'], 'codes': ['3095']}, 400
    else:
        return {'error': False, 'count': 1, 'data': pack_info(list)}


@route('/api/packs/<gameshort>/<packid>/addmod', methods=['POST'])
@user_has('packs-edit', params=['gameshort', 'packid'])
@with_session
def packs_add_mod(gameshort, packid):
    """
    Adds a new mod to the modlist. Required fields: modid
    """
    mod_id = request.json.get('modid')

    # Error check
    errors = ()
    codes = ()
    if not isinstance(mod_id, int) or not Mod.query.filter(Mod.id == mod_id).first():
        errors.append('The mod ID is invalid')
        codes.append('2130')
    elif isinstance(mod_id, int) and not Mod.query.filter(Mod.published).filter(Mod.id == mod_id).first():
        errors.append('The mod is not published.')
        codes.append('3020')
    if not packid.isdigit() or not ModList.query.filter(ModList.id == int(packid)).first():
        errors.append('The pack ID is invalid')
        codes.append('2135')
    if packid.isdigit() and not ModList.query.filter(ModList.id == int(packid)).filter(ModList.game_id == game_id(gameshort)).first():
        errors.append('The gameshort is invalid')
        codes.append('2125')
    if isinstance(mod_id, int) and packid.isdigit() and ModListItem.query.filter(ModListItem.mod_id == mod_id).filter(ModListItem.mod_list_id == int(packid)).first():
        errors.append('The specified mod was already added to the modlist')
        codes.append('2030')
    if any(errors):
        return {'error': True, 'reasons': errors, 'codes': codes}, 400

    # Get the list
    pack = ModList.query.filter(ModList.id == int(packid)).first()
    moditem = ModListItem(mod=Mod.query.filter(Mod.id == mod_id).first(), modlist=pack)

    db.add(moditem)
    db.flush()

    return {'error': False}

@route('/api/packs/<gameshort>/<packid>/delmod', methods=['POST'])
@user_has('packs-edit', params=['gameshort', 'packid'])
@with_session
def packs_del_mod(gameshort, packid):
    """
    Removes a mod from the modlist. Required fields: modid
    """
    mod_id = request.json.get('modid')

    # Error check
    errors = ()
    codes = ()
    if not isinstance(mod_id, int) or not Mod.query.filter(Mod.id == mod_id).first():
        errors.append('The mod ID is invalid')
        codes.append('2130')
    if not packid.isdigit() or not ModList.query.filter(ModList.id == int(packid)).first():
        errors.append('The pack ID is invalid')
        codes.append('2135')
    if not ModList.query.filter(ModList.id == int(packid)).filter(ModList.game_id == game_id(gameshort)).first():
        errors.append('The gameshort is invalid')
        codes.append('2125')
    if not ModListItem.query.filter(ModListItem.mod_id == mod_id).filter(ModListItem.mod_list_id == int(packid)).first():
        errors.append('The specified mod is not included in the modlist')
        codes.append('3053')
    if any(errors):
        return {'error': True, 'reasons': errors, 'codes': codes}, 400

    # Get the list
    pack = ModList.query.filter(ModList.id == int(packid)).first()
    moditem = ModListItem.query.filter(ModListItem.mod_list_id == int(packid)).filter(ModListItem.mod_id == mod_id).first()

    db.delete(moditem)

    return {'error': False}