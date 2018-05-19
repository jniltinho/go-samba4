# -*- encoding: utf-8 -*-
"""
Python Aplication Template
Licence: GPLv3
"""

import samba
from samba import credentials
from samba.param import LoadParm

from ldb import SCOPE_SUBTREE
from samba.auth import AUTH_SESSION_INFO_DEFAULT_GROUPS as df_gp1
from samba.auth import AUTH_SESSION_INFO_AUTHENTICATED as df_gp2


class auth_base(object):

    def __init__(self, user, password):
        self.user = user
        self.password = password
        self.lp = LoadParm()
        self.lp.load_default()
        self.ip = '127.0.0.1'
        self.WorkGroup = str(self.lp.get("workgroup"))
        self.creds = credentials.Credentials()
        self.creds.set_username(self.user)
        self.creds.set_password(self.password)
        self.creds.set_domain(self.WorkGroup)
        self.creds.set_workstation("")

    def autenticate(self):
        try:
            session_info_flags = (df_gp1 | df_gp2)

            LdapConn = samba.Ldb("ldap://%s" %
                                 self.ip, lp=self.lp, credentials=self.creds)
            DomainDN = LdapConn.get_default_basedn()
            search_filter = "sAMAccountName=%s" % self.user
            res = LdapConn.search(
                base=DomainDN, scope=SCOPE_SUBTREE, expression=search_filter, attrs=["dn"])
            if len(res) == 0:
                return False

            user_dn = res[0].dn
            session = samba.auth.user_session(
                LdapConn, lp_ctx=self.lp, dn=user_dn, session_info_flags=session_info_flags)
            token = session.security_token

            if (token.has_builtin_administrators()):
                return True

            if(token.is_system()):
                return True

        except Exception:
            return False

        return False
