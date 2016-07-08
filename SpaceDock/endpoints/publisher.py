from flask import jsonify, request
from flask_login import login_required, current_user
from sqlalchemy import desc
from SpaceDock.objects import *
from SpaceDock.formatting import publisher_info

class PublisherEndpoints:

    def publishers_list(self):
        """
        Outputs all publishers known by the application
        """
        results = dict()
        for v in Publisher.query.order_by(desc(Publisher.id)).all():
            results[v.id] = v.name
        return {"error": False, 'count': len(results), 'data': results}

    publishers_list.api_path = "/api/publishers"

    def publisher_info(self, pubid):
        """
        Outputs detailed infos for one publisher
        """
        if not pubid.isdigit() or not Publisher.query.filter(Publisher.id == int(pubid)).first():
            return {'error': True, 'reasons': ['Invalid publisher ID']}, 400
        # Return the info
        pub = Publisher.query.filter(Publisher.id == int(pubid)).first()
        return {'error': False, 'count': 1, 'data': publisher_info(pub)}

    publisher_info.api_path = '/api/publishers/<pubid>'