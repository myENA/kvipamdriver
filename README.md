# Experimental docker ipam remote driver.

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
- glide install to populate the vendor directory
- cd enaipandriver
- go build ./...
- if all went well, you should have a binary in your current directory  called enaipamdriver
- ???
- PROFIT!
