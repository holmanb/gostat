package main
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

func main(){
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
	table := bpf.NewTable(m.TableId("counts"), m)

	fmt.Printf("%24v","add_to_page_cache_lru")
	fmt.Printf("%24v","mark_buffer_dirty")
	fmt.Printf("%24v","account_page_dirtied")
	fmt.Printf("%24v\n","mark_page_accessed")
	init := true
	for {
		it := table.Iter()
		for ;it.Next(); {
		        kInt := bpf.GetHostByteOrder().Uint64(it.Leaf())
			if err != nil {
				panic(err)
			}
		        strval := strconv.FormatUint(bpf.GetHostByteOrder().Uint64(it.Key())-1, 16)
			kname, err := ksym.Ksym(strval)
			if err != nil {
				panic(err)
			}
			kmap[kname] = kInt
		}
		// don't print unchanged values
		if ! init {
			fmt.Printf("%24v",kmap["add_to_page_cache_lru"])
			fmt.Printf("%24v",kmap["mark_buffer_dirty"])
			fmt.Printf("%24v",kmap["account_page_dirtied"])
			fmt.Printf("%24v\n",kmap["mark_page_accessed"])
		}
		init = false
		time.Sleep(time.Second)
	}
}
