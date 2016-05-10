from flask import jsonify
from SpaceDock.endpoints.accounts import AccountEndpoints
from SpaceDock.endpoints.admin import AdminEndpoints
from SpaceDock.endpoints.game import GameEndpoints

class API:
    def __init__(self, flask, documentation, cfg, db, email, profiler):
        self.flask = flask
        self.documentation = documentation
        self.cfg = cfg
        self.db = db
        self.email = email
        self.profiler = profiler
        self.register_endpoints()
        
    def register_api_endpoint(self, endpoints):
        for member_name in dir(endpoints):
            member = getattr(endpoints, member_name)
            if 'api_path' in dir(member):
                if self.cfg.getb('profiler-histogram') or self.cfg.geti('profiler') != 0:
                    member = self.profiler.profile_method(member)
                print("Registered " + member.api_path)
                self.flask.add_url_rule(member.api_path, member.__name__, member)
                if member.api_path.endswith('/'):
                    self.flask.add_url_rule(member.api_path[:-1], member.__name__, member)
                else:
                    self.flask.add_url_rule(member.api_path + '/', member.__name__, member)
                self.documentation.add_documentation(member.api_path, member)

    def register_endpoints(self):
        if self.cfg.getb('profiler-histogram'):
            self.register_api_endpoint(self.profiler)
        admin_endpoints = AdminEndpoints(self.db, self.email)
        self.register_api_endpoint(admin_endpoints)
        account_endpoints = AccountEndpoints(self.cfg, self.db, self.email)
        self.register_api_endpoint(account_endpoints)
        game_endpoints = GameEndpoints(self.cfg, self.db)
        self.register_api_endpoint(game_endpoints)