#!/usr/bin/python
# -*- coding: utf-8 -*-

"""
https://realpython.com/introduction-to-flask-part-1-setting-up-a-static-site/
https://www.digitalocean.com/community/tutorials/how-to-structure-large-flask-applications
rm -rf dist
pyinstaller -w -F --add-data "app:app" go_samba4.py
cp -aR app/templates app/static dist/
rm -rf build go_samba4.spec *.pyc app/*.pyc

"""

import os
from gevent.wsgi import WSGIServer


from app import app
app.secret_key = os.urandom(12)
# app.run(host='0.0.0.0', port=8080)

http_server = WSGIServer(('', 8080), app)
http_server.serve_forever()
