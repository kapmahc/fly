#!/bin/sh
if [ ! -f "/tmp/fly.tar.gz" ]; then
  echo "file /tmp/fly.tar.gz not exist."
  exit -1
fi

{{range .shared_directories}}
if [ ! -d "/{{$.root}}/shared/{{.}}" ]; then
  mkdir -pv /{{$.root}}/shared/{{.}}
fi
{{end}}

{{range .shared_files}}
if [ ! -f "/{{$.root}}/shared/{{.}}" ]; then
  echo "file /{{$.root}}/shared/{{.}} not exist."
  exit -1
fi
{{end}}

VERSION=`date +%Y%m%d%H%M%S`
mkdir -pv /{{$.root}}/${VERSION}
cd /{{$.root}}/${VERSION}
tar xvf /tmp/fly.tar.gz
{{range .shared_directories}}
ln -sv /{{$.root}}/shared/{{.}} /{{$.root}}/${VERSION}/{{.}}
{{end}}
{{range .shared_files}}
rm -fv /{{$.root}}/${VERSION}/{{.}}
ln -sv /{{$.root}}/shared/{{.}} /{{$.root}}/${VERSION}/{{.}}
{{end}}

npm install
fuser -k {{.port}}/tcp
nohup ./fly > /dev/null 2>&1 &
rm -f tmp/nginx.conf && ./fly g ng
sudo nginx -s reload

echo 'DONE!'
