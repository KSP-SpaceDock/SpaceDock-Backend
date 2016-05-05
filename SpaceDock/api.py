from flask import jsonify
from SpaceDock.endpoints.game import GameEndpoints

class API:
    def __init__(self, flask, documentation, cfg, db):
        self.flask = flask
        self.documentation = documentation
        self.cfg = cfg
        self.db = db
        self.register_endpoints()
        
    def register_api_endpoint(self, url, method):
        self.flask.add_url_rule("/api/" + url, method.__name__, method)
        self.documentation.add_documentation(url, method)
        
    def register_endpoints(self):
        game_endpoints = GameEndpoints(self.cfg, self.db)
        self.register_api_endpoint("listgame", game_endpoints.listgame)