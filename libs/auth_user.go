package main

// https://github.com/go-ldap/ldap/blob/master/example_test.go

import (
	"fmt"
	"gopkg.in/ldap.v2"
	"log"
)

func main() {

	Example_userAuthentication("username", "password")

}

func Example_userAuthentication(username, password string) {

	bindusername := "bindusername"
	bindpassword := "bindpassword"

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", 389))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	// First bind with a read only user
	err = l.Bind(bindusername, bindpassword)
	if err != nil {
		log.Fatal(err)
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		"dc=linuxpro,dc=net",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=user)(sAMAccountName=%s))", username),
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	if len(sr.Entries) != 1 {
		log.Fatal("User does not exist")
	} else {
		log.Printf("User %s authenticade", username)
	}

	userdn := sr.Entries[0].DN

	//log.Print(userdn)

	// Bind as the user to verify their password
	err = l.Bind(userdn, password)
	if err != nil {
		log.Fatal(err)
	}

}
