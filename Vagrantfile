# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure(2) do |config|
  # The most common configuration options are documented and commented below.
  # For a complete reference, please see the online documentation at
  # https://docs.vagrantup.com.

  # Every Vagrant development environment requires a box.
  config.vm.box = "ubuntu/xenial64"

  # Setup private network to make nfs work
  config.vm.network :private_network, ip: "10.11.12.13"

  # Enable agent forwarding over SSH connections
  config.ssh.forward_agent = true

  # Share an additional folder to the guest VM.
  config.vm.synced_folder ".", "/home/ubuntu/src/github.com/afrolovskiy/udoit"

  # Enable provisioning with a shell script.
  config.vm.provision "shell", path: 'scripts/provision.sh'
end
