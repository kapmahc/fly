# FLY

A complete open source e-commerce solution by The Go Programming Language.

## Install nodejs

```bash
curl -o- https://raw.githubusercontent.com/creationix/nvm/v0.33.2/install.sh | zsh
nvm install node
nvm alias default node
```

## Install go

```bash
zsh < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
gvm install go1.9 -B
gvm use go1.9 --default
```

## Atom plugins

- go-plus
- git-plus
- file-icons

## Notes

- Create database

```
psql -U postgres
CREATE DATABASE db-name WITH ENCODING = 'UTF8';
CREATE USER user-name WITH PASSWORD 'change-me';
GRANT ALL PRIVILEGES ON DATABASE db-name TO user-name;
```

- Database migrations.

```
go get -u -d github.com/mattes/migrate/cli github.com/lib/pq
go build -ldflags '-s -w' -tags 'postgres' -o $GOPATH/bin/migrate github.com/mattes/migrate/cli
migrate -path db/migrate -database postgres://postgres@localhost:5432/fly?sslmode=disable up
```

- Chrome browser: F12 => Console settings => Log XMLHTTPRequests

- Rabbitmq Management Plugin(<http://localhost:15612>)

  ```
  rabbitmq-plugins enable rabbitmq_management
  rabbitmqctl change_password guest change-me
  rabbitmqctl add_user who-am-i change-me
  rabbitmqctl set_user_tags who-am-i administrator
  rabbitmqctl list_vhosts
  rabbitmqctl add_vhost /v-host
  rabbitmqctl set_permissions -p /v-host who-am-i ".*" ".*" ".*"
  ```

- "RPC failed; HTTP 301 curl 22 The requested URL returned error: 301"

  ```
  git config --global http.https://gopkg.in.followRedirects true
  ```

- 'Peer authentication failed for user', open file "/etc/postgresql/9.5/main/pg_hba.conf" change line:

  ```
  local   all             all                                     peer  
  TO:
  local   all             all                                     md5
  ```

- Generate openssl certs

  ```
  openssl genrsa -out www.change-me.com.key 2048
  openssl req -new -x509 -key www.change-me.com.key -out www.change-me.com.crt -days 3650 # Common Name:*.change-me.com
  ```

- [For gmail smtp](http://stackoverflow.com/questions/20337040/gmail-smtp-debug-error-please-log-in-via-your-web-browser)
