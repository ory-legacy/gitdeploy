# ory-am/gitdeploy

[Gitdeploy.io](http://gitdeploy.io), the first time bounded deployment for your apps out there.  
Hooked? Read the [project description](https://github.com/ory-am/gitdeploy/wiki)!

GitDeploy is built on top of the next-gen PaaS [Flynn](http://flynn.io).

## Deploy your application

**Preamble**:
* The `.gitdeploy.yml` file configures [Gitdeploy.io](http://gitdeploy.io) and must be saved to the projects root directory.
* Although your app does not have to be [12factor](http://12factor.net/) compliant, the web process needs to listen on
the `$PORT` and `$HOST` environment variables:
[Example 1](https://github.com/ory-am/gitdeploy-go-example/blob/master/main.go#L22-L23)
[Example 2](https://github.com/ory-am/gitdeploy-go-example/blob/master/main.go#L124-L125).  

### [Golang](http://golang.org/) example

See a Go example in action: 
[![Deploy gitdeploy-go-example via gitdeploy.io](https://img.shields.io/badge/gitdeploy.io-deploy%20gitdeploy--go--example/master-green.svg)](https://www.gitdeploy.io/deploy?repository=https%3A%2F%2Fgithub.com%2Fory-am%2Fgitdeploy-go-example.git)

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

## Local gitdeploy

You need [Go](http://golang.org/) and [NodeJS](http://nodejs.org/) installed on your machine. Additionally, a MongoDB instance is required. You can set up a MongoDB instance using Docker:

**On Windows and Max OS X**, download and install [Virtualbox](https://www.virtualbox.org/) and [Boot2Docker](http://boot2docker.io/). Next run

```
> boot2docker init
> VBoxManage modifyvm "boot2docker-vm" --natpf1 "guestmongodb,tcp,127.0.0.1,27017,,27017"
> boot2docker start
> boot2docker ssh
> docker run -d -p 27017:27017 library/mongo
```

**On Linux** download and install [Docker](https://www.docker.com/) and run `$ docker run -d -p 27017:27017 library/mongo`

**IMPORTANT:** If you reboot the boot2docker-vm or the host you need to restart the container as well. You can get the container id by doing `docker ps -l` and start it by doing `docker start {id}` (replace {id} with id from `docker ps -l`).

Next thing you need is a flynn cluster. First, install the [Flynn cli](https://github.com/flynn/flynn/tree/master/cli), second run `flynn install` in your console and follow the instructions.

Now you're almost done. Run in two separate terminals:

```
$ go run main.go
```

and

```
$ cd app
$ npm install -g grunt-cli bower yo generator-karma generator-angular
$ npm install
$ bower install
$ grunt serve
```

A window with Gitdeploy should open up automatically. If not, go to [localhost:9000](http://localhost:9000)

### Production

To run a production build, do:

```
$ cd app
$ npm install -g grunt-cli bower yo generator-karma generator-angular
$ npm install
$ bower install
$ grunt build
$ cd ..
$ go run main.go
```

Go to [localhost:7654](http://localhost:7654)
