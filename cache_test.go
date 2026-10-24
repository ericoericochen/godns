package godns

import (
	"reflect"
	"testing"
	"time"
)

// Verifies that looking up a missing key returns nil without changing the cache size or existing entries.
func TestDnsCacheGetMissingKey(t *testing.T) {
	cache := NewDnsCache(1)
	missing := DnsCacheKey{Qname: "missing.example.com", Qtype: TypeA, Qclass: ClassIN}
	if entry := cache.Get(missing); entry != nil {
		t.Fatalf("Get on an empty cache = %+v, want nil", entry)
	}
	if size := cache.Size(); size != 0 {
		t.Fatalf("Size after missing Get on an empty cache = %d, want 0", size)
	}

	key := DnsCacheKey{Qname: "example.com", Qtype: TypeA, Qclass: ClassIN}
	value := DnsCacheValue{Answers: []DnsRecord{{TTL: 60, RData: "192.0.2.1"}}}
	if !cache.Put(key, value) {
		t.Fatal("Put failed to insert the cache entry")
	}
	if entry := cache.Get(missing); entry != nil {
		t.Fatalf("Get for a missing key in a populated cache = %+v, want nil", entry)
	}
	if size := cache.Size(); size != 1 {
		t.Fatalf("Size after missing Get in a populated cache = %d, want 1", size)
	}
	entry := cache.Get(key)
	if entry == nil {
		t.Fatal("Get for the existing key returned nil")
	}
	if !reflect.DeepEqual(entry.Value, value) {
		t.Fatalf("Existing value after missing Get = %+v, want %+v", entry.Value, value)
	}
	if entry.Hits != 2 {
		t.Fatalf("Existing entry hits after missing Get and one successful Get = %d, want 2", entry.Hits)
	}
}

// Verifies that capacity zero disables caching by rejecting insertions and keeping all lookups empty.
func TestDnsCacheZeroCapacityDisablesCache(t *testing.T) {
	cache := NewDnsCache(0)
	value := DnsCacheValue{Answers: []DnsRecord{{TTL: 60}}}
	keys := []DnsCacheKey{
		{Qname: "a.example.com", Qtype: TypeA, Qclass: ClassIN},
		{Qname: "b.example.com", Qtype: TypeA, Qclass: ClassIN},
	}
	for _, key := range keys {
		if cache.Put(key, value) {
			t.Errorf("Put(%s) succeeded with caching disabled, want false", key.Qname)
		}
		if entry := cache.Get(key); entry != nil {
			t.Errorf("Get(%s) with caching disabled = %+v, want nil", key.Qname, entry)
		}
		if size := cache.Size(); size != 0 {
			t.Errorf("Size with caching disabled = %d, want 0", size)
		}
	}
}

// Verifies that inserting a duplicate live key fails and preserves the original value and cache size.
func TestDnsCachePutRejectsDuplicateLiveKey(t *testing.T) {
	cache := NewDnsCache(1)
	key := DnsCacheKey{Qname: "example.com", Qtype: TypeA, Qclass: ClassIN}
	original := DnsCacheValue{Answers: []DnsRecord{{TTL: 60, RData: "192.0.2.1"}}}
	replacement := DnsCacheValue{Answers: []DnsRecord{{TTL: 120, RData: "192.0.2.2"}}}
	if !cache.Put(key, original) {
		t.Fatal("Put failed to insert the original entry")
	}
	if cache.Put(key, replacement) {
		t.Fatal("Put succeeded for a duplicate live key, want false")
	}
	if size := cache.Size(); size != 1 {
		t.Fatalf("Size after duplicate Put = %d, want 1", size)
	}
	entry := cache.Get(key)
	if entry == nil {
		t.Fatal("Get returned nil after duplicate Put")
	}
	if !reflect.DeepEqual(entry.Value, original) {
		t.Fatalf("Value after duplicate Put = %+v, want %+v", entry.Value, original)
	}
	if entry.Hits != 2 {
		t.Fatalf("Hits after duplicate Put and one Get = %d, want 2", entry.Hits)
	}
}

