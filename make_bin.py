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


def compile_py(file_name):
    if os.path.exists(file_name):
        bin_name = file_name.split(".")[0]
        os.system("find . -type f -iname *.pyc -exec rm  -f {} \;")
        os.system("pyinstaller %s --clean --add-data 'app:app' -F -n %s" %
                  (file_name, bin_name))
        os.system("rm -rf build *.spec *.pyc")
        os.system("find . -type f -iname *.pyc -exec rm  -f {} \;")
        os.system("rm -rf dist/ssl dist/static dist/templates")
        os.system("cp -aR ssl themes/AdminLTE/* dist/")
    else:
        print "File %s NotFound !!!" % (file_name)


if len(sys.argv) == 1:
    print "%s + O nome do Arquivo Python para Compilar" % (sys.argv[0])
    sys.exit(1)


if len(sys.argv) == 2:
    file_name = sys.argv[1]
    compile_py(file_name)
