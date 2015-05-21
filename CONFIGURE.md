```yml
#This file defines the GitDeploy configuration settings. 
#Declare the configuration's version number (MANDATORY)
GDVersion: 0.1

# Declare custom Buildpack (OPTIONAL)
GDBuildpack: https://github.com/ddollar/heroku-buildpack-multi.git

# Set environment variables (OPTIONAL) @arekkas: are these just taken from local env, or are they specific to GD?
# Warning: Variables called HOST and PORT will be ignored
GDEnv:
    # $FOO=bar
    foo: bar
    # $BAZ=foo
    baz: foo
    
# Attach services (OPTIONAL)
GDAddons:
    MongoDB:
        # Requires a specific MongoDB version (RECOMMENDED) (DEFAULT: 3.0)
        version: 3.0
        # Set up environment variable bindings
        # Warning: Variables called HOST and PORT will be ignored
        user: MGO_USER
        password: MGO_PW
        host: MGO_HOST
        port: MGO_PORT
        db: MGO_DB
        url: MGO_URL
        
    Postgres:
        # Require a specific MongoDB version (RECOMMENDED) (DEFAULT: 9.4)
        version: 9.4
        # Set up environment variable bindings here
        # Warning: Variables called HOST and PORT will be ignored
        user: PG_USER
        password: PG_PW
        host: PG_HOST
        port: PG_PORT
        db: PG_DB
        url: PG_URL
        
# Specify processes (RECOMMENDED)
# MANDATORY for Go applications
# OPTIONAL for Node applications
GDProcs:
    # The web process accessible through http (RECOMMENDED)
    web: myexample
    # Specify additional processes. Keys are arbitrary (OPTIONAL)
    worker: myworker
    clock: myclock

# Specify the go directory (RECOMMENDED for Go Applications). Default is export GDGodir = "$GOPATH"
GDGodir: github.com/user/myexample
```
