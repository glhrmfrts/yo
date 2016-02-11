OUT = elo

ifdef SystemRoot
	OUT = elo.exe
endif

$(OUT):
	@go build -o $@

clean:
	@rm $(OUT)

.PHONY: clean