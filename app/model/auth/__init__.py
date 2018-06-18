import os
import sys

lib_samba = ['/opt/samba4/lib/python2.7/site-packages',
             '/opt/samba4/lib64/python2.7/site-packages',
             '/usr/local/samba/lib/python2.7/site-packages',
             '/usr/local/samba/lib64/python2.7/site-packages'
]

for i in lib_samba:
    if (os.path.exists(i)):
        sys.path.append(i)
