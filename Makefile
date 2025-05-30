
# GOOS=darwin GOARCH=amd64

appname=MYAPP
appnamelong=This Is My App

# appname in lower case
buildname=$(shell echo $(appname) | tr '[:upper:]' '[:lower:]')
arch=$(shell go env GOARCH)
os=$(shell go env GOOS)

build: clean

	# Build for local platform
	CGO_ENABLED=1 go build -o bin/$(buildname)_$(os).$(arch) -ldflags="-s -w -X 'app/app_conf.AppName=$(appname)' -X 'app/app_conf.AppNameLong=$(appnamelong)'" cmd/main.go

run:
	go run cmd/main.go

runapp:
	./bin/$(buildname)_$(os).$(arch)

clean:
 	#if file exists, delete Inventory
	@if [ -f bin/$(buildname)_$(os).$(arch) ]; then \
		rm bin/$(buildname)_$(os).$(arch); \
	fi
	
