from flask import Flask, jsonify
from flask_json import FlaskJSON
from flask_login import LoginManager
from SpaceDock.config import Config
from SpaceDock.database import Database
from SpaceDock.documentation import Documentation
from SpaceDock.profiler import Profiler


app = Flask(__name__)
json = FlaskJSON(app)
cfg = Config()
db = Database(cfg)
documentation = Documentation(app)
profiler = Profiler(cfg)

#Import of search must come after DB
from SpaceDock.search import Search
search = Search(db)

#Import of email must come after cfg/db load
from SpaceDock.email import Email
email = Email(cfg)

login_manager = LoginManager()
login_manager.init_app(app)

#Import must come after DB loads
from SpaceDock.api import API
api = API(app, documentation, cfg, db, email, profiler, search)

@app.before_first_request
def prepare():
    db.init_db()

from SpaceDock.objects import User
@login_manager.user_loader
def load_user(username):
    return User.query.filter(User.username == username).first()

login_manager.anonymous_user = lambda: None

if __name__ == '__main__':
    if cfg.get_environment() == 'dev':
        app.debug = True
    app.config['JSON_ADD_STATUS'] = False
    app.config['JSON_JSONP_OPTIONAL'] = False
    app.secret_key = cfg['secret-key']
    app.run(host = cfg['debug-host'], port = cfg.geti('debug-port'))
