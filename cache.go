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
	TTL        int
	Answers    []DnsRecord
	Authority  []DnsRecord
	Additional []DnsRecord
}

// func getTTL(answers []DnsRecord, authority []DnsRecord, additional []DnsRecord) int {
// 	ttl := math.MaxInt
// 	records := []DnsRecord{}
// 	records = append(records, cv.Answers...)
// 	records = append(records, cv.Authority...)
// 	records = append(records, cv.Additional...)
// 	for _, record := range records {
// 		ttl = min(ttl, int(record.TTL))
// 	}

// 	return ttl
// }

type DnsCacheEntry struct {
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
}

func (vn *valueCacheNode) expiresAt() time.Time {
	return vn.timestamp.Add(time.Duration(vn.value.TTL) * time.Second)
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
	fmt.Println("putting ", key, "=", value)

	c.mu.Lock()
	defer c.mu.Unlock()

	_, exists := c.hashmap[key]
	if exists {
		return false
	}

	// create value node
	valueNode := &valueCacheNode{key: key, value: value, timestamp: time.Now()}

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
	freqNode := freqHead.Value.(*freqCacheNode)
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

	valueNode := valueEl.Value.(*valueCacheNode)

	// check if value has expired, if so evict from cache
	if valueNode.expired() {
		// TODO: add evict logic here
		return nil
	}

	// increment frequency
	freqEl := valueNode.freqParent
	freqNode := freqEl.Value.(*freqCacheNode)
	neighborFreqEl := freqEl.Next()
	nextFreq := freqNode.freq + 1

	insertAfter := false
	var nextFreqNode *freqCacheNode
	var nextFreqEl *list.Element

	if neighborFreqEl == nil {
		nextFreqNode = newFreqNode(nextFreq)
		insertAfter = true
	} else {
		neighborFreqNode := neighborFreqEl.Value.(*freqCacheNode)
		if neighborFreqNode.freq == nextFreq {
			nextFreqNode = neighborFreqNode
			nextFreqEl = neighborFreqEl
		} else {
			nextFreqNode = newFreqNode(nextFreq)
			insertAfter = true
		}
	}

	if insertAfter {
		nextFreqEl = c.freqs.InsertAfter(nextFreqNode, freqEl)
	}

	// add the value node to the next frequency node
	freqNode.values.Remove(valueEl)
	valueNode.freqParent = nextFreqEl

	entry := &DnsCacheEntry{}
	entry.Value = valueNode.value
	entry.Timestamp = valueNode.timestamp
	entry.Hits = nextFreq
	return entry
}

func (c *DnsCache) evict() {
	// get the head of the frequency doubly linked list
	freqHead := c.freqs.Front()
	if freqHead == nil {
		return
	}

	freqNode := freqHead.Value.(*freqCacheNode)
	valueTail := freqNode.values.Back()
	if valueTail == nil {
		return
	}

	// remove least recently used item at the least frequently used node
	freqNode.values.Remove(valueTail)

	// remove key from hashmap
	valueEl := valueTail.Value.(*valueCacheNode)
	delete(c.hashmap, valueEl.key)

	c.size--
}

func (c *DnsCache) Size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.size
}
