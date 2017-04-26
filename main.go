package main

import (
	"fmt"
	"io"
	"os"
	"time"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

var version = "0.0.1"

func init() {
	
}

func main() {

	if len(os.Args) != 2 {
		println("need port to listen on")
		os.Exit(111)
	}

	port := os.Args[1]
	println("starting ssh server on port:", port)
	err := ssh.ListenAndServe(

		// interface+port
		"0.0.0.0:"+port,

		// ssh handler
		newuserhandler,

		// pubkey always true option
		ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true // always return true, keeping key
		}),

		// host key option -> $HOME/.ssh/id_sshd
	//	ssh.HostKeyFile(os.Getenv("HOME")+"/.ssh/id_sshd"),
		ssh.HostKeyFile("id_sshd"),
	)

	// only fatal errors
	if err != nil {
		println(err.Error())
		os.Exit(111)
	}
}

func goodbye(s ssh.Session) {
	io.WriteString(s, "Goodbye! 8-)\n")
	s.Exit(1)
}

func newuserhandler(s ssh.Session) {
	// log
	println(time.Now().String(), s.RemoteAddr().String(), s.Environ(), s.Command(),"\n")

	// get pubkey or die
	pubkey := s.PublicKey()
	if pubkey == nil {
		goodbye(s)
		return
	}
	pkey := gossh.MarshalAuthorizedKey(pubkey)
	// get username or die
	username := s.User()
	if username == "" {
		goodbye(s)
		return
	}



	//
	// send pubkey and username to API
	//
	newuser(username, string(pkey), "sf1.hashbang.sh")

	// tell user some info

	io.WriteString(s, fmt.Sprintf("public key used by %s:\n", s.User()))
	s.Write(pkey)

	io.WriteString(s, "Hello, Goodbye!\n")
	s.Exit(1)
}
