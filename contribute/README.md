## Criando a VM no VirtualBox
Criando o Servidor com Debian 9.1 64Bits



## Instalando, Compilando e Configurando o Samba4 4.8.0

```bash
apt-get install -y libreadline-dev git build-essential libattr1-dev libblkid-dev libpam0g-dev
apt-get install -y autoconf python-dev python-dnspython libacl1-dev gdb pkg-config libpopt-dev libldap2-dev 
apt-get install -y dnsutils acl attr libbsd-dev docbook-xsl libcups2-dev libgnutls28-dev


cd /usr/src
get_samba4=https://download.samba.org/pub/samba/stable/samba-4.8.0.tar.gz
wget -c ${get_samba4}
tar xvfz $(basename ${get_samba4})
cd $(basename ${get_samba4}|sed "s/.tar.gz//")
./configure --with-ads --with-shared-modules=idmap_ad --enable-debug --enable-selftest --with-systemd --prefix=/opt/samba4
make
make install



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
WantedBy=multi-user.target' > /etc/systemd/system/samba4.service

### Add PATH
echo 'export PATH=$PATH:/opt/samba4/bin:/opt/samba4/sbin' >> /etc/profile
source /etc/profile


### Create Domain Samba4 like AD
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
pip install flask pyinstaller gevent psutil python-ldap Flask-Caching

```

## Clonando o go-samba4 do GitHub

```bash

cd /opt/
git clone https://github.com/jniltinho/go-samba4.git

```

