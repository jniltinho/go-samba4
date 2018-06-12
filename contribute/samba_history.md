
## Sobre Samba

O Samba é um software livre, licenciado pela GPL (Generic Public License) criado em 1992, que permite o compartilhamento de arquivos/impressão, entre máquinas Windows e Linux, além de outras funções.

O Samba permite o compartilhamento de impressão "segura" e algumas de suas características são estabilidade e velocidade para garantir a total interoperabilidade entre os Sistemas Operacionais Unix e Windows.
Samba é um componente importante para integrar servidores e desktops Linux/Unix em ambientes Active Directory. Ele pode funcionar tanto como um controlador de domínio ou como um membro de domínio.

## Sobre Samba 4

No dia [11 de Dezembro de 2012 Foi lançado o Samba 4.0](https://www.samba.org/samba/history/samba-4.0.0.html), em desenvolvimento desde 2006, o primeiro release estável desta série com grandes aprimoramentos.

O Samba 4 trará novas melhorias e de acordo com a desenvolvedora, uma maior flexibilidade em sua configuração. O programa já estava em fase beta há um bom tempo, porém, é a primeira vez que ele ganha uma data de lançamento concreta.

## Alguns recursos do Servidor Samba 4

Ele possui recursos como um servidor de DNS dinâmico (que pode ser implementado a partir de um servidor DNS privado ou por meio do plugin BIND), um servidor de diretório LDAP, recursos para a implementação de um active directory e um servidor de autenticação Kerberos.

É possível criar o compartilhamento de arquivos e impressoras para ambientes que possuem máquinas Windows e Linux.

Além disso, todos os arquivos para a criação de um Active Directory Domain Controller que seja compatível com as versões do Windows 2000 (2003, 2008, 2012) foram disponibilizados nessa versão.

A compatibilidade com o Active Directory foi possível graças à documentação oficial e os testes de interoperabilidade da Microsoft. (Ano de 2008)

Nesse sentido, a empresa atuou lado a lado com os desenvolvedores do projeto Samba para garantir que o software fosse totalmente compatível com dispositivos Windows, eliminando a necessidade de a comunidade utilizar sniffers em rede e engenharia reversa para o Samba ser funcional.

Os aplicativos que levaram à criação do Samba foram originalmente criados por Andrew Tridgell. O projeto foi iniciado em 1991, quando Andrew precisou criar um protocolo de rede próprio
para trabalhar com o programa Pathworks.

Diante da ausência de documentação sobre o SMB, Andrew realizou a engenharia reversa do protocolo e publicou o resultado do seu trabalho na internet com o nome Samba.

Para quem não conhece, Samba é um software Servidor para Linux e outros sistemas unix. Sua função principal é gerenciar grupos e usuários, compartilhamento de recursos em redes formadas por PCs com o Windows. Assim, o usuário que utilizar o Samba poderá manter o Linux como servidor de arquivos.

## Cenários de implantação do Samba 4

O Samba 4 será a solução para quem precisa de controladores de domínios e ou Servidores de arquivos para o seguinte cenários:

- Como um servidor único sendo Domain Controller e Servidor arquivos em um mesmo servidor.
- Um controlador de Domínio Principal abrigando apenas essa função.
- Como servidor de arquivo totalmente dedicado para essa carga de trabalho
- Como um servidor de arquivos em ambiente que você já tenha um Domain Controller Windows Server e gostaria de servidores de arquivos adicionais, para balancear a carga, sem preocupar com custo com licenças.
- Como Domain Controllers adicionais e para outro Samba 4 ou Windows Server.
- Como Servidores de arquivos para namespace DFS Windows.
- Como Domain Controller Somente leitura para as filiais que não tem a mesma segurança física e logica que tem na matriz.

Esses são apenas alguns exemplos de uso do Samba 4, as possibilidades são enormes, podendo até coloca-lo em um ambiente de alta disponibilidade.

## 20 razões para usar o Samba 4

- É possível criar um AD Completo com Samba 4.
- É possível criar um controlador de domínio Principal.
- É possível criar um Controlador de domínio somente leitura (RODC).
- É possível criar um Controlador de domínio Adicional.
- Pode ser administrado usando interface Gráfica do próprio Windows, como usuários e computadores do Active Directory.
- Posso Migrar de forma fácil de AD Windows para um AD Linux e vice versa.
- É possível trabalhar com perfil móvel.
- Trabalhar com Pasta Base.
- Lixeira de Servidor de Arquivos (Tipo copia de Sombra).
- Auditoria de Acesso.
- Trabalhar com permissões como a do Windows.
- Trabalhar com GPO.
- Fazer replicação de Servidores  (TIPO DFS).
- Trabalhar com dados em camadas.
- Triagem de Arquivos  (Proibir gravação de arquivos pela extensão).
- É software Livre não precisa de licença.
- Não precisa de CALs de acesso para as estações.
- Posso fazer o SAMBA 4 trabalhar como controlador de domínio adicional do Windows server e vice versa.
- Já vem com DNS, kerberos, LDAP integrado.
- Posso fazer a integração do Samba 4 com o proxy Squid, pfSense e etc.

## Sobre o Go-Samba4

O projeto Go-Samba4 tem como solução a gestão de usuários e grupos no ambiente server do sistema Samba 4, com uma interface web simples, leve e intuitiva, desenvolvida para promover um ambiente controlado onde o administrador possa gerenciar um sistema que antes só podia ser gerenciado por linha de comando usando o samba-tool no ambiente Linux ou por uma interface Desktop Windows.

Esse trabalho apresenta um projeto que vai criar uma solução web desenvolvida na linguagem Python usando o framework Flask para gestão de usuários no software Samba 4, preocupando-se em criar algo que seja simples e intuitivo, de fácil execução em um servidor com o software Samba na versão 4.7.0 ou superior instalado.

É um projeto open source (aberto), com seu código fonte disponível no [GitHub](https://github.com/jniltinho/go-samba4) possibilite que os desenvolvedores possam baixar, usar e contribuir de alguma forma na melhoria e crescimento do projeto.

## Links

https://www.baboo.com.br/arquivo/internet/samba-4-sera-lancado-em-11-de-dezembro/

http://blog.astreinamentos.com.br/2015/08/20-razoes-porque-eu-uso-samba-4-ao-inves-de-ad-e-file-server-pirata.html

http://e-tinet.com/linux/por-que-samba-4-como-servidor/

https://blogs.technet.microsoft.com/openchoice/2008/03/03/testemunhos-sobre-o-anuncio-da-semana-passada-sobre-interoperabilidade/

http://www.theregister.co.uk/2008/02/21/microsoft_goes_open/