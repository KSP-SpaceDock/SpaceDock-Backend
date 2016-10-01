from flask import request
from sqlalchemy import desc
from SpaceDock.common import boolean, edit_object, game_id, user_has, with_session
from SpaceDock.database import db
from SpaceDock.formatting import game_info, game_version_info
from SpaceDock.objects import Game, GameVersion, Publisher
from SpaceDock.routing import route


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
    if any(errors):
        return {'error': True, 'reasons': errors}, 400

    # Get the matching game and edit it
    game = Game.query.filter(Game.short == gameshort).first()
    code = edit_object(game, request.json)

    # Error check
    if code == 2:
        return {'error': True, 'reasons': ['You tried to edit a value that doesn\'t exist.']}, 400
    elif code == 1:
        return {'error': True, 'reasons': ['You tried to edit a value that is marked as read-only.']}, 400
    else:
        return {'error': False, 'count': 1, 'data': game_info(game)}

@route('/api/games/add', methods=['POST'])
@user_has('game-add', params=['pubid'])
@with_session
def add_game():
    """
    Adds a new game based on the request parameters. Required fields: name, pubid, short
    """
    name = request.json.get('name')
    pubid = request.json.get('pubid')
    short = request.json.get('short')

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
    db.commit()
    return {'error': False, 'count': 1, 'data': game_info(game)}

@route('/api/games/remove', methods=['POST'])
@user_has('game-remove', params=['short'])
@with_session
def remove_game():
    """
    Removes a game from existence. Required fields: short
    """
    short = request.json.get('short')

    # Check if the gameshort is valid
    if not Game.query.filter(Game.short == short).first():
        return {'error': True, 'reasons': ['The gameshort is invalid.']}, 400

    # Get the game and remove it
    game = Game.query.filter(Game.short == short).first()
    db.delete(game)
    return {'error': False}