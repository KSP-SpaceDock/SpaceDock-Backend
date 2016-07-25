from flask import Flask, jsonify
from flask_json import FlaskJSON
from flask_login import LoginManager
from SpaceDock.config import cfg
from SpaceDock.database import init_db
from SpaceDock.objects import User
from SpaceDock.plugins import load_plugins

# Create Flask
app = Flask(__name__)
json = FlaskJSON(app)
login_manager = LoginManager(app)
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

from SpaceDock.common import json_output
from SpaceDock.routing import add_wrapper

# Register JSON output
add_wrapper(json_output)

# Load plugins
load_plugins()

# Register endpoints
import SpaceDock.endpoints.access
import SpaceDock.endpoints.accounts
import SpaceDock.endpoints.admin
import SpaceDock.endpoints.featured
import SpaceDock.endpoints.game
import SpaceDock.endpoints.mods
import SpaceDock.endpoints.publisher
import SpaceDock.endpoints.user


