# -*- coding: utf-8 -*-

import os
import psutil
import commands
from functools import wraps
from flask import Flask, redirect
from flask import url_for, session, flash


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
    data = {}
    users = cmd("samba-tool user list")
    groups = cmd("samba-tool group list")
    data.update({'load': os.getloadavg(), 'jobs': len(psutil.pids())})
    data.update({'users': len(users), 'groups': len(groups)})
    return data

# user_create("nicolas", "Nicolas@2018", "Nicolas", "de Oliveira Silva")
# group_list()
# user_list()
