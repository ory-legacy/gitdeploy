```yml
version: 0.2 (MANDATORY)

buildpack: https://github.com/ddollar/heroku-buildpack-multi.git (OPTIONAL)

# Set environment variables (OPTIONAL)
# Warning: Variables called HOST and PORT will be ignored
env:
    # $FOO=bar
    foo: bar
    # $BAZ=foo
    baz: foo
    
# Attach services (OPTIONAL)
addons:
    MongoDB:
        # Require a specific MongoDB version (OPTIONAL)
        # Warning: Variables called HOST and PORT will be ignored
        version: 3.0
        # set up environment variable bindings
        user: MGO_USER
        password: MGO_PW
        host: MGO_HOST
        port: MGO_PORT
        db: MGO_DB
        url: MGO_URL
        
    Postgres-9.4
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
procs:
    web: myexample

# Specify the go directory (RECOMMENDED for Go Applications)
godir: github.com/user/myexample
```
