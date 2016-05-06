from flask import jsonify
from SpaceDock.endpoints.accounts import AccountEndpoints
from SpaceDock.endpoints.game import GameEndpoints

class API:
    def __init__(self, flask, documentation, cfg, db, email):
        self.flask = flask
        self.documentation = documentation
        self.cfg = cfg
        self.db = db
        self.email = email
        self.register_endpoints()
        
    def register_api_endpoint(self, endpoints):
        for member_name in dir(endpoints):
            member = getattr(endpoints, member_name)
            if 'api_path' in dir(member):
                print("Registered " + member.__name__)
                self.flask.add_url_rule(member.api_path, member.__name__, member)
                self.documentation.add_documentation(member.api_path, member)
        
    def register_endpoints(self):
        account_endpoints = AccountEndpoints(self.cfg, self.db, self.email)
        self.register_api_endpoint(account_endpoints)
        game_endpoints = GameEndpoints(self.cfg, self.db)
        self.register_api_endpoint(game_endpoints)