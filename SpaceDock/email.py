from SpaceDock.celery import send_mail
from SpaceDock.config import cfg
from werkzeug.utils import secure_filename

def send_confirmation(user, followMod=None):
    with open("emails/confirm-account") as f:
        confirmation = ""
        if followMod != None:
            confirmation = user.confirmation + "?f=" + followMod
        else:
            confirmation = user.confirmation
        message = f.read().format(
            **{ 
                'site_name': cfg['site-name'], 
                'username': user.username, 
                'domain': cfg['domain'],
                'confirmation': confirmation
            })
        send_mail.delay(cfg['support-mail'], [ user.email ], "Welcome to " + cfg['site-name'] + "!", message, important=True)

def send_reset(user):
    with open("emails/password-reset") as f:
        message = f.read().format(
            **{ 
                'site_name': cfg['site-name'], 
                'username': user.username, 
                'domain': cfg['domain'], 
                'confirmation': user.passwordReset 
             })
        send_mail.delay(cfg['support-mail'], [ user.email ], "Reset your password on " + cfg['site-name'], message, important=True)

def send_grant_notice(mod, user):
    with open("emails/grant-notice") as f:
        message = f.read().format(
            **{ 
                'username': user.username, 
                'mod_username': mod.user.username,
                'mod_name': mod.name,
                'site_name': cfg['site-name'], 
                "domain": cfg["domain"],
                'url': create_mod_url(mod.id, mod.name) 
            })
        send_mail.delay(cfg['support-mail'], [ user.email ], "You've been asked to co-author a mod on " + cfg['site-name'], message, important=True)

def send_update_notification(mod, version, user):
    followers = [u.email for u in mod.followers]
    changelog = version.changelog
    if changelog:
        changelog = '\n'.join(['    ' + l for l in changelog.split('\n')])

    targets = list()
    for follower in followers:
        targets.append(follower)
    if len(targets) == 0:
            return
    with open("emails/mod-updated") as f:
        message = f.read().format(
            **{
                'username': user.username,
                'friendly_version': version.friendly_version,
                'mod_name': mod.name,
                'site_name': cfg['site-name'],
                'changelog': changelog,
                'domain': cfg["domain"],
                'url': create_mod_url(mod.id, secure_filename(mod.name)[:64]),
                'game_name': version.gameversion.game.name,
                'gameversion': version.gameversion.friendly_version
            })
        subject = user.username + " has just updated " + mod.name + "!"
        send_mail.delay(cfg['support-mail'], targets, subject, message)

def send_autoupdate_notification(mod):
    followers = [u.email for u in mod.followers]
    changelog = mod.default_version().changelog
    if changelog:
        changelog = '\n'.join(['    ' + l for l in changelog.split('\n')])

    targets = list()
    for follower in followers:
        targets.append(follower)
    if len(targets) == 0:
        return
    with open("emails/mod-autoupdated") as f:
        message = f.read().format(
            **{
                'username': mod.user.username,
                'friendly_version': mod.default_version().friendly_version,
                'mod_name': mod.name,
                'game_name': mod.game.name,
                'gameversion': mod.default_version().gameversion.friendly_version,
                'domain': cfg["domain"],
                'url': create_mod_url(mod.id, secure_filename(mod.name)[:64])
            })
        subject = mod.name + " is compatible with " + mod.game.name + " " + mod.versions[0].gameversion.friendly_version + "!"
        send_mail.delay(cfg['support-mail'], targets, subject, message)

def send_bulk_email(users, subject, body):
    targets = list()
    for u in users:
        targets.append(u)
    send_mail.delay(cfg['support-mail'], targets, subject, body)

def create_mod_url(id, name):
    route = cfg['mod-url']
    return route.replace('{id}', str(id)).replace('{name}', name) # Using manual replacement here, so users dont need to use both values
