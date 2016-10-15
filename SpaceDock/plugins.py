from importlib.machinery import SourceFileLoader

import os

# A list of all plugin modules
plugins = []

def load_file(path):
    # Read control header
    name = ''
    deps = []
    with open(os.getcwd() + '/plugins/' + path) as f:
        lines = f.readlines()
        for line in lines:
            if not line.startswith('#'):
                break
            if line[1:].strip().lower().startswith('name:'):
                name = line[1:].strip()[5:].strip()
            if line[1:].strip().lower().startswith('depends:'):
                deps.append(line[1:].strip()[8:].strip())

    # Dont load the file if it has no name
    if not name:
        print('Plugin file "' + path + '" has no name declaration. Please add one using the control header.')
        return
            
    # Load deps
    for d in deps:
        if d.endswith('.py'):
            load_file(d)
    plugins.append(SourceFileLoader(name, 'plugins/' + path).load_module())
    print(' * Loaded ' + name)

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