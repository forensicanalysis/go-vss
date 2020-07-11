.PHONY: test

LIBVERSION=libvshadow-20191221

CC=/usr/bin/x86_64-w64-mingw32-gcc
CXX=/usr/bin/x86_64-w64-mingw32-g++
AR=/usr/bin/x86_64-w64-mingw32-ar
OBJDUMP=/usr/bin/x86_64-w64-mingw32-objdump
RANLIB=/usr/bin/x86_64-w64-mingw32-ranlib
STRIP=/usr/bin/x86_64-w64-mingw32-strip
MINGWFLAGS=
CFLAGS=$(MINGWFLAGS)
CXXFLAGS=$(MINGWFLAGS)
LDFLAGS=$(LDFLAGS):/usr/include/wine-development/wine/windows

prepare_windows:
	sudo apt-get install -y mingw-w64
	tar -xf libvshadow-alpha-20191221.tar.gz

	cd $(LIBVERSION); CC=$(CC) CXX=$(CXX) AR=$(AR) OBJDUMP=$(OBJDUMP) RANLIB=$(RANLIB) STRIP=$(STRIP) ./configure --host=x64-mingw32msvc --prefix=/usr/bin/x86_64-w64-mingw32 --enable-winapi=yes --enable-multi-threading-support=no
	cd $(LIBVERSION); CC=$(CC) CXX=$(CXX) AR=$(AR) OBJDUMP=$(OBJDUMP) RANLIB=$(RANLIB) STRIP=$(STRIP) CFLAGS=$(CFLAGS) CXXFLAGS=$(CXXFLAGS) make

prepare_linux:
	tar -xf libvshadow-alpha-20191221.tar.gz
	cd $(LIBVERSION); ./configure --enable-multi-threading-support=no #--enable-shared=no
	cd $(LIBVERSION); make

prepare_macos: prepare_linux

build:
	go build $(FLAGS) .

test:
	go test -count=1 $(FLAGS) .

coverage:
	go-acc ./... -- $(FLAGS)

clean:
	rm -r $(LIBVERSION)

