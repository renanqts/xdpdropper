// +build ignore

#include "bpf_endian.h"
#include "common.h"

char __license[] SEC("license") = "GPL";

#define MAX_MAP_ENTRIES 1000 // It allows 1k entries

/* Define a hash map for storing packet by source IPv4 address */
struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, MAX_MAP_ENTRIES);
	__type(key, __u32);   // source IPv4 address
	__type(value, __u32); // packet count
} drop_map SEC(".maps");

/*
Attempt to parse the IPv4 source address from the packet.
Returns 0 if there is no IPv4 header field; otherwise returns non-zero.
We care only about IPv4 for the time being.
*/
static __always_inline int parse_ip_src_addr(struct xdp_md *ctx, __u32 *ip_src_addr) {
	void *data_end = (void *)(long)ctx->data_end;
	void *data     = (void *)(long)ctx->data;

	// First, parse the ethernet header.
	struct ethhdr *eth = data;
	if ((void *)(eth + 1) > data_end) {
		// use to debug
		// bpf_printk("Failed to parse the ethernet header");
		return 0;
	}

	// The protocol is not IPv4, so we can't parse an IPv4 source address.
	if (eth->h_proto != bpf_htons(ETH_P_IP)) {
		// use to debug
		// bpf_printk("The protocol is not IPv4");	
		return 0;
	}

	// Then parse the IP header.
	struct iphdr *ip = (void *)(eth + 1);
	if ((void *)(ip + 1) > data_end) {
		// use to debug
		// bpf_printk("Failed to parse the IP header");
		return 0;
	}

	// Return the source IP address in network byte order.
	*ip_src_addr = (__u32)(ip->saddr);
	return 1;
}

SEC("xdp")
int xdp_drop_func(struct xdp_md *ctx) {
	__u32 ip;
	if (!parse_ip_src_addr(ctx, &ip)) {
		// Not an IPv4 packet, so do nothing.
		goto done;
	}

	__u32 *pkt_count = bpf_map_lookup_elem(&drop_map, &ip);
	if (pkt_count) {
		// IP was found blacklisted.
		// so increment the counters atomically using an LLVM built-in
		__sync_fetch_and_add(pkt_count, 1);
		// Drop it
		return XDP_DROP;
	}

done:
	return XDP_PASS;
}
