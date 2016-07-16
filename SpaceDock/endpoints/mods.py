from SpaceDock.objects import *
from SpaceDock.common import *
from SpaceDock.formatting import mod_info
from sqlalchemy import desc
from flask_login import current_user
import json

class ModEndpoints:
    def __init__(self, cfg, db):
        self.db = db.get_database()
        self.cfg = cfg

    def mod_list(self):
        """
        Returns a list of all mods
        """
        results = list()
        for mod in Mod.query.order_by(desc(Mod.id)).filter(Mod.published):
            results.append(mod_info(mod))
        return {'error': False, 'count': len(results), 'data': results}

    mod_list.api_path = '/api/mods'

    def mod_info(self, modid):
        """
        Returns information for one mod
        """
        if not modid.isdigit() or not Mod.query.filter(Mod.id == int(modid)).first():
            return {'error': True, 'reasons': ['The modid is invalid']}, 400
        # Get the mod
        mod = Mod.query.filter(Mod.id == int(modid)).first()
        return {'error': False, 'count': 1, 'data': mod_info(mod)}

    mod_info.api_path = '/api/mods/<modid>'

    @with_session
    @user_has('mods-edit', params=['modid'])
    def mod_edit(self, modid):
        """
        Edits a mod, based on the request parameters. Required fields: data
        """
        errors = list()
        if not modid.isdigit() or not Mod.query.filter(Mod.id == int(modid)).first():
            errors.append('The Mod ID is invalid.')
        if not request.form.get('data') or not is_json(request.form.get('data')):
            errors.append('The patch data is invalid.')
        if any(errors):
            return {'error': True, 'reasons': errors}, 400

        # Get variables
        parameters = json.loads(request.form.get('data'))

        # Get the matching mod and edit it
        mod = Mod.query.filter(Mod.id == int(modid)).first()
        edit_object(mod, parameters)
        return {'error': False}

    mod_edit.api_path = '/api/mods/<modid>/edit'
    mod_edit.methods = ['POST']

    @with_session
    @user_has('mods-add', params=['gameshort'])
    def add_mod(self):
        """
        Adds a mod, based on the request parameters. Required fields: name, gameshort, license
        """
        # Get variables
        name = request.form.get('name')
        short = request.form.get('gameshort')
        license = request.form.get('license')

        # Check the vars
        errors = list()
        if not name:
            errors.append('Invalid mod name.')
        if Mod.query.filter(Mod.name == name).first():
            errors.append('A mod with this name does already exist.')
        if not short:
            errors.append('Invalid gameshort.')
        if not license:
            errors.append('Invalid License.')
        if any(errors):
            return {'error': True, 'reasons': errors}, 400

        # Add new mod
        mod = Mod(name, current_user.id, game_id(short), license)
        return {'error': False}

    add_mod.api_path = '/api/mods/add'
    add_mod.methods = ['POST']

    @with_session
    @user_has('mods-remove', params=['gameshort', 'name']) # We might want to allow deletion of own mods. Gameshort is here to allow per-game moderators.
    def remove_mod(self):
        """
        Removes a mod, based on the request parameters. Required fields: name, gameshort
        """
        # Get variables
        name = request.form.get('name')
        short = request.form.get('gameshort')

        # Check the vars
        errors = list()
        if not name:
            errors.append('Invalid mod name.')
        if not short:
            errors.append('Invalid gameshort.')
        if name and short and not Mod.query.filter(Mod.name == name).filter(Mod.game_id == game_id(short)).first():
            errors.append('A mod with theese parameters does not exist.')
        if any(errors):
            return {'error': True, 'reasons': errors}, 400

        # Add new mod
        mod = Mod.query.filter(Mod.name == name).filter(Mod.game_id == game_id(short)).first()
        self.db.remove(mod)
        return {'error': False}

    remove_mod.api_path = '/api/mods/remove'
    remove_mod.methods = ['POST']