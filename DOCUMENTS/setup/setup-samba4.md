# Setup — Samba 4 como Active Directory Domain Controller
## Ubuntu 24.04 LTS

> **Sistema:** Ubuntu Server 24.04 LTS (Noble Numbat)
> **Função:** Active Directory Domain Controller (AD DC) primário
> **Referência oficial:** [wiki.samba.org — Setting up Samba as an AD DC](https://wiki.samba.org/index.php/Setting_up_Samba_as_an_Active_Directory_Domain_Controller)

---

## Sumário

1. [Pré-requisitos](#1-pré-requisitos)
2. [Configuração de Rede (IP Estático)](#2-configuração-de-rede-ip-estático)
3. [Hostname e /etc/hosts](#3-hostname-e-etchosts)
4. [Desabilitar systemd-resolved](#4-desabilitar-systemd-resolved)
5. [Instalar Pacotes](#5-instalar-pacotes)
6. [Preparar Ambiente](#6-preparar-ambiente)
7. [Provisionar o Domínio](#7-provisionar-o-domínio)
8. [Configurar Kerberos](#8-configurar-kerberos)
9. [Configurar Serviços Systemd](#9-configurar-serviços-systemd)
10. [Iniciar e Verificar](#10-iniciar-e-verificar)
11. [Configurar Firewall (UFW)](#11-configurar-firewall-ufw)
12. [Habilitar LDAPS (porta 636)](#12-habilitar-ldaps-porta-636)
13. [Criar Usuário Bind para go-samba4](#13-criar-usuário-bind-para-go-samba4)
14. [Criar Grupos de Acesso para go-samba4](#14-criar-grupos-de-acesso-para-go-samba4)
15. [Verificações Pós-Instalação](#15-verificações-pós-instalação)
16. [Configuração de NTP/Chrony](#16-configuração-de-ntpchrony)
17. [Gerenciamento com samba-tool](#17-gerenciamento-com-samba-tool)
18. [Troubleshooting](#18-troubleshooting)

---

## 1. Pré-requisitos

### Hardware mínimo recomendado

| Recurso | Mínimo | Recomendado |
|---|---|---|
| CPU | 1 vCPU | 2 vCPU |
| RAM | 1 GB | 4 GB |
| Disco | 10 GB | 40 GB |
| Rede | 1 interface | 1 interface (IP estático obrigatório) |

### Requisitos de software

- Ubuntu Server 24.04 LTS — instalação limpa, **sem** outros serviços DNS (bind9, dnsmasq)
- Acesso root ou sudo
- Conexão com a internet para instalar pacotes
- Sistema de arquivos **ext4 ou xfs** — ACLs e xattrs habilitados por padrão

### Definir variáveis do ambiente (use os seus valores reais)

Ao longo deste guia, os seguintes valores de exemplo serão usados. **Substitua por seus dados reais antes de executar qualquer comando:**

```
Hostname curto:   dc1
FQDN:             dc1.empresa.local
Realm (Kerberos): EMPRESA.LOCAL
Domínio NetBIOS:  EMPRESA
IP do servidor:   192.168.1.10
Gateway:          192.168.1.1
Interface de rede: ens18   (verifique com: ip a)
DNS forwarder:    8.8.8.8
```

> ⚠️ **Atenção:** O Realm Kerberos **deve ser em maiúsculas**. O domínio DNS (`empresa.local`) deve ser um domínio que você controla e que **não** seja roteável na internet. Evite usar `.local` em redes onde mDNS (Avahi/Bonjour) esteja presente — prefira `.lan`, `.intranet` ou um subdomínio real.

---

## 2. Configuração de Rede (IP Estático)

O DC **deve** ter IP estático. Descubra o nome da sua interface de rede:

```bash
ip a
```

Edite o arquivo de configuração do Netplan:

```bash
sudo nano /etc/netplan/01-netcfg.yaml
```

Conteúdo (substitua `ens18` pela sua interface):

```yaml
network:
  version: 2
  renderer: networkd
  ethernets:
    ens18:
      dhcp4: no
      addresses:
        - 192.168.1.10/24
      routes:
        - to: default
          via: 192.168.1.1
      nameservers:
        # O DC deve apontar para si mesmo como DNS primário
        addresses:
          - 192.168.1.10
          - 8.8.8.8
        search:
          - empresa.local
```

> ⚠️ **Importante:** O DC deve apontar para o seu **próprio IP** como DNS primário após o provisionamento. Durante a instalação, pode usar `8.8.8.8` como primário.

Aplicar configuração:

```bash
sudo chmod 600 /etc/netplan/01-netcfg.yaml
sudo netplan apply
```

Verificar:

```bash
ip a show ens18
ip route
```

---

## 3. Hostname e /etc/hosts

Definir o hostname curto:

```bash
sudo hostnamectl set-hostname dc1
```

Editar `/etc/hosts` para que o FQDN resolva corretamente para o IP estático:

```bash
sudo nano /etc/hosts
```

O arquivo deve conter (remova qualquer linha duplicada para `127.0.1.1`):

```
127.0.0.1       localhost
192.168.1.10    dc1.empresa.local dc1
```

> ⚠️ **Nunca** coloque o FQDN do DC apontando para `127.0.0.1` ou `127.0.1.1`. Deve apontar para o IP estático real da interface de rede.

Verificar resolução:

```bash
hostname
# dc1

hostname -f
# dc1.empresa.local

ping -c1 dc1.empresa.local
# deve retornar 192.168.1.10
```

---

## 4. Desabilitar systemd-resolved

O Samba precisa controlar a porta 53 (DNS). O `systemd-resolved` ocupa essa porta por padrão no Ubuntu 24.04 e **deve ser desativado**.

```bash
# Parar e desabilitar o serviço
sudo systemctl disable --now systemd-resolved

# Remover o symlink do resolv.conf gerenciado pelo systemd
sudo rm /etc/resolv.conf

# Criar resolv.conf estático apontando para o próprio DC
sudo tee /etc/resolv.conf << 'EOF'
nameserver 192.168.1.10
nameserver 8.8.8.8
search empresa.local
EOF

# Tornar o arquivo imutável para evitar sobrescrita
sudo chattr +i /etc/resolv.conf
```

Verificar que a porta 53 está livre:

```bash
sudo ss -tlnp | grep ':53'
# Não deve retornar nada
```

---

## 5. Instalar Pacotes

Atualizar o sistema e instalar os pacotes necessários:

```bash
sudo apt update && sudo apt upgrade -y

# Pacote principal do AD DC (Ubuntu 24.04 usa samba-ad-dc separado)
sudo apt install -y \
    samba-ad-dc \
    krb5-user \
    winbind \
    smbclient \
    bind9-dnsutils \
    ldb-tools \
    chrony
```

> **Nota:** Durante a instalação do `krb5-user`, o instalador perguntará sobre o Realm e servidores Kerberos. Pode pressionar `Enter` para aceitar os defaults — o arquivo `/etc/krb5.conf` será substituído corretamente após o provisionamento.

---

## 6. Preparar Ambiente

### 6.1 Desabilitar os serviços Samba padrão

O Ubuntu instala `smbd`, `nmbd` e `winbind` como serviços separados. Para o modo AD DC, eles devem estar **desabilitados e mascarados** — o serviço `samba-ad-dc` os substitui:

```bash
sudo systemctl disable --now smbd nmbd winbind
sudo systemctl mask smbd nmbd winbind
```

Habilitar o serviço AD DC (sem iniciar ainda):

```bash
sudo systemctl unmask samba-ad-dc
sudo systemctl enable samba-ad-dc
```

### 6.2 Remover configuração antiga do Samba

O `samba-tool domain provision` recusa-se a sobrescrever um `smb.conf` existente:

```bash
# Fazer backup do smb.conf original
sudo mv /etc/samba/smb.conf /etc/samba/smb.conf.original

# Limpar arquivos de banco de dados anteriores (se existirem)
sudo rm -f /var/lib/samba/*.tdb
sudo rm -f /var/cache/samba/*.tdb
```

---

## 7. Provisionar o Domínio

Execute o provisionamento interativo com suporte a RFC 2307 (atributos POSIX para UID/GID no AD):

```bash
sudo samba-tool domain provision \
    --use-rfc2307 \
    --interactive
```

Respostas esperadas para cada pergunta:

```
Realm []:                          EMPRESA.LOCAL
Domain [EMPRESA]:                  EMPRESA
Server Role (dc, member, standalone) [dc]:   dc
DNS backend (SAMBA_INTERNAL, BIND9_FLATFILE, BIND9_DLZ, NONE) [SAMBA_INTERNAL]:   SAMBA_INTERNAL
DNS forwarder IP address [10.x.x.x]:   8.8.8.8
Administrator password:            (senha forte — mín. 8 chars, complexidade)
Retype password:
```

> **`--use-rfc2307`**: Habilita atributos POSIX (uidNumber, gidNumber, loginShell, unixHomeDirectory) no schema do AD. Necessário para integração de clientes Linux ao domínio.
>
> **DNS backend `SAMBA_INTERNAL`**: Mais simples para começar. O Samba inclui seu próprio servidor DNS interno. Para ambientes maiores, considere `BIND9_DLZ`.

Ao final do provisionamento bem-sucedido, você verá:

```
Server Role:           active directory domain controller
Hostname:              dc1
NetBIOS Domain:        EMPRESA
DNS Domain:            empresa.local
DOMAIN SID:            S-1-5-21-XXXXXXXXXX-XXXXXXXXXX-XXXXXXXXXX
```

O arquivo `/etc/samba/smb.conf` terá sido criado com a configuração do AD DC.

---

## 8. Configurar Kerberos

O provisionamento gera um arquivo `krb5.conf` correto em `/var/lib/samba/private/`. Copie-o para o local padrão:

```bash
sudo cp -f /var/lib/samba/private/krb5.conf /etc/krb5.conf
```

Verificar o conteúdo gerado:

```bash
cat /etc/krb5.conf
```

Deve conter algo como:

```ini
[libdefaults]
    default_realm = EMPRESA.LOCAL
    dns_lookup_realm = false
    dns_lookup_kdc = true

[realms]
    EMPRESA.LOCAL = {
        default_domain = empresa.local
    }

[domain_realm]
    empresa.local = EMPRESA.LOCAL
    .empresa.local = EMPRESA.LOCAL
```

---

## 9. Configurar Serviços Systemd

### 9.1 Verificar smb.conf gerado

```bash
sudo testparm -s
```

A saída deve mostrar `Server role: ROLE_ACTIVE_DIRECTORY_DC` sem erros.

### 9.2 Verificar configuração do samba-ad-dc

```bash
# Confirmar que o serviço está habilitado
sudo systemctl is-enabled samba-ad-dc
# enabled

# Confirmar que smbd/nmbd/winbind estão mascarados
sudo systemctl is-masked smbd nmbd winbind
# masked (x3)
```

---

## 10. Iniciar e Verificar

### 10.1 Iniciar o Samba AD DC

```bash
sudo systemctl start samba-ad-dc
sudo systemctl status samba-ad-dc
```

A saída deve mostrar `active (running)`.

### 10.2 Verificar DNS

```bash
# Resolução do hostname do DC
host -t A dc1.empresa.local localhost
# dc1.empresa.local has address 192.168.1.10

# Registros SRV do LDAP (necessário para domain join)
host -t SRV _ldap._tcp.empresa.local localhost
# _ldap._tcp.empresa.local has SRV record 0 100 389 dc1.empresa.local.

# Registros SRV do Kerberos UDP
host -t SRV _kerberos._udp.empresa.local localhost
# _kerberos._udp.empresa.local has SRV record 0 100 88 dc1.empresa.local.

# Registros SRV do Kerberos TCP
host -t SRV _kerberos._tcp.empresa.local localhost
# _kerberos._tcp.empresa.local has SRV record 0 100 88 dc1.empresa.local.
```

### 10.3 Verificar autenticação Kerberos

```bash
kinit administrator@EMPRESA.LOCAL
# Password for administrator@EMPRESA.LOCAL: (senha definida no provisionamento)

klist
# Ticket cache: FILE:/tmp/krb5cc_XXXX
# Default principal: administrator@EMPRESA.LOCAL
# Valid starting   Expires            Service principal
# xx/xx/xx xx:xx   xx/xx/xx xx:xx    krbtgt/EMPRESA.LOCAL@EMPRESA.LOCAL
```

### 10.4 Verificar portas abertas

```bash
sudo ss -tlnp | grep samba
# Deve mostrar: 88, 135, 389, 445, 464, 636, 3268, 3269
```

### 10.5 Verificar nível funcional do domínio

```bash
sudo samba-tool domain level show
# Domain and forest function level for domain 'DC=empresa,DC=local'
# Forest function level: (Windows) 2008 R2
# Domain function level: (Windows) 2008 R2
# Lowest function level of a DC: (Windows) 2008 R2
```

---

## 11. Configurar Firewall (UFW)

Habilitar e configurar o UFW com todas as portas necessárias para um AD DC:

```bash
# Habilitar UFW
sudo ufw enable

# SSH (não esqueça de liberar antes de ativar o UFW!)
sudo ufw allow 22/tcp

# Kerberos
sudo ufw allow 88/tcp
sudo ufw allow 88/udp

# DNS (Samba DNS interno)
sudo ufw allow 53/tcp
sudo ufw allow 53/udp

# LDAP e LDAPS
sudo ufw allow 389/tcp
sudo ufw allow 389/udp
sudo ufw allow 636/tcp

# SMB / CIFS
sudo ufw allow 445/tcp
sudo ufw allow 139/tcp

# NetBIOS
sudo ufw allow 137/udp
sudo ufw allow 138/udp

# Kerberos password change
sudo ufw allow 464/tcp
sudo ufw allow 464/udp

# RPC Endpoint Mapper
sudo ufw allow 135/tcp

# Global Catalog LDAP
sudo ufw allow 3268/tcp
sudo ufw allow 3269/tcp

# RPC dinâmico (necessário para domain join Windows)
sudo ufw allow 49152:65535/tcp

# Recarregar
sudo ufw reload
sudo ufw status verbose
```

---

## 12. Habilitar LDAPS (porta 636)

O Samba gera automaticamente um certificado TLS autoassinado durante o provisionamento. O LDAPS (porta 636) já está ativo por padrão.

### 12.1 Verificar certificados gerados

```bash
ls -la /var/lib/samba/private/tls/
# ca.pem      — CA autoassinada
# cert.pem    — Certificado do servidor
# key.pem     — Chave privada (protegida)
```

### 12.2 Verificar LDAPS funcionando

```bash
# Teste de conexão LDAPS
ldapsearch -H ldaps://dc1.empresa.local \
    -x \
    -D "CN=Administrator,CN=Users,DC=empresa,DC=local" \
    -W \
    -b "DC=empresa,DC=local" \
    "(objectClass=domain)" \
    -LLL \
    --no-verify-cert
```

### 12.3 Configurar TLS no smb.conf (opcional — para certificado próprio)

Se quiser usar um certificado próprio (CA interna ou Let's Encrypt), adicione ao `[global]` do `/etc/samba/smb.conf`:

```ini
[global]
    tls enabled  = yes
    tls keyfile  = /etc/ssl/private/samba-dc.key
    tls certfile = /etc/ssl/certs/samba-dc.crt
    tls cafile   = /etc/ssl/certs/ca.crt
    # Desabilitar TLS 1.0 e 1.1 (manter apenas TLS 1.2 e 1.3)
    tls priority = NORMAL:-VERS-SSL3.0:-VERS-TLS1.0:-VERS-TLS1.1
```

Após qualquer mudança no `smb.conf`:

```bash
sudo systemctl restart samba-ad-dc
```

### 12.4 Configurar LDAP para aceitar conexões simples (necessário para go-samba4)

Por padrão, o Samba 4 exige autenticação forte (TLS ou SASL) em conexões LDAP. Para permitir bind simples com senha via LDAPS (porta 636), adicione ao `[global]` do `smb.conf`:

```ini
[global]
    ldap server require strong auth = allow_sasl_over_tls
```

> **Nota de segurança:** O valor `allow_sasl_over_tls` permite bind simples **somente sobre TLS** (LDAPS/porta 636 ou STARTTLS). Nunca use `no`, pois isso permitiria senhas em texto claro via porta 389.

Reiniciar após a mudança:

```bash
sudo systemctl restart samba-ad-dc
```

---

## 13. Criar Usuário Bind para go-samba4

O `go-samba4` usa uma conta de serviço dedicada para se conectar ao AD via LDAP. **Não use o Administrator** para isso.

### 13.1 Criar a conta de serviço

```bash
sudo samba-tool user create samba4admin \
    --description="Bind user for go-samba4 web panel" \
    --mail-address="samba4admin@empresa.local"

# A senha será solicitada
```

### 13.2 Definir senha que não expire

```bash
sudo samba-tool user setexpiry samba4admin --noexpiry
```

### 13.3 Adicionar ao grupo Domain Admins (necessário para ler/escrever todos os atributos AD)

```bash
sudo samba-tool group addmembers "Domain Admins" samba4admin
```

> **Alternativa mais segura para produção:** Em vez de adicionar ao `Domain Admins`, configure permissões granulares via `dsacls` ou delegue apenas as OUs necessárias. Para desenvolvimento, `Domain Admins` é suficiente.

### 13.4 Verificar o DN do usuário criado

```bash
sudo samba-tool user show samba4admin
# dn: CN=samba4admin,CN=Users,DC=empresa,DC=local
```

Este DN será usado no `config.toml` do `go-samba4`:

```toml
[ldap]
bind_user = "CN=samba4admin,CN=Users,DC=empresa,DC=local"
```

---

## 14. Criar Grupos de Acesso para go-samba4

O `go-samba4` usa RBAC mapeado a grupos AD. Criar os grupos necessários:

```bash
# Grupo de administradores do painel web (acesso total)
sudo samba-tool group add SambaWebAdmins \
    --description="Full access to go-samba4 web panel"

# Grupo de operadores (CRUD de usuários/grupos, sem configurações)
sudo samba-tool group add SambaWebOperators \
    --description="Operator access to go-samba4 web panel"

# Grupo helpdesk (reset de senha, habilitar/desabilitar contas)
sudo samba-tool group add SambaWebHelpdesk \
    --description="Helpdesk access to go-samba4 web panel"

# Grupo somente leitura
sudo samba-tool group add SambaWebReadOnly \
    --description="Read-only access to go-samba4 web panel"
```

Adicionar usuários aos grupos:

```bash
# Adicionar o Administrator ao grupo de admins do painel
sudo samba-tool group addmembers SambaWebAdmins Administrator

# Verificar membros
sudo samba-tool group listmembers SambaWebAdmins
```

---

## 15. Verificações Pós-Instalação

### 15.1 Cheklist completo de verificação

```bash
# 1. Status do serviço
sudo systemctl status samba-ad-dc

# 2. Versão do Samba
samba --version

# 3. Nível funcional do domínio
sudo samba-tool domain level show

# 4. Listar DCs
sudo samba-tool domain info 192.168.1.10

# 5. Verificar políticas de senha
sudo samba-tool domain passwordsettings show

# 6. Listar usuários
sudo samba-tool user list

# 7. Listar grupos
sudo samba-tool group list

# 8. Verificar todos os registros DNS SRV
for srv in _ldap._tcp _kerberos._tcp _kerberos._udp _kpasswd._tcp _kpasswd._udp; do
    echo -n "${srv}.empresa.local: "
    host -t SRV ${srv}.empresa.local localhost 2>&1 | grep "has SRV" || echo "FALHOU"
done

# 9. Testar autenticação LDAP via smbclient
smbclient -L dc1.empresa.local -U Administrator
```

### 15.2 Verificar bind LDAP com o usuário de serviço

```bash
ldapsearch -H ldaps://dc1.empresa.local \
    -x \
    -D "CN=samba4admin,CN=Users,DC=empresa,DC=local" \
    -W \
    -b "CN=Users,DC=empresa,DC=local" \
    "(objectClass=user)" \
    sAMAccountName displayName mail \
    -LLL \
    --no-verify-cert
```

A saída deve listar os usuários do domínio sem erros de autenticação.

### 15.3 Verificar registro de auditoria do Samba

```bash
sudo journalctl -u samba-ad-dc -f
```

---

## 16. Configuração de NTP/Chrony

O Kerberos é **extremamente sensível à diferença de horário** entre cliente e servidor (tolerância padrão: 5 minutos). O DC deve ser a fonte de tempo autoritativa para todos os membros do domínio.

### 16.1 Configurar Chrony como servidor NTP

```bash
sudo nano /etc/chrony/chrony.conf
```

Adicione ao final do arquivo:

```
# Fontes NTP externas
pool pool.ntp.org iburst

# Permitir que clientes do domínio sincronizem com este DC
allow 192.168.1.0/24

# Servir tempo mesmo sem sincronização externa (para ambientes isolados)
local stratum 10
```

Reiniciar o Chrony:

```bash
sudo systemctl restart chrony
sudo chronyc sources -v
sudo chronyc tracking
```

### 16.2 Integração NTP com Samba (MS-SNTP)

O Samba possui suporte a MS-SNTP (autenticado) para clientes Windows. Para habilitá-lo, adicione ao `smb.conf`:

```ini
[global]
    ntp signd socket directory = /var/lib/samba/ntp_signd
```

Criar o diretório e ajustar permissões:

```bash
sudo mkdir -p /var/lib/samba/ntp_signd
sudo chown root:_chrony /var/lib/samba/ntp_signd
sudo chmod 750 /var/lib/samba/ntp_signd
sudo systemctl restart samba-ad-dc chrony
```

---

## 17. Gerenciamento com samba-tool

Referência rápida dos comandos `samba-tool` mais usados no dia a dia:

### Usuários

```bash
# Criar usuário
sudo samba-tool user create joao.silva \
    --given-name="João" \
    --surname="Silva" \
    --mail-address="joao.silva@empresa.local" \
    --department="TI" \
    --title="Analista"

# Listar usuários
sudo samba-tool user list

# Ver detalhes de um usuário
sudo samba-tool user show joao.silva

# Desabilitar conta
sudo samba-tool user disable joao.silva

# Habilitar conta
sudo samba-tool user enable joao.silva

# Resetar senha
sudo samba-tool user setpassword joao.silva

# Deletar usuário
sudo samba-tool user delete joao.silva

# Definir que senha não expira
sudo samba-tool user setexpiry joao.silva --noexpiry

# Mover usuário para outra OU
sudo samba-tool computer move "CN=joao.silva,CN=Users,DC=empresa,DC=local" \
    "OU=TI,DC=empresa,DC=local"
```

### Grupos

```bash
# Criar grupo
sudo samba-tool group add "Financeiro" \
    --description="Equipe Financeira"

# Adicionar membro ao grupo
sudo samba-tool group addmembers "Financeiro" joao.silva

# Remover membro do grupo
sudo samba-tool group removemembers "Financeiro" joao.silva

# Listar membros do grupo
sudo samba-tool group listmembers "Financeiro"

# Listar todos os grupos
sudo samba-tool group list
```

### OUs (Organizational Units)

```bash
# Criar OU
sudo samba-tool ou create "OU=TI,DC=empresa,DC=local"
sudo samba-tool ou create "OU=Financeiro,DC=empresa,DC=local"
sudo samba-tool ou create "OU=Servidores,DC=empresa,DC=local"

# Listar OUs
sudo samba-tool ou list
```

### Políticas de senha

```bash
# Ver políticas atuais
sudo samba-tool domain passwordsettings show

# Definir comprimento mínimo de 12 caracteres
sudo samba-tool domain passwordsettings set --min-pwd-length=12

# Definir histórico de senhas (24 senhas anteriores)
sudo samba-tool domain passwordsettings set --history-length=24

# Definir idade máxima (90 dias)
sudo samba-tool domain passwordsettings set --max-pwd-age=90

# Definir lockout após 5 tentativas falhas
sudo samba-tool domain passwordsettings set --account-lockout-threshold=5
sudo samba-tool domain passwordsettings set --account-lockout-duration=30
```

### DNS

```bash
# Adicionar registro A
sudo samba-tool dns add localhost empresa.local \
    servidor01 A 192.168.1.20 -U administrator

# Listar registros de uma zona
sudo samba-tool dns query localhost empresa.local @ ALL -U administrator

# Remover registro
sudo samba-tool dns delete localhost empresa.local \
    servidor01 A 192.168.1.20 -U administrator
```

---

## 18. Troubleshooting

### Porta 53 já em uso

```bash
sudo ss -tlnp | grep ':53'
# Se systemd-resolved aparecer:
sudo systemctl disable --now systemd-resolved
sudo rm /etc/resolv.conf
echo "nameserver 192.168.1.10" | sudo tee /etc/resolv.conf
sudo chattr +i /etc/resolv.conf
sudo systemctl restart samba-ad-dc
```

### Erro: "An existing smb.conf file exists"

```bash
sudo mv /etc/samba/smb.conf /etc/samba/smb.conf.bak
# Re-executar o provisionamento
```

### Erro de Kerberos: "Clock skew too great"

```bash
# Sincronizar horário imediatamente
sudo chronyc makestep
sudo chronyc tracking
# Verificar diferença de horário com cliente Windows
# w32tm /query /status (no Windows)
```

### LDAP retorna "Strong(er) authentication required"

Adicionar ao `[global]` do `smb.conf`:

```ini
ldap server require strong auth = allow_sasl_over_tls
```

```bash
sudo systemctl restart samba-ad-dc
```

### DNS SRV records não encontrados

```bash
# Verificar se o Samba DNS está respondendo
dig @127.0.0.1 _ldap._tcp.empresa.local SRV

# Se falhar, verificar logs do samba
sudo journalctl -u samba-ad-dc --since "10 minutes ago"

# Readicionar registro A do DC manualmente
sudo samba-tool dns add localhost empresa.local \
    dc1 A 192.168.1.10 -U administrator
```

### Verificar logs em tempo real

```bash
sudo journalctl -u samba-ad-dc -f --no-pager
```

### Reiniciar completamente

```bash
sudo systemctl stop samba-ad-dc
sudo systemctl start samba-ad-dc
sudo systemctl status samba-ad-dc
```

---

## Referências

- [Samba Wiki — Setting up Samba as an AD DC](https://wiki.samba.org/index.php/Setting_up_Samba_as_an_Active_Directory_Domain_Controller)
- [Samba Wiki — AD DC Port Usage](https://wiki.samba.org/index.php/Samba_AD_DC_Port_Usage)
- [Samba Wiki — Configuring LDAPS](https://wiki.samba.org/index.php/Configuring_LDAP_over_SSL_(LDAPS)_on_a_Samba_AD_DC)
- [Ubuntu Server Docs — Provisioning a Samba AD DC](https://documentation.ubuntu.com/server/how-to/samba/provision-samba-ad-controller/)
- [Samba GitLab — Código-fonte oficial](https://gitlab.com/samba-team/samba)

---

*Setup Guide v1.0 — Samba4 AD DC em Ubuntu 24.04 LTS — Março 2026*
