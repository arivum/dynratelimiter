
BPF        = dynratelimit.o
BPFCC      = clang
BPFCFLAGS  = -target bpf -Wall -O2 -emit-llvm -g -I $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))/include
BPFLC      = llc
BPFLCFLAGS = -march=bpf -mcpu=probe -filetype=obj

all: $(BPF)

%.o: %.c
	echo $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
	$(BPFCC) $(BPFCFLAGS) -c $< -o - | \
	$(BPFLC) $(BPFLCFLAGS) -o $@

clean:
	rm -f *.o
