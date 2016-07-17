from SpaceDock.common import json_output
from SpaceDock.endpoints.access import AccessEndpoints
from SpaceDock.endpoints.accounts import AccountEndpoints
from SpaceDock.endpoints.admin import AdminEndpoints
from SpaceDock.endpoints.game import GameEndpoints
from SpaceDock.endpoints.mods import ModEndpoints
from SpaceDock.endpoints.publisher import PublisherEndpoints
from SpaceDock.endpoints.user import UserEndpoints

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

    # Decorate all API endpoints
    def decorate_function(self, member):
        if self.cfg.getb('profiler-histogram') or self.cfg.geti('profiler') != 0:
            member = self.profiler.profile_method(member)
        member = json_output(member)
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
        self.register_api_endpoint(AccessEndpoints(self.cfg, self.db))
        self.register_api_endpoint(AccountEndpoints(self.cfg, self.db, self.email))
        self.register_api_endpoint(AdminEndpoints(self.db, self.email))
        # self.register_api_endpoint(AnonymousEndpoints(self.cfg, self.db, self.search))
        # self.register_api_endpoint(ApiEndpoints(self.cfg, self.db, self.email, self.search))
        self.register_api_endpoint(GameEndpoints(self.cfg, self.db))
        self.register_api_endpoint(ModEndpoints(self.cfg, self.db))
        self.register_api_endpoint(PublisherEndpoints(self.db))
        self.register_api_endpoint(UserEndpoints(self.cfg, self.db))