### 2. **Makefile**

CC = gcc
CFLAGS = -fPIC
LIBNAME = particles.so
PYTHON = python3

all: $(LIBNAME)

$(LIBNAME): particles.c particles.h
	$(CC) -shared -o $(LIBNAME) $(CFLAGS) particles.c

run: $(LIBNAME)
	$(PYTHON) particles_app.py

clean:
	rm -f $(LIBNAME)
