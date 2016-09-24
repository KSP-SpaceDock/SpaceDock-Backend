from SpaceDock.common import *
from SpaceDock.database import db
from SpaceDock.email import *
from SpaceDock.formatting import feature_info
from SpaceDock.objects import *
from SpaceDock.routing import route


@route('/api/mods/featured')
def list_featured():
    """
    Returns a list of featured mods.
    """
    result = list()
    for feature in Featured.query.all():
        result.append(feature_info(feature))
    return {'error': False, 'count': len(result), 'data': result}

@route('/api/mods/featured/<gameshort>')
def list_featured_game(gameshort):
    """
    Returns a list of featured mods for a specific game.
    """
    if not Game.query.filter(Game.short == gameshort).first():
        return {'error': True, 'reasons': ['The gameshort is invalid.']}, 400

    # Get the features
    result = list()
    for feature in Featured.query.all():
        mod = Mod.query.filter(Mod.id == Featured.mod_id).first()
        if mod.game.short == gameshort:
            result.append(feature_info(feature))
    return {'error': False, 'count': len(result), 'data': result}

@route('/api/mods/featured/add/<gameshort>', methods=['POST'])
@user_has('mods-feature', params=['gameshort'])
@with_session
def add_feature(gameshort):
    """
    Features a mod for this game. Required fields: modid
    """
    modid = request.json.get('modid')

    # Errorcheck
    errors = list()
    if not modid.isdigit() or not Mod.query.filter(Mod.id == int(modid)).first():
        errors.append('The Mod ID is invalid.')
    if not any(errors) and not Mod.query.filter(Mod.id == int(modid)).filter(Mod.game_id == game_id(gameshort)).first():
       errors.append('The gameshort is invalid.')
    if any(errors):
        return {'error': True, 'reasons': errors}, 400

    # Everything's fine, let's feature the mod    
    if Featured.query.filter(Featured.mod_id == int(modid)).first():
        return {'error': True, 'reasons': ['The mod is already featured']}, 400
    feature = Featured(int(modid))
    db.add(feature)
    db.commit()
    return {'error': False, 'count': 1, 'data': feature_info(feature)}

@route('/api/mods/featured/remove/<gameshort>', methods=['POST'])
@user_has('mods-feature', params=['gameshort'])
@with_session
def remove_feature(gameshort):
    """
    Unfeatures a mod for this game. Required fields: modid
    """
    modid = request.json.get('modid')

    # Errorcheck
    errors = list()
    if not modid.isdigit() or not Mod.query.filter(Mod.id == int(modid)).first():
        errors.append('The Mod ID is invalid.')
    if not any(errors) and not Mod.query.filter(Mod.id == int(modid)).filter(Mod.game_id == game_id(gameshort)).first():
       errors.append('The gameshort is invalid.')
    if modid.isdigit() and not Featured.query.filter(Featured.mod_id == int(modid)).first():
        errors.append('This mod isn\'t featured.')
    if any(errors):
        return {'error': True, 'reasons': errors}, 400

    # Unfeature the mod
    feature = Featured.query.filter(Featured.mod_id == int(modid)).first()
    db.delete(feature)
    return {'error': False}