// Verifies that inserting an expired key succeeds without a preceding Get and stores a fresh entry.
func TestDnsCachePutReplacesExpiredKey(t *testing.T) {
	cache := NewDnsCache(1)
	key := DnsCacheKey{Qname: "example.com", Qtype: TypeA, Qclass: ClassIN}
	expiringValue := DnsCacheValue{Answers: []DnsRecord{{TTL: 1, RData: "192.0.2.1"}}}
	freshValue := DnsCacheValue{Answers: []DnsRecord{{TTL: 60, RData: "192.0.2.2"}}}
	if !cache.Put(key, expiringValue) {
		t.Fatal("Put failed to insert the original entry")
	}

	time.Sleep(1100 * time.Millisecond)

	// Put must handle expiration itself; a Get here would clean up the expired key.
	if !cache.Put(key, freshValue) {
		t.Fatal("Put failed to replace an expired key")
	}
	if size := cache.Size(); size != 1 {
		t.Fatalf("Size after replacing expired key = %d, want 1", size)
	}
	entry := cache.Get(key)
	if entry == nil {
		t.Fatal("Get returned nil after replacing expired key")
	}
	if !reflect.DeepEqual(entry.Value, freshValue) {
		t.Fatalf("Value after replacing expired key = %+v, want %+v", entry.Value, freshValue)
	}
	if entry.TTL != 60 {
		t.Fatalf("TTL after replacing expired key = %d, want 60", entry.TTL)
	}
	if entry.Hits != 2 {
		t.Fatalf("Hits after replacing expired key and one Get = %d, want 2", entry.Hits)
	}
}

// Verifies that repeated reads preserve the cached value and size while incrementing hits on every read.
func TestDnsCacheRepeatedReads(t *testing.T) {
	cache := NewDnsCache(1)
	key := DnsCacheKey{Qname: "example.com", Qtype: TypeA, Qclass: ClassIN}
	value := DnsCacheValue{
		Answers: []DnsRecord{{Name: "example.com", Type: TypeA, Class: ClassIN, TTL: 60, RData: "192.0.2.1"}},
	}
	if !cache.Put(key, value) {
		t.Fatal("Put failed to insert the cache entry")
	}

	for read := 1; read <= 100; read++ {
		entry := cache.Get(key)
		if entry == nil {
			t.Fatalf("Get number %d returned nil", read)
		}
		if !reflect.DeepEqual(entry.Value, value) {
			t.Fatalf("Value on read %d = %+v, want %+v", read, entry.Value, value)
		}
		if entry.Hits != read+1 {
			t.Fatalf("Hits on read %d = %d, want %d", read, entry.Hits, read+1)
		}
		if size := cache.Size(); size != 1 {
			t.Fatalf("Size on read %d = %d, want 1", read, size)
		}
	}
}

// Verifies that expiration cleanup allows reinserting the same key and subsequent eviction at capacity.
func TestDnsCacheExpirationCleanupAllowsReinsertionAndEviction(t *testing.T) {
	cache := NewDnsCache(1)
	a := DnsCacheKey{Qname: "a.example.com", Qtype: TypeA, Qclass: ClassIN}
	b := DnsCacheKey{Qname: "b.example.com", Qtype: TypeA, Qclass: ClassIN}
	expiringValue := DnsCacheValue{
		Answers: []DnsRecord{{TTL: 1, RData: "192.0.2.1"}},
	}
	if !cache.Put(a, expiringValue) {
		t.Fatal("Put failed to insert A")
	}
	// Read A before expiration so cleanup must also handle an entry promoted by Get.
	if entry := cache.Get(a); entry == nil {
		t.Fatal("Get(A) returned nil before expiration")
	}

	time.Sleep(1100 * time.Millisecond)

	if entry := cache.Get(a); entry != nil {
		t.Fatalf("Get(A) after expiration = %+v, want nil", entry)
	}
	if size := cache.Size(); size != 0 {
		t.Fatalf("Size after expiration cleanup = %d, want 0", size)
	}

	freshValue := DnsCacheValue{
		Answers: []DnsRecord{{TTL: 60, RData: "192.0.2.2"}},
	}
	if !cache.Put(a, freshValue) {
		t.Fatal("Put failed to reinsert A after expiration cleanup")
	}
	entry := cache.Get(a)
	if entry == nil {
		t.Fatal("Get(A) returned nil after reinsertion")
	}
	if !reflect.DeepEqual(entry.Value, freshValue) {
		t.Fatalf("Value after reinsertion = %+v, want %+v", entry.Value, freshValue)
	}
	if entry.Hits != 2 {
		t.Fatalf("Hits after reinsertion and one Get = %d, want 2", entry.Hits)
	}
	if size := cache.Size(); size != 1 {
		t.Fatalf("Size after reinsertion = %d, want 1", size)
	}

	if !cache.Put(b, freshValue) {
		t.Fatal("Put failed to insert B")
	}
	if entry := cache.Get(a); entry != nil {
		t.Fatalf("Get(A) after inserting B = %+v, want nil", entry)
	}
	if entry := cache.Get(b); entry == nil {
		t.Fatal("Get(B) returned nil, want a surviving entry")
	}
	if size := cache.Size(); size != 1 {
		t.Fatalf("Size after eviction = %d, want 1", size)
	}
}

