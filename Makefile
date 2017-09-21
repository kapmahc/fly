dist=dist
pkg=github.com/kapmahc/fly/plugins/nut/app
theme=bootstrap

VERSION=`git rev-parse --short HEAD`
BUILD_TIME=`date -R`
AUTHOR_NAME=`git config --get user.name`
AUTHOR_EMAIL=`git config --get user.email`
COPYRIGHT=`head -n 1 LICENSE`
USAGE=`sed -n '3p' README.md`


build: backend frontend
	tar jcf $(dist).tar.bz2 $(dist)


frontend:
	cd dashboard && npm run build
	-cp -r dashboard/dist $(dist)/dashboard



backend:
	go build -ldflags "-s -w -X ${pkg}.Version=${VERSION} -X '${pkg}.BuildTime=${BUILD_TIME}' -X '${pkg}.AuthorName=${AUTHOR_NAME}' -X ${pkg}.AuthorEmail=${AUTHOR_EMAIL} -X '${pkg}.Copyright=${COPYRIGHT}' -X '${pkg}.Usage=${USAGE}'" -o ${dist}/fly main.go
	-cp -r locales templates $(dist)/
	mkdir -p $(dist)/themes/$(theme)
	-cp -r themes/$(theme)/{assets,views,package.json,package-lock.json} $(dist)/themes/$(theme)


init:
	govendor sync
	npm install
	cd dashboard && npm install

clean:
	-rm -r $(dist) $(dist).tar.bz2 dashboard/dist
