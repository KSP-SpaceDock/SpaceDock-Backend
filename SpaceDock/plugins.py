from importlib.machinery import SourceFileLoader
import os

# A list of all plugin modules
plugins = []

def load_file(path):
    # Read control header
    name = ''
    deps = []
    with open(path) as f:
        lines = f.readlines()
        for line in lines:
            if not line.startswith('#'):
                break
            if line[1:].trim().lower().startswith('name:'):
                name = line[1:].trim()[5:].trim()
            if line[1:].trim().lower().startswith('depends:'):
                deps.append(line[1:].trim()[8:].trim())

    # Dont load the file if it has no name
    if not name:
        print('Plugin file "' + path + '" has no name declaration. Please add one using the control header.')
            
    # Load deps
    for d in deps:
        if d.endswith('.py'):
            load_file(d)
    plugins.append(SourceFileLoader(name, 'plugins/' + file).load_module())

def load_plugins():
    """
    Loads all python files inside the "plugins/" directory.
    """
    if not os.path.exists('plugins/'):
        return
    for file in os.listdir('plugins/'):
        if file.endswith('.py'):
            load_file(file)

## ADD CALLBACK API HERE ##