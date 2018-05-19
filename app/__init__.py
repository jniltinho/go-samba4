# -*- coding: utf-8 -*-

# Import flask and template operators
from datetime import timedelta
from flask import Flask, render_template, request
from flask import session


# Define the WSGI application object
app = Flask(__name__)


@app.before_request
def log_request():
    app.logger.debug("Request Headers %s", request.headers)
    return None


@app.before_request
def make_session_permanent():
    session.permanent = True
    app.permanent_session_lifetime = timedelta(minutes=5)


@app.errorhandler(404)
def not_found(error):
    return render_template('404.html'), 404


## Init App
from views import default, users, groups

app.register_blueprint(default.mod)
app.register_blueprint(users.mod)
app.register_blueprint(groups.mod)
