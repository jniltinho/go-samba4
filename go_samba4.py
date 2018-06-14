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
import optparse
import traceback
from gevent import reinit
from gevent.wsgi import WSGIServer
from gevent.monkey import patch_all
from flask_caching import Cache


os.environ["PATH"] += ':/opt/samba4/bin:/opt/samba4/sbin'

from app import app
app.config['CACHE_TYPE'] = 'simple'
app.cache = Cache(app)
app.secret_key = os.urandom(12)


reinit()
patch_all(dns=False)


def server_prod(host="0.0.0.0", port=8088, ssl=True, debug=True):
    if sys.argv[0].split(".")[-1] != 'py':
        debug = False
    app.debug = debug
    try:
        if ssl:
            print('Starting Gevent HTTP server on https://%s:%s' % (host, port))
            server = WSGIServer(
                (host, port), app, keyfile='ssl/server.key', certfile='ssl/server.crt')
        else:
            print('Starting Gevent HTTP server on http://%s:%s' % (host, port))
            server = WSGIServer((host, port), app)
        server.serve_forever()
    except KeyboardInterrupt:
        print "Shutdown requested...exiting"
    except Exception:
        traceback.print_exc(file=sys.stdout)
    sys.exit(0)


def server_dev(host="0.0.0.0", port=8088, ssl=True, debug=True):
    context = ('ssl/server.crt', 'ssl/server.key')
    if sys.argv[0].split(".")[-1] != 'py':
        debug = False
    try:
        if ssl:
            app.run(host=host, port=port, debug=debug, ssl_context=context)
        else:
            app.run(host=host, port=port, debug=debug)
    except KeyboardInterrupt:
        print "Shutdown requested...exiting"
    except Exception:
        traceback.print_exc(file=sys.stdout)
    sys.exit(0)


def main():
    usage = "Usage: %prog --server-prod|--server-dev"
    parser = optparse.OptionParser(usage)
    parser.add_option("--server-prod", action="store_true",
                      dest="SRV_PROD", default=False, help="Server Gevent Prod")
    parser.add_option("--server-dev", action="store_true",
                      dest="SRV_DEV", default=False, help="Server Flask Desenv")
    parser.add_option("--SSL", action="store_true",
                      dest="SSL", default=False, help="Enable SSL")


    options, args = parser.parse_args()

    if len(sys.argv) == 1:
        parser.print_help()
        sys.exit(1)

    if (options.SRV_PROD):
        server_prod(ssl=options.SSL)

    if (options.SRV_DEV):
        server_dev(ssl=options.SSL)


if __name__ == "__main__":
    main()
