from flask import request
from flask_login import current_user
from sqlalchemy import desc
from SpaceDock.common import *
from SpaceDock.formatting import user_info, admin_user_info
from SpaceDock.objects import *
from SpaceDock.routing import route

import json


@route('/api/users')
def get_users():
    """
    Returns a list of users
    """
    users = list()
    for user in User.query.order_by(desc(User.id)).all():
        if has_ability('view-users-full') or (user.public and user.confirmation == None):
            users.append(user_info(user) if not has_ability('view-users-full') else admin_user_info(user))
    return {'error': False, 'count': len(users), 'data': users}

@route('/api/users/<userid>')
def get_user_info(userid):
    """
    Returns more data for one user
    """
    user = None
    if userid == 'current':
        user = current_user
    elif not userid.isdigit() or not User.query.filter(User.id == int(userid)).first():
        return {'error': True, 'reasons': ['The userid is invalid']}, 400
    if not user:
        user = User.query.filter(User.id == int(userid)).first()
    if has_ability('view-users-full') or (user.public and user.confirmation == None):
        return {'error': False, 'count': 1, 'data': (user_info(user) if not has_ability('view-users-full') else admin_user_info(user))}
    else:
        return {'error': True, 'reasons': ['The userid is invalid']}, 400

@route('/api/users/<userid>/edit', methods=['POST'])
@user_has('user-edit', params=['userid'], public=False)
@with_session
def edit_user(userid):
    """
    Edits a user, based on the request parameters. Required fields: data
    """
    if not userid.isdigit() or not User.query.filter(User.id == int(userid)).first():
        return {'error': True, 'reasons': ['The userid is invalid.']}, 400

    # Get variables
    parameters = json.loads(request.json['data'])

    # Get the matching user and edit it
    user = User.query.filter(User.id == int(userid)).first()
    code = edit_object(user, parameters)

    # Error check
    if code == 2:
        return {'error': True, 'reasons': ['You tried to edit a value that doesn\'t exist.']}, 400
    elif code == 1:
        return {'error': True, 'reasons': ['You tried to edit a value that is marked as read-only.']}, 400
    else:
        return {'error': False, 'count': 1, 'data': user_info(user)}

@route('/api/users/<userid>/update-bg', methods=['POST'])
@user_has('user-edit', params=['userid'], public=False)
@with_session
def user_updateBG(userid):
    """
    Updates a users background. Required fields: image
    """
    errors = list()
    if not userid.isdigit() or not User.query.filter(User.id == int(userid)).first():
        rerrors.append('The userid is invalid.')
    if not request.files.get('image'):
        errors.append('The background is invalid.')
    if any(errors):
        return {'error': True, 'reasons': errors}, 400

    # Find the user
    user = User.query.filter(User.id == int(userid)).first()

    # Get the file and save it to disk
    f = request.files['image']
    filetype = os.path.splitext(os.path.basename(f.filename))[1]
    if not filetype in ['.png', '.jpg']:
        return {'error': True, 'reasons': ['This file type is not acceptable.']}, 400
    filename = secure_filename(user.username) + filetype
    base_path = os.path.join(secure_filename(user.username) + '-' + str(time.time()) + '_' + str(user.id))
    full_path = os.path.join(cfg['storage'], base_path)
    if not os.path.exists(full_path):
        os.makedirs(full_path)
    path = os.path.join(full_path, filename)
    try:
        os.remove(os.path.join(cfg['storage'], user.backgroundMedia))
    except:
        pass # who cares
    f.save(path)
    user.backgroundMedia = os.path.join(base_path, filename)
    return {'error': False, 'count': 1, 'data': user_info(user)}

    f = request.files['image']
    filetype = os.path.splitext(os.path.basename(f.filename))[1]
    if not filetype in ['.png', '.jpg']:
        return {'error': True, 'reasons': ['This file type is not acceptable.']}, 400
    filename = secure_filename(mod.name) + '-' + str(time.time()) + filetype
    base_path = os.path.join(secure_filename(user.username) + '_' + str(mod.user.id), secure_filename(mod.name))
    full_path = os.path.join(cfg['storage'], base_path)
    if not os.path.exists(full_path):
        os.makedirs(full_path)
    path = os.path.join(full_path, filename)
    try:
        os.remove(os.path.join(cfg['storage'], user.backgroundMedia))
    except:
        pass # who cares
    f.save(path)
    user.backgroundMedia = os.path.join(base_path, filename)
    return {'error': False, 'count': 1, 'data': user_info(user)}