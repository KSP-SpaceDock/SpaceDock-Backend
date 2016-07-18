from importlib.machinery import SourceFileLoader
import os

# A list of all plugin modules
plugins = []

def load_plugins():
    """
    Loads all python files inside the "plugins/" directory.
    """
    if not os.path.exists('plugins/'):
        return
    for file in os.listdir('plugins/'):
        if file.endswith('.py'):
            plugins.append(SourceFileLoader(os.path.splitext(os.path.basename(file))[0], 'plugins/' + file).load_module())

## ADD CALLBACK API HERE ##