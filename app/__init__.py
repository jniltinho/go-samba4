# Import flask and template operators
from flask import Flask, render_template, request


# Define the WSGI application object
app = Flask(__name__)


from app import views


@app.before_request
def log_request():
    app.logger.debug("Request Headers %s", request.headers)
    return None


@app.errorhandler(404)
def not_found(error):
    return render_template('404.html'), 404
