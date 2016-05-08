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
        'created': game.created.isotime() if not game.created == None else None,
        'updated': game.updated.isotime() if not game.updated == None else None,
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