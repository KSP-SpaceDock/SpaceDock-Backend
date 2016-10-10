# Sets up a barebone instance of SpaceDock Backend

from werkzeug.utils import secure_filename
from SpaceDock.config import cfg
from SpaceDock.database import db
from SpaceDock.common import *
from SpaceDock.objects import *
from zipfile import ZipFile

import os.path
import SpaceDock.app

# Makes a new User
def new_user(name, password, email, admin):
    user = User(name, email, password)
    db.add(user)

    # Setup roles
    user.add_roles(name)
    db.commit()
    role = Role.query.filter(Role.name == name).first()
    role.add_abilities('user-edit', 'mods-add', 'packs-add', 'logged-in')
    role.add_param('user-edit', 'userid', user.id)
    role.add_param('mods-add', 'gameshort', '.*')
    role.add_param('packs-add', 'gameshort', '.*')

    # Admin Roles
    if admin:
        user.add_roles('admin')
        admin_role = Role.query.filter(Role.name == 'admin').first()
        admin_role.add_abilities_re('.*')
        admin_role.add_abilities('mods-invite')

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
        admin_role.add_param('user-edit', 'userid', '.*')

        db.add(admin_role)
        

    # Confirmation
    user.confirmation = None
    user.public = True
    db.commit()
    return user

# Makes a new game
def new_game(name, short, publisher):
    if Publisher.query.filter(Publisher.name == publisher).first():
        pub = Publisher.query.filter(Publisher.name == publisher).first()
    else:
        pub = Publisher(publisher)
    db.add(pub)

    # Create the game
    game = Game(name, pub, short)
    game.active = True
    db.add(game)

    # Commit and return
    db.commit()
    return game

# Makes a new gameversion
def new_version(game, name, beta):
    version = GameVersion(name, Game.query.filter(Game.short == game).first(), beta)
    db.add(version)
    db.commit()
    return version

# Creates a new user to admin a game
def new_game_admin(name, password, email, game):
    user = new_user(name, password, email, False)

    # Add game specific stuff
    user.add_roles(game)
    game_role = Role.query.filter(Role.name == game).first()
    game_role.add_abilities('game-edit')
    game_role.add_abilities_re('mods-.*')
    game_role.add_abilities('mods-invite')
    game_role.add_abilities_re('packs-.*')

    # Params
    game_role.add_param('mods-feature', 'gameshort', game)
    game_role.add_param('game-edit', 'gameshort', game)
    game_role.add_param('mods-edit', 'gameshort', game)
    game_role.add_param('mods-add', 'gameshort', game)
    game_role.add_param('mods-remove', 'gameshort', game)
    game_role.add_param('packs-add', 'gameshort', game)
    game_role.add_param('packs-remove', 'gameshort', game)

    # Assign
    db.add(game_role)

    # Commit and return
    db.commit()
    return user

# Creates a new mod
def new_mod(name, user, game, license):
    user = User.query.filter(User.username == user).first()
    game = Game.query.filter(Game.short == game).first()

    # Create new object
    mod = Mod(name, user, game, license)
    mod.published = True

    # Roles
    user.add_roles(name)
    role = Role.query.filter(Role.name == name).first()
    role.add_abilities('mods-edit', 'mods-remove')
    role.add_param('mods-edit', 'modid', str(mod.id))
    role.add_param('mods-remove', 'name', name)

    # Assign
    db.add(role)
    db.add(mod)

    # Commit and return
    db.commit()
    return mod

# Creates a new version for the mod
def new_mod_version(modname, friendly_version, game, gameversion, beta):
    mod = Mod.query.filter(Mod.name == modname).first()
    game = GameVersion.query.filter(Game.query.filter(Game.short == game).first().id == GameVersion.game_id).filter(GameVersion.friendly_version == gameversion).first()

    # Path
    filename = secure_filename(mod.name) + '-' + secure_filename(friendly_version) + '.zip'
    base_path = os.path.join(secure_filename(mod.user.username) + '_' + str(mod.user_id), secure_filename(mod.name))
    full_path = os.path.join(cfg['storage'], base_path)
    if not os.path.exists(full_path):
        os.makedirs(full_path)
    path_ = os.path.join(full_path, filename)

    # Create the object
    version = ModVersion(mod, friendly_version, game, '/content/' + base_path.replace("\\", "/") + "/" + filename, beta)

    # Save data
    zip = ZipFile(path_, 'w')
    zip.writestr('SUPRISE.txt', ('As it seems, you downloaded ' + modname + ' ' + friendly_version).encode('utf-8'))
    zip.close()

    # add and commit
    db.add(version)
    db.flush()
    if not beta:
        mod.default_version_id = version.id
        print(mod.default_version_id)
    db.commit()
    return version
    


# Only run this when the file is run directly
if __name__ == '__main__':

    # Setup the DB
    admin = new_user('Administrator', 'admin', 'admin@noname.net', True)
    user = new_user('SpaceDockUser', 'user', 'user@noname.net', False)

    # Game 1
    game_KSP = new_game('Kerbal Space Program', 'kerbal-space-program', 'Squad MX')
    version_112 = new_version('kerbal-space-program', '1.1.2', False)
    version_113 = new_version('kerbal-space-program', '1.1.3', False)
    version_12 = new_version('kerbal-space-program', '1.2', True)

    # Game 2
    game_Factorio = new_game('Factorio', 'factorio', 'Wube Software')
    version_012 = new_version('factorio', '0.12', False)

    # Game admins
    ksp_admin = new_game_admin('GameAdminKSP', 'gameadminksp', 'gameadminksp@noname.net', 'kerbal-space-program')
    fac_admin = new_game_admin('GameAdminFAC', 'gameadminfac', 'gameadminfac@noname.net', 'factorio')

    # Mods
    mod_ksp_1 = new_mod('DarkMultiPlayer', 'SpaceDockUser', 'kerbal-space-program', 'MIT')
    mod_ksp_2 = new_mod('CookieEngine', 'GameAdminKSP', 'kerbal-space-program', 'GPL')

    # Versions
    mod_ksp_1_1 = new_mod_version('DarkMultiPlayer', '0.1', 'kerbal-space-program', '1.1.2', False)
    mod_ksp_1_2 = new_mod_version('DarkMultiPlayer', '0.2', 'kerbal-space-program', '1.1.3', False)
    mod_ksp_1_3 = new_mod_version('DarkMultiPlayer', '0.3', 'kerbal-space-program', '1.2', True)
    mod_ksp_2_1 = new_mod_version('CookieEngine', '1.2', 'kerbal-space-program', '1.1.3', False)
    mod_ksp_2_2 = new_mod_version('CookieEngine', '1.4', 'kerbal-space-program', '1.2', True)


    
             