package godns

import (
	"container/list"
	"math"
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

func newDnsCacheKeyFromMessage(message DnsMessage) DnsCacheKey {
	qname := message.Qname()
	qtype := message.Qtype()
	qclass := message.Qclass()
	return DnsCacheKey{
		Qname:  qname,
		Qtype:  qtype,
		Qclass: qclass,
	}
}

type DnsCacheValue struct {
	Answers    []DnsRecord
	Authority  []DnsRecord
	Additional []DnsRecord
}

type DnsCacheEntry struct {
	TTL       int
	Timestamp time.Time
	Hits      int
	Value     DnsCacheValue
}

// internal nodes
type freqCacheNode struct {
	freq   int
	values *list.List
}

type valueCacheNode struct {
	freqParent *list.Element
	key        DnsCacheKey
	value      DnsCacheValue
	timestamp  time.Time
	ttl        int
}

func (vn *valueCacheNode) expiresAt() time.Time {
	return vn.timestamp.Add(time.Duration(vn.ttl) * time.Second)
}

func (vn *valueCacheNode) expired() bool {
	return !time.Now().Before(vn.expiresAt())
}

func newFreqNode(freq int) *freqCacheNode {
	freqNode := &freqCacheNode{}
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

func (c *DnsCache) Put(key DnsCacheKey, value DnsCacheValue) bool {
	if c.capacity == 0 {
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	existingValueEl, exists := c.hashmap[key]
	if exists {
		valueNode := existingValueEl.Value.(*valueCacheNode)
		// evict cache entry
		if valueNode.expired() {
			c.evictKey(key)
		} else {
			return false
		}
	}

	// check if cache is at capacity and if so evict
	if c.size == c.capacity {
		c.evict()
	}

	// create value node
	ttl := getTTL(value.Answers, value.Authority, value.Additional)
	valueNode := &valueCacheNode{key: key, value: value, timestamp: time.Now(), ttl: ttl}

	// value node should point to frequency node 1 at the head
	// create the node if it doesn't exist
	freqHead := c.freqs.Front()
	if freqHead == nil {
		freqNode := newFreqNode(1)
		freqHead = c.freqs.PushFront(freqNode)
	} else {
		freqNode := freqHead.Value.(*freqCacheNode)
		if freqNode.freq != 1 {
			freqNode := newFreqNode(1)
			freqHead = c.freqs.PushFront(freqNode)
		}
	}

	// add value node to the front of doubly linked list of freq node
	valueNode.freqParent = freqHead
	freqNode := freqHead.Value.(*freqCacheNode)
	valueEl := freqNode.values.PushFront(valueNode)
	c.hashmap[key] = valueEl

	// increment cache size
	c.size++

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

	valueNode := valueEl.Value.(*valueCacheNode)

	// check if value has expired, if so evict from cache
	if valueNode.expired() {
		c.evictValueEl(valueEl)
		return nil
	}

	// increment frequency
	// get the current freq node and the next freq node
	freqEl := valueNode.freqParent
	freqNode := freqEl.Value.(*freqCacheNode)
	neighborFreqEl := freqEl.Next()
	nextFreq := freqNode.freq + 1

	var nextFreqNode *freqCacheNode
	var nextFreqEl *list.Element

	if neighborFreqEl == nil {
		// no next freq node
		// create one thats the current freq + 1
		nextFreqNode = newFreqNode(nextFreq)
		nextFreqEl = c.freqs.InsertAfter(nextFreqNode, freqEl)
	} else {
		neighborFreqNode := neighborFreqEl.Value.(*freqCacheNode)
		if neighborFreqNode.freq == nextFreq {
			// next freq node is already current + 1
			nextFreqNode = neighborFreqNode
			nextFreqEl = neighborFreqEl
		} else {
			// next freq node is not current + 1 so create it
			nextFreqNode = newFreqNode(nextFreq)
			nextFreqEl = c.freqs.InsertAfter(nextFreqNode, freqEl)
		}
	}

	// add the value node to the next frequency node
	freqNode.values.Remove(valueEl)
	newValueEl := nextFreqNode.values.PushFront(valueNode)

	// no cache entries with current freq hits
	// remove the freq node
	if freqNode.values.Len() == 0 {
		c.freqs.Remove(freqEl)
	}
	valueNode.freqParent = nextFreqEl
	c.hashmap[key] = newValueEl

	entry := &DnsCacheEntry{}
	entry.Value = valueNode.value
	entry.Timestamp = valueNode.timestamp
	entry.Hits = nextFreq
	entry.TTL = valueNode.ttl
	return entry
}

func (c *DnsCache) getEvictionCandidate() *list.Element {
	freqHead := c.freqs.Front()
	if freqHead == nil {
		return nil
	}
	freqNode := freqHead.Value.(*freqCacheNode)
	valueTail := freqNode.values.Back()
	if valueTail == nil {
		return nil
	}
	return valueTail
}

func (c *DnsCache) evict() {
	// get the eviction candidate
	candidate := c.getEvictionCandidate()

	// remove the candidate
	c.evictValueEl(candidate)
}

func (c *DnsCache) evictKey(key DnsCacheKey) {
	valueEl, _ := c.hashmap[key]
	c.evictValueEl(valueEl)
}

func (c *DnsCache) evictValueEl(valueEl *list.Element) {
	if valueEl == nil {
		return
	}

	// get the freq parent node
	valueNode := valueEl.Value.(*valueCacheNode)
	freqEl := valueNode.freqParent
	freqNode := freqEl.Value.(*freqCacheNode)

	// remove the value node
	freqNode.values.Remove(valueEl)

	// remove key from the hashmap
	delete(c.hashmap, valueNode.key)

	// if freq node has no cache entries, remove the freq node
	if freqNode.values.Len() == 0 {
		c.freqs.Remove(freqEl)
	}

	c.size--
}

func (c *DnsCache) Size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.size
}

// utils
func getTTL(answers []DnsRecord, authority []DnsRecord, additional []DnsRecord) int {
	ttl := math.MaxInt
	records := []DnsRecord{}
	records = append(records, answers...)
	records = append(records, authority...)
	records = append(records, additional...)
	for _, record := range records {
		ttl = min(ttl, int(record.TTL))
	}
	return ttl
}
