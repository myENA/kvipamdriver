package driver

import (
	// "github.com/docker/libnetwork/ipams/remote"
	// "github.com/docker/libnetwork/ipams/remote/api"
	"github.com/docker/go-plugins-helpers/ipam"
	"github.com/docker/libnetwork/datastore"
	"net"
)

type IPAMDriver struct {
	Addresses ipam.AddressSpacesResponse
	Alloc     *Allocator
}

func (ipd *IPAMDriver) GetCapabilities() (*ipam.CapabilitiesResponse, error) {
	return &ipam.CapabilitiesResponse{RequiresMACAddress: true}, nil
}

// GetDefaultAddressSpaces returns the default local and global address spaces for this ipam
func (ipd *IPAMDriver) GetDefaultAddressSpaces() (*ipam.AddressSpacesResponse, error) {

	return &ipd.Addresses, nil
}

// RequestPool returns an address pool along with its unique id. Address space is a mandatory field
// which denotes a set of non-overlapping pools. pool describes the pool of addresses in CIDR notation.
// subpool indicates a smaller range of addresses from the pool, for now it is specified in CIDR notation.
// Both pool and subpool are non mandatory fields. When they are not specified, Ipam driver may choose to
// return a self chosen pool for this request. In such case the v6 flag needs to be set appropriately so
// that the driver would return the expected ip version pool.
func (ipd *IPAMDriver) RequestPool(pool *ipam.RequestPoolRequest) (*ipam.RequestPoolResponse, error) {
	key, nw, data, err := ipd.Alloc.RequestPool(pool.AddressSpace, pool.Pool, pool.SubPool, pool.Options, pool.V6)
	if err != nil {
		return nil, err
	}
	resp := new(ipam.RequestPoolResponse)
	resp.Pool = nw.String()
	resp.PoolID = key
	resp.Data = data
	return resp, nil
	
}

// Release the address from the specified pool ID
func (ipd *IPAMDriver) ReleasePool(req *ipam.ReleasePoolRequest) error {
	return ipd.Alloc.ReleasePool(req.PoolID)
}

// Request an Address
func (ipd *IPAMDriver) RequestAddress(req *ipam.RequestAddressRequest) (*ipam.RequestAddressResponse, error) {
	ip := net.ParseIP(req.Address)
	newip, data, err := ipd.Alloc.RequestAddress(req.PoolID, ip, req.Options)
	if err != nil {
		return nil, err
	}
	resp := new(ipam.RequestAddressResponse)
	resp.Address = newip.String()
	resp.Data = data
	return resp, nil
}

func (ipd *IPAMDriver) ReleaseAddress(req *ipam.ReleaseAddressRequest) error {
	return ipd.Alloc.ReleaseAddress(req.PoolID, net.ParseIP(req.Address))
}

func Init(Addresses *ipam.AddressSpacesResponse, cfg *datastore.ScopeCfg) (*IPAMDriver, error) {
	var err error
	ipd := new(IPAMDriver)
	
//	cfg := new(datastore.ScopeCfg)
//	cfg.Client.Address = "cnllab01.dev.ena.net"
//	cfg.Client.Provider = "consul"
//	cfg.Client.Config = new(store.Config)
//	cfg.Client.Config.ConnectionTimeout = 10 * time.Second
	
	ds, err := datastore.NewDataStore(datastore.GlobalScope, cfg)
	if err != nil {
		return nil, err
	}
	 
	ipd.Alloc, err = NewAllocator(ds, ds)
	return ipd, nil
	
}
