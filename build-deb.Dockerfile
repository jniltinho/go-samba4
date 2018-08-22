ARG DISTRO_IMAGE=debian:9
FROM $DISTRO_IMAGE
LABEL maintainer="Nilton OS <jniltinho@gmail.com>"

# docker build -t build-samba4 -f duild-deb.Dockerfile
# docker build -t build-samba4 --build-arg DISTRO_IMAGE=ubuntu:xenial -f duild-deb.Dockerfile
# docker run build-samba4
# docker run -it build-samba4 /bin/bash
# docker run --rm -it -v "${PWD}:/src" build-samba4 /bin/bash


# docker tag build-samba4 jniltinh/build-samba4
# docker push jniltinh/build-samba4

# docker stop $(docker ps -a -q)
# docker rm $(docker ps -a -q)
# docker rmi -f $(docker images -q)


# https://gist.github.com/jniltinho/7a59467a8a4e5e88a8166f9e7e679e4d
# http://sig9.hatenablog.com/entry/2017/12/04/000000



ADD scripts/build_deb.sh /

## Install base packages
RUN chmod +x /build_deb.sh && sync && sleep 2 && /build_deb.sh

VOLUME /src
WORKDIR /src

CMD ls /src/ /root/