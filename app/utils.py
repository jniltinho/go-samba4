# -*- coding: utf-8 -*-

import os
import psutil
import commands
import platform
from functools import wraps
from flask import Flask, redirect, current_app
from flask import url_for, session, flash

dist = platform.dist()[0].lower()

# Define the WSGI application object
app = Flask(__name__)


def cmd(shell_cmd):
    """
    Classe para enviar comandos shell
    ddDdD
    sdgsdgsdggs
    """
    status, output = commands.getstatusoutput(shell_cmd)
    if status == 0:
        return output.splitlines()


def login_required(f):
    @wraps(f)
    def wrap(*args, **kwargs):
        if 'logged_in' in session:
            return f(*args, **kwargs)
        else:
            flash("You need to login first")
            return redirect(url_for('home'))

    return wrap


def user_create(username, password, given_name, surname):
    cli = "samba-tool user create %s %s --given-name='%s' --surname='%s'" % (
        username, password, given_name, surname)
    res = cmd(cli)
    print res


def get_pkgs():
    cached = current_app.cache.get('get_pkgs')
    if cached:
        return cached
    if dist == 'ubuntu' or dist == 'debian':
        output = cmd("dpkg-query -f '${Package};${Version};deb\n' -W")
        pkg_type = 'deb'
    if dist == 'centos':
        output = cmd("rpm -qa --qf '%{NAME}.%{ARCH};%{VERSION}-%{RELEASE};rpm\n'|sort")
        pkg_type = 'rpm'
    pkgs = []
    for value in output:
        pak = value.split(";")
        pkgs.append([pak[0], pak[1], pak[2]])

    result = {'pkgs': pkgs, 'pkg_count': len(pkgs), 'pkg_type': pkg_type}
    current_app.cache.set('get_pkgs', result, timeout=300)
    return result


def user_delete(username):
    res = cmd("samba-tool user delete %s" % (username))
    return res


def group_list():
    res = cmd("samba-tool group list")
    print res


def get_users():
    users = {'result': []}
    res = cmd("samba-tool user list|egrep -v 'Guest|krbtgt'")
    cli = "|".join(res)
    get_users = cmd("pdbedit -Lf|egrep '%s'" % cli)
    # print get_users
    for user in get_users:
        if user.split(':')[2]:
            users['result'].append({'username': user.split(':')[0], 'full_name': user.split(':')[2]})
        else:
            users['result'].append({'username': user.split(':')[0], 'full_name': user.split(':')[0].title()})

    return users


def get_groups():
    ls_groups = sorted(cmd("samba-tool group list"))
    return ls_groups


def get_cpu_stats():
    if current_app.cache.get('get_cpu_stats'):
        return current_app.cache.get('get_cpu_stats')
    data = {}
    users = cmd("samba-tool user list")
    groups = cmd("samba-tool group list")
    data.update({'load': os.getloadavg(), 'jobs': len(psutil.pids())})
    data.update({'users': len(users), 'groups': len(groups)})
    current_app.cache.set('get_cpu_stats', data, timeout=300)
    return data

# user_create("nicolas", "Nicolas@2018", "Nicolas", "de Oliveira Silva")
# group_list()
# user_list()
