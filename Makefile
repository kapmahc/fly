dst=fly

build:
	go build -ldflags '-s -w' -tags 'postgres' -o migrate github.com/mattes/migrate/cli
	bee pack -v -ba="-ldflags '-s'" -exp=tmp:node_modules:.git -exs=.un~:.swp:.tmp:.go:Makefile

clean:
	-rm migrate $(dst) $(dst).tar.gz
