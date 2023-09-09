package ssh

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// SSH stores required SSH configuration.
type SSH struct {
	conf *ssh.ClientConfig
	Host string
	Key  string
	Port string
	User string
}

// NewSSH return a SSH struct.
func NewSSH(host string, port string, user string, keyPath string) (*SSH, error) {
	if err := checkSSHArgs(host, port, user, keyPath); err != nil {
		return nil, fmt.Errorf("check ssh args: %w", err)
	}

	if _, err := os.Stat(keyPath); err != nil {
		return nil, fmt.Errorf("ssh key stat %q: %w", keyPath, err)
	}

	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("read ssh key: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(keyData)
	if err != nil {
		return nil, fmt.Errorf("parse key: %w", err)
	}

	hostKeyCallback, err := knownhosts.New("/Users/" + user + "/.ssh/known_hosts")
	if err != nil {
		return nil, fmt.Errorf("host key callback: %w", err)
	}

	conf := &ssh.ClientConfig{
		User:              user,
		HostKeyCallback:   hostKeyCallback,
		HostKeyAlgorithms: []string{ssh.KeyAlgoED25519},
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	return &SSH{
		conf: conf,
		Host: host,
		Key:  keyPath,
		Port: port,
		User: user,
	}, nil
}

// checkSSHArgs ensures that complete SSH configuration is provided.
func checkSSHArgs(host string, port string, user string, keyPath string) error {
	if host == "" {
		return fmt.Errorf("host value is empty")
	}
	if port == "" {
		return fmt.Errorf("port value is empty")
	}
	if user == "" {
		return fmt.Errorf("user value is empty")
	}
	if keyPath == "" {
		return fmt.Errorf("key path value is empty")
	}

	return nil
}

// CMD executes shell commands over SSH connection.
func (s *SSH) CMD(cmd string) error {
	dest := fmt.Sprintf("%s:%s", s.Host, s.Port)

	conn, err := ssh.Dial("tcp", dest, s.conf)
	if err != nil {
		return fmt.Errorf("ssh dial: %w", err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("new session: %w", err)
	}
	defer session.Close()

	stdOut, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("STDOUT pipe: %w", err)
	}

	stdErr, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("STDERR pipe: %w", err)
	}

	if err := session.Run(cmd); err != nil {
		scannerErr := bufio.NewScanner(stdErr)
		for scannerErr.Scan() {
			fmt.Println(scannerErr.Text())
		}

		return fmt.Errorf("run cmd: %w", err)
	}

	scannerOut := bufio.NewScanner(stdOut)
	for scannerOut.Scan() {
		fmt.Println(scannerOut.Text())
	}

	return nil
}
