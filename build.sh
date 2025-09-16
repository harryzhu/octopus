CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o dist/macos_arm64/octopus -ldflags "-w -s" main.go
zip dist/octopus_macos_arm64.zip dist/macos_arm64/octopus

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/macos_intel/octopus -ldflags "-w -s" main.go
zip dist/octopus_macos_intel.zip dist/macos_intel/octopus


CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/linux_amd64/octopus -ldflags "-w -s" main.go
zip dist/octopus_linux_amd64.zip dist/linux_amd64/octopus


CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/windows_amd64/octopus.exe -ldflags "-w -s" main.go
zip dist/octopus_windows_amd64.zip dist/windows_amd64/octopus.exe
