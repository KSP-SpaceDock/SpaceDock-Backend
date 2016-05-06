from flask import jsonify, request
from sqlalchemy import desc, asc
from SpaceDock.objects import *
from SpaceDock.formatting import game_info

class GameEndpoints:
    def __init__(self, cfg, db):
        self.cfg = cfg
        self.db = db
        
    def list_games(self):
        """
        Displays a list of all games in the database.
        """
        results = list()
        for game in Game.query.order_by(desc(Game.name)).all():
            results.append(game_info(game))
        return jsonify(results)
        
    list_games.api_path = "/api/games"