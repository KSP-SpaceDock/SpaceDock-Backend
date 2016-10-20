from flask import request
from SpaceDock.common import game_id, user_has, with_session
from SpaceDock.database import db
from SpaceDock.formatting import feature_info
from SpaceDock.objects import Featured, Game, Mod
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
    if not Game.query.filter(Game.active).filter(Game.short == gameshort).first():
        return {'error': True, 'reasons': ['The gameshort is invalid.'], 'codes': ['2125']}, 400

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
    if not isinstance(modid, int) or not Mod.query.filter(Mod.id == modid).first():
        return {'error': True, 'reasons': ['The modid is invalid.'], 'codes': ['2130']}, 400
    elif not Mod.query.filter(Mod.id == modid).filter(Mod.game_id == game_id(gameshort)).first():
        return {'error': True, 'reasons': ['The gameshort is invalid.'], 'codes': ['2125']}, 400
    elif not Mod.query.filter(Mod.published).filter(Mod.id == modid).first():
        return {'error': True, 'reasons': ['The mod must be published first.'], 'codes': ['3022']}, 400
    if Featured.query.filter(Featured.mod_id == modid).first():
        return {'error': True, 'reasons': ['The mod is already featured'], 'codes': ['3015']}, 400

    # Everything's fine, let's feature the mod    
    feature = Featured(Mod.query.filter(Mod.id == modid).first())
    db.add(feature)
    db.flush()
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
    if not isinstance(modid, int) or not Mod.query.filter(Mod.id == modid).first():
        return {'error': True, 'reasons': ['The modid is invalid.'], 'codes': ['2130']}, 400
    elif not Mod.query.filter(Mod.id == modid).filter(Mod.game_id == game_id(gameshort)).first():
        return {'error': True, 'reasons': ['The gameshort is invalid.'], 'codes': ['2125']}, 400
    elif not Mod.query.filter(Mod.published).filter(Mod.id == modid).first():
        return {'error': True, 'reasons': ['The mod must be published first.'], 'codes': ['3022']}, 400
    elif not Featured.query.filter(Featured.mod_id == modid).first():
        return {'error': True, 'reasons': ['This mod isn\'t featured.']}, 400

    # Unfeature the mod
    feature = Featured.query.filter(Featured.mod_id == modid).first()
    db.delete(feature)
    return {'error': False}