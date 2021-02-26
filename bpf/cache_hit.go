package bpf

// based on bcc/cachestat.py
import (
	"os"
	"fmt"
	"strconv"
	bpf "github.com/iovisor/gobpf/bcc"
	ksym "github.com/iovisor/gobpf/pkg/ksym"
)

import "C"

type CacheHit struct {
	table *bpf.Table
	m *bpf.Module
	counters map[string]uint64
	addrCache map[string]string
	Debug bool
}

func (c *CacheHit) Init(){

	c.addrCache = make(map[string]string)
	c.counters = map[string]uint64{
		"add_to_page_cache_lru":0,
		"mark_buffer_dirty":0,
		"account_page_dirtied":0,
		"mark_page_accessed":0,
	};

	if c.m = bpf.NewModule(source, []string{}); c.m == nil {
		fmt.Fprintf(os.Stderr, "module init failed\n")
		os.Exit(1)
	}

	do_count_probe, err := c.m.LoadKprobe("do_count")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load do_count : %s\n", err)
		os.Exit(1)
	}

	for k,_ := range c.counters {

		if err := c.m.AttachKprobe(k, do_count_probe, -1); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to attach %s: %s\n",k, err)
			os.Exit(1)
		}
	}

	if c.table = bpf.NewTable(c.m.TableId("counts"), c.m); c.table == nil {
		fmt.Fprintf(os.Stderr, "table init failed\n")
		os.Exit(1)
	}
}

func (c *CacheHit) Update(s chan<- string){

		hits := float64(0)
		misses := uint64(0)
		total := uint64(0)
		ratio := float64(0)
		c.counter_lookup()
		if c.counters["mark_page_accessed"] > c.counters["mark_buffer_dirty"]{
			total = c.counters["mark_page_accessed"] - c.counters["mark_buffer_dirty"]
		}
		if c.counters["add_to_page_cache_lru"] > c.counters["account_page_dirtied"] {
			misses = c.counters["add_to_page_cache_lru"] - c.counters["account_page_dirtied"]
		}

		hits = float64(total - misses)
		if hits < 0 {
			misses = total
			hits = 0
		}
		if total > 0 {
			ratio = hits / float64(total)
		}

		if c.Debug {
			s <- fmt.Sprintf("%24v%24v%24v", hits, misses, ratio*100)
		} else {
			s <- fmt.Sprintf("%.1f",ratio * 100)
		}

		for k,_ := range c.counters{
			c.counters[k] = 0
		}
		// clear counters
		c.table.DeleteAll()
}

func (c *CacheHit) counter_lookup(){
	it := c.table.Iter()
	for {
		if(!it.Next()){
			break
		}
		counterVal := bpf.GetHostByteOrder().Uint64(it.Leaf())

		// fixme: this ugly hack shouldn't be necessary (bug in the lib??)
		// https://github.com/iovisor/gobpf/issues/273
		symbolAddr := strconv.FormatUint(bpf.GetHostByteOrder().Uint64(it.Key())-1, 16)

		// keeps local cache to lookup symbol name from symbol address
		if _, ok := c.addrCache[symbolAddr]; !ok {
			kname, err := ksym.Ksym(symbolAddr)
			if err != nil {
				panic(err)
			}
			c.addrCache[symbolAddr] = kname
		}
		c.counters[c.addrCache[symbolAddr]] = max(0, counterVal)
	}
}

func (c *CacheHit) Close(){
	c.m.Close()
}

const source string = `
#include <uapi/linux/ptrace.h>
struct key_t {
    u64 ip;
};
BPF_HASH(counts, struct key_t, u64, 4);
int do_count(struct pt_regs *ctx) {
    struct key_t key = {};
    u64 ip;
    key.ip = PT_REGS_IP(ctx);
    counts.increment(key); // update counter
    return 0;
}
`

type key_t struct {
    ip uint64
}

func max(a,b uint64) uint64 {
	if (a < b) {
		 return b
	}
	return a
}
