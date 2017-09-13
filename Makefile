dst=fly

define ANNOUNCE_BODY
package nut

const (
	// Version version
	Version = "$(shell git rev-parse --short HEAD)"
	// BuildTime build time
	BuildTime = "$(shell date -R)"
	// Usage usage
	Usage = "$(shell sed -n '3p' README.md)"
	// Copyright copyright
	Copyright = "$(shell head -n 1 LICENSE)"
	// AuthorName author's name
	AuthorName = "$(shell git config --get user.name)"
	// AuthorEmail author's email
	AuthorEmail = "$(shell git config --get user.email)"
)

endef


build: frontend backend
	mkdir dist
	tar -zxf fly.tar.gz -C dist
	-cp -r dashboard/build dist/dashboard
	mkdir -p dist/tmp
	cd dist && tar -jcf ../dist.tar.bz2 *


frontend:
	cd dashboard && npm run build


backend: ANNOUNCE.txt
	bee pack -v -ba="-ldflags '-s'" -exp=tmp:dist:node_modules:dashboard:.git -exs=.un~:.swp:.tmp:.go:Makefile



ANNOUNCE.txt: export ANNOUNCE_BODY:=$(ANNOUNCE_BODY)
ANNOUNCE.txt:
	echo "$${ANNOUNCE_BODY}" > plugins/nut/constants.go




clean:
	-rm -rv $(dst) $(dst).tar.gz lastupdate.tmp routers/commentsRouter_*.go dashboard/build dist*
