#!/usr/bin/env bash
set -e

# "ubuntu/xenial64" box doesn't have default "vagrant" user.
# Default user is "ubuntu"

# Fix pty/locale errors
export DEBIAN_FRONTEND=noninteractive
export LANGUAGE=en_US.UTF-8
export LANG=en_US.UTF-8
export LC_ALL=en_US.UTF-8
# locale-gen en_US.UTF-8
dpkg-reconfigure locales
echo 'LC_ALL="en_US.UTF-8"' >> /etc/environment

# Set timezone
ln -sf /usr/share/zoneinfo/UTC /etc/localtime

# Register postgresql repository
echo 'deb http://apt.postgresql.org/pub/repos/apt/ trusty-pgdg main' > /etc/apt/sources.list.d/pgdg.list
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
apt-get update

# Setup postgres
apt-get install -y git postgresql-9.5

# Init database
service postgresql start
sudo -u postgres createuser -d udoit
sudo -u postgres createdb -O udoit udoit
echo 'host all all 0.0.0.0/0 trust' > /etc/postgresql/9.5/main/pg_hba.conf
echo 'local all all trust' >> /etc/postgresql/9.5/main/pg_hba.conf
sed -i "s/#listen_addresses = 'localhost'/listen_addresses = '*'/" /etc/postgresql/9.5/main/postgresql.conf
service postgresql restart

# Setup Go
VERSION='1.7'
curl -sL "https://storage.googleapis.com/golang/go$VERSION.linux-amd64.tar.gz" | tar -C /usr/local -xzf -
echo 'export PATH=$PATH:/usr/local/go/bin' >> /home/ubuntu/.bashrc

# Setup the $GOPATH and append path to bin folder to the $PATH
echo 'export GOPATH=/home/ubuntu' >> /home/ubuntu/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> /home/ubuntu/.bashrc

# Auto cd to sources dir
echo 'cd ~/src/github.com/afrolovskiy/udoit' >> /home/ubuntu/.bashrc

# Install udoit
export PATH=$PATH:/usr/local/go/bin
export GOPATH=/home/ubuntu
cd $GOPATH/src/github.com/afrolovskiy/udoit
go install github.com/afrolovskiy/udoit

# Set ubuntu as owner of $GOPATH
chown ubuntu:ubuntu -R $GOPATH
