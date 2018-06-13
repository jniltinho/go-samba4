# Samba4 Web Manager - Python + Flask Framework

A web interface to manage samba4 server - Python, Flask Framework

## System requirements for development

- Debian 9 64Bits
- Samba4 4.8.2
- Python
- Flask Framework
- Pyinstaller

## System requirements for use

- Debian 9 64Bits
- Samba4 4.8.2

## Installing Samba4 on Debian 9

[Install Samba4](https://github.com/jniltinho/go-samba4/tree/master/contribute)

## Run go-samba4

```bash

cd /opt/
git clone https://github.com/jniltinho/go-samba4.git
mv go-samba4/dist /opt/
rm -rf go-samba4 && mv dist go-samba4

cd /opt/go-samba4/
chmod +x go-samba4/go_samba4
bash <(curl -Ss https://my-netdata.io/kickstart-static64.sh) --dont-wait --dont-start-it
cp /opt/netdata/system/netdata.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable netdata
systemctl start netdata

./go_samba4 --server-prod
## Run https://0.0.0.0:8088
```

## SystemD Daemon go-samba4

```bash
## Create daemon systemd
echo '[Unit]
Description=Go-Samba4 Daemon
After=syslog.target network.target
 
[Service]
WorkingDirectory=/opt/go-samba4
ExecStart=/opt/go-samba4/go_samba4 --server-prod
NonBlocking=true
 
[Install]
WantedBy=multi-user.target' > /etc/systemd/system/go_samba4.service

## Add start script on boot
systemctl daemon-reload
systemctl enable go_samba4.service
systemctl start go_samba4.service
```

## AdminLTE Admin Template

[AdminLTE](https://github.com/almasaeed2010/AdminLTE)

**AdminLTE** -- is a fully responsive admin template. Based on **[Bootstrap 3](https://github.com/twbs/bootstrap)** framework. Highly customizable and easy to use. Fits many screen resolutions from small mobile devices to large desktops. Check out the live preview now and see for yourself.

**Download & Preview on [AdminLTE.IO](https://adminlte.io)**

## Telas

![image](https://raw.github.com/jniltinho/go-samba4/master/screens/login.png)
![image](https://raw.github.com/jniltinho/go-samba4/master/screens/dashboard.png)
![image](https://raw.github.com/jniltinho/go-samba4/master/screens/users.png)
![image](https://raw.github.com/jniltinho/go-samba4/master/screens/grupos.png)
![image](https://raw.github.com/jniltinho/go-samba4/master/screens/add_user.png)
![image](https://raw.github.com/jniltinho/go-samba4/master/screens/add_group.png)
