# -*- coding: utf-8 -*-

import commands


def cmd(shell_cmd):
    """
    Classe para enviar comandos shell
    ddDdD
    sdgsdgsdggs
    """
    status, output = commands.getstatusoutput(shell_cmd)
    if status == 0:
        return output.splitlines()


def user_create(username, password, given_name, surname):
    cli = "samba-tool user create %s %s --given-name='%s' --surname='%s'" % (
        username, password, given_name, surname)
    res = cmd(cli)
    print res


def group_list():
    res = cmd("samba-tool group list")
    print res


def user_list():
    res = cmd("samba-tool user list")
    res.remove('krbtgt')
    print res


# user_create("nicolas", "Nicolas@2018", "Nicolas", "de Oliveira Silva")
# group_list()
user_list()
