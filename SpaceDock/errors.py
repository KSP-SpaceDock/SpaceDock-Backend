from SpaceDock.app import app
from flask import jsonify


# Declare any HTTP error you want
@app.errorhandler(404)
@app.errorhandler(500)
def handle_error(e):
    return jsonify(error=str(e)), e.code if hasattr(e, 'code') else 500
