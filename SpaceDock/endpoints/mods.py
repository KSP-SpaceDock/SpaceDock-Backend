from SpaceDock.objects import *
from SpaceDock.common import *
from SpaceDock.formatting import mod_info
from sqlalchemy import desc

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