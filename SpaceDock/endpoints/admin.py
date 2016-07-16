from sqlalchemy import desc
from SpaceDock.objects import User, Mod, GameVersion, Game, Publisher
from SpaceDock.common import *
from SpaceDock.formatting import *
from flask_login import current_user, login_user, logout_user

class AdminEndpoints:
    def __init__(self, db, email):
        self.db = db.get_database()
        self.email = email

    @user_has('admin-impersonate', params=['userid']) # I feel like adding this as a param could be useful
    def impersonate(self, userid):
        """
        Log into another persons account from an admin account
        """
        if not userid.isdigit() or not User.query.filter(User.id == int(userid)).first():
            return {'error': True, 'reasons': ['The userid is invalid']}, 400
        user = User.query.filter(User.id == int(userid)).first()
        login_user(user)
        return {'error': False}

    impersonate.api_path = "/api/admin/impersonate/<userid>"

    @user_has('admin-email')
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
            return {'error': True, 'reason': 'Required fields are missing'}
        if subject == '' or body == '':
            return {'error': True, 'reason': 'Required data is missing'}
        users = User.query.all()
        if modders_only:
            users = [u for u in users if len(u.mods) != 0 or u.username == current_user.username]
        self.email.send_bulk_email([u.email for u in users], subject, body)
        return {'error': False}

    email.api_path = "/admin/email"
    email.methods = ['POST']

    @user_has('admin-confirm')
    @with_session
    def manual_confirm(self, userid):
        if not userid.isdigit() or not User.query.filter(User.id == int(userid)).first():
            return {'error': True, 'reasons': ['The userid is invalid']}, 400
        user = User.query.filter(User.id == int(userid)).first()
        user.confirmation = None
        role.add_abilities('user-edit', 'mods-add')
        role.add_param('user-edit', 'userid', user.id)
        role.add_param('mods-add', 'gameshort', '*.')
        return {'error': False}

    manual_confirm.api_path = "/admin/manual-confirmation/<userid>"
