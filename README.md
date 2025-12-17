# Gator

`Gator` is a terminal based RSS feed blog aggregator written in Go, backed by a
PostgreSQL database for persistent storage. It allows users to register
accounts, subscribe to RSS feeds, periodically fetch posts, and browse
aggregated post metadata with links that can be opened in a web browser.

## Learning Project Disclaimer

This project was built as part of a guided learning project from [Boot.dev](https://www.boot.dev)
and is not intended for production use.

## Tech stack

- Go
- PostgreSQL
- Goose
- sqlc

## Supported Platforms

- macOS
- Linux
- Windows via WSL

## Requirements

The following must be installed:

- Go (v1.25+)
- PostgreSQL (v15)

### Installing Go

[https://go.dev/doc/install](https://go.dev/doc/install).

### Installing PostgreSQL

#### 1. Install using the terminal

On macOS:

```bash
brew install postgresql@15
```

On Linux/WSL:

```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```

#### 2. Ensure PostgreSQL is installed correctly

To ensure that version 15 of PostgreSQL was installed correctly:

```bash
psql --version
```

If installed correctly, the output should include a version number that begins
with 15.

#### 3. Set password for `postgres` user (Linux only)

When using Linux, it is necessary to set a password:

```bash
sudo passwd postgres
```

On macOS, setting a password isn't necessary.

## Setup

### Create the `gator` database

#### 1. Ensure the PostgreSQL service is running

On macOS:

```bash
brew services start postgresql@15
```

On Linux/WSL:

```bash
sudo service postgresql start
```

#### 2. Connect to PostgreSQL server

Use the `psql` CLI to connect to the PostgreSQL server.

On macOS:

```bash
psql postgres
```

On Linux/WSL:

```bash
sudo -u postgres psql
```

#### 3. Create database

From within the `psql` CLI, create a database for the `gator` application:

```sql
CREATE DATABASE gator;
```

#### 4. Set PostgreSQL server password (Linux only)

When on Linux, a password will need to be set for the PostgreSQL server. In the
`psql` CLI, use this command to change it:

```sql
ALTER USER postgres PASSWORD 'postgres';
```

> Remember this password, as it will be used in the database URL later.

This step is not required on macOS.

### Create gator config file

The `gator` application expects a config file (`.gatorconfig`) to be located in
the `$HOME` directory.

Create this file by running the following:

```bash
touch $HOME/.gatorconfig
```

Using an editor, paste the following into `.gatorconfig`:

```json
{
  "db_url": "connection_string_goes_here",
  "current_user_name": "username_goes_here"
}
```

The URL for the database created earlier must be added to this file. The URL
follows a format like this:

```txt
postgres://username:password@host:port/database?sslmode=disable
```

- `username`:
  - macOS: username of the current system user
  - Linux/WSL: `postgres`
- `password`:
  - macOS: leave blank (but be sure to include the `:` after the username)
  - Linux/WSL: use the password set in the `psql` CLI (the one used in the
  command like `ALTER USER postgres PASSWORD 'postgres';`)
- `host`: `localhost`
- `port`: `5432` (PostgreSQL default)
- `database`: `gator` (unless a different name with the
`CREATE DATABASE` command was used)
- `?sslmode=disable`: include at the end of the URL. This is
needed when running the database locally.

For example, on macOS with no password:

```txt
postgres://alice:@localhost:5432/gator?sslmode=disable
```

Update the `db_url` field in `.gatorconfig` with the URL for the newly created
database.

Example:

```json
{
  "db_url": "postgres://alice:@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

The `current_user_name` field will be set later from within the `gator`
application.

## Installing Gator

Install `gator` by using the following command:

```bash
go install github.com/sbrown3212/gator@latest
```

This will install the binary to `$GOPATH/bin`.

## Commands

### `register`

Create a new user with the `register` command.

This command requires one argument that will be used as the user's username.

Example:

```bash
gator register alice
```

### `login`

Use the `login` command to switch to another user that has already registered.

This command requires one argument, which should be the name of the user to be
switched to.

Example:

```bash
gator login bob
```

### `addfeed`

Add a feed to the current user with the `addfeed` command.

This command requires two arguments. The first is the feed name (use
quotes if it contains spaces). The second is the URL of the
feed.

Example:

```bash
gator addfeed "Hacker News" https://news.ycombinator.com/rss
```

### `agg`

The `agg` command is used to periodically fetch posts for all feeds. It will
run continuously and is intended to be left running in a separate terminal
instance from the one used to interact with the rest of the `gator` application.

This command requires one argument, which is used to set the time delay between
fetching posts for the next feed. This argument should be written with a number
directly followed by a letter (no spaces in between). The letter corresponds to
the unit of time, and the number corresponds to the quantity. For example, `5m`
would be parsed as five minutes, `30s` would be 30 seconds, and `1h` would be
one hour.

Example:

```bash
gator agg 30s
```

This example would fetch the posts of the feed that has gone the longest without
being fetched, and again every 30 seconds, until the terminal is closed or the
process is terminated.

### `browse`

Use the `browse` command to view the most recent posts from the feeds of the
current user.

This command optionally takes one argument, a number, which is used to specify
the number of posts to display. Defaults to 2 if no argument is provided.

Example:

```bash
gator browse # shows 2 posts
# or
gator browse 10 # shows 10 posts
```

## Uninstalling Gator

### Delete the binary

```bash
rm "(go env GOPATH)/bin/gator"
```

### Delete the `.gatorconfig`

```bash
rm ~/.gatorconfig
```

### Delete the `gator` database

From within the `psql` client, run:

```sql
DROP DATABASE gator;
```
