from flask import jsonify, request
from SpaceDock.objects import *

class GameEndpoints:
    def __init__(self, cfg, db):
        self.cfg = cfg
        self.db = db