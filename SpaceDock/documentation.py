from flask import redirect, url_for, jsonify

class Documentation:
    def __init__(self, flask):
        self.flask = flask
        self.flask.add_url_rule('/documentation', 'documentation', self.documentation_page)
        self.flask.add_url_rule('/', 'index', self.index_page)        
        self.methods = {}
        
    def index_page(self):
        return redirect(url_for('documentation'), 302)
    
    def documentation_page(self):
        return jsonify(self.methods)
    
    def add_documentation(self, url, method):
		if not method.__doc__ == None:
			self.methods[url] = method.__doc__.strip()
		else:
			self.methods[url] = 'No documentation available'
        self.flask.add_url_rule('/documentation/' + url, "doc_" + method.__name__, self.methods[url])