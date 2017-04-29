package main

import (
	"fmt"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
	term "golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var version = "0.0.3"

func init() {

}

func authPublicKey (ctx ssh.Context, key ssh.PublicKey) bool {
	return true // always return true, keeping key
}

func authKeyboardInteractive (ctx ssh.Context, password string) bool {
	return true // always return true, 
}

var DefaultServer = ssh.Server {
	Addr: "0.0.0.0:4444",
	Handler: handleEntrypoint,
	PublicKeyHandler: handlePublicKey,
	PasswordHandler: handlePassword,
}

func init(){
	ssh.DefaultServer.ssh.HostKeyFile("key.pem")
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
		ssh.PublicKeyAuth(),
		
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
	if pkey == nil {
		goodbye(s)
		return
	}

	// get username or die
	username := s.User()
	if username == "" {
		goodbye(s)
		return
	}

	s.Write([]byte("Creating account on sf1.hashbang.sh\n"))
	
	// log
	fmt.Fprintln(os.Stderr, time.Now().String(), username, "status")
	hosts := decodestatus(getstatus("https://hashbang.sh/server/stats"))
	io.WriteString(s, fmt.Sprintf("found %v hosts:\n", len(hosts)))

	for name, host := range hosts {
		host.hostname = name
	}

	//
	// send pubkey and username to API
	//

	t := term.NewTerminal(s, "> ")
	hostname, _ := os.Hostname()

	if hostname != "" {
		t.Write([]byte("\nYou are connected to: " + hostname + "\n\n"))
	}
	pstring := strings.TrimSuffix(string(pkey), "\n")

	var input, resp string
	var err error

	for {
		input, err = t.ReadLine()
		if err != nil {
			if err != io.EOF {
				log.Println(err.Error())
			}
			goodbye(s)
			s.Exit(1)
		}
		cmd := strings.TrimSuffix(input, "\n")
		switch cmd {
		case "status":
		
			hosts := decodestatus(getstatus("https://hashbang.sh/server/stats"))
			for hostname, host := range hosts {
				io.WriteString(s, "\n\n"+hostname+"\n")
				io.WriteString(s, host.String())
			}
		case "new":
			io.WriteString(s, fmt.Sprintf("Username: %s\n", username))
			io.WriteString(s, "Creating a new #! account..\n")
			hosts := decodestatus(getstatus("https://hashbang.sh/server/stats"))
			io.WriteString(s, "Which #! hostname?\nAvailable Hosts:\n\n")
			for name  := range hosts {
				io.WriteString(s, fmt.Sprintf("%s ", name))
			}
			io.WriteString(s, "\n\n")
			input, err = t.ReadLine()
			if err != nil {
				log.Println(err.Error())
				goodbye(s)
				s.Exit(1)
			}
			resp = newuser(username, pstring, hostname)
		case "exit", "", "EOF", "EOF\n":
			goto Done
		case "help":
			io.WriteString(s, "Commands: 'new', 'exit', 'help'\n")
		default:
			resp = fmt.Sprintf("error: command %q not found", cmd)
		}

		// print reply
		log.Println(username, resp)
		io.WriteString(s, resp+"\n")
		if resp == "success" {
			io.WriteString(s, fmt.Sprintf("now log in lik this: ssh %s@%s\n", username, hostname))
		}

	}

Done:

	s.Exit(1)
}
