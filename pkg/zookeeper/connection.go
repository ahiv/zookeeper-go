package zookeeper

import (
	"errors"
	"net"
	"time"
)

type Connection interface {
	Close() error
}

const networkType = "tcp"

type Host struct {
	address string
}

func (host *Host) Connect() (net.Conn, error) {
	return net.Dial(networkType, host.address)
}

func (host *Host) ConnectWithTimeout(timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(networkType, host.address, timeout)
}

type ConnectionConfig struct {
	connectTimeout time.Duration
	cluster Cluster
}

type ClusterConnection struct {
}

// ErrNoReachableHostFound is returned by a method after it fails to establish
// a connection to any host in a cluster.
var ErrNoReachableHostFound = errors.New(
	"no reachable host found in the configured cluster")

func (config *ConnectionConfig) connectAndInitialize() (net.Conn, error) {
	networkConnection, err := config.connect()
	if err != nil {
		return nil, err
	}
	err = config.initializeConnection(networkConnection)
	return networkConnection, err
}

func (config *ConnectionConfig) connect() (net.Conn, error) {
	hosts := config.cluster.IterateHosts()
	return config.connectToAnyHost(hosts)
}

// connectToAnyHost tries to establish a Connection to any host from the hosts
// iterator. If it fails, an ErrNoReachableHostFound is returned.
func (config *ConnectionConfig) connectToAnyHost(hosts HostIterator) (net.Conn, error) {
	for hosts.HasNext() {
		nextHost := hosts.Next()
		if networkConnection, err := config.connectToHost(nextHost); err == nil {
			return networkConnection, nil
		}
	}
	return nil, ErrNoReachableHostFound
}

func (config *ConnectionConfig) connectToHost(host *Host) (net.Conn, error) {
	return host.ConnectWithTimeout(config.connectTimeout)
}

func (config *ConnectionConfig) initializeConnection(networkConnection net.Conn) error {
	authentication := &authentication{connection:networkConnection}
	return authentication.run()
}

type authentication struct {
	connection net.Conn
}

func (authentication *authentication) run() error {
	return nil
}
