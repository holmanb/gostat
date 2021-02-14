package main

// based on bcc/cachestat.py
import (
	"os"
	"fmt"
	"time"
	"strconv"
	bpf "github.com/iovisor/gobpf/bcc"
	ksym "github.com/iovisor/gobpf/pkg/ksym"
)

import "C"

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

func get_counters(table *bpf.Table,kmap *map[string]uint64) {
	it := table.Iter()
	for ;it.Next(); {
		kInt := bpf.GetHostByteOrder().Uint64(it.Leaf())

		// fixme: this ugly hack shouldn't be necessary (bug in the lib??)
		// https://github.com/iovisor/gobpf/issues/273
		strval := strconv.FormatUint(bpf.GetHostByteOrder().Uint64(it.Key())-1, 16)
		kname, err := ksym.Ksym(strval)
		if err != nil {
			panic(err)
		}
		(*kmap)[kname] = max(0,kInt)
	}
}

func main(){
	debug := false
	kmap := map[string]uint64{
		"add_to_page_cache_lru":0,
		"mark_buffer_dirty":0,
		"account_page_dirtied":0,
		"mark_page_accessed":0,
	};
	m := bpf.NewModule(source, []string{})
	defer m.Close()

	do_count_probe, err := m.LoadKprobe("do_count")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load do_count : %s\n", err)
		os.Exit(1)
	}
	for k,_ := range kmap {

		err = m.AttachKprobe(k, do_count_probe, -1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to attach %s: %s\n",k, err)
			os.Exit(1)
		}
	}

	fmt.Printf("%24v%24v%24v\n", "HITS", "MISSES", "RATIO")
	if debug {
		fmt.Printf("%24v%24v%24v%24v%24v%24v%24v%24v\n",
			"mpa",
			"mbd",
			"apcl",
			"apd",
			"total",
			"misses",
			"hits",
			"ratio")
	}
	table := bpf.NewTable(m.TableId("counts"), m)

	for {
		hits := float64(0)
		misses := uint64(0)
		total := uint64(0)
		ratio := float64(0)
		get_counters(table, &kmap)
		if kmap["mark_page_accessed"] > kmap["mark_buffer_dirty"]{
			total = kmap["mark_page_accessed"] - kmap["mark_buffer_dirty"]
		}
		if kmap["add_to_page_cache_lru"] > kmap["account_page_dirtied"] {
			misses = kmap["add_to_page_cache_lru"] - kmap["account_page_dirtied"]
		}

		hits = float64(total - misses)
		if hits < 0 {
			misses = total
			hits = 0
		}
		if total > 0 {
			ratio = hits / float64(total)
		}

		if debug {
			fmt.Printf("%24v%24v%24v%24v%24v%24v%24v%24v\n",
				kmap["mark_page_accessed"],
				kmap["mark_buffer_dirty"],
				kmap["add_to_page_cache_lru"],
				kmap["account_page_dirtied"],
				total, misses, hits,ratio*100)
		} else {
			fmt.Printf("%24v%24v%24v\n", hits, misses, fmt.Sprintf("%.1f",ratio * 100))
		}

		for k,_ := range kmap {
			kmap[k] = 0
		}
		// clear counters
		table.DeleteAll()

		time.Sleep(time.Second)
	}
}
