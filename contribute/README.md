## Criando a VM no VirtualBox
Criando o Servidor com Debian 8.7 64Bits



## Instalando, Compilando e Configurando o Samba4 4.6.3

```bash
apt-get install -y libreadline-dev git build-essential libattr1-dev libblkid-dev 
apt-get install -y autoconf python-dev python-dnspython libacl1-dev gdb pkg-config libpopt-dev libldap2-dev 
apt-get install -y dnsutils acl attr libbsd-dev docbook-xsl libcups2-dev libgnutls28-dev


cd /usr/src
wget https://download.samba.org/pub/samba/stable/samba-4.6.3.tar.gz
tar -xzvf samba-4.6.3.tar.gz
cd samba-4.6.3
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

## Instalando Golang 1.8 64Bits

```bash

cd /usr/local/
wget https://storage.googleapis.com/golang/go1.8.1.linux-amd64.tar.gz
tar -xzf go1.8.1.linux-amd64.tar.gz && rm -f go1.8.1.linux-amd64.tar.gz

echo 'export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/GO
export PATH=$PATH:$GOPATH/bin' >> /etc/profile


source /etc/profile
mkdir -p $HOME/GO

```

## Clonando o go-samba4 do GitHub

```bash

cd /opt/
git clone https://github.com/jniltinho/go-samba4.git

```

## Instalando a WEB IDE PARA DESENV Golang

```bash

groupadd wide && useradd -M -s /bin/bash -g wide -d /opt/wide wide

cd /opt/
wget https://github.com/jniltinho/go-samba4/raw/master/contribute/wide-1.5.2-linux-amd64.tar.gz
tar -xvf wide-1.5.2-linux-amd64.tar.gz && rm -f wide-1.5.2-linux-amd64.tar.gz
chown -R wide:wide /opt/wide

cp /opt/wide/scripts/wide.service /etc/systemd/system/wide.service
chmod +x /opt/wide/scripts/start-wide.sh


systemctl daemon-reload
systemctl start wide.service
systemctl enable wide.service

su - wide
go get github.com/astaxie/beego
go get github.com/beego/bee
exit

```

