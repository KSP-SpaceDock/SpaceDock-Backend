Installation
============

Set up packages
---------------

apt-get install python3-dev virtualenv

Set up environment
------------------

cp alembic.ini.example alembic.ini
cp config.ini.example config.ini
virtualenv -p /usr/bin/python3 .
. bin/activate

Install requirements
--------------------

pip install -r requirements.txt
