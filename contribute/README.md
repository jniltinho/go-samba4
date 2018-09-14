
# Criando a VM no VirtualBox

Criando o Servidor com Debian 9 64Bits

## Instalando, Compilando e Configurando o Samba4 4.9.0

```bash
apt-get update
apt-get -yq install ruby-dev
apt-get -yq install libreadline-dev git build-essential libattr1-dev libblkid-dev libpam0g-dev
apt-get -yq install autoconf python-dev python-dnspython libacl1-dev gdb pkg-config libpopt-dev
apt-get -yq install libldap2-dev libtirpc-dev libxslt1-dev python-pycryptopp libgnutls28-dev
apt-get -yq install dnsutils acl attr libbsd-dev libcups2-dev libgnutls28-dev curl wget
apt-get -yq install docbook-xsl libacl1-dev gdb liblmdb-dev libjansson-dev libpam0g-dev libgpgme-dev
apt-get -yq install tracker libtracker-sparql-1.0-dev libavahi-client-dev libavahi-common-dev bison flex
apt-get -yq install libarchive-dev
gem install fpm


cd /usr/src
get_samba4=https://download.samba.org/pub/samba/stable/samba-4.9.0.tar.gz

PKG=$(basename ${get_samba4}|sed "s/.tar.gz//")
PKG_NAME=$(basename ${get_samba4}|sed "s/.tar.gz//"|cut -d- -f1)
PKG_VERSION=$(basename ${get_samba4}|sed "s/.tar.gz//"|cut -d- -f2)

wget -c ${get_samba4}
tar xvfz $(basename ${get_samba4})
cd $(basename ${get_samba4}|sed "s/.tar.gz//")
./configure --with-ads --systemd-install-services --with-shared-modules=idmap_ad --enable-debug --enable-selftest --with-systemd --enable-spotlight --prefix=/opt/samba4
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


fpm -s dir -t deb -n ${PKG_NAME} -v ${PKG_VERSION} -C /tmp/installdir \
  -d "python-minimal" \
  -d "libpython2.7" \
  -d "libbsd0" \
  -d "libpopt0" \
  -d "libgnutls30" \
  -d "libldap-2.4-2" \
  -d "libcups2" \
  -d "libjansson4" \
  -d "libtracker-sparql-1.0-0" \
  -d "libgpgme11" \
  -p ${PKG}+dfsg-1.amd64.deb .

mv ${PKG}+dfsg-1.amd64.deb /root/
dpkg -i /root/${PKG}+dfsg-1.amd64.deb

### Add PATH
echo 'export PATH=$PATH:/opt/samba4/bin:/opt/samba4/sbin' >> /etc/profile
source /etc/profile


### Create Domain Samba4 like AD
# hostnamectl set-hostname samba4.linuxpro.net 
# samba-tool domain provision --server-role=dc --use-rfc2307 --dns-backend=SAMBA_INTERNAL --realm=LINUXPRO.NET --domain=LINUXPRO --adminpass=Linuxpro123456
# or
# samba-tool domain provision --server-role=dc --use-rfc2307 --function-level=2008_R2 --use-xattrs=yes --dns-backend=SAMBA_INTERNAL --realm=LINUXPRO.NET --domain=LINUXPRO --adminpass=Linuxpro123456

### Add start script on boot
# systemctl daemon-reload
# systemctl enable samba4.service
# systemctl start samba4.service
```

## Instalando o Framework Flask

```bash

apt-get install -y python-pip git-core dos2unix
apt-get install -y libsasl2-dev python-dev libldap2-dev libssl-dev
```

## Clonando o go-samba4 do GitHub

```bash

cd /opt/
mkdir dev_go-samba4 && cd dev_go-samba4
git clone https://github.com/jniltinho/go-samba4.git
cd go-samba4
pip install -r requirements.txt
```