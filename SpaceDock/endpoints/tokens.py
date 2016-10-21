from flask import request
from SpaceDock.common import user_has, with_session
from SpaceDock.database import db
from SpaceDock.formatting import token_format
from SpaceDock.objects import Token
from SpaceDock.routing import route

@route('/api/tokens/generate', methods=['POST'])
@user_has('token-generate')
@with_session
def generate_token():
    """
    Generates a new API token
    """
    token = Token()
    db.add(token)
    db.flush()
    token['ips'] = list()
    return {'error': False, 'count': 1, 'data': token_format(token)}

@route('/api/tokens/edit', methods=['POST'])
@user_has('token-edit', params=['tokenid'])
@with_session
def edit_token():
    """
    Edits the IP-Adresses of a token
    """
    tokenid = request.json.get('tokenid')
    ips = request.json.get('ips')

    # Check for int
    if not isinstance(tokenid, int) or not Token.query.filter(Token.id == tokenid).first():
        return {'error': True, 'reasons': ['The token ID is invalid'], 'codes': ['2131']}, 400
    if not ips or not isinstance(ips, list):
        return {'error': True, 'reasons': ['The list of IP Addresses is invalid.'], 'codes': ['2132']}, 400
    token = Token.query.filter(Token.id == tokenid).first()

    # Edit the token
    token['ips'] = ips
    return {'error': False, 'count': 1, 'data': token_format(token)}

@route('/api/tokens/revoke', methods=['POST'])
@user_has('token-revoke', params=['tokenid'])
@with_session
def revoke_token():
    """
    Removes a token completely
    """
    tokenid = request.json.get('tokenid')

    # Check for int
    if not isinstance(tokenid, int) or not Token.query.filter(Token.id == tokenid).first():
        return {'error': True, 'reasons': ['The token ID is invalid'], 'codes': ['2131']}, 400
    token = Token.query.filter(Token.id == tokenid).first()
    db.delete(token)
    return {'error': False}
