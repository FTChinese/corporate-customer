config_file_name := api.toml
local_config_file := $(HOME)/config/$(config_file_name)

version := `git describe --tags`
build_time := `date +%FT%T%z`
commit := `git log --max-count=1 --pretty=format:%aI_%h`

ldflags := -ldflags "-w -s -X main.version=${version} -X main.build=${build_time} -X main.commit=${commit}"

app_name := ftacademy
go_version := go1.16

sys := $(shell uname -s)
hardware := $(shell uname -m)
build_dir := build
src_dir := .

default_exec := $(build_dir)/$(sys)/$(hardware)/$(app_name)
compile_default_exec := go build -o $(default_exec) $(ldflags) -tags production -v $(src_dir)

linux_x86_exec := $(build_dir)/linux/x86/$(app_name)
compile_linux_x86 := GOOS=linux GOARCH=amd64 go build -o $(linux_x86_exec) $(ldflags) -tags production -v $(src_dir)

linux_arm_exec := $(build_dir)/linux/arm/$(app_name)
compile_linux_arm := GOOS=linux GOARM=7 GOARCH=arm go build -o $(linux_arm_exec) $(ldflags) -tags production -v $(src_dir)

.PHONY: build
build :
	@echo "Build version $(version)"
	$(compile_default_exec)

.PHONY: run
run :
	$(default_exec)

.PHONY: amd64
amd64 :
	@echo "Build production linux version $(version)"
	$(compile_linux_x86)

.PHONY: arm
arm :
	@echo "Build production arm version $(version)"
	$(compile_linux_arm)

.PHONY: install-go
install-go:
	@echo "Install go version $(go_version)"
	gvm install $(go_version)
	gvm use $(go_version)

.PHONY: config
config :
	rsync -v tk11:/home/node/config/$(config_file_name) ./$(build_dir)

publish :
	ssh ucloud "rm -f /home/node/go/bin/$(app_name).bak"
	rsync -v ./$(default_exec) ucloud:/home/node/go/bin/$(app_name).bak

restart :
	ssh ucloud "cd /home/node/go/bin/ && \mv $(app_name).bak $(app_name)"
	ssh ucloud supervisorctl restart $(app_name)

clean :
	go clean -x
	rm build/*

