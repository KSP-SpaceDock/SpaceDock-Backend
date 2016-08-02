# Sets up a barebone instance of SpaceDock Backend

from SpaceDock.config import cfg
from SpaceDock.database import db
from SpaceDock.common import *
from SpaceDock.objects import *

# Makes a new User
def new_user(name, password, email, admin):
    user = User(name, email, password)
    db.add(user)

    # Setup roles
    user.add_roles(name)
    db.commit()
    role = Role.query.filter(Role.name == name.lower()).first()
    role.add_abilities('user-edit', 'mods-add', 'logged-in')
    role.add_param('user-edit', 'userid', user.id)
    role.add_param('mods-add', 'gameshort', '.*')
    role.add_param('packs-add', 'gameshort', '.*')

    # Admin Roles
    if admin:
        user.add_roles('admin')
        admin_role = Role.query.filter(Role.name == 'admin').first()
        admin_role.add_abilities_re('.*')

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
    game = Game(name, pub.id, short)
    game.active = True
    db.add(game)

    # Commit and return
    db.commit()
    return game

# Makes a new gameversion
def new_version(game, name, beta):
    version = GameVersion(name, Game.query.filter(Game.short == game).first().id, beta)
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



# Setup the DB
admin = new_user('Administrator', 'admin', 'admin@noname.net', True)
user = new_user('SpaceDockUser', 'user', 'user@noname.net', False)

# Game 1
game_KSP = new_game('Kerbal Space Program', 'kerbal-space-program', 'Squad MX')
version_112 = new_version('kerbal-space-program', '1.1.2', False)
version_113 = new_version('kerbal-space-program', '1.1.3', False)
version_12 = new_version('kerbal-space-program', '1.2', True)

# Game 2
game_Factorio = new_game('factorio', 'factorio', 'Wube Software')
version_012 = new_version('factorio', '0.12', False)

# Game admins
ksp_admin = new_game_admin('GameAdminKSP', 'gameadminksp', 'gameadminksp@noname.net', 'kerbal-space-program')
fac_admin = new_game_admin('GameAdminFAC', 'gameadminfac', 'gameadminfac@noname.net', 'factorio')


    
             