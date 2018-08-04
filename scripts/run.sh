#!/bin/bash

### Create Domain
/opt/samba4/bin/samba-tool domain provision --server-role=dc --use-rfc2307 \
 --dns-backend=SAMBA_INTERNAL --realm=LINUXPRO.NET --domain=LINUXPRO \
 --adminpass=Linuxpro123456

/opt/samba4/sbin/samba -D
/etc/init.d/netdata start

#cd /opt/go-samba4/
#exec ./go_samba4 --server-prod --ssl
