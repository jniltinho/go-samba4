# -*- coding: utf-8 -*-


from flask import Blueprint, render_template
from flask import url_for, session, current_app
from flask import flash, redirect, request


# from app import app
from app.model.users import get_users, get_pkgs
from app.model.users import get_cpu_stats
from app.model.auth.auth_base import auth_base


mod = Blueprint('default', __name__)


@mod.route('/')
def index():
    if not session.get('logged_in'):
        return render_template('default/login.html')
    else:
        get_proc = get_cpu_stats()
        get_proc.update({'ls_users': get_users()})
        get_proc.update(get_pkgs())
        return render_template('default/index.html', **get_proc)


@mod.route('/netdata/')
def netdata():
    if not session.get('logged_in'):
        return render_template('default/login.html')
    else:
        netdata_ip = request.host.split(':')[0] + ':19999'
        return render_template('default/netdata.html', netdata_ip=netdata_ip)  # render a template


@mod.route('/login', methods=['POST'])
def login():
    # app.logger.debug("Request Form %s", request.form)
    base = auth_base(request.form['username'], request.form['password'])
    if base.autenticate():
        session['logged_in'] = True
        session['username'] = request.form['username']
        session['password'] = request.form['password']
    else:
        flash('wrong password!')
    return redirect(url_for('default.index'))


@mod.route("/logout")
def logout():
    current_app.cache.clear()
    session.clear()
    # session['logged_in'] = False
    return redirect(url_for('default.index'))
