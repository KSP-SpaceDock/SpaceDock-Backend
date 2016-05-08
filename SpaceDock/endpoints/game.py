from flask import jsonify, request
from sqlalchemy import desc, asc
from SpaceDock.objects import *
from SpaceDock.formatting import game_info, game_version_info

class GameEndpoints:
    def __init__(self, cfg, db):
        self.cfg = cfg
        self.db = db
        
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
        if not gameid.isdigit():
            return jsonify({'error': 'True', 'message': 'The number you entered is not a valid ID.'}), 400
        if len(Game.query.filter(Game.id == int(gameid)).all()) == 0:
            return jsonify({'error': 'True', 'message': 'The number you entered is not a valid ID.'}), 400

        # Return gameinfo
        game = Game.query.filter(Game.id == int(gameid)).first()
        return jsonify(game_info(game))

    game_info.api_path = '/api/games/<gameid>'

    def game_versions(self, gameid):
        """
        Displays information about the versions of a game. Required parameters: gameid
        """
        if not gameid.isdigit():
            return jsonify({'error': 'True', 'message': 'The number you entered is not a valid ID.'}), 400
        if len(Game.query.filter(Game.id == int(gameid)).all()) == 0:
            return jsonify({'error': 'True', 'message': 'The number you entered is not a valid ID.'}), 400
   
        # get game versions        
        game = Game.query.filter(Game.id == int(gameid)).first()
        versions = GameVersion.query.filter(GameVersion.game_id == game.id).all()

        # Format them
        result = list()
        for version in versions:
            result.append(game_version_info(version))
        return jsonify({game.id: result})

    game_versions.api_path = '/api/games/<gameid>/versions'