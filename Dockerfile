FROM debian:stretch
LABEL maintainer="Nilton OS <jniltinho@gmail.com>"

## docker build -t debian-samba4 .
## docker run -it debian-samba4 /bin/bash
## docker run -it jniltinho/debian-samba4 /bin/bash
## docker run -d --restart=always -p 8088:8088 debian-samba4
## docker run -d --restart=always -p 8088:8088 jniltinho/debian-samba4
## https://github.com/titpetric/netdata/blob/master/releases/latest/Dockerfile

# docker stop $(docker ps -a -q)
# docker rm $(docker ps -a -q)
# docker rmi -f $(docker images -q)

## Install base packages
RUN apt-get update && \
    apt-get -yq install \
		libreadline-dev wget git build-essential libattr1-dev libblkid-dev libpam0g-dev \
		autoconf python-dev python-dnspython libacl1-dev gdb pkg-config libpopt-dev libldap2-dev \
		dnsutils acl attr libbsd-dev docbook-xsl libcups2-dev libgnutls28-dev ca-certificates && \
	apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /var/cache/apt/archive/*.deb

RUN wget https://download.samba.org/pub/samba/stable/samba-4.8.3.tar.gz && tar xvfz samba-4.8.3.tar.gz && \
    cd samba-4.8.3/ && \
    ./configure --with-ads --with-shared-modules=idmap_ad --enable-debug --enable-selftest \
    --with-systemd --prefix=/opt/samba4 && \
    make && make install && cd ../ && rm -rf samba-4.8.3*

RUN cd /tmp/ && git clone https://github.com/jniltinho/go-samba4.git && \
    mv go-samba4/dist /opt/go-samba4 && rm -rf go-samba4 && chmod +x /opt/go-samba4/go_samba4 && \
    wget https://my-netdata.io/kickstart-static64.sh && /bin/bash kickstart-static64.sh --dont-wait --dont-start-it && \
    cp /opt/netdata/system/netdata-lsb /etc/init.d/netdata && chmod +x /etc/init.d/netdata && \
    cd ../ && rm -rf /tmp/*

ADD docker-entrypoint.sh /usr/local/bin/
RUN ln -s /usr/local/bin/docker-entrypoint.sh /entrypoint.sh && chmod +x /usr/local/bin/docker-entrypoint.sh

EXPOSE 8088
ENTRYPOINT ["/entrypoint.sh"]