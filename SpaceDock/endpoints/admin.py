from flask import request
from flask_login import current_user, login_user
from SpaceDock.common import user_has, with_session
from SpaceDock.database import db
from SpaceDock.email import send_bulk_email
from SpaceDock.objects import Role, User
from SpaceDock.routing import route

@route('/api/admin/impersonate/<userid>')
@user_has('admin-impersonate', params=['userid']) # I feel like adding this as a param could be useful
def impersonate(userid):
    """
    Log into another persons account from an admin account
    """
    if not userid.isdigit() or not User.query.filter(User.id == int(userid)).first():
        return {'error': True, 'reasons': ['The userid is invalid'], 'codes': ['2145']}, 400
    user = User.query.filter(User.id == int(userid)).first()
    login_user(user)
    return {'error': False}

@route('/api/admin/email', methods=['POST'])
@user_has('admin-email')
def email():
    """
    Emails everyone. Required fields: subject, body. Optional fields: modders-only
    """
    subject = request.json.get('subject')
    body = request.json.get('body')
    modders_only = request.json.get('modders-only')
    if not isinstance(modders_only, bool):
        return {'error': True, 'reasons': ['"modders_only" is invalid']}, 400
    if not subject or not body:
        return {'error': True, 'reasons': ['Required fields are missing'], 'codes': ['2510']}
    if subject == '' or body == '':
        return {'error': True, 'reasons': ['Required data is missing'], 'codes': ['2500']}
    users = User.query.all()
    if modders_only:
        users = [u for u in users if len(u.mods) != 0 or u.username == current_user.username]
    send_bulk_email([u.email for u in users], subject, body)
    return {'error': False}

@route('/api/admin/manual-confirmation/<userid>')
@user_has('admin-confirm')
@with_session
def manual_confirm(userid):
    if not userid.isdigit() or not User.query.filter(User.id == int(userid)).first():
        return {'error': True, 'reasons': ['The userid is invalid'], 'codes': ['2145']}, 400
    user = User.query.filter(User.id == int(userid)).first()
    user.confirmation = None
    user.add_roles(user.username)
    role = Role.query.filter(Role.name == user.username).first()
    role.add_abilities('user-edit', 'mods-add', 'packs-add', 'logged-in')
    role.add_param('user-edit', 'userid', user.id)
    role.add_param('mods-add', 'gameshort', '*.')
    db.add(role)
    return {'error': False}
