#!/bin/bash

# exec custom command
if [[ $# -gt 0 ]] ; then
        exec "$@"
        exit
fi

### Create Domain
/opt/samba4/bin/samba-tool domain provision --server-role=dc --use-rfc2307 \
 --dns-backend=SAMBA_INTERNAL --realm=LINUXPRO.NET --domain=LINUXPRO \
 --adminpass=Linuxpro123456

/etc/init.d/netdata start

exec /usr/bin/supervisord -c /etc/supervisord.conf