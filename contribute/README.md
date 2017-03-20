## Criando a VM no VirtualBox
Criando o Servidor com Debian 8.7 64Bits



## Instalando, Compilando e Configurando o Samba4 4.6.0


## Instalando Golang 1.8 64Bits



## Clonando o go-samba4 do GitHub



## Instalando a WEB IDE PARA DESENV Golang

```bash

groupadd wide && useradd -M -s /bin/bash -g wide -d /opt/wide wide

cd /opt/
wget https://github.com/jniltinho/go-samba4/raw/master/contribute/wide-1.5.2-linux-amd64.tar.gz
tar -xvf wide-1.5.2-linux-amd64.tar.gz && rm -f wide-1.5.2-linux-amd64.tar.gz
chown -R wide:wide /opt/wide

cp /opt/wide/scripts/wide.service /etc/systemd/system/wide.service
chmod +x /opt/wide/scripts/start-wide.sh

```

