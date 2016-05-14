from flask import jsonify, request
from flask_login import login_required, current_user
from sqlalchemy import desc
from SpaceDock.objects import *
from SpaceDock.formatting import publisher_info

class PublisherEndpoints:

    def publishers_list(self):
        results = dict()
        for v in Publisher.query.order_by(desc(Publisher.id)).all():
            results[v.id] = publisher_info(v)
        return jsonify({"error": False, 'data': results})
    
    publishers_list.api_path = "/api/publishers"