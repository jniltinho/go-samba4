#!/bin/bash
set -e
DEBIAN_FRONTEND=noninteractive

apt-get update
apt-get -yq install wget git-core supervisor python-dev
apt-get -yq install python-minimal libpython2.7 libbsd0 libpopt0 libgnutls30 libldap-2.4-2 libcups2
apt-get -yq install ca-certificates nginx python-pip
apt-get -yq install libsasl2-dev libldap2-dev libssl-dev

dpkg -i /tmp/samba-*.amd64.deb

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