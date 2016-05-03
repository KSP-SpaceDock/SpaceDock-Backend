from flask import Flask, jsonify
from SpaceDock.config import Config
from SpaceDock.api import API
from SpaceDock.documentation import Documentation

cfg = Config()
app = Flask(__name__)
documentation = Documentation(app)
api = API(app, documentation)

@app.before_first_request
def prepare():
    pass

if __name__ == '__main__':
    if cfg.get_environment() == 'dev':
        app.debug = True
    app.run(host = cfg['debug-host'], port = cfg.geti('debug-port'))
