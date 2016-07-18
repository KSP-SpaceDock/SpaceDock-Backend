from SpaceDock.app import app
from SpaceDock.config import cfg

# Start up the app
if __name__ == '__main__':
    app.run(host = cfg['debug-host'], port = cfg.geti('debug-port'))