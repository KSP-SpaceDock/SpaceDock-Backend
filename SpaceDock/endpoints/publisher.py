from datetime import datetime
from flask import request
from sqlalchemy import desc
from SpaceDock.common import with_session, user_has
from SpaceDock.database import db
from SpaceDock.formatting import publisher_info
from SpaceDock.objects import *
from SpaceDock.routing import route


@route('/api/publishers')
def publishers_list():
    """
    Outputs all publishers known by the application
    """
    results = dict()
    for v in Publisher.query.order_by(desc(Publisher.id)).all():
        results[v.id] = v.name
    return {"error": False, 'count': len(results), 'data': results}

@route('/api/publishers/<pubid>')
def publisher_info(pubid):
    """
    Outputs detailed infos for one publisher
    """
    if not pubid.isdigit() or not Publisher.query.filter(Publisher.id == int(pubid)).first():
        return {'error': True, 'reasons': ['Invalid publisher ID']}, 400
    # Return the info
    pub = Publisher.query.filter(Publisher.id == int(pubid)).first()
    return {'error': False, 'count': 1, 'data': publisher_info(pub)}

@route('/api/publishers/<pubid>/edit', methods=['POST'])
@user_has('publisher-edit', params=['pubid'])
@with_session
def edit_publisher(publid):
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
    pub.updated = datetime.now()
    return {'error': False}

@route('/api/publishers/add', methods=['POST'])
@user_has('publisher-add')
@with_session
def add_publisher():
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
    db.add(pub)
    return {'error': False}

@route('/api/publishers/remove', methods=['POST'])
@user_has('publisher-remove')
@with_session
def remove_publisher():
    """
    Removes a game from existence. Required fields: pubid
    """
    pubid = request.form['pubid']

    # Check if the pubid is valid
    if not pubid.isdigit() or not Publisher.query.filter(Publisher.id == int(pubid)).first():
        return {'error': True, 'reasons': ['Invalid publisher ID']}, 400

    # Get the publisher and remove it
    pub = Publisher.query.filter(Publisher.id == int(pubid)).first()
    db.remove(pub)
    return {'error': False}