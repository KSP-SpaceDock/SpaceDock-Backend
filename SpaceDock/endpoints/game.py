from flask import request
from sqlalchemy import desc
from SpaceDock.common import *
from SpaceDock.database import db
from SpaceDock.formatting import game_info, game_version_info
from SpaceDock.objects import *
from SpaceDock.routing import route

import json


@route('/api/games')
def list_games():
    """
    Displays a list of all games in the database.
    """
    results = list()
    includeInactive = False
    if request.args.get('includeInactive'):
        includeInactive = boolean(request.args.get('includeInactive'))
    # Game.active or includeInactive refuses to work :-(
    f = Game.active
    if includeInactive:
        f = True
    for game in Game.query.order_by(desc(Game.name)).filter(f):
        results.append(game_info(game))
    return {'error': False, 'count': len(results), 'data': results}

@route('/api/games/<gameshort>')
def games_info(gameshort):
    """
    Displays information about a game.
    """
    # Get the games with the according gameshort
    filter = Game.query.filter(Game.short == gameshort)

    # Game doesn't exist
    if len(filter.all()) == 0:
        return {'error': True, 'reasons': ['The gameshort is invalid.']}, 400

    # Game does exist
    game = filter.first()
    return {'error': False, 'count': 1, 'data': game_info(game)}

@route('/api/games/<gameshort>/versions')
def game_versions(gameshort):
    """
    Displays information about the versions of a game.
    """
    if not Game.query.filter(Game.short == gameshort).first():
        return {'error': True, 'reasons': ['The gameshort is invalid.']}, 400

    # Get the ID
    gameid = game_id(gameshort)

    # get game versions
    versions = GameVersion.query.filter(GameVersion.game_id == gameid).all()

    # Format them
    results = list()
    for version in versions:
        results.append(game_version_info(version))
    return {'error': False, 'count': len(results), 'data': results}

@route('/api/games/<gameshort>/mods')
def game_mods(gameshort):
    """
    Displays a list of all mods added for this game.
    """
    if not Game.query.filter(Game.short == gameshort).first():
        return {'error': True, 'reasons': ['The gameshort is invalid.']}, 400

    # Get the ID
    gameid = game_id(gameshort)

    # Get mods
    mods = Mod.query.filter(Mod.game_id == int(gameid)).all()

    # Format
    result = dict()
    for mod in mods:
        result[str(mod.id)] = mod.name
    return {'error': False, 'count': len(result), 'data': result}

@route('/api/games/<gameshort>/modlists')
def game_modlists(gameshort):
    """
    Displays all mod lists this game knows.
    """
    if not Game.query.filter(Game.short == gameshort).first():
        return {'error': True, 'reasons': ['The gameshort is invalid.']}, 400

    # Get the ID
    gameid = game_id(gameshort)

    # Get mod lists
    game = Game.query.filter(Game.id == int(gameid)).first()
    modlists = ModList.query.filter(ModList.game_id == int(gameid)).all()

    # Format
    result = dict()
    for ml in modlists:
        result[str(ml.id)] = ml.name
    return {'error': False, 'count': len(result), 'data': result}


@route('/api/games/<gameshort>/edit', methods=['POST'])
@user_has('game-edit', params=['gameshort'])
@with_session
def edit_game(gameshort):
    """
    Edits a game, based on the request parameters. Required fields: data
    """
    errors = list()
    if not Game.query.filter(Game.short == gameshort).first():
        errors.append('The gameshort is invalid.')
    if not request.form.get('data') or not is_json(request.form.get('data')):
        errors.append('The patch data is invalid.')
    if any(errors):
        return {'error': True, 'reasons': errors}, 400

    # Get variables
    parameters = json.loads(request.form.get('data'))

    # Get the matching game and edit it
    game = Game.query.filter(Game.short == gameshort).first()
    edit_object(game, parameters)
    return {'error': False}

@route('/api/games/add', methods=['POST'])
@user_has('game-add', params=['pubid'])
@with_session
def add_game():
    """
    Adds a new game based on the request parameters. Required fields: name, pubid, short
    """
    name = request.form.get('name')
    pubid = request.form.get('pubid')
    short = request.form.get('short')

    errors = list()

    # Check if the publisher ID is valid
    if not pubid or not pubid.isdigit() or not Publisher.query.filter(Publisher.id == int(pubid)).first():
        errors.append('The pubid is invalid.')
    if not name:
        errors.append('The name is invalid.')
    if not short:
        errors.append('The gameshort is invalid.')

    # Check if the game already exists
    if Game.query.filter(Game.short == short).first():
        errors.append('The gameshort already exists.')
    if Game.query.filter(Game.name == name).first():
        errors.append('The game name already exists.')

    # Errors
    if len(errors) > 0:
        return {'error': True, 'reasons': errors}, 400

    # Make a new game
    game = Game(name, int(pubid), short)
    db.add(game)
    return {'error': False}

@route('/api/games/remove', methods=['POST'])
@user_has('game-remove', params=['short'])
@with_session
def remove_game():
    """
    Removes a game from existence. Required fields: short
    """
    short = request.form.get('short')

    # Check if the gameshort is valid
    if not Game.query.filter(Game.short == short).first():
        return {'error': True, 'reasons': ['The gameshort is invalid.']}, 400

    # Get the game and remove it
    game = Game.query.filter(Game.short == short).first()
    db.delete(game)
    return {'error': False}