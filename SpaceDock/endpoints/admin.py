from flask import Blueprint, render_template, abort, redirect
from flask.ext.login import current_user
from sqlalchemy import desc
from SpaceDock.objects import User, Mod, GameVersion, Game, Publisher
from SpaceDock.common import *
from SpaceDock.formatting import *
from flask.ext.login import current_user, login_user, logout_user

class AdminEndpoints:
    #TODO: Split this up maybe?.
    @adminrequired
    def backend(self):
        """
        Get all backend data
        """
        user_count = User.query.count()
        users_raw = User.query.order_by(desc(User.created)).all();
        users = []
        for user in users_raw:
            users.append(admin_user_info(user))
        mods = Mod.query.count()
        versions_raw = GameVersion.query.order_by(desc(GameVersion.id)).all()
        versions = []
        for version in versions_raw:
            versions.append(game_version_info(version))
        games_raw = Game.query.filter(Game.active == True).order_by(desc(Game.id)).all()
        games = []
        for game in games_raw:
            games.append(game_info(game))
        publishers_raw = Publisher.query.order_by(desc(Publisher.id)).all()
        publishers = []
        for publisher in publishers_raw:
            publishers.append(publisher_info(publisher))
        return jsonify({'user_count': user_count, 'mods': mods, 'users': users, 'versions': versions, 'games': games, 'publishers': publishers})
     
    backend.api_path="/api/admin"
    
    def __init__(self, db, email):
        self.db = db
        self.email = email
    
    @adminrequired
    def impersonate(self, username):
        """
        Log into another persons account from an admin account
        """
        user = User.query.filter(User.username == username).first()
        if not user:
            return jsonify({'error': True, 'reason': 'User does not exist'})
        login_user(user)
        return jsonify({'error': False})
    
    impersonate.api_path = "/admin/impersonate/<username>"
    
    @adminrequired
    @with_session
    def create_version(self):
        """
        Create a game version.
        Required fields: game_version, game_id, beta (True/False)
        """
        game_version = request.form.get("game_version")
        game_id = request.form.get("game_id")
        beta = request.form.get("beta")
        if not game_version or not game_id or not beta:
            return jsonify({'error': True, 'reason': 'Required fields are missing'})
        game = Game.query.filter(Game.name == game_id).first()
        if not game:
            return jsonify({'error': True, 'reason': 'Game does not exist'})
        gameversion = GameVersion.query.filter(Game.id == game_id).filter(GameVersion.friendly_version == game_version).first()
        if gameversion:
            return jsonify({'error': True, 'reason': 'Game version already exists'})
        isbeta = False
        if beta.lower() == "true" or beta == "1":
            isbeta = True
        version = GameVersion(game_version, game_id, isbeta)
        self.db.add(version)
        self.db.commit()
        return jsonify({'error': False})
    
    create_version.api_path = "/admin/versions/create"
    create_version.methods = ['POST']
    
    @adminrequired
    @with_session
    def create_game(self):
        """
        Creates a game.
        Required fields: game_name, game_short, publisher_id
        """
        game_name = request.form.get("game_name")
        game_short = request.form.get("game_short")
        publisher_id = request.form.get("publisher_id")
        if not game_name or not publisher_id or not game_short:
            return jsonify({'error': True, 'reason': 'Required fields are missing'})
        publisher = Publisher.query.filter(Publisher.id == publisher_id).first()
        if not publisher:
            return jsonify({'error': True, 'reason': 'Publisher does not exist'})
        game = Game.query.filter(Game.name == game_name).first()
        if game:
            return jsonify({'error': True, 'reason': 'Game already exists'})
        go = Game(game_name, publisher_id, game_short)
        self.db.add(go)
        self.db.commit()
        return jsonify({'error': False})
    
    create_game.api_path = "/admin/games/create"
    create_game.methods = ['POST']
    
    @adminrequired
    @with_session
    def create_publisher(self):
        """
        Creates a publisher.
        Required fields: publisher_name
        """
        publisher_name = request.form.get("publisher_name")
        if not publisher_name:
            return jsonify({'error': True, 'reason': 'Required fields are missing'})
        publisher = Publisher.query.filter(Publisher.name == publisher_name).first()
        if publisher:
            return jsonify({'error': True, 'reason': 'Publisher already exists'})
        gname = Publisher(publisher_name)
        self.db.add(gname)
        self.db.commit()
        return jsonify({'error': False})
    
    create_publisher.api_path = "/admin/publishers/create"
    create_publisher.methods = ['POST']
    
    
    @adminrequired
    def email(self):
        """
        Emails everyone
        Required fields: subject, body
        Optional fields: modders (on/off)
        """
        subject = request.form.get('subject')
        body = request.form.get('body')
        modders_only = request.form.get('modders-only') == 'on'
        if not subject or not body:
            return jsonify({'error': True, 'reason': 'Required fields are missing'})
        if subject == '' or body == '':
            return jsonify({'error': True, 'reason': 'Required data is missing'})
        users = User.query.all()
        if modders_only:
            users = [u for u in users if len(u.mods) != 0 or u.username == current_user.username]
        send_bulk_email([u.email for u in users], subject, body)
        return jsonify({'error': False})
    
    email.api_path = "/admin/email"
    email.methods = ['POST']
    
    @adminrequired
    @with_session
    def manual_confirm(self, user_id):
        user = User.query.filter(User.id == int(user_id)).first()
        if not user:
            return jsonify({'error': True, 'reason': 'User does not exist'})
        user.confirmation = None
        return jsonify({'error': False})
    
    manual_confirm.api_path = "/admin/manual-confirmation/<user_id>"
    