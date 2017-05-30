SRC = $(wildcard ast/*.go)
SRC += $(wildcard parse/*.go)
SRC += $(wildcard pretty/*.go)
SRC += $(wildcard run/*.go)
SRC += $(wildcard *.go)
OUT = yo

ifdef SystemRoot
	OUT = yo.exe
endif

$(OUT): $(SRC)
	@go build -o $@ ./run

clean:
	@rm $(OUT)

.PHONY: clean
