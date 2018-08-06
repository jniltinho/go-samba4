#!/bin/bash
set -e
DEBIAN_FRONTEND=noninteractive

apt-get update
apt-get -yq install libreadline-dev wget git build-essential libattr1-dev libblkid-dev libpam0g-dev supervisor
apt-get -yq install autoconf python-dev python-dnspython libacl1-dev gdb pkg-config libpopt-dev libldap2-dev
apt-get -yq install dnsutils acl attr libbsd-dev docbook-xsl libcups2-dev libgnutls28-dev ca-certificates nginx

cd /tmp/
wget https://download.samba.org/pub/samba/stable/samba-4.8.3.tar.gz
tar xvfz samba-4.8.3.tar.gz
cd samba-4.8.3/
./configure --with-ads --with-shared-modules=idmap_ad --enable-debug --enable-selftest --with-systemd --prefix=/opt/samba4
make && make install
cd ..

git clone https://github.com/jniltinho/go-samba4.git
mv go-samba4/dist /opt/go-samba4
rm -rf go-samba4
chmod +x /opt/go-samba4/go_samba4

wget https://my-netdata.io/kickstart-static64.sh
/bin/bash kickstart-static64.sh --dont-wait --dont-start-it
cp /opt/netdata/system/netdata-lsb /etc/init.d/netdata
chmod +x /etc/init.d/netdata
cd ..

apt-get clean
rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /var/cache/apt/archive/*.deb

echo '[supervisord] 
nodaemon=true

[program:go_samba4]
directory=/opt/go-samba4
autostart=true
autorestart=true
command=/opt/go-samba4/go_samba4 --server-prod --ssl

[program:nginx]
command=/usr/sbin/nginx -g "daemon off;"
autostart=true
autorestart=true
#user=nobody' > /etc/supervisord.conf

mv /etc/nginx/sites-available/default /etc/nginx/sites-available/default_old_$$