import os, os.path, sys
sys.path.append(os.getcwd())

from SpaceDock.app import *
from SpaceDock.database import db
from SpaceDock.objects import *

try:

    for user in User.query.filter(User.confirmation == None).all():
        # Normal role
        user.add_roles(user.username)
        db.flush()    
        role = Role.query.filter(Role.name == user.username).first()
        role.add_abilities('user-edit', 'mods-add', 'packs-add', 'logged-in')
        role.add_param('user-edit', 'userid', user.id)
        role.add_param('mods-add', 'gameshort', '.*')
        role.add_param('packs-add', 'gameshort', '.*')
        db.add(role)
    
        # Mods
        for mod in Mod.query.filter(Mod.user_id == user.id).all():
            user.add_roles(mod.name)
            db.flush()
            role = Role.query.filter(Role.name == mod.name).first()
            role.add_abilities('mods-edit', 'mods-remove')
            role.add_param('mods-edit', 'modid', str(mod.id))
            role.add_param('mods-remove', 'name', mod.name)
            db.add(role)
        
        for shared in SharedAuthor.query.filter(SharedAuthor.user_id == user.id).all():
            mod = Mod.query.filter(Mod.id == shared.mod_id).first()
            if not mod: continue # Moo
            user.add_roles(mod.name)
            db.flush()
            role = Role.query.filter(Role.name == mod.name).first()
            role.add_abilities('mods-edit', 'mods-remove')
            role.add_param('mods-edit', 'modid', str(mod.id))
            role.add_param('mods-remove', 'name', mod.name)
            db.add(role)
        
        # Modpacks
        for pack in ModList.query.filter(ModList.user_id == user.id).all():
            user.add_roles(pack.name)    
            db.flush()
            role = Role.query.filter(Role.name == pack.name).first()
            role.add_abilities('packs-edit', 'packs-remove')
            role.add_param('packs-edit', 'packid', str(pack.id))
            role.add_param('packs-remove', 'name', pack.name)
            db.add(role)
        
        # Admins
        if user.username in sys.argv:
            user.add_roles('admin')
            admin_role = Role.query.filter(Role.name == 'admin').first()
            admin_role.add_abilities_re('.*')
            admin_role.add_abilities('mods-invite')
            admin_role.add_abilities('view-users-full')
            admin_role.add_abilities('no-limits')

            # Params
            admin_role.add_param('admin-impersonate', 'userid', '.*')
            admin_role.add_param('mods-feature', 'gameshort', '.*')
            admin_role.add_param('game-edit', 'gameshort', '.*')
            admin_role.add_param('game-add', 'pubid', '.*')
            admin_role.add_param('game-remove', 'short', '.*')
            admin_role.add_param('mods-edit', 'gameshort', '.*')
            admin_role.add_param('mods-add', 'gameshort', '.*')
            admin_role.add_param('mods-remove', 'gameshort', '.*')
            admin_role.add_param('packs-add', 'gameshort', '.*')
            admin_role.add_param('packs-remove', 'gameshort', '.*')
            admin_role.add_param('publisher-edit', 'publid', '.*')
            admin_role.add_param('token-edit', 'tokenid', '.*')
            admin_role.add_param('token-remove', 'tokenid', '.*')
            admin_role.add_param('user-edit', 'userid', '.*')
    
            db.add(admin_role)
    db.commit()
except:
    db.rollback()
    raise
        