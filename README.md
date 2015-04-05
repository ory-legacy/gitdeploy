# ory-am/gitdeploy

## Run

```
VBoxManage modifyvm "boot2docker-vm" --natpf1 "guestpostgresql,tcp,127.0.0.1,5432,,5432"
docker run --name gitdeploy-postgres -e POSTGRES_USER=gitdeploy -e POSTGRES_PASSWORD=changeme -d -p 5432:5432 postgres
```

### Configuration

```

```

## .gitignore

```
process:
    web: node server.js
```