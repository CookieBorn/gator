# gator

# Requirements:
- GO
- Postgres

# Install:
You will need to run go install from the root of the directory
You will need to set up a config file with the name ".gatorconfig.json" at the root directory of the system.
The config file will need to include one json object of the form:

{"db_url":"example.postgres.url/gator","current_user_name":}

# Usage:
- login <username>
Login as existing user
- register <username>
Create new user
- reset
Remove all data from database
- addfeed <name> <url>
Add url to users feed
- agg <time between fetch>
Fetch posts from feed urls
- browse <limit>
See title and urls of posts
- users
List users
- feeds
List feeds
- following
List feeds followed by the current user
- follow <url>
Add a feed to current users following List
- unfollow <url>
Remove a feed to current users following List
