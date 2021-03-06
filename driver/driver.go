// This package is an experiment in making a remote ipam driver for docker.
// Most of the code here was lifted from github.com/docker/libnetwork/ipam,
// the default docker ipam, which is copyright Docker Inc.
package driver

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/ipam"
	"github.com/docker/libnetwork/datastore"
	"net"
)

// This implements the Ipam interface
type IPAMDriver struct {
	Addresses ipam.AddressSpacesResponse
	Alloc     *Allocator
}

// Callback for our driver to determine whether or not it requires
func (ipd *IPAMDriver) GetCapabilities() (*ipam.CapabilitiesResponse, error) {
	log.Debugln("GetCapabilities() called")
	return &ipam.CapabilitiesResponse{RequiresMACAddress: true}, nil
}

// GetDefaultAddressSpaces returns the default local and global address spaces for this ipam
func (ipd *IPAMDriver) GetDefaultAddressSpaces() (*ipam.AddressSpacesResponse, error) {
	log.Debugln("GetDefaultAddressSpaces() called")
	return &ipd.Addresses, nil
}

// RequestPool returns an address pool along with its unique id. Address space is a mandatory field
// which denotes a set of non-overlapping pools. pool describes the pool of addresses in CIDR notation.
// subpool indicates a smaller range of addresses from the pool, for now it is specified in CIDR notation.
// Both pool and subpool are non mandatory fields. When they are not specified, Ipam driver may choose to
// return a self chosen pool for this request. In such case the v6 flag needs to be set appropriately so
// that the driver would return the expected ip version pool.
func (ipd *IPAMDriver) RequestPool(pool *ipam.RequestPoolRequest) (*ipam.RequestPoolResponse, error) {
	log.Debugf("RequestPool called with argument %#v", pool)
	key, nw, data, err := ipd.Alloc.RequestPool(pool.AddressSpace, pool.Pool, pool.SubPool, pool.Options, pool.V6)

	if err != nil {
		log.Errorf("RequestPool returned error: %s", err.Error())
		return nil, err
	}

	resp := &ipam.RequestPoolResponse{
		Pool:   nw.String(),
		PoolID: key,
		Data:   data,
	}

	log.Debugf("RequestPool returning %#v", resp)
	return resp, nil

}

// Release the address from the specified pool ID
func (ipd *IPAMDriver) ReleasePool(req *ipam.ReleasePoolRequest) error {

	log.Debugf("ReleasePool called with %#v", req)
	err := ipd.Alloc.ReleasePool(req.PoolID)

	if err != nil {
		log.Errorf("Error returned from ReleasePool: %s", err.Error())
	}

	return err
}

// Request an Address
func (ipd *IPAMDriver) RequestAddress(req *ipam.RequestAddressRequest) (*ipam.RequestAddressResponse, error) {

	log.Debugf("RequestAddress called with %#v", req)
	ip := net.ParseIP(req.Address)
	newip, data, err := ipd.Alloc.RequestAddress(req.PoolID, ip, req.Options)

	if err != nil {
		log.Errorf("error returned from RequestAddress: %s", err.Error())
		return nil, err
	}

	resp := new(ipam.RequestAddressResponse)
	resp.Address = newip.String()
	resp.Data = data
	log.Debugf("RequestAddress returning %#v", resp)
	return resp, nil
}

// Releases an address (not just a clever name)
func (ipd *IPAMDriver) ReleaseAddress(req *ipam.ReleaseAddressRequest) error {
	log.Debugf("ReleaseAddress called with %#v", req)

	err := ipd.Alloc.ReleaseAddress(req.PoolID, net.ParseIP(req.Address))
	if err != nil {
		log.Errorf("Error returned from ReleaseAddress: %s", err.Error())
	}
	return err
}

// This creates a new instance of IPAMDriver
func NewIPAMDriver(Addresses *ipam.AddressSpacesResponse, cfg *datastore.ScopeCfg) (*IPAMDriver, error) {
	var err error
	log.Debugf("Init called")

	if Addresses == nil {
		err = fmt.Errorf("Invalid Addresses")
		log.Error(err)
		return nil, err
	}

	dsg, err := datastore.NewDataStore(Addresses.GlobalDefaultAddressSpace, cfg)

	if err != nil {
		log.Errorf("Error returned from init: %s", err.Error())
		return nil, err
	}

	dsl, err := datastore.NewDataStore(Addresses.LocalDefaultAddressSpace, cfg)

	if err != nil {
		log.Errorf("Error returned from init: %s", err.Error())
		return nil, err
	}

	ipd := &IPAMDriver{Addresses: *Addresses}
	ipd.Alloc, err = NewAllocator(dsl, dsg)

	if err != nil {
		log.Errorf("NewAllocator returned error: %s", err.Error())
		return nil, err
	}

	log.Debug("Init success")
	return ipd, nil

}