// Verifies that an entry is available before its TTL expires and returns nil afterward.
func TestDnsCacheGetReturnsNilAfterTTL(t *testing.T) {
	cache := NewDnsCache(1)
	key := DnsCacheKey{
		Qname:  "example.com",
		Qtype:  TypeA,
		Qclass: ClassIN,
	}
	value := DnsCacheValue{
		Answers: []DnsRecord{{TTL: 1}}, // DNS TTLs are in seconds.
	}

	if !cache.Put(key, value) {
		t.Fatal("Put failed to insert the cache entry")
	}
	if entry := cache.Get(key); entry == nil {
		t.Fatal("Get returned nil before the TTL expired")
	}

	time.Sleep(1100 * time.Millisecond)

	if entry := cache.Get(key); entry != nil {
		t.Fatalf("Get after TTL expiration = %+v, want nil", entry)
	}
}

// Verifies that N reads produce N+1 hits, including the initial insertion.
func TestDnsCacheHitsAfterGets(t *testing.T) {
	cache := NewDnsCache(1)
	key := DnsCacheKey{
		Qname:  "example.com",
		Qtype:  TypeA,
		Qclass: ClassIN,
	}
	value := DnsCacheValue{
		Answers: []DnsRecord{{TTL: 60}},
	}

	if !cache.Put(key, value) {
		t.Fatal("Put failed to insert the cache entry")
	}

	const n = 5
	var entry *DnsCacheEntry
	for i := 0; i < n; i++ {
		entry = cache.Get(key)
		if entry == nil {
			t.Fatalf("Get number %d returned nil", i+1)
		}
	}

	if entry.Hits != n+1 {
		t.Fatalf("Hits after %d Gets = %d, want %d", n, entry.Hits, n+1)
	}
}

// Verifies that reading an expired entry returns nil and reduces the cache size.
func TestDnsCacheGetRemovesExpiredEntry(t *testing.T) {
	cache := NewDnsCache(1)
	key := DnsCacheKey{
		Qname:  "example.com",
		Qtype:  TypeA,
		Qclass: ClassIN,
	}
	value := DnsCacheValue{
		Answers: []DnsRecord{{TTL: 1}},
	}

	if !cache.Put(key, value) {
		t.Fatal("Put failed to insert the cache entry")
	}
	if size := cache.Size(); size != 1 {
		t.Fatalf("Size after Put = %d, want 1", size)
	}

	time.Sleep(1100 * time.Millisecond)

	if entry := cache.Get(key); entry != nil {
		t.Fatalf("Get after TTL expiration = %+v, want nil", entry)
	}
	if size := cache.Size(); size != 0 {
		t.Fatalf("Size after Get of expired entry = %d, want 0", size)
	}
}

// Verifies that entries with equal hit counts are evicted in insertion order when none have been read.
func TestDnsCacheEvictsInLRUOrderWhenHitsAreEqual(t *testing.T) {
	const capacity = 3
	cache := NewDnsCache(capacity)
	keys := []DnsCacheKey{
		{Qname: "a.example.com", Qtype: TypeA, Qclass: ClassIN},
		{Qname: "b.example.com", Qtype: TypeA, Qclass: ClassIN},
		{Qname: "c.example.com", Qtype: TypeA, Qclass: ClassIN},
		{Qname: "d.example.com", Qtype: TypeA, Qclass: ClassIN},
		{Qname: "e.example.com", Qtype: TypeA, Qclass: ClassIN},
		{Qname: "f.example.com", Qtype: TypeA, Qclass: ClassIN},
	}
	value := DnsCacheValue{
		Answers: []DnsRecord{{TTL: 60}},
	}

	for i, key := range keys {
		if !cache.Put(key, value) {
			t.Fatalf("Put failed to insert %s", key.Qname)
		}
		if i >= capacity {
			evictedKey := keys[i-capacity]
			if entry := cache.Get(evictedKey); entry != nil {
				t.Fatalf("after inserting %s, Get(%s) = %+v, want nil", key.Qname, evictedKey.Qname, entry)
			}
			if size := cache.Size(); size != capacity {
				t.Fatalf("Size after inserting %s = %d, want %d", key.Qname, size, capacity)
			}
		}
	}

	// Check survivors only after all insertions so Gets cannot affect eviction order.
	for _, key := range keys[capacity:] {
		if entry := cache.Get(key); entry == nil {
			t.Errorf("Get(%s) returned nil, want a surviving entry", key.Qname)
		}
	}
}

