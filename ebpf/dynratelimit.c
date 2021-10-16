/*
 * Copyright (c) 2021, arivum.
 * All rights reserved.
 * SPDX-License-Identifier: MIT
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/MIT
 */

#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/in.h>
#include <linux/udp.h>
#include <linux/tcp.h>
#include <sys/resource.h>

#define DYNRATELIMIT_MAP_SIZE 2
#define CONNECTION_MAP_SIZE 1000

#ifndef htons
# define htons(X)		__constant_htons((X))
#endif


enum connection_state {
	conn_initiated = (__u8) 1,
	conn_established = (__u8) 2,
	conn_closed = (__u8) 4,
};


struct bpf_map_def SEC("maps") dynratelimit_map = {
	.type = BPF_MAP_TYPE_ARRAY,
	.key_size = sizeof(__u32),
	.value_size = sizeof(__u32),
	.max_entries = DYNRATELIMIT_MAP_SIZE,
	.map_flags       = 0,
};

struct bpf_map_def SEC("maps") connection_map = {
	.type = BPF_MAP_TYPE_LRU_HASH,
	.key_size = sizeof(__u64),
	.value_size = sizeof(enum connection_state),
	.max_entries = CONNECTION_MAP_SIZE,
	.map_flags       = 0,
};

enum _index {
	limit,
	accepted_requests,
};

__u32 get_ratelimit() {
	enum _index map_index = limit;
	int *entry;

	entry = bpf_map_lookup_elem(&dynratelimit_map, &map_index);
	if (!entry) {
		return -1;
	}
	return *entry;
}

SEC("xdp")
int xdp_dynratelimit(struct xdp_md *ctx)
{
	enum _index map_index = accepted_requests;
	int *acc_requests = bpf_map_lookup_elem(&dynratelimit_map, &map_index);
	if (!acc_requests) {
		return XDP_DROP;
	}


	__u64 nh_off = 0;
	void* data_end = (void*)(long)ctx->data_end;
    void* data = (void*)(long)ctx->data;

    // Handle data as an ethernet frame header
    // Check frame header size
    struct ethhdr *eth = data;
    nh_off = sizeof(*eth);
    if (data + nh_off > data_end) {
        return XDP_PASS;
    }

    // Check protocol
    if (eth->h_proto != htons(ETH_P_IP)) {
        return XDP_PASS;
    }

    // Check ip packet header size
 	struct iphdr *iph = data + nh_off;
    nh_off += sizeof(struct iphdr);
    if (data + nh_off > data_end) {
        return XDP_DROP;
    }

    // Check protocol
    if (iph->protocol != IPPROTO_TCP) {
        return XDP_PASS;
    }

    // Check tcp header size
	struct tcphdr *tcph = data + nh_off;
    nh_off += sizeof(struct tcphdr);
    if (data + nh_off > data_end) {
        return XDP_PASS;
    }


	// syn=1 ack=0 -> new entry in map with saddr
	// syn=0 ack=1 -> update entry in map with saddr (set timestamp)
	// fin=1 -> delete entry from map with saddr

	__u64 key = (((__u64)iph->saddr) << 32) | (__u64) tcph->source;

	if (tcph->syn == 1 && tcph->ack == 0 && tcph->fin == 0) {
		if (*acc_requests >= get_ratelimit()) {
			return XDP_DROP;
		}
		enum connection_state new_entry = conn_initiated;

		if (bpf_map_update_elem(&connection_map, &key, &new_entry, BPF_ANY) > 0) {
			return XDP_DROP;
		} else {
			__sync_fetch_and_add(acc_requests, 1);
		}
	} else if (tcph->fin == 1) {
		bpf_map_delete_elem(&connection_map, &key);
	} else if (tcph->syn == 0 && tcph->ack == 1 && tcph->fin == 0) {
		enum connection_state *entry = bpf_map_lookup_elem(&connection_map, &key);
		if (!entry) {
			return XDP_DROP;
		} else if (*entry != conn_established) {
			*entry = conn_established;
		}
	} else {
		if (!bpf_map_lookup_elem(&connection_map, &key)) {
			return XDP_DROP;
		}
	}

	return XDP_PASS;
}


char _license[] SEC("license") = "GPL";


