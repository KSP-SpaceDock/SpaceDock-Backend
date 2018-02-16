## SpaceDock
SpaceDock is an open source website software that can host modifications for multiple games. It started as a fork of the popular KerbalStuff software for Kerbal Space Program and evolved since then.

SpaceDock is split into two parts:
* The frontend, which is currently developed under the name [OpenDock](https://github.com/KSP-SpaceDock/OpenDock)
* The backend which developed here.

The backend handles all the content that SpaceDock stores, like mods and users. The program is designed to be as lightweight as possible to distribute it across multiple nodes. A frontend can then query the multiple nodes. Due to the modular nature of the framework, third party persons who are interested in mod hosting can either host their own backend together with a custom frontend, or create just a frontend and rely on SpaceDocks infrastructure internally to host their mods. This allows for flexible and easily customised mod repositories (once there are free frontends available)

The API implemented by SpaceDock Backend is the suceessor of the old KerbalStuff API that is still used in the current SpaceDock. It offers way more routes, but lacks an actual frontend. Instead of HTML code, the backend returns JSON formatted responses, or optionally JSONP. Other formats, like XML are **not** planned. The API is **not** compatible with the old KerbalStuff API. Some routes may be the same, but we don't support Applications built for the old API.

### Setting up SpaceDock-Backend
To setup a version of SpaceDock-Backend on your local computer (development) or server (production) you need to install Go 1.8. You can get it here: https://golang.org/dl/. The go executable should be in your $PATH.

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

SpaceDocks build process follows the golang specifications: You need to clone the SpaceDock repository into your $GOPATH, so that it looks like this ($GOPATH is set to `/home/user/gocode` in this example)

```
/home/user/gocode/src/github.com/KSP-SpaceDock/SpaceDock-Backend/
---- app/
---- build/
---- config/
---- emails/
---- middleware/
---- objects/
---- routes/
---- tools/
---- vendor/
---- sdb.go
```

#### Glide
We use the [Glide package manager](https://glide.sh) to manage the sdb dependencies. For installing glide, please execute the appropreate buildscript for your plattform (Linux/MacOS Users: build/install_glide.sh, Windows Users: build/install_glide.ps1). This will install glide into $GOPATH/bin, which requires a valid $GOPATH! It is also **required** to add $GOPATH/bin to your $PATH!

#### Installing the dependencies
Installing the SpaceDock-Backend dependencies is as simple as running `$ glide install` in the SDB root directory after installing glide.

#### Plugins
Well, "plugins"...

If you want to use plugins in your SpaceDock instance, you need to enter their goland dependency urls into build/plugins.txt. You can use full glide versioning syntax here. After doing so, you need to run `build/fetch_plugins.sh` if you are on Linux/MacOS or `build/fetch_plugins.ps1` on Windows. This will create a file called `build_sdb.go` that includes the specified plugins and fetch them using glide.

Even if you don't plan to use plugins, execute the `fetch_plugins` script. It will work without a `build/plugins.txt` file being present, and just creates the `build_sdb.go` file.

#### Building and starting the application
To build the app you need to call `go build` on the `build_sdb.go` file. The only difference between Linux and Windows is, that Windows users should use `sdb.exe` instead of `sdb`

```
go build -v -o ./sdb ./build_sdb.go
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
SpaceDock-Backend is a Golang Application that uses [iris](https://github.com/kataras/iris) for serving content and [gorm](https://github.com/jinzhu/gorm) for persistency. Even though we are developing and running SpaceDock using PostgreSQL, you can use any SQL based Database in combination with gorm. (That means MySQL, MariaDB). SQLite could work, but supporting it is a pain, because it uses cgo, which wouldn't allow us to crosscompile the program. At the moment, we only support Postgres.

### The team
* godarklight
* V1TA5
* StollD
* RockyTV

### Special Thanks
* SirCmpwn, testeddoughnut and andyleap for developing the original KerbalStuff API
* All SpaceDock contributors for their help during the migration from KerbalStuff to SpaceDock
* The iris, glide and gorm teams for their amazing frameworks

### License
This application is licensed under the MIT License. You are free to remix, adapt and redistribute it, however, please credit the original authors