// Verifies that equal-hit entries are evicted by their last read order rather than their insertion order.
func TestDnsCacheEvictsLRUTieAfterReads(t *testing.T) {
	cache := NewDnsCache(2)
	a := DnsCacheKey{Qname: "a.example.com", Qtype: TypeA, Qclass: ClassIN}
	b := DnsCacheKey{Qname: "b.example.com", Qtype: TypeA, Qclass: ClassIN}
	c := DnsCacheKey{Qname: "c.example.com", Qtype: TypeA, Qclass: ClassIN}
	value := DnsCacheValue{
		Answers: []DnsRecord{{TTL: 60}},
	}

	// Insert B first so insertion order differs from the subsequent read order.
	for _, key := range []DnsCacheKey{b, a} {
		if !cache.Put(key, value) {
			t.Fatalf("Put failed to insert %s", key.Qname)
		}
	}
	for _, key := range []DnsCacheKey{a, b} {
		entry := cache.Get(key)
		if entry == nil {
			t.Fatalf("Get(%s) returned nil before eviction", key.Qname)
		}
		if entry.Hits != 2 {
			t.Fatalf("Hits for %s = %d, want 2", key.Qname, entry.Hits)
		}
	}

	if !cache.Put(c, value) {
		t.Fatal("Put failed to insert C")
	}
	if entry := cache.Get(a); entry != nil {
		t.Fatalf("Get(A) after inserting C = %+v, want nil", entry)
	}
	for _, key := range []DnsCacheKey{b, c} {
		if entry := cache.Get(key); entry == nil {
			t.Errorf("Get(%s) returned nil, want a surviving entry", key.Qname)
		}
	}
	if size := cache.Size(); size != 2 {
		t.Fatalf("Size after eviction = %d, want 2", size)
	}
}

// Verifies that adding an entry to a full cache evicts the entry with the fewest hits.
func TestDnsCacheEvictsLeastFrequentlyUsedEntry(t *testing.T) {
	cache := NewDnsCache(3)
	a := DnsCacheKey{Qname: "a.example.com", Qtype: TypeA, Qclass: ClassIN}
	b := DnsCacheKey{Qname: "b.example.com", Qtype: TypeA, Qclass: ClassIN}
	c := DnsCacheKey{Qname: "c.example.com", Qtype: TypeA, Qclass: ClassIN}
	d := DnsCacheKey{Qname: "d.example.com", Qtype: TypeA, Qclass: ClassIN}
	value := DnsCacheValue{
		Answers: []DnsRecord{{TTL: 60}},
	}

	for _, key := range []DnsCacheKey{a, b, c} {
		if !cache.Put(key, value) {
			t.Fatalf("Put failed to insert %s", key.Qname)
		}
	}

	// A has 3 hits, B has 2, and C still has its initial 1 hit.
	for _, key := range []DnsCacheKey{a, a, b} {
		if entry := cache.Get(key); entry == nil {
			t.Fatalf("Get(%s) returned nil before eviction", key.Qname)
		}
	}

	if !cache.Put(d, value) {
		t.Fatal("Put failed to insert D")
	}
	if entry := cache.Get(c); entry != nil {
		t.Fatalf("Get(C) after inserting D = %+v, want nil", entry)
	}
	for _, key := range []DnsCacheKey{a, b, d} {
		if entry := cache.Get(key); entry == nil {
			t.Errorf("Get(%s) returned nil, want a surviving entry", key.Qname)
		}
	}
	if size := cache.Size(); size != 3 {
		t.Fatalf("Size after eviction = %d, want 3", size)
	}
}

