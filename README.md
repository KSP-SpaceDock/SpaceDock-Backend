## SpaceDock
SpaceDock is an open source website software that can host modifications for multiple games. It started as a fork of the popular KerbalStuff software for Kerbal Space Program and evolved since then.

SpaceDock is split into two parts:
* The frontend, which is not released under an open source license and unavailable to the public
* The backend which is licensed as MIT and free for everyone to use.

The backend handles all the content that SpaceDock stores, like mods, background images and users. The program is designed to be as lightweight as possible to distribute it across multiple nodes. A frontend can then query the multiple nodes. Due to the modular nature of the framework, third party persons who are interested in mod hosting can either host their own backend together with a custom frontend, or create just a frontend and rely on SpaceDocks infrastructure internally to host their mods. This allows for flexible and easily customised mod repositories (once there are free frontends available)

The API implemented by SpaceDock Backend is the suceessor of the old KerbalStuff API that is still used in the current SpaceDock. It offers way more routes, but lacks an actual frontend. Instead of HTML code, the backend returns JSON formatted responses, or optionally JSONP. Other formats, like XML are **not** planned. The API is **not** compatible with the old KerbalStuff API. Some routes may be the same, but we don't support Applications built for the old API.

### Setting up SpaceDock-Backend
To setup a version of SpaceDock-Backend on your local computer (development) or server (production) you need to install Python3 and Virtualenv first. If you are running on an isolated system specifically used for SDB you can drop virtualenv, however in multi purpose systems we dearly recommend it.

First you need to clone the repository
```
git clone https://github.com/KSP-SpaceDock/SpaceDock-Backend.git
cd SpaceDock-Backend
```

Now you need to create a new virtualenv. Skip this step if you dont use virtualenv:
```
virtualenv -p /usr/bin/python3 . # Windows: C:\Python34\python.exe or wherever your python lies
source bin/activate # Windows: .\Scripts\activate
```

Copy the configuration files:
```
cp config/alembic.ini.example config/alembic.ini
cp config/config.ini.example config/config.ini
```

And finally install all the requirements. Depending on your internet connection this could take a while
```
pip install -r requirements.txt
```

To start the backend simply type
```
python app.py
```

### Requirements
SpaceDock-Backend is a Python Application that uses flask for serving content and sqlalchemy for persistency. Even though we are developing and running SpaceDock using PostgreSQL, you can use any SQL based Database in combination with sqlalchemy. (That means MySQL, MariaDB, SQLite). At the moment, we only support Postgres though.

### The team
* godarklight
* V1TA5
* ThomasKerman

### Special Thanks
* SirCmpwn, testeddoughnut and andyleap for developing the original KerbalStuff API
* All SpaceDock contributors for their help during the migration from KerbalStuff to SpaceDock
* The flask and sqlalchemy teams for their amazing frameworks

### License
This application is licensed under the MIT License. You are free to remix, adapt and redistribute it, however, please credit the original authors