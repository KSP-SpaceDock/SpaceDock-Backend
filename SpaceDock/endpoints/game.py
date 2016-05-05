from flask import jsonify
from SpaceDock.objects import *

class GameEndpoints:
    def __init__(self, cfg, db):
        self.cfg = cfg
        self.db = db
    
    def listgame(self):
        """
        List the first game
        """
        gameList = {}
        for g in Game.query.all():
            gameList[g.id] = {'name': g.name, 'short_name': g.short}
        return jsonify(gameList)