func logDnsCache(t *testing.T, cache *DnsCache) {
	t.Helper()
	cache.mu.Lock()
	defer cache.mu.Unlock()

	t.Logf("cache: size=%d capacity=%d map entries=%d", cache.size, cache.capacity, len(cache.hashmap))
	t.Log("frequency buckets from lowest to highest; keys from MRU to LRU:")
	for freqEl := cache.freqs.Front(); freqEl != nil; freqEl = freqEl.Next() {
		node := freqEl.Value.(*freqCacheNode)
		keys := []string{}
		for valueEl := node.values.Front(); valueEl != nil; valueEl = valueEl.Next() {
			value := valueEl.Value.(*valueCacheNode)
			keys = append(keys, value.key.Qname)
		}
		t.Logf("  freq=%d: %v", node.freq, keys)
	}
}

type expectedFrequencyBucket struct {
	freq int
	keys []DnsCacheKey // Most recently used first.
}

func assertFrequencyBuckets(t *testing.T, cache *DnsCache, want []expectedFrequencyBucket) {
	t.Helper()
	cache.mu.Lock()
	defer cache.mu.Unlock()

	got := []expectedFrequencyBucket{}
	for freqEl := cache.freqs.Front(); freqEl != nil; freqEl = freqEl.Next() {
		node := freqEl.Value.(*freqCacheNode)
		bucket := expectedFrequencyBucket{freq: node.freq, keys: []DnsCacheKey{}}
		for valueEl := node.values.Front(); valueEl != nil; valueEl = valueEl.Next() {
			value := valueEl.Value.(*valueCacheNode)
			bucket.keys = append(bucket.keys, value.key)
		}
		got = append(got, bucket)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("frequency buckets (MRU first) = %+v, want %+v", got, want)
	}
}

// Verifies that reads move entries into increasing frequency buckets and order ties by read recency.
func TestDnsCacheFrequencyBucketsAfterReads(t *testing.T) {
	cache := NewDnsCache(3)
	a := DnsCacheKey{Qname: "a.example.com", Qtype: TypeA, Qclass: ClassIN}
	b := DnsCacheKey{Qname: "b.example.com", Qtype: TypeA, Qclass: ClassIN}
	value := DnsCacheValue{Answers: []DnsRecord{{TTL: 60}}}
	for _, key := range []DnsCacheKey{b, a} {
		if !cache.Put(key, value) {
			t.Fatalf("Put failed to insert %s", key.Qname)
		}
	}
	assertFrequencyBuckets(t, cache, []expectedFrequencyBucket{
		{freq: 1, keys: []DnsCacheKey{a, b}},
	})

	if cache.Get(a) == nil {
		t.Fatal("Get(A) returned nil")
	}
	assertFrequencyBuckets(t, cache, []expectedFrequencyBucket{
		{freq: 1, keys: []DnsCacheKey{b}},
		{freq: 2, keys: []DnsCacheKey{a}},
	})

	if cache.Get(b) == nil {
		t.Fatal("Get(B) returned nil")
	}
	assertFrequencyBuckets(t, cache, []expectedFrequencyBucket{
		{freq: 2, keys: []DnsCacheKey{b, a}},
	})

	if cache.Get(a) == nil {
		t.Fatal("second Get(A) returned nil")
	}
	assertFrequencyBuckets(t, cache, []expectedFrequencyBucket{
		{freq: 2, keys: []DnsCacheKey{b}},
		{freq: 3, keys: []DnsCacheKey{a}},
	})
}

// Verifies that a new entry starts at frequency one when existing entries have higher frequencies.
func TestDnsCacheNewEntryStartsInFrequencyOneBucket(t *testing.T) {
	cache := NewDnsCache(3)
	a := DnsCacheKey{Qname: "a.example.com", Qtype: TypeA, Qclass: ClassIN}
	b := DnsCacheKey{Qname: "b.example.com", Qtype: TypeA, Qclass: ClassIN}
	value := DnsCacheValue{Answers: []DnsRecord{{TTL: 60}}}
	if !cache.Put(a, value) {
		t.Fatal("Put failed to insert A")
	}
	if cache.Get(a) == nil {
		t.Fatal("Get(A) returned nil")
	}
	assertFrequencyBuckets(t, cache, []expectedFrequencyBucket{
		{freq: 2, keys: []DnsCacheKey{a}},
	})

	if !cache.Put(b, value) {
		t.Fatal("Put failed to insert B")
	}
	assertFrequencyBuckets(t, cache, []expectedFrequencyBucket{
		{freq: 1, keys: []DnsCacheKey{b}},
		{freq: 2, keys: []DnsCacheKey{a}},
	})
}
