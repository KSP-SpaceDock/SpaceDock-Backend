from flask import url_for
from SpaceDock.objects import *

def bulk(array, formatter):
    for f in array:
        yield formatter(f)

def game_info(game):
    return {
        'id': game.id,
        'name': game.name,
        'active': game.active,
        'fileformats': game.fileformats,
        'altname': game.altname,
        'rating': game.rating,
        'releasedate': game.releasedate.isoformat() if not game.releasedate == None else None,
        'short': game.short,
        'publisher': game.publisher_id,
        'description': game.description,
        'short_description': game.short_description,
        'created': game.created.isoformat() if not game.created == None else None,
        'updated': game.updated.isoformat() if not game.updated == None else None,
        'background': game.background,
        'bgOffsetX': game.bgOffsetX,
        'bgOffsetY': game.bgOffsetY,
        'link': game.link
    }

def game_version_info(version):
    return {
        'id': version.id,
        'friendly_version': version.friendly_version,
        'is_beta': version.is_beta
    }

def publisher_info(publisher):
    # get games manually
    games = list()
    for g in Game.query.filter(Game.publisher_id == publisher.id).all():
        games.append(game_info(g))
    return {
        'id': publisher.id,
        'name': publisher.name,
        'short_description': publisher.short_description,
        'description': publisher.description,
        'created': publisher.created.isoformat() if not publisher.created == None else None,
        'updated': publisher.updated.isoformat() if not publisher.updated == None else None,
        'background': publisher.background,
        'bgOffsetX': publisher.bgOffsetX,
        'bgOffsetY': publisher.bgOffsetY,
        'link': publisher.link,
        'games': games
    }

def mod_info(mod):
    return {
        'id': mod.id,
        'author': mod.user.username,
        'user_id': mod.user_id,
        'game_id': mod.game_id,
        'game': mod.game.name,
        'game_short': mod.game.short,
        'shared_authors': [s.user.username for s in mod.shared_authors],
        'name': mod.name,
        'description': mod.description,
        'short_description': mod.short_description,
        'approved': mod.approved,
        'published': mod.published,
        'donation_link': mod.donation_link,
        'external_link': mod.external_link,
        'license': mod.license,
        'votes': mod.votes,
        'created': mod.created.isoformat() if not mod.created == None else None,
        'updated': mod.updated.isoformat() if not mod.updated == None else None,
        'background': mod.background,
        'bgOffsetX': mod.bgOffsetX,
        'bgOffsetY': mod.bgOffsetY,
        'medias': mod.medias,
        'default_version_id': mod.default_version_id,
        #'downloads': mod.downloads,
        #'follow_events': mod.follow_events,
        #'referrals': mod.referrals,
        'source_link': mod.source_link,
        'follower_count': mod.follower_count,
        'download_count': mod.download_count,
        'followers': mod.followers,
        #'rating': mod.rating,
        #'review': mod.review,
        'total_score': mod.total_score,
        'rating_count': mod.rating_count,
    }

def mod_version_info(modversion):
    return {
        'id': modversion.id,
        'mod_id': modversion.mod_id,
        'is_beta': modversion.is_beta,
        'friendly_version': modversion.friendly_version,
        'gameversion_id': modversion.gameversion_id,
        'game_version': modversion.gameversion.friendly_version,
        'created': modversion.created.isoformat() if not modversion.created == None else None,
        #TODO: Fix
        'download_path': modversion.download_path,
        #'download_path': url_for('mods.download', mod_id=modversion.mod.id, mod_name=modversion.mod.name, version=modversion.friendly_version),
        'changelog': modversion.changelog,
        'sort_index': modversion.sort_index,
        'file_size': modversion.file_size
    }

#WARNING: Some of this stuff is sensitive, make sure it comes from admin only access!
def admin_user_info(user):
    return {
    'id': user.id,
    'username': user.username,
    'email': user.email,
    'showEmail': user.showEmail,
    'public': user.public,
    'admin': user.admin,
    #Password skipped
    'description': user.description,
    'created': user.created.isoformat() if not user.created == None else None,
    'showCreated': user.showCreated,
    'forumUsername': user.forumUsername,
    'showForumName': user.showForumName,
    'forumId': user.forumId,
    'ircNick': user.ircNick,
    'showIRCName': user.showIRCName,
    'twitterUsername': user.twitterUsername,
    'showTwitterName': user.showTwitterName,
    'redditUsername': user.redditUsername,
    'showRedditName': user.showRedditName,
    'youtubeUsername': user.youtubeUsername,
    'showYoutubeName': user.showYoutubeName,
    'twitchUsername': user.twitchUsername,
    'showTwitchName': user.showTwitchName,
    'facebookUsername': user.facebookUsername,
    'showFacebookName': user.showFacebookName,
    'location': user.location,
    'showLocation': user.showLocation,
    'backgroundMedia': user.backgroundMedia,
    #Password reset skipped
    'bgOffsetX': user.bgOffsetX,
    'bgOffsetY': user.bgOffsetY,
    'dark_theme': user.dark_theme
    }

def user_info(user):
    return {
    'id': user.id,
    'username': user.username,
    'public': user.public,
    'admin': user.admin,
    #Email skipped
    #Password skipped
    'description': user.description,
    'created': user.created.isoformat() if not user.created == None else None,
    'showCreated': user.showCreated,
    'forumUsername': user.forumUsername,
    'showForumName': user.showForumName,
    'forumId': user.forumId,
    'ircNick': user.ircNick,
    'showIRCName': user.showIRCName,
    'twitterUsername': user.twitterUsername,
    'showTwitterName': user.showTwitterName,
    'redditUsername': user.redditUsername,
    'showRedditName': user.showRedditName,
    'youtubeUsername': user.youtubeUsername,
    'showYoutubeName': user.showYoutubeName,
    'twitchUsername': user.twitchUsername,
    'showTwitchName': user.showTwitchName,
    'facebookUsername': user.facebookUsername,
    'showFacebookName': user.showFacebookName,
    'location': user.location,
    'showLocation': user.showLocation,
    'backgroundMedia': user.backgroundMedia,
    #Password reset skipped
    'bgOffsetX': user.bgOffsetX,
    'bgOffsetY': user.bgOffsetY,
    'dark_theme': user.dark_theme,
    'roles': roles_format(user._roles)
    }

def roles_format(roles):
    for role in roles:
        yield role.name

def ability_format(ability):
    return {
        'id': ability.id,
        'name': ability.name
    }