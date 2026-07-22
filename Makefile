VERSION ?= 0.1.0
NAME ?= lisdb

INSTDIR ?= /usr/bin
OUTDIR ?= _build

GO ?= go

UNAME_ARCH ?= $(shell uname -m)

ifeq ($(UNAME_ARCH),x86_64)
    ARCH ?= amd64
else ifeq ($(UNAME_ARCH),aarch64)
    ARCH ?= arm64
else ifeq ($(UNAME_ARCH),i686)
    ARCH ?= i386
else
    ARCH ?= $(UNAME_ARCH)
endif

export OUTDIR NAME VERSION UNAME_ARCH ARCH

all: $(NAME) deb

$(NAME): test
	@mkdir -vp $(OUTDIR)
	$(GO) build -o $(OUTDIR)/$(NAME)

deb: $(NAME)
	scripts/make-deb.sh $(OUTDIR) $(NAME)

test:
	@$(GO) test -C llis
	@echo "\n\033[38;2;0;223;0m\033[1mPASS\033[0m"

install:
	install -m 0755 $(OUTDIR)/$(NAME) -t $(PREFIX)/$(INSTDIR)
	$(PREFIX)/$(INSTDIR)/$(NAME) --setup

clean:
	rm -rf .generated $(OUTDIR)
