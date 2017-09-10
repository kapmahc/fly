# fly
A complete open source e-commerce solution by spring-boot.

## Install java and gradle
```
curl -s "https://get.sdkman.io" | zsh
sdk install java
sdk install gradle
```

## Install nodejs
```
wget -qO- https://raw.githubusercontent.com/creationix/nvm/v0.33.2/install.sh | zsh
nvm install node
nvm alias default node
```

## Create database

```bash
psql -U postgres
CREATE DATABASE db-name WITH ENCODING = 'UTF8';
CREATE USER user-name WITH PASSWORD 'change-me';
GRANT ALL PRIVILEGES ON DATABASE db-name TO user-name;
```


## Build
```
git clone https://github.com/kapmahc/fly.git
cd fly
npm install
gradle build
```

## Notes

- Chrome browser: F12 => Console settings => Log XMLHTTPRequests

- Rabbitmq Management Plugin(<http://localhost:15612>)

  ```bash
  rabbitmq-plugins enable rabbitmq_management
  rabbitmqctl add_user test test
  rabbitmqctl set_user_tags test administrator
  rabbitmqctl set_permissions -p / test ".*" ".*" ".*"
  ```

- "RPC failed; HTTP 301 curl 22 The requested URL returned error: 301"

  ```bash
  git config --global http.https://gopkg.in.followRedirects true
  ```

- 'Peer authentication failed for user', open file "/etc/postgresql/9.5/main/pg_hba.conf" change line:

  ```
  local   all             all                                     peer  
  TO:
  local   all             all                                     md5
  ```

- Generate openssl certs

  ```bash
  openssl genrsa -out www.change-me.com.key 2048
  openssl req -new -x509 -key www.change-me.com.key -out www.change-me.com.crt -days 3650 # Common Name:*.change-me.com
  ```


- [For gmail smtp](http://stackoverflow.com/questions/20337040/gmail-smtp-debug-error-please-log-in-via-your-web-browser)


## Documents
 - [application.properties](https://docs.spring.io/spring-boot/docs/current/reference/html/common-application-properties.html)
 - [nginx](https://www.nginx.com/resources/deployment-guides/load-balance-apache-tomcat/)