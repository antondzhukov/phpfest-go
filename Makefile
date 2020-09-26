kSOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=phpfestgo

.DEFAULT_GOAL: $(BINARY)

$(BINARY): protofiles $(SOURCES)
	go build -o ${BINARY}

.PHONY: install
install: $(SOURCES)
	go install

.PHONY: protofiles
protofiles: ./phpfestproto/*.proto
	rm -f ./phpfestproto/*.go && protoc --go_out=:. --go-grpc_out=:. ./phpfestproto/*.proto

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

