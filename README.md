![autobd] (logo.png?raw=true "autobd")
# Backing you up since whenever...

[![Go Report Card](https://goreportcard.com/badge/github.com/tywkeene/autobd)](https://goreportcard.com/report/github.com/tywkeene/autobd)
[![Stories in Ready](https://badge.waffle.io/tywkeene/autobd.svg?label=ready&title=Ready)](http://waffle.io/tywkeene/autobd)
[![Build Status](https://travis-ci.org/tywkeene/autobd.svg)](https://travis-ci.org/tywkeene/autobd)
[![GoDoc](https://godoc.org/github.com/tywkeene/autobd?status.svg)](https://godoc.org/github.com/tywkeene/autobd)
[![Gitter](https://badges.gitter.im/autobd/Lobby.svg)](https://gitter.im/autobd/Lobby)

## What is it?

Autobd is an automatic backup daemon. What is that? A daemon (derived from the evil word 'demon') is a process or program that
runs (mostly) silently in the background, and handles certain tasks. This 'Demons' task is to back up a directory tree.

Say you have three servers. A, B and C. You have files on server A that you want on servers B and C. You don't really want to
some convoluted and hacky script, you need something that is just going to work in the background, require almost no configuration
or installation of any additional software, and is super easy to get running via docker.

### Enter Autobd.

All you need to do on server A is set the directory you want to watch via the ```DATA_DIR``` variable in ```scripts/docker/deploy-server.sh```, then you can 
start the daemon by running this script. This will start a single autobd server instance, running inside of docker. You can ```curl
http://0.0.0.0:8080/version```, and you will get version information from the server.

Now you just need to get your nodes going. This is just like the server, except you're using another script, pre-written. All
you need to do is let the node know which directory you want it to put synced files into, via the ```DATA_DIR``` variable in
the ```scripts/docker/deploy-nodes.sh``` script and run the it. You're good to go.

Once the nodes and server are running, the nodes will identify with the server, and request an initial file index from there server.
The first sync of course will be quite large, since the node is getting everything, but subsequent syncs will usually be much 
smaller, unless you dump a lot of files onto the server at once.

The nodes will do this again and again in an interval, that you can define in etc/config.toml.node. Additionally, each node will
send the server a 'heartbeat' in an interval. This lets the server know of the status of each node, if it's synced, and if it's
still online. So you need only access the server, to get information on the entire 'horde' of 'demons' :)

That's autobd in a nutshell.

See the [tutorial](./Documentation/TUTORIAL.md) for a more in-depth guide.


## Getting autobd

[![Join the chat at https://gitter.im/autobd/Lobby](https://badges.gitter.im/autobd/Lobby.svg)](https://gitter.im/autobd/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
Golang is required, so [get that](https://golang.org/doc/install). Most versions should work but it's probably best to
stick with at least 1.4.2 as that's what I'm developing on.

After you have go installed, it's as simple as either cloning the git repo:

`git clone https://github.com/tywkeene/autobd`

or using go get:

`go get github.com/tywkeene/autobd`

and to build

```
$ cd autobd
$ ./build.sh
```

## Getting it running

To run autobd , simply do: `./autobd -config etc/config.toml`

Autobd ships with two configuration files, config.toml.server and config.toml.node, to get you started running both


### Dockerfile

Autobd ships with two Dockerfiles and the scripts to deploy both server and node containers. You'll of course need [docker](https://docs.docker.com/engine/installation/)

#### scripts/docker/deploy-server.sh
Usage: ./scripts/docker/deploy-server.sh

Removes old autobd-server container, builds a new image, and deploys a new autobd-server container using the arguments:
```
DATA_DIR="/home/$USER/data/server-data"
SECRET_DIR="/home/$USER/secret"
ETC_DIR="/home/$USER/etc/autobd"
PORT=8080
```
These variables may be modified to suit your needs.


#### scripts/docker/deploy-nodes.sh
Usage: ./scripts/docker/deploy-nodes.sh n


Builds new autobd-node images, and deploys n nodes

#### scripts/docker/kill-nodes.sh
Usage: ./scripts/docker/kill-nodes.sh n

Runs ```docker rm -f```for n nodes

#### scripts/docker/setup-network.sh
Usage ./scripts/docker/setup-network.sh

Sets up a docker network for server and nodes to communicate. 
(I basically just use this testing, may or may not be useful to you)

## The API
See [API.md](./Documentation/API.md)

## Contributing

Open a pull request or issue with your idea or bug. Really. Anything wrong with anything anywhere, open a ticket and let me know,
then we can work on it together :) (Just be sure to check the [story board](https://waffle.io/tywkeene/autobd) before creating a new ticket) and checkout the [Development Process](./Documentation/development_process.png) flowchart.
