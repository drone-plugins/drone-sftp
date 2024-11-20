// Copyright (c) 2022, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/drone/drone-go/drone"
	"github.com/dustin/go-humanize"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

// Args provides plugin execution arguments.
type (
	Args struct {
		Pipeline

		// Level defines the plugin log level.
		Level string `envconfig:"PLUGIN_LOG_LEVEL"`

		// Skip verification of certificates
		SkipVerify bool `envconfig:"PLUGIN_SKIP_VERIFY"`

		// Plugin specific
		Host        string   `envconfig:"PLUGIN_HOST"`
		Port        int      `envconfig:"PLUGIN_PORT"`
		Username    string   `envconfig:"PLUGIN_USERNAME"`
		Password    string   `envconfig:"PLUGIN_PASSWORD"`
		Key         string   `envconfig:"PLUGIN_KEY"`
		KeyPath     string   `envconfig:"PLUGIN_KEY_PATH"`
		Passphrase  string   `envconfig:"PLUGIN_PASSPHRASE"`
		Files       []string `envconfig:"PLUGIN_FILES"`
		Destination string   `envconfig:"PLUGIN_DESTINATION_PATH"`
	}

	UploadInfo struct {
		Name string `json:"name"`
		Size string `json:"size"`
	}
)

const (
	defaultDirectory = "/"
	defaultPort      = 22
)

var errConfiguration = errors.New("configuration error")

// Exec executes the plugin.
func Exec(args *Args) error {
	err := verifyArgs(args)
	if err != nil {
		return fmt.Errorf("error in the configuration: %w", err)
	}

	files, err := findFileUploads(args)
	if err != nil {
		return fmt.Errorf("could not get the files to upload: %w", err)
	}

	fileCount := len(files)
	if fileCount == 0 {
		return fmt.Errorf("no files were found to upload: %w", errConfiguration)
	}

	logrus.Infof("found %d files to upload\n", fileCount)

	client, err := createSftpClient(args)
	if err != nil {
		return fmt.Errorf("could not connect to sftp server: %w", err)
	}

	err = createDirectory(client, args.Destination)
	if err != nil {
		return fmt.Errorf("unable to create the destination directory %s: %w", args.Destination, err)
	}

	destination := uploadDestination(args)
	logrus.Infof("ready to upload to %s\n", destination)

	uploads := []UploadInfo{}
	for _, file := range files {
		localFile := path.Base(file)
		remoteFile := path.Join(args.Destination, localFile)
		bytes, err := uploadFile(client, file, remoteFile)
		if err != nil {
			return fmt.Errorf("could not upload file %s: %w", file, err)
		}

		logrus.Infof("uploaded file %s\n", localFile)

		uploadInfo := UploadInfo{
			Name: localFile,
			Size: humanize.Bytes(bytes),
		}
		uploads = append(uploads, uploadInfo)
	}

	// Create the card data
	cardData := struct {
		Destination string       `json:"uploadTo"`
		URL         string       `json:"url"`
		Uploads     []UploadInfo `json:"uploads"`
	}{
		Destination: destination,
		URL:         fmt.Sprintf("sftp://%s", destination),
		Uploads:     uploads,
	}

	data, _ := json.Marshal(cardData)
	card := drone.CardInput{
		Schema: "https://drone-plugins.github.io/drone-sftp/card.json",
		Data:   data,
	}
	writeCard(args.Card.Path, &card)

	return nil
}

func verifyArgs(args *Args) error {
	if args.Username == "" {
		return fmt.Errorf("no username provided: %w", errConfiguration)
	}

	if args.Password == "" && args.Key == "" && args.KeyPath == "" {
		return fmt.Errorf("no password or key provided: %w", errConfiguration)
	}

	if args.Host == "" {
		return fmt.Errorf("no hostname provided: %w", errConfiguration)
	}

	if len(args.Files) == 0 {
		return fmt.Errorf("no files to upload provided: %w", errConfiguration)
	}

	if args.Port == 0 {
		args.Port = defaultPort
	}

	if args.Destination == "" {
		args.Destination = defaultDirectory
	}

	return nil
}

