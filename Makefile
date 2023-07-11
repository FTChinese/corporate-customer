config_file_name := api.toml
local_config_file := $(HOME)/config/$(config_file_name)

app_name := ftacademy
go_version := go1.16

current_dir := $(shell pwd)
sys := $(shell uname -s)
hardware := $(shell uname -m)
src_dir := $(current_dir)
out_dir := $(current_dir)/out
build_dir := $(current_dir)/build

default_exec := $(out_dir)/$(sys)/$(hardware)/$(app_name)

linux_x86_exec := $(out_dir)/linux/x86/$(app_name)

linux_arm_exec := $(out_dir)/linux/arm/$(app_name)

server_dir := /data/node/go/bin

.PHONY: build
build :
	go build -o $(default_exec) -tags production -v $(src_dir)

.PHONY: run
run :
	$(default_exec) -production=false -livemode=false

.PHONY: builddir
builddir :
	mkdir -p $(build_dir)

.PHONY: devenv
devenv : builddir
	rsync $(HOME)/config/env.dev.toml $(build_dir)/$(config_file_name)

.PHONY: version
version :
	git describe --tags > build/version
	date +%FT%T%z > build/build_time

# Use production db, sandbox api.
# For online production. use -production=true -livemode=true
.PHONY: sandbox
sandbox :
	$(default_exec) -production=true -livemode=false

.PHONY: amd64
amd64 :
	@echo "Build production linux version $(version)"
	GOOS=linux GOARCH=amd64 go build -o $(linux_x86_exec) -tags production -v $(src_dir)

.PHONY: arm
arm :
	@echo "Build production arm version $(version)"
	GOOS=linux GOARM=7 GOARCH=arm go build -o $(linux_arm_exec) -tags production -v $(src_dir)

.PHONY: install-go
install-go:
#	@echo "Install go version $(go_version)"
#	gvm install $(go_version)
	/data/opt/server/jenkins/jenkins/.gvm/bin/gvm use $(go_version)

.PHONY: config
config : outdir
	# Download configuration file
	rsync -v node@tk11:/home/node/config/$(config_file_name) $(build_dir)/$(config_file_name)
	ls ./$(build_dir)

.PHONY: publish
publish :
	# Remove the .bak file
	ssh ucloud "rm -f $(server_dir)/$(app_name).bak"
	# Sync binary to the xxx.bak file
	rsync -v ./$(default_exec) ucloud:$(server_dir)/$(app_name).bak

.PHONY: restart
restart :
	# Rename xxx.bak to app name
	ssh ucloud "cd $(server_dir)/ && \mv $(app_name).bak $(app_name)"
	ssh ucloud supervisorctl restart $(app_name)

.PHONY: clean
clean :
	go clean -x
	rm -r build/*


.PHONY: outdir
outdir :
	mkdir -p ./$(build_dir)