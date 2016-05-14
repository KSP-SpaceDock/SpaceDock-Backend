from flask import request
from flask.ext.login import current_user, login_user, logout_user
from datetime import datetime, timedelta
from SpaceDock.objects import *
from SpaceDock.common import *
from SpaceDock.formatting import user_info

import bcrypt
import re
import random
import base64
import binascii
import os

class AccountEndpoints:
    def __init__(self, cfg, db, email):
        self.cfg = cfg
        self.db = db.get_database()
        self.email = email
    
    @with_session
    def register(self):
        """
        Required parameters: email, username, password, repeatPassword
        Optional parameters: follow-mod
        """
        if not self.cfg.getb('registration'):
            return {'error': True, 'configError': 'Registrations are disabled'}, 400

        followMod = request.form.get('follow-mod')
        email = request.form.get('email')
        username = request.form.get('username')
        password = request.form.get('password')
        confirmPassword = request.form.get('repeatPassword')
        
        usernameErrors = []
        passwordErrors = []
        emailErrors = []
    
        emailError = self.check_email_for_registration(email)
        if emailError:
            emailErrors.append(emailError)
    
        usernameError = self.check_username_for_registration(username)
        if usernameError:
            usernameErrors.append(usernameError)
    
        if not password:
            passwordErrors.append('Password is required.')
        else:
            if password != confirmPassword:
                passwordErrors.append('Passwords do not match.')
            if len(password) < 5:
                passwordErrors.append('Your password must be greater than 5 characters.')
            if len(password) > 256:
                passwordErrors.append('We admire your dedication to security, but please use a shorter password.')
        
        if len(usernameErrors) > 0 or len(passwordErrors) > 0 or len(emailErrors) > 0:
            return {'error': True, 'usernameErrors': usernameErrors, 'passwordErrors': passwordErrors, 'emailErrors': emailErrors}, 400

        # All valid, let's make them an account
        user = User(username, email, password)
        user.confirmation = binascii.b2a_hex(os.urandom(20)).decode("utf-8")
        self.db.add(user)
        self.db.commit() # We do this manually so that we're sure everything's hunky dory before the email leaves
        if followMod:
            self.email.send_confirmation(user, followMod)
        else:
            self.email.send_confirmation(user)
        return {'error': False}

    register.api_path = "/api/register"
    register.methods = ['POST']    
    
    def check_username_for_registration(self, username):
        if not username:
            return 'Username is required.'
        if not re.match(r"^[A-Za-z0-9_]+$", username):
            return 'Please only use letters, numbers, and underscores.'
        if len(username) < 3 or len(username) > 24:
            return 'Usernames must be between 3 and 24 characters.'
        if User.query.filter(User.username.ilike(username)).first():
            return 'A user by this name already exists.'
        return None
    
    def check_email_for_registration(self, email):
        if not email:
            return 'Email is required.'
        if not re.match(r"^[^@]+@[^@]+\.[^@]+$", email):
            return 'Please specify a valid email address.'
        elif User.query.filter(User.email == email).first():
            return 'A user with this email already exists.'
        return None
    
    @with_session
    def confirm(self, username, confirmation):
        """
        Confirms the user. The confirmation must match what was sent to the user.
        """
        user = User.query.filter(User.username == username).first()
        if not user:
            return {"error": True, "confirmError": "User does not exist"}, 400
        if user.confirmation == None:
            return {"error": True, "confirmError": "User already confirmed"}, 400
        if user.confirmation != confirmation:
            return {"error": True, "confirmError": "Confirmation does not match"}, 400

        user.confirmation = None
        login_user(user)
        f = request.args.get('f')
        if f:
            mod = Mod.query.filter(Mod.id == int(f)).first()
            mod.follower_count += 1
            user.following.append(mod)
        return {"error": False}

    confirm.api_path = "/api/confirm/<username>/<confirmation>"
    
    
    def login(self):
        """
        Required fields: username, password
        """
        username = request.form['username']
        password = request.form['password']
        if not username or not password:
            return { 'error': True, 'reason': 'Missing username or password' }, 400
        user = User.query.filter(User.username.ilike(username)).first()
        if not user:
            return { 'error': True, 'reason': 'Username or password is incorrect' }, 400
        if not bcrypt.hashpw(password.encode('utf-8'), user.password.encode('utf-8')) == user.password.encode('utf-8'):
            return { 'error': True, 'reason': 'Username or password is incorrect' }, 400
        if user.confirmation != '' and user.confirmation != None:
            return { 'error': True, 'reason': 'User is not confirmed' }, 400
        login_user(user)
        return {'error': False}
    
    login.api_path = "/api/login"
    login.methods = ['POST']
    
    def logout(self):
        logout_user()
        return { 'error': False }
    
    logout.api_path = "/api/logout"
    
    def get_current_user(self):
        if not user:
            return { 'error': True, 'reason': 'User is not confirmed' }, 400
        return { 'error': False, 'current_user': user_info(current_user) }
    
    login.api_path = "/api/login"
    login.methods = ['POST']
    
    @with_session
    def forgot_password(self):
        """
        Required fields: email
        """
        email = request.form.get('email')
        if not email:
            return { 'error': True, 'reason': 'No email address' }, 400
        user = User.query.filter(User.email == email).first()
        if not user:
            return { 'error': True, 'reason': 'No user for provided email address' }, 400
        user.passwordReset = binascii.b2a_hex(os.urandom(20)).decode("utf-8")
        user.passwordResetExpiry = datetime.now() + timedelta(days=1)
        self.db.commit()
        send_reset(user)
        return {'error': False}
        
    forgot_password.api_path = "/api/reset"
    forgot_password.methods = ['POST']
    
    @with_session
    def reset_password(self, username, confirmation):
        """
        Required parameters: password, password2
        """
        user = User.query.filter(User.username == username).first()
        if not user:
            return { 'error': True, 'reason': 'Username is incorrect' }, 400
        
        if user.passwordResetExpiry == None or user.passwordResetExpiry < datetime.now():
            return { 'error': True, 'reason': 'Password reset invalid' }, 400
        if user.passwordReset != confirmation:
            return { 'error': True, 'reason': 'Password reset invalid' }, 400
        password = request.form.get('password')
        password2 = request.form.get('password2')
        if not password or not password2:
            return { 'error': True, 'reason': 'Passwords not provided' }, 400
        if password != password2:
            return { 'error': True, 'reason': 'Passwords do not match' }, 400
        user.set_password(password)
        user.passwordReset = None
        user.passwordResetExpiry = None
        self.db.commit()
        return { 'error': False }
    
    reset_password.api_path = "/reset/<username>/<confirmation>"
    reset_password.methods = ['POST']