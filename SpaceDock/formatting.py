from SpaceDock.objects import *

def game_info(game):
    return {
        "id": game.id,
        "name": game.name,
        "publisher_id": game.publisher_id,
        "short_description": game.short_description,
        "description": game.description,
        "created": game.created.isoformat(),
        "background": game.background,
        "bg_offset_x": game.bgOffsetX,
        "bg_offset_y": game.bgOffsetY,
        "link": game.link
    }