package main

import (
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"os"
)

var version = "0.0.3"
var motd = "Hello!\n\n"
func init() {

}

func authPublicKey(ctx ssh.Context, key ssh.PublicKey) bool {
	return true // always return true, keeping key
}

func authPassword(ctx ssh.Context, password string) bool {
	return true // always return true
}

func aauthKeyboardInteractive(ctx ssh.Context, challenge gossh.KeyboardInteractiveChallenge) bool {
	return true // always return true,
}
func authKeyboardInteractive(ctx ssh.Context, challenge gossh.KeyboardInteractiveChallenge) bool {
	ans, err := challenge("user",
		"instruction\n",
		[]string{"What color is grass? ALL CAPS\n", "What color is sky? ALL CAPS\n"},
		[]bool{true, true})
	if err != nil {
		log.Println(err)
		return false
	}
	ok := ans[0] == "GREEN" && ans[1] == "BLUE"
	if ok {
		challenge("user", motd, nil, nil)
		return true
	}
	return false
}

var DefaultServer = ssh.Server{
	Addr:                       "0.0.0.0:4444",
	Handler:                    handleEntrypoint,
	PublicKeyHandler:           authPublicKey,
	PasswordHandler:            authPassword,
	KeyboardInteractiveHandler: authKeyboardInteractive,
}

func init() {
	err := DefaultServer.SetOption(ssh.HostKeyFile("key.pem"))
	if err != nil {
		println(err.Error())
		os.Exit(111)
	}
}

func main() {

	if len(os.Args) != 2 {
		println("need port to listen on")
		os.Exit(111)
	}

	port := os.Args[1]

	println("starting ssh server on port:", port)
	listener, err := net.Listen("tcp", DefaultServer.Addr)
	if err != nil {
		println(err.Error())
		os.Exit(111)
	}
	err = DefaultServer.Serve(listener)
	if err != nil {
		println(err.Error())
		os.Exit(111)
	}

}

func goodbye(s ssh.Session) {
	io.WriteString(s, "Goodbye! 8-)\n")
	s.Exit(1)
}
