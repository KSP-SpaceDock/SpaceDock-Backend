from flask import jsonify, request
from flask_login import login_required, current_user
from sqlalchemy import desc
from SpaceDock.objects import *
from SpaceDock.formatting import publisher_info
from SpaceDock.common import with_session, user_has

class PublisherEndpoints:
    def __init__(self, db):
        self.db = db

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

    @with_session
    @user_has('publisher-edit', params=['pubid'])
    def edit_publisher(self, publid):
        """
        Edits a publisher, based on the request parameters. Required fields: data
        """
        if not pubid.isdigit() or not Publisher.query.filter(Publisher.id == int(pubid)).first():
            return {'error': True, 'reasons': ['Invalid publisher ID']}, 400

        # Get variables
        parameters = json.loads(request.form['data'])

        # Get the matching game and edit it
        pub = Publisher.query.filter(Publisher.id == int(pubid)).first()
        edit_object(pub, parameters)
        return {'error': False}

    edit_publisher.api_path = '/api/publishers/<pubid>/edit'
    edit_publisher.methods = ['POST']

    @with_session
    @user_has('publisher-add')
    def add_publisher(self):
        """
        Adds a publisher, based on the request parameters. Required fields: name
        """
        # Get variables
        name = request.form['name']

        # Check for existence
        if Publisher.query.filter(Publisher.name == name).first():
            return {'error': True, 'reasons': ['A publisher with this name already exists.']}, 400

        # Get the matching game and edit it
        pub = Publisher(name)
        self.db.add(pub)
        return {'error': False}

    add_publisher.api_path = '/api/publishers/add'
    add_publisher.methods = ['POST']

    @with_session
    @user_has('publisher-remove')
    def remove_publisher(self):
        """
        Removes a game from existence. Required fields: pubid
        """
        pubid = request.form['pubid']

        # Check if the pubid is valid
        if not pubid.isdigit() or not Publisher.query.filter(Publisher.id == int(pubid)).first():
            return {'error': True, 'reasons': ['Invalid publisher ID']}, 400

        # Get the publisher and remove it
        pub = Publisher.query.filter(Publisher.id == int(pubid)).first()
        self.db.remove(pub)
        return {'error': False}

    remove_publisher.api_path = '/api/publishers/remove'
    remove_publisher.methods = ['POST']