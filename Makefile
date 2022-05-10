# List all supported dist, run below command
#  go tool dist list

BINDIR="bin"
BINARY_NAME=goracle

all: clean dep windows_amd64 linux_amd64

windows_amd64:
	@echo
	@echo "** Build windows/amd64 binary **"
	@echo
	CGO_ENABLED=1 \
	GOOS=windows \
	GOARCH=amd64 \
	CC="zig cc -target x86_64-windows-gnu" \
	CXX="zig c++ -target x86_64-windows-gnu" \
	go build -o ${BINDIR}/goracle-windows-amd64 github.com/JohnWongCHN/goracle

linux_amd64:
	@echo
	@echo "** Build linux/amd64 binary **"
	@echo
	CGO_ENABLED=1 \
	GOOS=linux \
	GOARCH=amd64 \
	CC="zig cc -target x86_64-linux-gnu" \
	CXX="zig c++ -target x86_64-linux-gnu" \
	go build -o ${BINDIR}/goracle-linux-amd64 github.com/JohnWongCHN/goracle

.PHONY : clean

# clean files and code
# with minux ignore errors and continue
clean:
	@echo
	@echo "** Clean files and code **"
	@echo
	go clean
	@echo
	rm -f ${BINDIR}/*

# install dependency packages
dep:
	@echo
	@echo "** Installing dependencies **"
	@echo
	go mod download
	@echo
	go mod tidy

vet:
	go vet

test:
	go test ./...

test_verbose:
	go test ./... -v

test_coverage:
	go test ./... -coverprofile=coverage.out