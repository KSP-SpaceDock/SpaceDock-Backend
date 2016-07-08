from flask import jsonify
from functools import wraps
from SpaceDock.common import json
from SpaceDock.endpoints.accounts import AccountEndpoints
from SpaceDock.endpoints.api import ApiEndpoints
from SpaceDock.endpoints.anonymous import AnonymousEndpoints
from SpaceDock.endpoints.admin import AdminEndpoints
from SpaceDock.endpoints.game import GameEndpoints
from SpaceDock.endpoints.publisher import PublisherEndpoints

class API:
    def __init__(self, flask, documentation, cfg, db, email, profiler, search):
        self.flask = flask
        self.documentation = documentation
        self.cfg = cfg
        self.db = db
        self.email = email
        self.profiler = profiler
        self.search = search
        self.register_endpoints()

    #Decorate all API endpoints
    def decorate_function(self, member):
        if self.cfg.getb('profiler-histogram') or self.cfg.geti('profiler') != 0:
            member = self.profiler.profile_method(member)
        member = json(member)
        return member

    def register_api_endpoint(self, endpoints):
        for member_name in dir(endpoints):
            member = getattr(endpoints, member_name)
            if 'api_path' in dir(member):
                member = self.decorate_function(member)
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
        anonymous_endpoints = AnonymousEndpoints(self.cfg, self.db, self.search)
        self.register_api_endpoint(anonymous_endpoints)
        api_endpoints = ApiEndpoints(self.cfg, self.db, self.email, self.search)
        self.register_api_endpoint(api_endpoints)
        account_endpoints = AccountEndpoints(self.cfg, self.db, self.email)
        self.register_api_endpoint(account_endpoints)
        game_endpoints = GameEndpoints(self.cfg, self.db)
        self.register_api_endpoint(game_endpoints)
        publisher_endpoints = PublisherEndpoints()
        self.register_api_endpoint(publisher_endpoints)