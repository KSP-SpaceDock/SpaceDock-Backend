from flask import request
from sqlalchemy import desc
from SpaceDock.common import *
from SpaceDock.objects import *
from SpaceDock.formatting import game_info, game_version_info

import json

class GameEndpoints:
    def __init__(self, cfg, db):
        self.cfg = cfg
        self.db = db.get_database()

    def list_games(self):
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

    list_games.api_path = '/api/games'

    def game_info(self, gameshort):
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

    game_info.api_path = '/api/games/<gameshort>'

    def game_versions(self, gameshort):
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

    game_versions.api_path = '/api/games/<gameshort>/versions'

    def game_mods(self, gameshort):
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

    game_mods.api_path = '/api/games/<gameshort>/mods'

    def game_modlists(self, gameshort):
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

    game_modlists.api_path = '/api/games/<gameshort>/modlists'

    @with_session
    @user_has('game-edit', params=['gameshort'])
    def edit_game(self, gameshort):
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

    edit_game.api_path = '/api/games/<gameshort>/edit'
    edit_game.methods = ['POST']

    @with_session
    @user_has('game-add', params=['pubid'])
    def add_game(self):
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
        self.db.add(game)
        return {'error': False}

    add_game.api_path = '/api/games/add'
    add_game.methods = ['POST']

    @with_session
    @user_has('game-remove', params=['short'])
    def remove_game(self):
        """
        Removes a game from existence. Required fields: short
        """
        short = request.form['short']

        # Check if the gameshort is valid
        if not Game.query.filter(Game.short == short).first():
            return {'error': True, 'reasons': ['The gameshort is invalid.']}, 400

        # Get the game and remove it
        game = Game.query.filter(Game.short == short).first()
        self.db.remove(game)
        return {'error': False}

    remove_game.api_path = '/api/games/remove'
    remove_game.methods = ['POST']