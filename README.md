## SpaceDock
SpaceDock is an open source website software that can host modifications for multiple games. It started as a fork of the popular KerbalStuff software for Kerbal Space Program and evolved since then.

SpaceDock is split into two parts:
* The frontend, which is not released under an open source license and unavailable to the public
* The backend which is licensed as MIT and free for everyone to use.

The backend handles all the content that SpaceDock stores, like mods and users. The program is designed to be as lightweight as possible to distribute it across multiple nodes. A frontend can then query the multiple nodes. Due to the modular nature of the framework, third party persons who are interested in mod hosting can either host their own backend together with a custom frontend, or create just a frontend and rely on SpaceDocks infrastructure internally to host their mods. This allows for flexible and easily customised mod repositories (once there are free frontends available)

The API implemented by SpaceDock Backend is the suceessor of the old KerbalStuff API that is still used in the current SpaceDock. It offers way more routes, but lacks an actual frontend. Instead of HTML code, the backend returns JSON formatted responses, or optionally JSONP. Other formats, like XML are **not** planned. The API is **not** compatible with the old KerbalStuff API. Some routes may be the same, but we don't support Applications built for the old API.

### Setting up SpaceDock-Backend
To setup a version of SpaceDock-Backend on your local computer (development) or server (production) you need to install Go 1.8. You can get it here: https://golang.org/dl/

First you need to clone the repository
```
git clone https://github.com/KSP-SpaceDock/SpaceDock-Backend.git
cd SpaceDock-Backend
```

SpaceDock requires some apt-get packages to be installed
```
apt-get install postgresql-server-9.3 postgresql-server-dev-9.3
apt-get install libmysqlclient-dev mysql-server # Only if you want to use MySQL
```

SpaceDocks build process involves creating a virtual environment to set your GOPATH properly. If you are new to golang, and don't fully understand it's build / dependency concept, this is the solution for you. Running this code will start the virtual environment, fetch the SpaceDock dependencies, and build the executable file. Depending on your internet connection, this could take a while. The build process isn't verbose, it might look like the command line froze. Just leave it running, after ~5 minutes it will finish. 

If you want to use plugins in your SpaceDock instance, you need to enter their goland dependency urls into build/plugins.txt. The build process will fetch and embed them into the main app.
```
. build/activate.sh # On Windows: . .\build\activate.ps1
build sdb
```

Copy the configuration files:
```
cp config/config.example.yml config/config.yml
```

To start the backend simply type
```
./sdb # sdb.exe on Windows
```

### Requirements
SpaceDock-Backend is a Golang Application that uses [iris](https://github.com/kataras/iris) for serving content and [gorm](https://github.com/jinzhu/gorm) for persistency. Even though we are developing and running SpaceDock using PostgreSQL, you can use any SQL based Database in combination with gorm. (That means MySQL, MariaDB, SQLite). At the moment, we only support Postgres though.

### The team
* godarklight
* V1TA5
* StollD
* RockyTV

### Special Thanks
* SirCmpwn, testeddoughnut and andyleap for developing the original KerbalStuff API
* All SpaceDock contributors for their help during the migration from KerbalStuff to SpaceDock
* The iris and gorm teams for their amazing frameworks

### License
This application is licensed under the MIT License. You are free to remix, adapt and redistribute it, however, please credit the original authors