Installation
============

Set up packages
---------------

apt-get install python3-dev virtualenv

Set up environment
------------------

- git clone https://github.com/KSP-SpaceDock/SpaceDock-Backend.git
- cd SpaceDock-Backend
- cp alembic.ini.example alembic.ini
- cp config.ini.example config.ini
- virtualenv -p /usr/bin/python3 .
- source bin/activate
- pip install -r requirements.txt

Edit your settings
------------------
nano config.ini

Start the app
------------------
python app.py


