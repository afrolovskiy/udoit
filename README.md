# U DO IT (work in progress)

## Installation

Was tested on OSX. Should also work on Unix systems.

Requires Go 1.7+. Make sure you have Go properly installed, including setting up your GOPATH.

Create directory in your GOPATH and move folder with project there:

    $ mkdir -p $GOPATH/src/github.com/udoit
    $ mv udoit $GOPATH/src/github.com/udoit

Now you can install it:

    $ go install github.com/afrolovskiy/udoit

## Usage

    $ UDOIT_API_TOKEN='changeme' udoit
    2017/02/11 15:22:58 Authorized on account udoittestbot

## Development

### Environment

[Vagrant](https://www.vagrantup.com/) can be used for setting up development environment.

To create a new VM instance execute:

    $ vagrant up

To login trough SSH to VM machine execute:

    $ vagrant ssh

To run udoit execute:

    $ UDOIT_API_TOKEN='changeme' go run udoit.go

### Vendoring

All external code must be vendored.
[govendor](https://github.com/kardianos/govendor) used for dependency management.
