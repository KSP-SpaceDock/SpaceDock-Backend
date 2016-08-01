from flask import redirect, make_response
from SpaceDock.routing import route
from SpaceDock.config import cfg

import mimetypes
import os.path

@route('/content/<url>')
def download(url):
    """
    Downloads a file from the storage.
    """
    # Check for a CDN
    if cfg['cdn-domain']:
        return redirect('http://' + cfg['cdn-domain'] + '/' + url, code=302)

    # Check for X-Sendfile
    response = None
    if _cfg("use-x-accel") == 'nginx':
        response = make_response("")
        response.headers['Content-Type'] = mimetypes.guess_type(os.path.join(_cfg('storage'), url))[0]
        response.headers['Content-Disposition'] = 'attachment; filename=' + os.path.basename(url)
        response.headers['X-Accel-Redirect'] = '/internal/' + url
    if _cfg("use-x-accel") == 'apache':
        response = make_response("")
        response.headers['Content-Type'] = mimetypes.guess_type(os.path.join(_cfg('storage'), url))[0]
        response.headers['Content-Disposition'] = 'attachment; filename=' + os.path.basename(url)
        response.headers['X-Sendfile'] = os.path.join(_cfg('storage'), url)
    if response is None:
        response = make_response(send_file(os.path.join(_cfg('storage'), url), as_attachment = True))
    return response