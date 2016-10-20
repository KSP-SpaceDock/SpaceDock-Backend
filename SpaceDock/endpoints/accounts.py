from datetime import datetime, timedelta
from flask import request
from flask_login import current_user, login_user, logout_user
from SpaceDock.common import with_session
from SpaceDock.config import cfg
from SpaceDock.database import db
from SpaceDock.email import send_confirmation, send_reset
from SpaceDock.formatting import user_info
from SpaceDock.objects import Mod, Role, User
from SpaceDock.routing import route

import bcrypt
import re
import binascii
import os

@route('/api/register', methods=['POST'])
@with_session
def register():
    """
    Registers a new useraccount. Required parameters: email, username, password, repeatPassword. Optional parameters: follow-mod
    """
    if not cfg.getb('registration'):
        return {'error': True, 'reasons': ['Registrations are disabled'], 'codes': ['3010']}, 400

    followMod = request.json.get('follow-mod')
    email = request.json.get('email')
    username = request.json.get('username')
    password = request.json.get('password')
    confirmPassword = request.json.get('repeatPassword')

    errors = ()
    codes = ()
    emailError = check_email_for_registration(email)
    if emailError:
        errors.append(emailError)
        codes.append('4000')

    usernameError = check_username_for_registration(username)
    if usernameError:
        errors.append(usernameError)
        codes.append('4000')

    if not password:
        errors.append('Password is required.')
        codes.append('2515')
    else:
        if password != confirmPassword:
            errors.append('Passwords do not match.')
            codes.append('3005')
        if len(password) < 5:
            errors.append('Your password must be greater than 5 characters.')
            codes.append('2101')
        if len(password) > 256:
            errors.append('We admire your dedication to security, but please use a shorter password.')
            codes.append('2102')

    if len(errors) > 0:
        return {'error': True, 'reasons': errors, 'codes': codes}, 400

    # All valid, let's make them an account
    user = User(username, email, password)
    user.confirmation = binascii.b2a_hex(os.urandom(20)).decode("utf-8")
    db.add(user)
    db.commit() # We do this manually so that we're sure everything's hunky dory before the email leaves
    if followMod:
        send_confirmation(user, followMod)
    else:
        send_confirmation(user)
    return {'error': False, 'count': 1, 'data': user_info(user)}

def check_username_for_registration(username):
    if not username:
        return 'Username is required.'
    if not re.match(r"^[A-Za-z0-9_]+$", username):
        return 'Please only use letters, numbers, and underscores.'
    if len(username) < 3 or len(username) > 24:
        return 'Usernames must be between 3 and 24 characters.'
    if User.query.filter(User.username.ilike(username)).first():
        return 'A user by this name already exists.'
    return None

def check_email_for_registration(email):
    if not email:
        return 'Email is required.'
    if not re.match(r"^[^@]+@[^@]+\.[^@]+$", email):
        return 'Please specify a valid email address.'
    elif User.query.filter(User.email == email).first():
        return 'A user with this email already exists.'
    return None

@route('/api/confirm/<username>/<confirmation>')
@with_session
def confirm(username, confirmation):
    """
    Confirms the user. The confirmation must match what was sent to the user.
    """
    user = User.query.filter(User.username == username).first()
    if not user:
        return {'error': True, 'reasons': ['User does not exist'], 'codes': ['2165']}, 400
    if user.confirmation == None:
        return {'error': True, 'reasons': ['User already confirmed'], 'codes': ['3045']}, 400
    if user.confirmation != confirmation:
        return {'error': True, 'reasons': ['Confirmation does not match'], 'codes': ['2100']}, 400

    user.confirmation = None
    login_user(user)
    user.add_roles(username)
    role = Role.query.filter(Role.name == username).first()
    role.add_abilities('user-edit', 'mods-add', 'packs-add', 'logged-in')
    role.add_param('user-edit', 'userid', user.id)
    role.add_param('mods-add', 'gameshort', '.*')
    role.add_param('packs-add', 'gameshort', '.*')
    db.add(role)
    f = request.args.get('f')
    if f:
        mod = Mod.query.filter(Mod.id == int(f)).first()
        mod.follower_count += 1
        user.following.append(mod)
    return {'error': False}

@route('/api/login', methods=['POST'])
def login():
    """
    Login the user to use additional features. Required fields: username, password
    """
    username = request.json.get('username')
    password = request.json.get('password')
    if not username or not password:
        return {'error': True, 'reasons': ['Missing username or password'], 'codes': ['2515']}, 400
    if current_user:
        return {'error': True, 'reasons': ['You are already logged in'], 'codes': ['3060']}, 400
    user = User.query.filter(User.username.ilike(username)).first()
    if not user:
        return {'error': True, 'reasons': ['Username or password is incorrect'], 'codes': ['2175']}, 400
    if not bcrypt.hashpw(password.encode('utf-8'), user.password.encode('utf-8')) == user.password.encode('utf-8'):
        return {'error': True, 'reasons': ['Username or password is incorrect'], 'codes': ['2175']}, 400
    if user.confirmation == '' and user.confirmation == None:
        return {'error': True, 'reasons': ['User is not confirmed'], 'codes': ['3055']}, 400
    login_user(user)
    return {'error': False}

@route('/api/logout')
def logout():
    """
    Closes the session and logs the user out
    """
    if not current_user:
        return {'error': True, 'reasons': ['You are not logged in. Logging out now would be a bit difficult, right?'], 'codes': ['3070']}, 400 # Thats my daily portion of irony
    logout_user()
    return {'error': False}

@route('/api/reset', methods=['POST'])
@with_session
def forgot_password():
    """
    Sends you a confirmation key for resetting your password in case you forgot it. Required fields: email
    """
    email = request.json.get('email')
    if not email:
        return {'error': True, 'reasons': ['No email address'], 'codes': ['2520']}, 400
    user = User.query.filter(User.email == email).first()
    if not user:
        return {'error': True, 'reasons': ['No user for provided email address'], 'codes': ['2115']}, 400
    user.passwordReset = binascii.b2a_hex(os.urandom(20)).decode("utf-8")
    user.passwordResetExpiry = datetime.now() + timedelta(days=1)
    send_reset(user)
    return {'error': False}

@route('/api/reset/<username>/<confirmation>', methods=['POST'])
@with_session
def reset_password(username, confirmation):
    """
    Allows you to reset your password with the key you got through /api/reset. Required parameters: password, password2
    """
    user = User.query.filter(User.username == username).first()
    if not user:
        return {'error': True, 'reasons': ['Username is incorrect'], 'codes': ['2170']}, 400

    if user.passwordResetExpiry == None or user.passwordResetExpiry < datetime.now():
        return {'error': True, 'reasons': ['Password reset invalid'], 'codes': ['3000']}, 400
    if user.passwordReset != confirmation:
        return {'error': True, 'reasons': ['Password reset invalid'], 'codes': ['3000']}, 400
    password = request.json.get('password')
    password2 = request.json.get('password2')
    if not password or not password2:
        return {'error': True, 'reasons': ['Passwords not provided'], 'codes': ['2525']}, 400
    if password != password2:
        return {'error': True, 'reasons': ['Passwords do not match'], 'codes': ['3005']}, 400
    user.set_password(password)
    user.passwordReset = None
    user.passwordResetExpiry = None
    return {'error': False}