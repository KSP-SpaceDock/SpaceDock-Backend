from functools import wraps
from SpaceDock.app import app

# Wrappers used for routing
wrappers = []
wrappers_immideate = []

# Registers a new wrapper
def add_wrapper(wrap, immideate=False):
    if immideate:
        wrappers_immideate.append(wrap)
    else:
        wrappers.append(wrap)

# Removes a wrapper
def remove_wrapper(wrap):
    if wrap in wrappers_immideate:
        wrappers_immideate.remove(wrap)
    else:
        wrappers.remove(wrap)

# A decorator that is used to register a view function for a given URL rule.
def route(rule, **options):
    def wrapper(f):
        f.api_path = rule
        for wrapi in wrappers_immideate:
            f = wrapi(f)
        @wraps(f)
        def inner(*args, **kwargs):
            func = f
            for wrap in wrappers:
                func = wrap(func)
            return func(*args, **kwargs)
        return app.route(rule, **options)(inner)
    return wrapper