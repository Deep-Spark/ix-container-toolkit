all:
	mkdir -p build
	go build -o build/ix-container-runtime cmd/ix-container-runtime/main.go
	go build -o build/ix-ctk cmd/ix-ctk/main.go

install:
	mkdir -p /var/log/iluvatarcorex/ix-container-toolkit/
	install -Dm755 build/ix-container-runtime /usr/local/bin/ix-container-runtime
	install -Dm755 build/ix-ctk /usr/local/bin/ix-ctk

uninstall:
	rm -rf /var/log/iluvatarcorex/ix-container-toolkit/
	rm -f /usr/local/bin/ix-container-runtime
	rm -f /usr/local/bin/ix-ctk

clean:
	rm -rf build
