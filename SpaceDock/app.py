from flask import Flask
from flask_cors import CORS
from flask_json import FlaskJSON
from flask_limiter import Limiter
from flask_limiter.util import get_remote_address
from flask_login import LoginManager
from werkzeug.contrib.fixers import ProxyFix
from SpaceDock.config import cfg
from SpaceDock.database import init_db
from SpaceDock.objects import User
from SpaceDock.plugins import load_plugins

# Create Flask
app = Flask(__name__)
json = FlaskJSON(app)
login_manager = LoginManager(app)
limiter = Limiter(app, key_func=get_remote_address, headers_enabled=cfg.getb('limit-headers'), 
                  storage_uri=cfg["redis-connection"] if cfg.get_environment() != 'dev' else None)
if cfg.getb('disable-same-origin'):
    cors = CORS(app, supports_credentials=True)
init_db()

# Config
if cfg.get_environment() == 'dev':
    app.debug = True
app.config['JSON_ADD_STATUS'] = False
app.config['JSON_JSONP_OPTIONAL'] = False
app.secret_key = cfg['secret-key']
app.url_map.strict_slashes = False

@login_manager.user_loader
def load_user(username):
    return User.query.filter(User.username == username).first()

login_manager.anonymous_user = lambda: None

from SpaceDock.common import json_output, cache, limit
from SpaceDock.routing import add_wrapper
from SpaceDock import errors

# Register JSON output
add_wrapper(json_output)
add_wrapper(cache)
add_wrapper(limit)

# Load plugins
load_plugins()

# Register endpoints
import SpaceDock.endpoints.access
import SpaceDock.endpoints.accounts
import SpaceDock.endpoints.admin
import SpaceDock.endpoints.featured
import SpaceDock.endpoints.game
import SpaceDock.endpoints.general
import SpaceDock.endpoints.mods
import SpaceDock.endpoints.packs
import SpaceDock.endpoints.publisher
import SpaceDock.endpoints.tokens
import SpaceDock.endpoints.user

# Proxy fix
app.wsgi_app = ProxyFix(app.wsgi_app)