from datetime import datetime
from flask import request
from sqlalchemy import desc
from SpaceDock.common import edit_object, user_has, with_session
from SpaceDock.database import db
from SpaceDock.formatting import publisher_info
from SpaceDock.objects import Publisher
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
def publishers_info(pubid):
    """
    Outputs detailed infos for one publisher
    """
    if not pubid.isdigit() or not Publisher.query.filter(Publisher.id == int(pubid)).first():
        return {'error': True, 'reasons': ['Invalid publisher ID'], 'codes': ['2110']}, 400
    # Return the info
    pub = Publisher.query.filter(Publisher.id == int(pubid)).first()
    return {'error': False, 'count': 1, 'data': publisher_info(pub)}

@route('/api/publishers/<pubid>/edit', methods=['POST'])
@user_has('publisher-edit', params=['pubid'])
@with_session
def edit_publisher(pubid):
    """
    Edits a publisher, based on the request parameters. Required fields: data
    """
    if not pubid.isdigit() or not Publisher.query.filter(Publisher.id == int(pubid)).first():
        return {'error': True, 'reasons': ['Invalid publisher ID'], 'codes': ['2110']}, 400

    # Get the matching game and edit it
    pub = Publisher.query.filter(Publisher.id == int(pubid)).first()
    code = edit_object(pub, request.json)
    
    # Error check
    if code == 3:
        return {'error': True, 'reasons': ['The value you submitted is invalid'], 'codes': ['2180']}, 400
    elif code == 2:
        return {'error': True, 'reasons': ['You tried to edit a value that doesn\'t exist.'], 'codes': ['3090']}, 400
    elif code == 1:
        return {'error': True, 'reasons': ['You tried to edit a value that is marked as read-only.'], 'codes': ['3095']}, 400
    else:
        pub.updated = datetime.now()
        return {'error': False, 'count': 1, 'data': publisher_info(pub)}

@route('/api/publishers/add', methods=['POST'])
@user_has('publisher-add')
@with_session
def add_publisher():
    """
    Adds a publisher, based on the request parameters. Required fields: name
    """
    # Get variables
    name = request.json.get('name')

    # Check for existence
    if Publisher.query.filter(Publisher.name == name).first():
        return {'error': True, 'reasons': ['A publisher with this name already exists.'], 'codes': ['2000']}, 400

    # Get the matching game and edit it
    pub = Publisher(name)
    db.add(pub)
    db.flush()
    return {'error': False, 'count': 1, 'data': publisher_info(pub)}

@route('/api/publishers/remove', methods=['POST'])
@user_has('publisher-remove')
@with_session
def remove_publisher():
    """
    Removes a game from existence. Required fields: pubid
    """
    pubid = request.json.get('pubid')

    # Check if the pubid is valid
    if not isinstance(pubid, int) or not Publisher.query.filter(Publisher.id == pubid).first():
        return {'error': True, 'reasons': ['Invalid publisher ID'], 'codes': ['2110']}, 400

    # Get the publisher and remove it
    pub = Publisher.query.filter(Publisher.id == pubid).first()
    db.delete(pub)
    return {'error': False}