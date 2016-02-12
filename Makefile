SRC = $(wildcard elo/ast/*.go)
SRC += $(wildcard elo/parse/*.go)
SRC += $(wildcard elo/pretty/*.go)
SRC += $(wildcard elo/*.go)
OUT = elo

ifdef SystemRoot
	OUT = elo.exe
endif

$(OUT): $(SRC)
	@go build -o $@

clean:
	@rm $(OUT)

.PHONY: clean