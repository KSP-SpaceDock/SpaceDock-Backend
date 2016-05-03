from flask import jsonify

endpoints = {}

class API:
    def __init__(self, flask, documentation):
        self.flask = flask
        self.documentation = documentation
        self.register_api_endpoint("test", self.test)
        
    def register_api_endpoint(self, url, method):
        self.flask.add_url_rule("/api/" + url, method.__name__, method)
        self.documentation.add_documentation(url, method)
    
    def test(self):
        """
        An API endpoint
        """
        return jsonify(test="Some test")