package main

import (
	"flag"
	"log"
	"net"
	"path"
	"strings"

	"github.com/hpcloud/tail"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var (
	host      = flag.String("h", "127.0.0.1:22", "remoteHost:port")
	user      = flag.String("u", "root:123456", "user:password")
	localFile = flag.String("f", "", "local file path")
	remoteDir = flag.String("d", "/tmp", "remote dir")
)

// 建立ssh隧道
func sshConnect(h, u string) *ssh.Client {
	log.Printf("Connecting to remote server %s\n", h)

	auth := strings.Split(u, ":")
	user, passwd := auth[0], auth[1]
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(passwd),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", h, config)
	if err != nil {
		log.Fatalf("Connect to server fail, %s\n", err.Error())
	}

	return client
}

// 建立sftp客户端
func sftpClient(sshClient *ssh.Client) *sftp.Client {
	log.Printf("Creating sftp file client\n")

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		log.Fatalf("Sftp client error, %s\n", err.Error())
	}

	return sftpClient
}

func main() {
	flag.Parse()

	sshClient := sshConnect(*host, *user)
	defer sshClient.Close()

	sftpClient := sftpClient(sshClient)
	defer sftpClient.Close()

	remoteDir := *remoteDir
	remoteFileName := path.Base(*localFile)

	dstFile, err := sftpClient.Create(path.Join(remoteDir, remoteFileName))
	if err != nil {
		log.Fatalf("Create remote file fail, %s\n", err.Error())
	}
	defer dstFile.Close()

	t, err := tail.TailFile(*localFile, tail.Config{Follow: true})
	log.Println("Starting pushing ..............................")
	for line := range t.Lines {
		dstFile.Write([]byte(line.Text + "\n"))
	}
}
