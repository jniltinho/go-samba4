# -*- coding: utf-8 -*-


from flask import Blueprint, render_template
from flask import url_for, session, current_app
from flask import flash, redirect, request


# from app import app
from app.utils import get_users, get_pkgs
from app.utils import get_cpu_stats
from app.AuthSMB4 import AuthSMB4


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


@mod.route('/welcome')
def welcome():
    return render_template('default/welcome.html')  # render a template


@mod.route('/login', methods=['POST'])
def login():
    # app.logger.debug("Request Form %s", request.form)
    base = AuthSMB4(request.form['username'], request.form['password'])
    if base.Autenticate():
        session['logged_in'] = True
        session['username'] = request.form['username']
    else:
        flash('wrong password!')
    return redirect(url_for('default.index'))


@mod.route("/logout")
def logout():
    current_app.cache.clear()
    session.clear()
    # session['logged_in'] = False
    return redirect(url_for('default.index'))
