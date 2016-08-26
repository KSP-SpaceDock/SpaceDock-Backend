from flask import request, session, url_for, current_app
from flask.ext.login import current_user, login_user
from sqlalchemy import desc
from SpaceDock.objects import *
from SpaceDock.common import *
from SpaceDock.celery import notify_ckan
from SpaceDock.formatting import *
from datetime import datetime
from flask_json import FlaskJSON, JsonError, json_response, as_json, as_json_p
from flask.ext.cache import Cache

import time
import os
import zipfile
import urllib
import math
import json
import urllib.parse

class ApiEndpoints:

    def __init__(self, cfg, db, email, search):
        self.cfg = cfg
        self.db = db.get_database()
        self.search = search
        self.default_description = """This is your mod listing! You can edit it as much as you like before you make it public.

        To edit **this** text, you can click on the "**Edit this Mod**" button up there.

        By the way, you have a lot of flexibility here. You can embed YouTube videos or screenshots. Be creative.

        You can check out the SpaceDock [markdown documentation](/markdown) for tips.

        Thanks for hosting your mod on SpaceDock!"""

    @as_json
    def browse_mod(self,gameid,filterby):
        count = 30
        if filterby == "top":
            mods = Mod.query.filter(Mod.published,Mod.game_id == gameid).order_by(desc(Mod.created)).limit(count)
        elif filterby == "new":
            mods = Mod.query.filter(Mod.published, Mod.game_id == gameid).order_by(desc(Mod.created)).limit(count)
        else:
            mods = Mod.query.filter(Mod.published, Mod.game_id == gameid).order_by(desc(Mod.created)).limit(count)
        results = list()
        for mod in mods:
            results.append({
                "name": mod.name,
                "id": mod.id,
                "game": mod.game.name,
                "game_id": mod.game_id,
                "game_short": mod.game.short,
                "short_description": mod.short_description,
                "downloads": mod.download_count,
                "followers": mod.follower_count,
                "author": mod.user.username,
                "default_version_id": mod.default_version().id,
                "shared_authors": list(),
                "background": mod.background,
                "license": mod.license,
                "website": mod.external_link,
                "donations": mod.donation_link,
                "source_code": mod.source_link,
                "url": "/mod/"+ str(mod.id) + "/" + urllib.parse.quote_plus(mod.name)
            })
        return {"count":len(results),"data":results}

    browse_mod.api_path = "/api/game/<int:gameid>/browse/mod/<filterby>"


    @as_json
    def typeahead_mod(self):
        query = request.args.get('query')
        page = request.args.get('page')
        query = '' if not query else query
        page = 1 if not page or not page.isdigit() else int(page)
        results = list()
        for m in self.search.typeahead_mods(query):
            a = mod_info(m)
            results.append(a)
        return results

    typeahead_mod.api_path = "/api/typeahead/mod"


    @as_json
    def search_mod(self):
        query = request.args.get('query')
        page = request.args.get('page')
        query = '' if not query else query
        page = 1 if not page or not page.isdigit() else int(page)
        results = list()
        for m in self.search.search_mods(None, query, page, 30)[0]:
            a = mod_info(m)
            results.append(a)
        return results

    search_mod.api_path = "/api/search/mod"


    @as_json
    def search_user(self):
        query = request.args.get('query')
        page = request.args.get('page')
        query = '' if not query else query
        page = 0 if not page or not page.isdigit() else int(page)
        results = list()
        for u in self.search.search_users(query, page):
            a = user_info(u)
            a['mods'] = list()
            mods = Mod.query.filter(Mod.user == u, Mod.published == True).order_by(Mod.created)
            for m in mods:
                a['mods'].append(mod_info(m))
            results.append(a)
        return results

    search_user.api_path = "/api/search/user"


    @as_json
    def mod(self, modid):
        mod = Mod.query.filter(Mod.id == modid).first()
        if not mod:
            return { 'error': True, 'reason': 'Mod not found.' }, 404
        if not mod.published:
            return { 'error': True, 'reason': 'Mod not published.' }, 401
        info = mod_info(mod)
        #TODO: Fixme
        #info["description_html"] = str(current_app.jinja_env.filters['markdown'](mod.description))
        return info

    mod.api_path = "/api/mod/<int:modid>"


    @as_json
    def mod_version(self, modid, version):
        mod = Mod.query.filter(Mod.id == modid).first()
        if not mod:
            return { 'error': True, 'reason': 'Mod not found.' }, 404
        if not mod.published:
            return { 'error': True, 'reason': 'Mod not published.' }, 401
        if version == "latest" or version == "latest_version":
            v = mod.default_version()
        elif version.isdigit():
            v = ModVersion.query.filter(ModVersion.mod == mod,
                                        ModVersion.id == int(version)).first()
        else:
            return { 'error': True, 'reason': 'Invalid version.' }, 400
        if not v:
            return { 'error': True, 'reason': 'Version not found.' }, 404
        info = version_info(mod, v)
        return info

    mod_version.api_path = "/api/mod/<int:modid>/<version>"


    @as_json
    def user(self, username):
        user = User.query.filter(User.username == username).first()
        if not user:
            return { 'error': True, 'reason': 'User not found.' }, 404
        if not user.public:
            return { 'error': True, 'reason': 'User not public.' }, 401
        mods = Mod.query.filter(Mod.user == user, Mod.published == True).order_by(
            Mod.created)
        info = user_info(user)
        info['mods'] = list()
        for m in mods:
            info['mods'].append(mod_info(m))
        return info

    user.api_path = "/api/user/<username>"


    @as_json
    @with_session
    def update_mod_background(self, mod_id):
        if current_user == None:
            return { 'error': True, 'reason': 'You are not logged in.' }, 401
        mod = Mod.query.filter(Mod.id == mod_id).first()
        if not mod:
            return { 'error': True, 'reason': 'Mod not found.' }, 404
        editable = False
        if current_user:
            if current_user.admin:
                editable = True
            if current_user.id == mod.user_id:
                editable = True
            if any([u.accepted and u.user == current_user for u in mod.shared_authors]):
                editable = True
        if not editable:
            return { 'error': True, 'reason': 'Not enought rights.' }, 401
        f = request.files['image']
        filetype = os.path.splitext(os.path.basename(f.filename))[1]
        if not filetype in ['.png', '.jpg']:
            return { 'error': True, 'reason': 'This file type is not acceptable.' }, 400
        filename = secure_filename(mod.name) + '-' + str(time.time()) + filetype
        base_path = os.path.join(secure_filename(mod.user.username) + '_' + str(mod.user.id), secure_filename(mod.name))
        full_path = os.path.join(self.cfg['storage'], base_path)
        if not os.path.exists(full_path):
            os.makedirs(full_path)
        path = os.path.join(full_path, filename)
        try:
            os.remove(os.path.join(self.cfg['storage'], mod.background))
        except:
            pass # who cares
        f.save(path)
        mod.background = os.path.join(base_path, filename)
        return { 'path': '/content/' + mod.background }

    update_mod_background.methods = ['POST']
    update_mod_background.api_path = "/api/mod/<mod_id>/update-bg"


    @as_json
    @with_session
    def update_user_background(self, username):
        if current_user == None:
            return { 'error': True, 'reason': 'You are not logged in.' }, 401
        user = User.query.filter(User.username == username).first()
        if not current_user.admin and current_user.username != user.username:
            return { 'error': True, 'reason': 'You are not authorized to edit this user\'s background' }, 403
        f = request.files['image']
        filetype = os.path.splitext(os.path.basename(f.filename))[1]
        if not filetype in ['.png', '.jpg']:
            return { 'error': True, 'reason': 'This file type is not acceptable.' }, 400
        filename = secure_filename(user.username) + filetype
        base_path = os.path.join(secure_filename(user.username) + '-' + str(time.time()) + '_' + str(user.id))
        full_path = os.path.join(self.cfg['storage'], base_path)
        if not os.path.exists(full_path):
            os.makedirs(full_path)
        path = os.path.join(full_path, filename)
        try:
            os.remove(os.path.join(self.cfg['storage'], user.backgroundMedia))
        except:
            pass # who cares
        f.save(path)
        user.backgroundMedia = os.path.join(base_path, filename)
        return { '/content/' + user.backgroundMedia }

    update_user_background.methods = ['POST']
    update_user_background.api_path = "/api/user/<username>/update-bg"


    @as_json
    @with_session
    def grant_mod(self, mod_id):
        mod = Mod.query.filter(Mod.id == mod_id).first()
        if not mod:
            return { 'error': True, 'reason': 'Mod not found.' }, 404
        editable = False
        if current_user:
            if current_user.admin:
                editable = True
            if current_user.id == mod.user_id:
                editable = True
        if not editable:
            return { 'error': True, 'reason': 'Not enought rights.' }, 401
        new_user = request.json.get('user')
        new_user = User.query.filter(User.username.ilike(new_user)).first()
        if new_user == None:
            return { 'error': True, 'reason': 'The specified user does not exist.' }, 400
        if mod.user == new_user:
            return { 'error': True, 'reason': 'This user has already been added.' }, 400
        if any(m.user == new_user for m in mod.shared_authors):
            return { 'error': True, 'reason': 'This user has already been added.' }, 400
        if not new_user.public:
            return { 'error': True, 'reason': 'This user has not made their profile public.' }, 400
        author = SharedAuthor()
        author.mod = mod
        author.user = new_user
        mod.shared_authors.append(author)
        db.add(author)
        db.commit()
        self.email.send_grant_notice(mod, new_user)
        return { 'error': False }, 200

    grant_mod.methods = ['POST']
    grant_mod.api_path = "/api/mod/<mod_id>/grant"


    @as_json
    @with_session
    def accept_grant_mod(self, mod_id):
        if current_user == None:
            return { 'error': True, 'reason': 'You are not logged in.' }, 401
        mod = Mod.query.filter(Mod.id == mod_id).first()
        if not mod:
            return { 'error': True, 'reason': 'Mod not found.' }, 404
        author = [a for a in mod.shared_authors if a.user == current_user]
        if len(author) == 0:
            return { 'error': True, 'reason': 'You do not have a pending authorship invite.' }, 200
        author = author[0]
        if author.accepted:
            return { 'error': True, 'reason': 'You do not have a pending authorship invite.' }, 200
        author.accepted = True
        return { 'error': False }, 200

    accept_grant_mod.methods = ['POST']
    accept_grant_mod.api_path = "/api/mod/<mod_id>/accept_grant"

    @as_json
    @with_session
    def reject_grant_mod(self, mod_id):
        if current_user == None:
            return { 'error': True, 'reason': 'You are not logged in.' }, 401
        mod = Mod.query.filter(Mod.id == mod_id).first()
        if not mod:
            return { 'error': True, 'reason': 'Mod not found.' }, 404
        author = [a for a in mod.shared_authors if a.user == current_user]
        if len(author) == 0:
            return { 'error': True, 'reason': 'You do not have a pending authorship invite.' }, 200
        author = author[0]
        if author.accepted:
            return { 'error': True, 'reason': 'You do not have a pending authorship invite.' }, 200
        mod.shared_authors = [a for a in mod.shared_authors if a.user != current_user]
        db.delete(author)
        return { 'error': False }, 200

    reject_grant_mod.methods = ['POST']
    reject_grant_mod.api_path = "/api/mod/<mod_id>/reject_grant"


    @as_json
    @with_session
    def revoke_mod(self, mod_id):
        if current_user == None:
            return { 'error': True, 'reason': 'You are not logged in.' }, 401
        mod = Mod.query.filter(Mod.id == mod_id).first()
        if not mod:
            return { 'error': True, 'reason': 'Mod not found.' }, 404
        editable = False
        if current_user:
            if current_user.admin:
                editable = True
            if current_user.id == mod.user_id:
                editable = True
        if not editable:
            return { 'error': True, 'reason': 'Not enought rights.' }, 401
        new_user = request.json.get('user')
        new_user = User.query.filter(User.username.ilike(new_user)).first()
        if new_user == None:
            return { 'error': True, 'reason': 'The specified user does not exist.' }, 404
        if mod.user == new_user:
            return { 'error': True, 'reason': 'You can\'t remove yourself.' }, 400
        if not any(m.user == new_user for m in mod.shared_authors):
            return { 'error': True, 'reason': 'This user is not an author.' }, 400
        author = [a for a in mod.shared_authors if a.user == new_user][0]
        mod.shared_authors = [a for a in mod.shared_authors if a.user != current_user]
        db.delete(author)
        return { 'error': False }, 200

    revoke_mod.methods = ['POST']
    revoke_mod.api_path = "/api/mod/<mod_id>/revoke"


    @as_json
    @with_session
    def set_default_version(self, mid, vid):
        mod = Mod.query.filter(Mod.id == mid).first()
        if not mod:
            return { 'error': True, 'reason': 'The specified mod does not exist.' }, 404
        editable = False
        if current_user:
            if current_user.admin:
                editable = True
            if current_user.id == mod.user_id:
                editable = True
            if any([u.accepted and u.user == current_user for u in mod.shared_authors]):
                editable = True
        if not editable:
            return { 'error': True, 'reason': 'You do not have permission to do this.' }, 400
        if not any([v.id == vid for v in mod.versions]):
            return { 'error': True, 'reason': 'This mod does not have the specified version.' }, 404
        mod.default_version_id = vid
        return { 'error': False }, 200

    set_default_version.methods = ['POST']
    set_default_version.api_path = "/api/mod/<int:mid>/set-default/<int:vid>"


    @as_json
    @with_session
    def create_list(self):
        if not current_user:
            return { 'error': True, 'reason': 'You are not logged in.' }, 401
        if not current_user.public:
            return { 'error': True, 'reason': 'Only users with public profiles may create mod packs.' }, 403
        name = request.json.get('name')
        if not name:
            return { 'error': True, 'reason': 'All fields are required.' }, 400
        game = request.json.get('game')
        if not game:
            return {'error': True, 'reason': 'Please select a game.'}, 400
        if len(name) > 100:
            return { 'error': True, 'reason': 'Fields exceed maximum permissible length.' }, 400
        mod_list = ModList()
        mod_list.name = name
        mod_list.user = current_user
        mod_list.game_id = game
        db.add(mod_list)
        db.commit()
        return { 'url': url_for("lists.view_list", list_id=mod_list.id, list_name=mod_list.name) }

    create_list.methods = ['POST']
    create_list.api_path = "/api/pack/create"


    @as_json
    @with_session
    def create_mod(self):
        if not current_user:
            return { 'error': True, 'reason': 'You are not logged in.' }, 401
        if not current_user.public:
            return { 'error': True, 'reason': 'Only users with public profiles may create mods.' }, 403
        name = request.json.get('name')
        game = request.json.get('game')
        short_description = request.json.get('short-description')
        version = request.json.get('version')
        game_version = request.json.get('game-version')
        license = request.json.get('license')
        ckan = request.json.get('ckan')
        zipball = request.files.get('zipball')
        # Validate
        if not name \
            or not short_description \
            or not version \
            or not game \
            or not game_version \
            or not license \
            or not zipball:
            return { 'error': True, 'reason': 'All fields are required.' }, 400
        # Validation, continued
        if len(name) > 100 \
            or len(short_description) > 1000 \
            or len(license) > 128:
            return { 'error': True, 'reason': 'Fields exceed maximum permissible length.' }, 400
        if ckan == None:
            ckan = False
        else:
            ckan = (ckan.lower() == "true" or ckan.lower() == "yes" or ckan.lower() == "on")
        test_game = Game.query.filter(Game.id == game).first()
        if not test_game:
            return { 'error': True, 'reason': 'Game does not exist.' }, 400
        test_gameversion = GameVersion.query.filter(GameVersion.game_id == test_game.id).filter(GameVersion.friendly_version == game_version).first()
        if not test_gameversion:
            return { 'error': True, 'reason': 'Game version does not exist.' }, 400
        game_version_id = test_gameversion.id
        mod = Mod()
        mod.user = current_user
        mod.name = name
        mod.game_id = game
        mod.short_description = short_description
        mod.description = self.default_description
        mod.ckan = ckan
        mod.license = license
        # Save zipball
        filename = secure_filename(name) + '-' + secure_filename(version) + '.zip'
        base_path = os.path.join(secure_filename(current_user.username) + '_' + str(current_user.id), secure_filename(name))
        full_path = os.path.join(self.cfg['storage'], base_path)
        if not os.path.exists(full_path):
            os.makedirs(full_path)
        path = os.path.join(full_path, filename)
        if os.path.isfile(path):
            # We already have this version
            # We'll remove it because the only reason it could be here on creation is an error
            os.remove(path)
        zipball.save(path)
        if not zipfile.is_zipfile(path):
            os.remove(path)
            return { 'error': True, 'reason': 'This is not a valid zip file.' }, 400
        version = ModVersion(secure_filename(version), game_version_id, os.path.join(base_path, filename))
        mod.versions.append(version)
        db.add(version)
        # Save database entry
        db.add(mod)
        db.commit()
        mod.default_version_id = version.id
        db.commit()
        ga = Game.query.filter(Game.id == game).first()
        session['game'] = ga.id;
        notify_ckan.delay(mod.id, 'create')
        return { 'url': url_for("mods.mod", id=mod.id, mod_name=mod.name), "id": mod.id, "name": mod.name }

    create_mod.methods = ['POST']
    create_mod.api_path = "/api/mod/create"


    @as_json
    @with_session
    def update_mod(self, mod_id):
        if current_user == None:
            return { 'error': True, 'reason': 'You are not logged in.' }, 401
        mod = Mod.query.filter(Mod.id == mod_id).first()
        if not mod:
            return { 'error': True, 'reason': 'Mod not found.' }, 404
        editable = False
        if current_user:
            if current_user.admin:
                editable = True
            if current_user.id == mod.user_id:
                editable = True
            if any([u.accepted and u.user == current_user for u in mod.shared_authors]):
                editable = True
        if not editable:
            return { 'error': True, 'reason': 'Not enought rights.' }, 401
        version = request.json.get('version')
        changelog = request.json.get('changelog')
        game_version = request.json.get('game-version')
        notify = request.json.get('notify-followers')
        zipball = request.files.get('zipball')
        if not version \
            or not game_version \
            or not zipball:
            # Client side validation means that they're just being pricks if they
            # get here, so we don't need to show them a pretty error reason
            # SMILIE: this doesn't account for "external" API use --> return a json error
            return { 'error': True, 'reason': 'All fields are required.' }, 400
        test_gameversion = GameVersion.query.filter(GameVersion.game_id == Mod.game_id).filter(GameVersion.friendly_version == game_version).first()
        if not test_gameversion:
            return { 'error': True, 'reason': 'Game version does not exist.' }, 400
        game_version_id = test_gameversion.id
        if notify == None:
            notify = False
        else:
            notify = (notify.lower() == "true" or notify.lower() == "yes")
        filename = secure_filename(mod.name) + '-' + secure_filename(version) + '.zip'
        base_path = os.path.join(secure_filename(current_user.username) + '_' + str(current_user.id), secure_filename(mod.name))
        full_path = os.path.join(self.cfg['storage'], base_path)
        if not os.path.exists(full_path):
            os.makedirs(full_path)
        path = os.path.join(full_path, filename)
        for v in mod.versions:
            if v.friendly_version == secure_filename(version):
                return { 'error': True, 'reason': 'We already have this version. Did you mistype the version number?' }, 400
        if os.path.isfile(path):
            os.remove(path)
        zipball.save(path)
        if not zipfile.is_zipfile(path):
            os.remove(path)
            return { 'error': True, 'reason': 'This is not a valid zip file.' }, 400
        version = ModVersion(secure_filename(version), game_version_id, os.path.join(base_path, filename))
        version.changelog = changelog
        # Assign a sort index
        if len(mod.versions) == 0:
            version.sort_index = 0
        else:
            version.sort_index = max([v.sort_index for v in mod.versions]) + 1
        mod.versions.append(version)
        mod.updated = datetime.now()
        if notify:
            self.email.send_update_notification(mod, version, current_user)
        db.add(version)
        db.commit()
        mod.default_version_id = version.id
        db.commit()
        notify_ckan.delay(mod_id, 'update')
        return { 'url': url_for("mods.mod", id=mod.id, mod_name=mod.name), "id": version.id  }

    update_mod.methods = ['POST']
    update_mod.api_path = "/api/mod/<mod_id>/update"