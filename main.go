package main

import (
	"io"
	"log"
	"fmt"
	"os"
	"time"
	"strings"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
	term "golang.org/x/crypto/ssh/terminal"
)

var version = "0.0.2"

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
	fmt.Fprintln(os.Stderr, time.Now().String(), s.RemoteAddr().String(), s.Environ(), s.Command(), "\n")

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

	
	s.Write([]byte("Creating account on sf1.hashbang.sh\n"))
	b := getstatus("https://hashbang.sh/server/stats")

	// log
	fmt.Fprintln(os.Stderr, time.Now().String(), username, "status")
	io.WriteString(s, fmt.Sprintln(decodestatus(b)))
	<-time.After(3*time.Second)
	
	//
	// send pubkey and username to API
	//

	t := term.NewTerminal(s, "> ")

	oldState, err := term.MakeRaw(0)
	if err != nil {
	        panic(err)
	        s.Exit(1)
	}
	defer term.Restore(0, oldState)


	
	pstring := strings.TrimSuffix(string(pkey), "\n")
	shellhost, err :=  t.ReadLine()
	if err != nil {
		log.Println(err.Error())
		goodbye(s)
		s.Exit(1)
	}
	println(username,pstring,shellhost)

	
//	resp := newuser(username, pstring, hostname)
	resp := "net disabled"
	// print reply
	fmt.Println(time.Now().String(), username, resp)
	// tell user response
	io.WriteString(s, resp+"\n")

	

	s.Exit(1)
}
