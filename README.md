## Licensing

```
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

## Index

* [Installation](https://github.com/Natnael-Alemayehu/service?tab=readme-ov-file#installation)
* [Create Your Own Version](https://github.com/Natnael-Alemayehu/service?tab=readme-ov-file#create-your-own-version)
* [Running The Project](https://github.com/Natnael-Alemayehu/service?tab=readme-ov-file#running-the-project)
* [Joining the Go Slack Community](https://github.com/Natnael-Alemayehu/service?tab=readme-ov-file#joining-the-go-slack-community)

## Installation

To clone the project, create a folder and use the git clone command. Then please read the [makefile](makefile) file to learn how to install all the tooling and docker images.

```
$ cd $HOME
$ mkdir code
$ cd code
$ git clone https://github.com/ardanlabs/service or git@github.com:ardanlabs/service.git
$ cd service
```

## Create Your Own Version

If you want to create a version of the project for your own use, use the new gonew command.

```
$ go install golang.org/x/tools/cmd/gonew@latest

$ cd $HOME
$ mkdir code
$ cd code
$ gonew github.com/ardanlabs/service github.com/mydomain/myproject
$ cd myproject
$ go mod vendor
```

Now you have a copy with your own module name. Now all you need to do is initialize the project for git.

## Running The Project

To run the project use the following commands.

```
# Install Tooling
$ make dev-gotooling
$ make dev-brew
$ make dev-docker

# Run Tests
$ make test

# Shutdown Tests
$ make test-down

# Run Project
$ make dev-up
$ make dev-update-apply
$ make token
$ export TOKEN=<COPY TOKEN>
$ make users

# Run Load
$ make load

# Run Tooling
$ make grafana
$ make statsviz

# Shut Project
$ make dev-down
```

## Joining the Go Slack Community

We use a Slack channel to share links, code, and examples during the training.  This is free.  This is also the same Slack community you will use after training to ask for help and interact with may Go experts around the world in the community.

1. Using the following link, fill out your name and email address: https://invite.slack.gobridge.org
1. Check your email, and follow the link to the slack application.
1. Join the training channel by clicking on this link: https://gophers.slack.com/messages/training/
1. Click the “Join Channel” button at the bottom of the screen.
___
All material is licensed under the [Apache License Version 2.0, January 2004](http://www.apache.org/licenses/LICENSE-2.0).
