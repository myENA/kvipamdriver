# Experimental docker ipam remote driver.

Emphasis on experimental.  This is basically a repackaging of the
default docker ipam driver as a remote version, but with one important
difference: both local and global all use a distributed kv store
backend.  Currently this only supports consul, but it would be trivial
to support zookeeper and etcd since this is using libkv and
libnetwork/datastore to implement the consul functionality and these
readily support zookeeper and etcd.

## Why?

Currently there isn't a good way to distribute IP assignments in the
same subnet across multiple docker host without using swarm and
overlay networks.  Overlay networks are difficult to expose outside of
the docker swarm, but other options exists such as creating a custom
docker bridge and having that bridge exist as part of a larger
network.  In this case, however, if you aren'tcusing swarm / overlay,
the IPAM for bridge networks is always local to the machine, even
though the functionality clearly exists to use a distributed KV store.
The ultimate goal of this project is to allow the distributed nature
of the default IPAM docker driver to work for non-overlay networks.

##Build Requirements:
- linux
- go 1.5.x (tested on go 1.5.3)
- glide 0.8.3

##Installation Requirements:
- docker 1.9+ (tested on 1.10.1, might work on 1.9)

##Installation Instructions:
- ensure that the GO15VENDOREXPERIMENT=1 environment variable is set  and exported
- clone the repo into the appropriate location in your GOPATH
- cd into the repo
- glide install
- cd kvipamdriver
- go build ./...
- if all went well, you should have a binary in your current directory
  called kvipamdriver.  Copy this to each docker machine that will
  share the network.  This will require write access to
  /etc/docker/plugins directory.

##Copyright

As mentioned previously, this is mostly a copy of the Docker
github.com/libnetwork/ipam package, with some minor modifications to
allow it to run as a standalone remote IPAM driver and to always use a
distributed KV store.  Those portions are copyright Docker Inc. under
the Apache 2.0 license.  The rest is copyright Education Networks of
America and licensed under the Apache 2.0 license.
