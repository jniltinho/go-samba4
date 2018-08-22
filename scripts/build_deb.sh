#!/bin/bash
## Create RPM package Samba4 4.8.4 (Debian 9)
## http://sig9.hatenablog.com/entry/2017/12/04/000000
## Run as root (sudo su)

apt-get update
apt-get -yq install ruby-dev
apt-get -yq install libreadline-dev git build-essential libattr1-dev libblkid-dev libpam0g-dev
apt-get -yq install autoconf python-dev python-dnspython libacl1-dev gdb pkg-config libpopt-dev libldap2-dev 
apt-get -yq install dnsutils acl attr libbsd-dev docbook-xsl libcups2-dev libgnutls28-dev curl

gem install fpm


mkdir -p /build && cd /build

curl -L -O http://www.samba.org/samba/ftp/samba-latest.tar.gz
tar zxvf samba-latest.tar.gz
cd samba-4.8*
./configure --with-ads --with-shared-modules=idmap_ad --enable-debug --enable-selftest --with-systemd --prefix=/opt/samba4
make -j 2
make install install DESTDIR=/tmp/installdir

mkdir -p /tmp/installdir/etc/systemd/system

echo '[Unit]
Description=Samba4 AD Daemon
After=syslog.target network.target
 
[Service]
Type=forking
PIDFile=/opt/samba4/var/run/samba.pid
LimitNOFILE=16384
EnvironmentFile=-/etc/sysconfig/samba4
ExecStart=/opt/samba4/sbin/samba $SAMBAOPTIONS
ExecReload=/usr/bin/kill -HUP $MAINPID
 
[Install]
WantedBy=multi-user.target' > /tmp/installdir/etc/systemd/system/samba4.service

fpm -s dir -t deb -n samba -v 4.8.4 -C /tmp/installdir \
  -d "python-minimal" \
  -d "libpython2.7" \
  -d "libbsd0" \
  -d "libpopt0" \
  -d "libgnutls30" \
  -d "libldap-2.4-2" \
  -d "libcups2" \
  -p samba-4.8.4+dfsg-1.amd64.deb .

mv samba-4.8.4+dfsg-1.amd64.deb /root/

cd /
apt-get clean
rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /var/cache/apt/archive/*.deb /build


## Install DEB
## apt-get update
## apt-get -yq install python-minimal libpython2.7 libbsd0 libpopt0 libgnutls30 libldap-2.4-2 libcups2
## dpkg -i /root/samba-4.8.4+dfsg-1.amd64.deb

### Add PATH
# echo 'export PATH=$PATH:/opt/samba4/bin:/opt/samba4/sbin' >> /etc/profile
# source /etc/profile


### Create Domain Samba4 like AD
# hostnamectl set-hostname samba4.linuxpro.net 
# samba-tool domain provision --server-role=dc --use-rfc2307 --dns-backend=SAMBA_INTERNAL --realm=LINUXPRO.NET --domain=LINUXPRO --adminpass=Linuxpro123456
# or
# samba-tool domain provision --server-role=dc --use-rfc2307 --function-level=2008_R2 --use-xattrs=yes --dns-backend=SAMBA_INTERNAL --realm=LINUXPRO.NET --domain=LINUXPRO --adminpass=Linuxpro123456

### Add start script on boot
# systemctl daemon-reload
# systemctl enable samba4.service
# systemctl start samba4.service