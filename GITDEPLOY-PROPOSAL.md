**.gitdeploy.yml** configuration options (proposal)

```yml
# This file defines the GitDeploy configuration settings. 
# Declare the configuration's version number (MANDATORY)
version: 0.1

# Declare custom Buildpack (OPTIONAL)
buildpack: https://github.com/ddollar/heroku-buildpack-multi.git

# Set up environment variables (OPTIONAL)
# Warning: Variables called HOST and PORT will be ignored
env:
    FOO: bar
    SESSION_SECRET: notasecret
    
# Attach services (OPTIONAL)
addons:
    MongoDB:
        # Requires a specific MongoDB version (RECOMMENDED) (DEFAULT: 3.0)
        version: 3.0
        # Set up environment variable bindings
        # Warning: Variables called HOST and PORT will be overriden
        # INFO: In the following case, the mongodb instance's hostname is going
        # to be bound to the environment variable called $MGO_HOST.
        # The environment variable's name is arbitrary, e.g. $MONGODB_HOSTNAME or $DATABASE_HOST.
        host: MGO_HOST
        port: MGO_PORT
        db: MGO_DB
        url: MGO_URL
        # In this case, the environment variables could look like (these are examplary values!):
        # <app-id>.gitdeploy.io:$ env
        # > MGO_HOST=mongodb.gitdeploy.io
        # > MGO_PORT=27017
        # > MGO_DB=uihg0w9834tpi
        # > MGO_URL=mongodb://fgu38gak4584:gjk259agjg8@mongodb.gitdeploy.io:27017/uihg0w9834tpi
        
    Postgres:
        # Require a specific MongoDB version (RECOMMENDED) (DEFAULT: 9.4)
        version: 9.4
        # Set up environment variable bindings
        # Warning: Variables called HOST and PORT will be overriden
        # INFO: In the following case, the mongodb instance's hostname is going
        # to be bound to the environment variable called $PG_HOST.
        # The environment variable's name is arbitrary, e.g. $POSTGRES_HOSTNAME or $DATABASE_HOST.
        user: PG_USER
        password: PG_PW
        host: PG_HOST
        port: PG_PORT
        db: PG_DB
        url: PG_URL
        # In this case, the environment variables could look like (these are examplary values!):
        # <app-id>.gitdeploy.io:$ env
        # > PG_USER=gh57whe456ujdxj
        # > PG_PW=5zrur6kf
        # > PG_HOST=postgres.gitdeploy.io
        # > PG_PORT=5432
        # > PG_DB=kg0923ruhe45
        # > PG_URL=mongodb://gh57whe456ujdxj:5zrur6kf@postgres.gitdeploy.io:5432/kg0923ruhe45
        
# Specify processes (RECOMMENDED)
# MANDATORY for Go applications
# OPTIONAL for Node applications
procs:
    # The web process accessible through http (RECOMMENDED)
    web: myexample
    # If it's a node app, this could look like:
    # web: node server.js
    # Specify additional processes. Keys are arbitrary (OPTIONAL)
    worker: myworker
    clock: myclock

# REQUIRED when not using [godep](https://github.com/tools/godep)
# SHOULD NOT be used when using [godep](https://github.com/tools/godep)
godir: github.com/user/myexample
```
