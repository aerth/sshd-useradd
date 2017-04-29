package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	log.SetPrefix("")
	log.SetFlags(log.Ltime)
}

func handleEntrypoint(s ssh.Session) {
	// log
	log.Println(s.RemoteAddr().String(), s.Environ(), s.Command())

	// get pubkey or die
	pubkey := s.PublicKey()
	if pubkey == nil {
		s.Write([]byte("Please come back with a public key. Any key will do. To create one, run 'ssh-keygen'\n"))
		goodbye(s)
		return
	}

	pkey := gossh.MarshalAuthorizedKey(pubkey)
	if pkey == nil {
		s.Write([]byte("Please come back with a public key. Any key will do. To create one, run 'ssh-keygen'\n"))
		goodbye(s)
		return
	}

	// get username or die
	username := s.User()
	if username == "" {
		goodbye(s)
		return
	}

	s.Write(([]byte("\033[H\033[2J" + "NEW ACCOUNT\n")))

	// log
	log.Println(username, "status")

	t := terminal.NewTerminal(s, "> ")
	hostname, _ := os.Hostname()

	if hostname != "" {
		t.Write([]byte("\nYou are connected to: " + hostname + "\n\n"))
	}

	pstring := strings.TrimSuffix(string(pkey), "\n")
	t.Write([]byte(pstring + "\n\n"))
	t.Write([]byte("This server has passwordless authentication.\n"))
	t.Write([]byte("The key you used to log in is the only way to enter.\n"))
	t.Write([]byte("Don't lose the key. If this is the wrong key, exit now.\n"))
	t.Write([]byte("\n"))

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

			hosts := decodestatus(getstatus(apistatusurl))
			for hostname, host := range hosts {
				io.WriteString(s, "\n\n"+hostname+"\n")
				io.WriteString(s, host.String())
			}
		case "new":
			io.WriteString(s, fmt.Sprintf("Username: %s\n", username))
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
	}

Done:

	s.Exit(1)
}

func ListHosts(s ssh.Session) {
	hosts := decodestatus(getstatus(apistatusurl))
	io.WriteString(s, fmt.Sprintf("found %v hosts:\n", len(hosts)))
	for name, host := range hosts {
		host.hostname = name
	}
	io.WriteString(s, "Which #! hostname?\nAvailable Hosts:\n\n")
	for name := range hosts {
		io.WriteString(s, fmt.Sprintf("%s ", name))
	}

}
