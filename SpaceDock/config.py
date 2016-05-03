from configparser import ConfigParser

class Config:
    def __init__(self):
        self._config = ConfigParser()
        self._config.readfp(open('config.ini'))
        self._env = self._config['meta']['environment']
        
    def get_environment(self):
        """
        Returns the current environment
        """
        return self._env
    
    def get(self, key):
        """
        Returns a string
        """
        return self._config[self._env][key]
    
    def geti(self, key):
        """
        Returns an integer
        """
        return int(self.get(key))
    
    def getb(self, key):
        """
        Returns a boolean
        """
        return bool(self.get(key))
    
    def __getitem__(self, key):
        """
        Returns a string, indexer access
        """
        return self.get(key)
    