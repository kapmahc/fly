dst=fly

build:
	bee pack -v -ba="-ldflags '-s'" -exp=tmp:node_modules:.git -exs=.un~:.swp:.go:Makefile


clean:
	-rm -r $(dst) $(dst).tar.gz


