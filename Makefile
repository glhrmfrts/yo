SRC = $(wildcard ast/*.go)
SRC += $(wildcard parse/*.go)
SRC += $(wildcard pretty/*.go)
SRC += $(wildcard run/*.go)
SRC += $(wildcard *.go)
OUT = went

ifdef SystemRoot
	OUT = went.exe
endif

$(OUT): $(SRC)
	@go build -o $@ ./run

clean:
	@rm $(OUT)

.PHONY: clean