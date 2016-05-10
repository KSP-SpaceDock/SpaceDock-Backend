from SpaceDock.objects import *

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
    'created': user.created,
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