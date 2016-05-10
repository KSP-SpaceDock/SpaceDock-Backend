from flask import jsonify, request
from flask_login import login_required, current_user
from sqlalchemy import desc, asc
from SpaceDock.objects import *
from SpaceDock.formatting import game_info, game_version_info
from SpaceDock.common import with_session, edit_object

class GameEndpoints:
    def __init__(self, cfg, db):
        self.cfg = cfg
        self.db = db.get_database()
        
    def list_games(self):
        """
        Displays a list of all games in the database.
        """
        results = dict()
        for game in Game.query.order_by(desc(Game.name)).filter(Game.active):
            results[str(game.id)] = game.short;
        return jsonify(results)
        
    list_games.api_path = '/api/games'

    def game_info(self, gameid):
        """
        Displays information about a game. Required parameters: gameid
        """
        if not gameid.isdigit() or not Game.query.filter(Game.id == int(gameid)).first():
            return jsonify({'error': True, 'idErrors': 'The number you entered is not a valid ID.'}), 400

        # Return gameinfo
        game = Game.query.filter(Game.id == int(gameid)).first()
        return jsonify({'error': False, 'game_info': game_info(game)})

    game_info.api_path = '/api/games/<gameid>'

    def game_versions(self, gameid):
        """
        Displays information about the versions of a game. Required parameters: gameid
        """
        if not gameid.isdigit() or not Game.query.filter(Game.id == int(gameid)).first():
            return jsonify({'error': True, 'idErrors': 'The number you entered is not a valid ID.'}), 400
   
        # get game versions        
        versions = GameVersion.query.filter(GameVersion.game_id == int(gameid)).all()

        # Format them
        result = list()
        for version in versions:
            result.append(game_version_info(version))
        return jsonify({'error': False, 'game_versions': result})

    game_versions.api_path = '/api/games/<gameid>/versions'

    def game_mods(self, gameid):
        """
        Displays a list of all mods added for this game. Required parameters: gameid
        """
        if not gameid.isdigit() or not Game.query.filter(Game.id == int(gameid)).first():
            return jsonify({'error': True, 'idErrors': 'The number you entered is not a valid ID.'}), 400

        # Get mods
        mods = Mod.query.filter(Mod.game_id == int(gameid)).all()

        # Format
        result = dict()
        for mod in mods:
            result[str(mod.id)] = mod.name
        return jsonify({'error': False, 'mods': result})

    game_mods.api_path = '/api/games/<gameid>/mods'

    def game_modlists(self, gameid):
        """
        Displays all mod lists this game knows. Required parameters: gameid
        """
        if not gameid.isdigit() or len(Game.query.filter(Game.id == int(gameid)).all()) == 0:
            return jsonify({'error': True, 'idErrors': 'The number you entered is not a valid ID.'}), 400

        # Get mod lists
        game = Game.query.filter(Game.id == int(gameid)).first()
        modlists = ModList.query.filter(ModList.game_id == int(gameid)).all()

        # Format
        result = dict()
        for ml in modlists:
            result[str(ml.id)] = ml.name
        return jsonify({'error': False, 'modlists': result})

    game_modlists.api_path = '/api/games/<gameid>/modlists'

    #TODO: Move this to admin maybe?
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
