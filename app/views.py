# -*- coding: utf-8 -*-


from flask import url_for, session, abort
from flask import flash, redirect, render_template, request


from app import app
from app.utils import login_required, get_users, get_pkgs
from app.utils import get_cpu_stats, get_groups, user_delete
from app.AuthSMB4 import AuthSMB4


@app.route('/')
def home():
    if not session.get('logged_in'):
        return render_template('login.html')
    else:
        get_proc = get_cpu_stats()
        get_proc.update({'ls_users': get_users()})
        get_proc.update(get_pkgs())
        return render_template('index.html', **get_proc)


@app.route('/welcome')
def welcome():
    return render_template('welcome.html')  # render a template


@app.route('/users')
@login_required
def users():
    ls_users = get_users()
    return render_template('users.html', users=ls_users)


@app.route('/del/<username>')
@login_required
def users_del(username):
    user_get = request.args.get('user')
    res = user_delete(user_get)
    print res


@app.route('/groups')
@login_required
def groups():
    ls_groups = get_groups()
    return render_template('groups.html', groups=ls_groups)


@app.route('/login', methods=['POST'])
def login():
    # app.logger.debug("Request Form %s", request.form)
    base = AuthSMB4(request.form['username'], request.form['password'])
    if base.Autenticate():
        session['logged_in'] = True
        session['username'] = request.form['username']
    else:
        flash('wrong password!')
    return redirect(url_for('home'))


@app.route("/logout")
def logout():
    session.clear()
    # session['logged_in'] = False
    return redirect(url_for('home'))
