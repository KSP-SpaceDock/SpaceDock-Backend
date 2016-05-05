from flask import Flask, jsonify
from SpaceDock.config import Config
from SpaceDock.database import Database
from SpaceDock.documentation import Documentation

app = Flask(__name__)
cfg = Config()
db = Database(cfg)
documentation = Documentation(app)

#Import must come after DB loads
from SpaceDock.api import API
api = API(app, documentation, cfg, db)

@app.before_first_request
def prepare():
    pass

if __name__ == '__main__':
    if cfg.get_environment() == 'dev':
        app.debug = True
    app.run(host = cfg['debug-host'], port = cfg.geti('debug-port'))
