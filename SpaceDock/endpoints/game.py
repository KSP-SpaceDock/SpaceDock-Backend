from flask import request
from sqlalchemy import desc
from SpaceDock.common import *
from SpaceDock.objects import *
from SpaceDock.formatting import game_info, game_version_info

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
        Displays information about a game. Required parameters: gameshort
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
        Displays information about the versions of a game. Required parameters: gameshort
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
        Displays a list of all mods added for this game. Required parameters: gameshort
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
        Displays all mod lists this game knows. Required parameters: gameshort
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


    #TODO: Move this to admin maybe?
    #TODO: Redo with propper access control
    #@accessrequired
    #@with_session
    #def edit_game(self, gameid):
    #    """
    #    Edits a game, based on the request parameters. Required parameters: gameid
    #    """
    #    if not gameid.isdigit() or len(Game.query.filter(Game.id == int(gameid)).all()) == 0:
    #        return jsonify({'error': True, 'idErrors': 'The number you entered is not a valid ID.'}), 400
    #
    #    # Get variables
    #    parameters = request.form.to_dict()
    #
    #    # Get the matching game and edit it
    #    game = Game.query.filter(Game.id == int(gameid)).first()
    #    edit_object(game, parameters)
    #    return jsonify({'error': False})
    #
    #edit_game.api_path = '/api/games/<gameid>/edit'
    #edit_game.methods = ['POST']
