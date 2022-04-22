# List all supported dist, run below command
#  go tool dist list

BIN="bin"
BINARY_NAME=goracle
PLATFORMS = windows/amd64 \
			linux/amd64 \

# loop in recipe is simpliy as shell script itself
# use backslash to split long lines
build: clean dep
	@echo
	@echo "** Build binary **"
	@echo
	for platform in ${PLATFORMS}; \
	do \
	  IFS='/' read -r -a array <<< "$${platform}"; \
	  GOOS=$${array[0]}; \
	  GOARCH=$${array[1]}; \
	  OUTPUT_NAME="${BIN}/${BINARY_NAME}-$${GOOS}-$${GOARCH}"; \
	  if [ $${GOOS} = "windows" ]; then \
	    OUTPUT_NAME="$${OUTPUT_NAME}.exe"; \
	  fi; \
	  GOARCH=$${GOARCH} GOOS=$${GOOS} go build -o $${OUTPUT_NAME} github.com/JohnWongCHN/goracle; \
	done

.PHONY : clean

# clean files and code
# with minux ignore errors and continue
clean:
	@echo
	@echo "** Clean files and code **"
	@echo
	go clean
	@echo
	-for platform in ${PLATFORMS}; \
	do \
	  filename=${BIN}/${BINARY_NAME}-$$(echo $$platform | sed -r 's/\//-/g'); \
	  if [[ $${filename} == *"windows"* ]]; then \
	    filename="$${filename}.exe"; \
	  fi; \
	  rm $${filename}; \
	done

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