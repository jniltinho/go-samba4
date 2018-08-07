#!/bin/bash
set -e
DEBIAN_FRONTEND=noninteractive

apt-get update
apt-get -yq install libreadline-dev wget git-core build-essential libattr1-dev libblkid-dev libpam0g-dev supervisor
apt-get -yq install autoconf python-dev python-dnspython libacl1-dev gdb pkg-config libpopt-dev libldap2-dev
apt-get -yq install dnsutils acl attr libbsd-dev docbook-xsl libcups2-dev libgnutls28-dev ca-certificates nginx
apt-get -yq install python-pip dos2unix libsasl2-dev libldap2-dev libssl-dev

cd /tmp/
get_samba4=https://download.samba.org/pub/samba/rc/samba-4.9.0rc2.tar.gz
wget ${get_samba4}
tar xvfz $(basename ${get_samba4})
cd $(basename ${get_samba4}|sed "s/.tar.gz//")
./configure --with-ads --with-shared-modules=idmap_ad --with-systemd --prefix=/opt/samba4
make && make install

cd /tmp/
git clone https://github.com/jniltinho/go-samba4.git
cd go-samba4
rm -rf dist/*
pip install -r requirements.txt
python make_bin.py go_samba4.py
mv /tmp/go-samba4/dist /opt/go-samba4
chmod +x /opt/go-samba4/go_samba4

cd /tmp/
wget https://my-netdata.io/kickstart-static64.sh
/bin/bash kickstart-static64.sh --dont-wait --dont-start-it
cp /opt/netdata/system/netdata-lsb /etc/init.d/netdata
chmod +x /etc/init.d/netdata

cd /
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