package godns

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type DnsCache struct {
	mu       sync.Mutex
	hashmap  map[DnsCacheKey]*list.Element // map of cache key to value nodes
	freqs    *list.List                    // doubly linked list of frequency nodes
	capacity int
	size     int
}

type DnsCacheKey struct {
	Qname  string
	Qtype  Type
	Qclass Class
}

type DnsCacheValue struct {
	Answers    []DnsRecord
	Authority  []DnsRecord
	Additional []DnsRecord
}

type DnsCacheEntry struct {
	Timestamp time.Time
	TTL       int
	Hits      int
	Value     DnsCacheValue
}

// internal nodes
type freqNode struct {
	freq   int
	values *list.List
}

type valueNode struct {
	freqParent *list.Element
	key        DnsCacheKey
	value      DnsCacheValue
	timestamp  time.Time
	ttl        int
}

func (vn *valueNode) expiresAt() time.Time {
	return vn.timestamp.Add(time.Duration(vn.ttl) * time.Second)
}

func (vn *valueNode) expired() bool {
	return !time.Now().Before(vn.expiresAt())
}

func newFreqNode(freq int) *freqNode {
	freqNode := &freqNode{}
	freqNode.freq = freq
	freqNode.values = list.New()
	return freqNode
}

func NewDnsCache(capacity int) *DnsCache {
	cache := &DnsCache{}
	cache.capacity = capacity
	cache.size = 0
	cache.freqs = list.New()
	cache.hashmap = map[DnsCacheKey]*list.Element{}
	return cache
}

func (c *DnsCache) Put(key DnsCacheKey, value DnsCacheValue, ttl int) bool {
	fmt.Println("putting ", key, "=", value)

	c.mu.Lock()
	defer c.mu.Unlock()

	_, exists := c.hashmap[key]
	if exists {
		return false
	}

	// create value node
	valueNode := &valueNode{key: key, value: value, timestamp: time.Now(), ttl: ttl}

	// value node should point to frequency node 1 at the head
	// create the node if it doesn't exist
	freqHead := c.freqs.Front()
	if freqHead == nil {
		fmt.Println("freq head is nil")
		freqNode := newFreqNode(1)
		freqHead = c.freqs.PushFront(freqNode)
	}

	// add value node to the front of doubly linked list of freq node
	valueNode.freqParent = freqHead
	freqNode := freqHead.Value.(*freqNode)
	valueEl := freqNode.values.PushFront(valueNode)
	c.hashmap[key] = valueEl

	// increment cache size
	c.size++

	// overfilled cache, evict
	if c.size > c.capacity {
		c.evict()
	}

	return true
}

func (c *DnsCache) Get(key DnsCacheKey) *DnsCacheEntry {
	c.mu.Lock()
	defer c.mu.Unlock()

	// O(1) lookup of value node in hashmap
	valueEl, hit := c.hashmap[key]
	if !hit {
		return nil
	}

	valueNode := valueEl.Value.(*valueNode)

	// check if value has expired, if so evict from cache
	if valueNode.expired() {
		// TODO: add evict logic here
		return nil
	}

	// increment frequency

	entry := &DnsCacheEntry{}
	entry.Value = valueNode.value
	entry.TTL = valueNode.ttl
	entry.Timestamp = valueNode.timestamp
	entry.Hits = 1
	return entry
}

func (c *DnsCache) evict() {
	// get the head of the frequency doubly linked list
	freqHead := c.freqs.Front()
	if freqHead == nil {
		return
	}

	freqNode := freqHead.Value.(*freqNode)
	valueTail := freqNode.values.Back()
	if valueTail == nil {
		return
	}

	// remove least recently used item at the least frequently used node
	freqNode.values.Remove(valueTail)

	// remove key from hashmap
	valueEl := valueTail.Value.(*valueNode)
	delete(c.hashmap, valueEl.key)

	c.size--
}

func (c *DnsCache) Size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.size
}
