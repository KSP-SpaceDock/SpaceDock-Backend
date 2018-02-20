# SpaceDock

SpaceDock is an open source website software that can host modifications for multiple games. It started as a fork of the popular KerbalStuff software for Kerbal Space Program and evolved since then.

SpaceDock is split into two parts:
* The frontend, which is currently developed under the name [OpenDock](https://github.com/KSP-SpaceDock/OpenDock)
* The backend which is developed here

The backend handles all the content that SpaceDock stores, like mods and users. The program is designed to be as lightweight as possible to distribute it across multiple nodes. A frontend can then query the multiple nodes. Due to the modular nature of the framework, third party persons who are interested in mod hosting can either host their own backend together with a custom frontend, or create just a frontend and rely on SpaceDocks infrastructure internally to host their mods. This allows for flexible and easily customised mod repositories (once there are free frontends available)

The API implemented by SpaceDock-Backend is the suceessor of the old KerbalStuff API that is still used in the current SpaceDock. It offers way more routes, but lacks an actual frontend. Instead of HTML code, the backend returns JSON formatted responses, or optionally JSONP. Other formats, like XML are **not** planned. The API is **not** compatible with the old KerbalStuff API. Some routes may be the same, but we don't support Applications built for the old API.

## Building SpaceDock-Backend

### Get the code

First you need to set up a GOPATH environment variable and directory. Edit your `~/.bashrc` and add:

```sh
export GOPATH="$HOME/gocode"
export PATH="$PATH:$GOPATH/bin"
```

Then clone the repository. SpaceDock's build process follows the golang specifications: You need to clone the SpaceDock repository into your `$GOPATH`. From a fresh command prompt (so your `.bashrc` changes can take effect):

```sh
mkdir -p $GOPATH/src/github.com/KSP-SpaceDock
cd $GOPATH/src/github.com/KSP-SpaceDock
git clone https://github.com/KSP-SpaceDock/SpaceDock-Backend.git
cd SpaceDock-Backend
```

After these steps, the repository should look like this:

```
$GOPATH/src/github.com/KSP-SpaceDock/SpaceDock-Backend/
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

### Get the dependencies

#### Install Go

Install Go 1.9. You can get it from https://golang.org/dl/ or from apt-get:

```sh
apt-get install golang-1.9
```

The `go` executable should be in your `$PATH`.

#### Install Glide

We use the [Glide package manager](https://glide.sh) to manage the backend's build dependencies. Execute the appropriote build script for your platform:

- Linux/MacOS: `build/install_glide.sh`
- Windows: `build\install_glide.ps1`

This will install `glide` into `$GOPATH/bin`.

#### Run glide

Install SpaceDock-Backend's build dependencies:

```
cd $GOPATH/src/github.com/KSP-SpaceDock/SpaceDock-Backend
glide install
```

### Get plugins

Several plugins are available at the [SpaceDock-Extras](https://github.com/KSP-SpaceDock/SpaceDock-Extras) repo. These are required for OpenDock.

1. Edit `build/plugins.txt` and add your plugins' goland dependency URLs. You can use full glide versioning syntax here.
2. Run the plugin build script for your platform:
   - Linux/MacOS: `build/fetch_plugins.sh`
   - Windows: `build/fetch_plugins.ps1`

Even if you don't plan to use plugins, execute the `fetch_plugins` script. It will work without a `build/plugins.txt` file being present, and just creates the `build_sdb.go` file.

### Build

To build the app you need to call `go build` on the `build_sdb.go` file. The only difference between Linux and Windows is, that Windows users should use `sdb.exe` instead of `sdb`

```sh
go build -v -o ./sdb ./build_sdb.go
```

### Runtime dependencies

#### SQL database

The backend requires an SQL database. This is used to store permanent data such as users and mods.

If you want to use PostGreSQL (recommended):

```sh
apt-get install postgresql-server-9.3 postgresql-server-dev-9.3
```

If you want to use MySQL:

```sh
apt-get install libmysqlclient-dev mysql-server 
```

Once your database is installed, you need to create the user and database for SpaceDock-Backend.

```sh
sudo -u postgres psql
postgres=# CREATE USER spacedock WITH PASSWORD 'spacedock';
postgres=# CREATE DATABASE spacedockbackend;
```

SpaceDock-Backend is a Golang Application that uses [iris](https://github.com/kataras/iris) for serving content and [gorm](https://github.com/jinzhu/gorm) for persistency. Even though we are developing and running SpaceDock using PostgreSQL, you can use any SQL based Database in combination with gorm. (That means MySQL, MariaDB). SQLite could work, but supporting it is a pain, because it uses cgo, which wouldn't allow us to crosscompile the program. At the moment, we only support Postgres.

#### Store

The backend also requires a place to store temporary data such as sessions of logged in users. For a scalable server, you should install `redis`:

```sh
apt-get install redis-server
```

However, for a small dev/test backend, you can use a memory store instead.

### Configuration

Copy the configuration file template:

```sh
cp config/config.example.yml config/config.yml
```

If this is your first time running the backend, there are several values you should change.

First, disable same-origin checking to avoid getting in trouble with CORS with errors like `Response to preflight request doesn't pass access control check: No 'Access-Control-Allow-Origin' header is present on the requested resource.`:

```yml
disable-same-origin: true
```

Then, tell the backend how to access your SQL database:

```yml
dialect: "postgres"
connection-data: "postgresql://spacedock:spacedock@127.0.0.1:5432/spacedockbackend"
```

If you chose not to install redis, switch to the memory store:

```yml
store-type: "memory"
```

### Starting the backend

To start the backend:

```sh
./sdb # sdb.exe on Windows
```

This will start a backend process listening on TCP/IP port 5000. You are now ready to start [OpenDock](https://github.com/KSP-SpaceDock/OpenDock)!

### The team

* godarklight
* V1TA5
* StollD
* RockyTV

### Special thanks

* SirCmpwn, testeddoughnut and andyleap for developing the original KerbalStuff API
* All SpaceDock contributors for their help during the migration from KerbalStuff to SpaceDock
* The iris, glide and gorm teams for their amazing frameworks

### License

This application is licensed under the MIT License. You are free to remix, adapt and redistribute it, however, please credit the original authors.
