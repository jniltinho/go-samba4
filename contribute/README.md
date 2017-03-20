# go-samba4
A web interface to manage a remote samba4 server for the Go programming language.




```bash

mkdir /opt/wide
groupadd wide && useradd -M -s /bin/bash -g wide -d /opt/wide wide


cp /opt/wide/scripts/wide.service /etc/systemd/system/wide.service
chmod +x /opt/wide/scripts/start-wide.sh

```

