package zookeeper

// Cluster is a set of hosts that share a ZookeeperCluster. Each host in is a
// viable connection target for the client library. Clusters may be reconfigured
// or reconfigure themselves after initialization. If this is the case, they
// have to ensure that updates are always fully visible to reader threads.
type Cluster interface {
	// IterateHosts creates a HostIterator that is iterating the clusters hosts.
	// Subsequent modifications of the clusters host list are not visible to the
	// iterator. There are no guarantees for a certain order.
	IterateHosts() HostIterator
}

// StaticCluster is a cluster of members that are known at initialization time.
// It does not add or remove members after initialization. The host-list may
// therefor contain failed nodes.
type StaticCluster struct {
	hosts []*Host
}

// CreateStaticCluster creates a cluster from a static set of hosts. The hosts
// can not be changed after initialization.
func CreateStaticCluster(hosts []*Host) Cluster {
	copiedHosts := make([]*Host, len(hosts))
	copy(copiedHosts, hosts)
	return &StaticCluster{
		hosts: copiedHosts,
	}
}

// IterateHosts iterates the clusters hosts.
func (cluster *StaticCluster) IterateHosts() HostIterator {
	return CreateHostIteratorFromSlice(cluster.hosts)
}

// HostIterator iterates over a finite set of hosts.
type HostIterator struct {
	items  []*Host
	index  int
	length int
}

const initialIteratorIndex = 0

// CreateHostIteratorFromSlice creates a HostIterator that iterates the hosts
// in the passed slice. This method creates a defensive copy of the hosts array.
func CreateHostIteratorFromSlice(hosts []*Host) HostIterator {
	length := len(hosts)
	items := make([]*Host, length)
	copy(items, hosts)
	return HostIterator{
		items:  items,
		length: length,
		index:  initialIteratorIndex,
	}
}

// Next returns the next Host of the iterator. This method panics if there is
// no next node, it should therefor only be called in combination with the
// HasNext() method. Items are removed from the iterator after this method was
// called. Subsequent calls doe not return the same item.
func (iterator *HostIterator) Next() *Host {
	iterator.ensureHasNext()
	iterator.removeLastItem()
	return iterator.lookupNextItemAndAdvance()
}

func (iterator *HostIterator) lookupNextItemAndAdvance() *Host {
	nextHostIndex := iterator.index
	iterator.index++
	return iterator.items[nextHostIndex]
}

// HasNext returns whether the iterator has a next item. It is only safe to
// call Next(), if this method returns true. Subsequent calls to HasNext()
// always return the same result until Next() is called.
func (iterator *HostIterator) HasNext() bool {
	return iterator.index < iterator.length
}

func (iterator *HostIterator) ensureHasNext() {
	if !iterator.HasNext() {
		panic("HostIterator does not have any remaining items")
	}
}

func (iterator *HostIterator) removeLastItem() {
	lastIndex := iterator.index - 1
	if iterator.isIndexInBounds(lastIndex) {
		iterator.removeItemAtIndex(lastIndex)
	}
}

// removeItemAtIndex removes the item at the passed index from the array.
// This is done to delete the reference to the host object.
func (iterator *HostIterator) removeItemAtIndex(index int) {
	iterator.items[index] = nil
}

func (iterator *HostIterator) isIndexInBounds(index int) bool {
	return index >= 0 && index < iterator.length
}