func findFileUploads(args *Args) ([]string, error) {
	var files []string
	for _, glob := range args.Files {
		globed, err := filepath.Glob(glob)

		if err != nil {
			return nil, fmt.Errorf("failed to glob %s: %w", glob, err)
		}

		if globed != nil {
			files = append(files, globed...)
		}
	}

	return files, nil
}

func createSftpClient(args *Args) (*sftp.Client, error) {
	authMethods, err := createSftpAuthMethods(args)
	if err != nil {
		return nil, err
	}

	server := fmt.Sprintf("%s:%d", args.Host, args.Port)
	config := &ssh.ClientConfig{
		User:            args.Username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if args.SkipVerify {
		config.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

	conn, err := ssh.Dial("tcp", server, config)
	if err != nil {
		return nil, fmt.Errorf("could not connect to %s: %w", server, err)
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("could not create sftp client for %s: %w", server, err)
	}

	return client, nil
}

func createSftpAuthMethods(args *Args) ([]ssh.AuthMethod, error) {
	var methods []ssh.AuthMethod

	if args.Password != "" {
		methods = append(methods, ssh.Password(args.Password))
	}

	if args.Key != "" {
		signer, err := parsePrivateKey([]byte(args.Key), []byte(args.Passphrase))
		if err != nil {
			return nil, err
		}
		methods = append(methods, ssh.PublicKeys(signer))
	}

	if args.KeyPath != "" {
		buffer, err := os.ReadFile(args.KeyPath)
		if err != nil {
			return nil, fmt.Errorf("could not read private key: %w", err)
		}
		signer, err := parsePrivateKey(buffer, []byte(args.Passphrase))
		if err != nil {
			return nil, err
		}
		methods = append(methods, ssh.PublicKeys(signer))
	}

	if len(methods) > 0 {
		return methods, nil
	} else {
		return nil, fmt.Errorf("could not determinate an sftp auth method")
	}
}

func parsePrivateKey(key []byte, passphrase []byte) (ssh.Signer, error) {
	var signer ssh.Signer
	var err error
	if len(passphrase) > 0 {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(key, passphrase)
	} else {
		signer, err = ssh.ParsePrivateKey(key)
	}
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}
	return signer, nil
}

func createDirectory(client *sftp.Client, directory string) error {
	if directory == defaultDirectory {
		return nil
	}

	stat, err := client.Lstat(directory)
	if err == nil {
		if !stat.IsDir() {
			return fmt.Errorf("path is not a directory %s: %w", directory, errConfiguration)
		}
		return nil
	}

	parent := path.Dir(directory)
	err = createDirectory(client, parent)
	if err != nil {
		return err
	}

	err = client.Mkdir(directory)
	if err != nil {
		return fmt.Errorf("error creating directory %s: %w", directory, err)
	}
	return nil
}

func uploadFile(client *sftp.Client, localFile, remoteFile string) (uint64, error) {
	local, err := os.Open(localFile)
	if err != nil {
		return 0, fmt.Errorf("could not open file %s: %w", localFile, err)
	}

	remote, err := client.Create(remoteFile)
	if err != nil {
		return 0, fmt.Errorf("unable to create destination file %s: %w", remoteFile, err)
	}

	bytes, err := io.Copy(remote, local)
	if err != nil {
		return 0, fmt.Errorf("could not copy to destination file %s: %w", remoteFile, err)
	}

	return uint64(bytes), nil
}

func uploadDestination(args *Args) string {
	var server string
	if args.Port == defaultPort {
		server = args.Host
	} else {
		server = fmt.Sprintf("%s:%d", args.Host, args.Port)
	}

	if args.Destination == defaultDirectory {
		return server
	}

	return fmt.Sprintf("%s%s", server, args.Destination)
}
