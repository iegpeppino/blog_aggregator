# blog_aggregator / Gator

## ENG

### Description 
This go program is a simple cli called gator, that connects to a postgres database 
and allows the user to register and login as different users that can use a custom
set of commands in order to follow RSS feeds from blogs, parse their posts and 
print their contents to the screen.

### Requirements

To run this program you must have Go and Postgresql insalled on your system.
Make sure you have the latests Go and Postgres versions

### Set-up

- Clone the repo to your system and install: 

```bash
git clone https://github.com/iegpeppino/gator
cd path/to/gator
go install ...
```

- Create a Database config file in your home directory, ~/.gatorconfig.json

Be sure to replace the username, password and database name with the ones
set up in your postgres server.

```json
{
  "db_url": "postgres://username:@localhost:5432/database?sslmode=disable",
}
```

### Usage

- First, you'll have to create a new user:

```bash
gator register <username>
```

- Add a feed to the database: 

```bash
gator addfeed <rssURL>
```

- Aggregate posts every time interval (eg: 1m = 1 minute):

```bash
gator agg <time_duration>
```

- Print up to n posts to terminal:

```bash
gator browse <n>
```

### Other commands you can use

- `gator login <username>` - Log in as an existing user
- `gator users` - Lists all registered users' names
- `gator reset` - Deletes all users from the db
- `gator feeds` - Lists all feeds
- `gator follow <rssURL>` - Makes current user follow a desired feed
- `gator unfollow <rssURL>` - Makes current user unfollow a feed