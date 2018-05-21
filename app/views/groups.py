# -*- coding: utf-8 -*-


from flask import Blueprint, render_template
from flask import session, current_app, jsonify
from flask import request, redirect, url_for

from app.model.users import login_required
from app.model.users import get_groups, group_create, group_delete


mod = Blueprint('groups', __name__, url_prefix='/groups')


@mod.route('/')
@login_required
def index():
    ls_groups = get_groups()
    return render_template('groups/index.html', groups=ls_groups)


@mod.route('/add/', methods=['POST', 'GET'])
@login_required
def add():
    if request.method == 'POST':
        groupname = request.form['groupname'].lstrip().rstrip()
        group_create(groupname)
        return redirect(url_for('groups.index'))
    return render_template('groups/add_group.html')


@mod.route('/del/', methods=['POST'])
@login_required
def groups_del():
    groupname = request.json['groupname']
    res = group_delete(groupname)
    if res:
        current_app.cache.clear()
        return jsonify(message=res, deleted=True)
    return jsonify(message="No Deleted Group", deleted=False)