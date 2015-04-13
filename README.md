# ory-am/gitdeploy

[Gitdeploy.io](http://gitdeploy.io), the first 1-click deployment for your apps out there.

[Gitdeploy.io](http://gitdeploy.io) is built on top of the next-gen PaaS [Flynn](http://flynn.io).

Try it yourself:  
[![Deploy gitdeploy-go-example via gitdeploy.io](https://img.shields.io/badge/gitdeploy.io-deploy%20gitdeploy--go--example/master-green.svg)](http://localhost:8124/deploy?repository=https%3A%2F%2Fgithub.com%2Fory-am%2Fgitdeploy-go-example.git)

## Deploy your application

The `.gitdeploy.yml` file configures [Gitdeploy.io](http://gitdeploy.io) and must be saved to the projects root directory.

### [Golang](http://golang.org/)

See a Go example in action:  
[![Deploy gitdeploy-go-example via gitdeploy.io](https://img.shields.io/badge/gitdeploy.io-deploy%20gitdeploy--go--example/master-green.svg)](http://localhost:8124/deploy?repository=https%3A%2F%2Fgithub.com%2Fory-am%2Fgitdeploy-go-example.git)

Flynn uses Heroku-like buildpacks to deploy Go applications: [Deploy Go](https://flynn.io/docs/how-to-deploy-go)  
To deploy your app via [Gitdeploy.io](http://gitdeploy.io), you'll need a `.gitdeploy.yml` file which combines
`Procfile` and `.godir`.

**Use [Godep](https://github.com/tools/godep):** As suggested in the [deploy Go on Fylnn](https://flynn.io/docs/how-to-deploy-go) docs you *should* use
[Godep](https://github.com/tools/godep) for your dependencies to significantly reduce deployment time.

```yml
# .gitdeploy.yml

# Learn more about Procfile: https://devcenter.heroku.com/articles/procfile
Procfile:
    web: myexample


# Learn more about .godir: https://github.com/kr/heroku-buildpack-go#godir-and-godeps
Godir:
    github.com/user/myexample
```

### 

## Contribute

Windows / Mac OSX

* [Virtualbox](https://www.virtualbox.org/)
* [Boot2Docker](http://boot2docker.io/)

```
> boot2docker init
> VBoxManage modifyvm "boot2docker-vm" --natpf1 "guestmongodb,tcp,127.0.0.1,27017,,27017"
> boot2docker start
> boot2docker ssh
> docker run -d -p 27017:27017 dockerfile/mongodb
```

Linux

* docker https://www.docker.com/

```
$ docker run -d -p 27017:27017 dockerfile/mongodb
```

### Configuration

tbd

### Things to note

**Why we're using MongoDB**

* Short lifetime of database items
* No 100% consistency needed (neither logs nor apps)
* Many writes to logs