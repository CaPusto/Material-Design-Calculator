BINARY=calculator

.PHONY: all build build-linux build-windows run check-deps clean tidy test

all: check-deps build

PKG_MANAGER := $(shell command -v dnf >/dev/null 2>&1 && echo "dnf" || (command -v yum >/dev/null 2>&1 && echo "yum" || echo "apt"))

# Packages for different app managers
DEB_FYNEDEPS := gcc pkg-config libgl1-mesa-dev libegl1-mesa-dev \
    libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev \
    libxkbcommon-dev libwayland-dev
RPM_FYNEDEPS := gcc pkgconfig mesa-libGL-devel mesa-libEGL-devel \
    libXrandr-devel libXcursor-devel libXinerama-devel libXi-devel \
    libxkbcommon-devel wayland-devel

check-deps:
	@echo "Checking system dependencies for Fyne..."
ifeq ($(PKG_MANAGER),dnf)
	@for pkg in $(RPM_FYNEDEPS); do \
		if rpm -q $$pkg >/dev/null 2>&1; then \
			echo "  [OK] $$pkg"; \
		else \
			echo "  [MISSING] $$pkg — install: sudo dnf install $$pkg"; \
		fi \
	done
else ifeq ($(PKG_MANAGER),yum)
	@for pkg in $(RPM_FYNEDEPS); do \
		if rpm -q $$pkg >/dev/null 2>&1; then \
			echo "  [OK] $$pkg"; \
		else \
			echo "  [MISSING] $$pkg — install: sudo yum install $$pkg"; \
		fi \
	done
else
	@for pkg in $(DEB_FYNEDEPS); do \
		if dpkg -s $$pkg >/dev/null 2>&1; then \
			echo "  [OK] $$pkg"; \
		else \
			echo "  [MISSING] $$pkg — install: sudo apt install $$pkg"; \
		fi \
	done
endif
	@command -v go >/dev/null 2>&1 && echo "  [OK] go" || echo "  [MISSING] go"

build: build-linux

build-linux:
	go build -ldflags="-s -w" -o $(BINARY) .

build-windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 \
	CC=x86_64-w64-mingw32-gcc \
	go build -ldflags="-s -w -H=windowsgui" -o $(BINARY).exe .

run: build-linux
	./$(BINARY)

clean:
	rm -f $(BINARY) $(BINARY).exe

tidy:
	go mod tidy

test:
	go test ./...
