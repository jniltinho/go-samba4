#!/usr/bin/env python
# -*- coding: utf-8 -*-


"""
Usando o pyinstaller para gerar um binario
pip install pyinstaller
pyinstaller file.py -F
rm -rf build *.spec *.pyc
"""

import os
import sys
import platform


os_plat = platform.system().lower()
bits = platform.machine()[-2:]


def compile_32(file_name):
    if os.path.exists(file_name):
        os.system('pyinstaller %s -F --add-data "app:app" --distpath=dist_32' % (file_name))
        os.system("rm -rf build *.spec *.pyc app/*.pyc")
        os.system("rm -rf dist_32/templates dist_32/static dist_32/ssl")
        os.system("cp -aR app/templates app/static ssl dist_32/")
    else:
        print "File %s NotFound !!!" % (file_name)


def compile_64(file_name):
    if os.path.exists(file_name):
        os.system('pyinstaller %s -F --add-data "app:app" --distpath=dist' % (file_name))
        os.system("rm -rf build *.spec *.pyc app/*.pyc")
        os.system("rm -rf dist/templates dist/static dist/ssl")
        os.system("cp -aR app/templates app/static ssl dist/")
    else:
        print "File %s NotFound !!!" % (file_name)


if len(sys.argv) == 1:
    print "%s + O nome do Arquivo Python para Compilar" % (sys.argv[0])
    sys.exit(1)


if len(sys.argv) == 2:
    file_name = sys.argv[1]
    if bits == '64':
        compile_64(file_name)
    else:
        compile_32(file_name)
