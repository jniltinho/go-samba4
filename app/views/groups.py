# -*- coding: utf-8 -*-


from flask import Blueprint, render_template


from app.model.users import login_required
from app.model.users import get_groups


mod = Blueprint('groups', __name__, url_prefix='/groups')


@mod.route('/')
@login_required
def index():
    ls_groups = get_groups()
    return render_template('groups/index.html', groups=ls_groups)
