#!/usr/bin/env python
# -*- coding: utf-8 -*-


"""
https://realpython.com/introduction-to-flask-part-1-setting-up-a-static-site/
https://www.digitalocean.com/community/tutorials/how-to-structure-large-flask-applications
https://github.com/stgraber/samba4-manager
https://github.com/baboons/samba4-gaps

pyinstaller -F --add-data "app:app" go_samba4.py
cp -aR app/templates app/static dist/
rm -rf build *.spec *.pyc app/*.pyc
openssl req -nodes -new -x509 -keyout ssl/server.key -out ssl/server.crt

"""

import os
import sys
import traceback
from gevent import reinit
from gevent.wsgi import WSGIServer
from gevent.monkey import patch_all
reinit()
patch_all(dns=False)


host = "0.0.0.0"
port = 8088


def main():
    from app import app
    app.secret_key = os.urandom(12)
    # app.run(host='0.0.0.0', port=8080)
    try:
        print('Starting Gevent HTTP server on https://%s:%s' % (host, port))
        WSGIServer((host, port), app, keyfile='ssl/server.key',
                   certfile='ssl/server.crt').serve_forever()
    except KeyboardInterrupt:
        print "Shutdown requested...exiting"
    except Exception:
        traceback.print_exc(file=sys.stdout)
    sys.exit(0)


if __name__ == "__main__":
    main()
