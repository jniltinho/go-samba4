#!/usr/bin/python
# -*- coding: utf-8 -*-


from flask import url_for, session, abort
from flask import flash, redirect, render_template, request

from app import app


@app.route('/')
def home():
    if not session.get('logged_in'):
        return render_template('login.html')
    else:
        return render_template('index.html')


@app.route('/welcome')
def welcome():
    return render_template('welcome.html')  # render a template


@app.route('/login', methods=['POST'])
def login():
    # app.logger.debug("Request Form %s", request.form)
    if request.form['password'] == '123456' and request.form['username'] == 'admin':
        session['logged_in'] = True
    else:
        flash('wrong password!')
    return redirect(url_for('home'))


@app.route("/logout")
def logout():
    session['logged_in'] = False
    return redirect(url_for('home'))
