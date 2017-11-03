package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/ricottadeploy/ricotta/comms"
	"github.com/ricottadeploy/ricotta/security"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	basePath string
	cfgFile  string
	certFile string
	keyFile  string
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringVar(&basePath, "path", "C:/ricotta/agent", "")
	cfgFile = basePath + "/conf/config.yaml"
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		viper.ReadInConfig()
	}
}

var rootCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Ricotta Agent")
		generateCertIfNotExist()
		for {
			c, conerr := connect()
			if conerr != nil {
				fmt.Printf("Connection Error: %s\n", conerr)
			} else {
				err := listen(c)
				fmt.Printf("Communication Error: %s\n", err)
			}
			time.Sleep(1 * time.Second)
		}
	},
}

func generateCertIfNotExist() {
	certsPath := basePath + "/certs"
	certFile = certsPath + "/certificate.pem"
	keyFile = certsPath + "/private.pem"
	created := security.GenerateDefaultX509SelfSignedCertIfNotExist(certFile, keyFile)
	if created {
		fmt.Println("Generated certificate")
	}
	fingerprint := security.GetX509CertSHA1Fingerprint(certFile)
	fmt.Printf("Fingerprint: %s\n", fingerprint)
}

func connect() (net.Conn, error) {
	masterAddr := viper.GetString("master.address")
	masterFingerprint := viper.GetString("master.fingerprint")
	fmt.Printf("Connecting to master at %s\n", masterAddr)

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal("Error loading certificate. ", err)
	}
	tlsCfg := &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
	}
	tlsCfg.BuildNameToCertificate()
	conn, err := tls.Dial("tcp4", masterAddr, tlsCfg)
	if err != nil {
		return conn, err
	}

	fingPrint := security.GetPeerFingerprint(conn)
	if fingPrint != masterFingerprint {
		log.Fatalf("Trusted master fingerprint is: %s\nFingerprint of master at %s is: %s\nExiting...", masterFingerprint, masterAddr, fingPrint)
	}
	fmt.Printf("Master fingerprint verified. Connection successful.\n")
	return conn, nil
}

func listen(conn net.Conn) error {
	fmt.Println("Listening to commands from master")
	cc := comms.NewConn(conn)
	for {
		msgFrame, err := cc.ReadMsgFrame()
		if err != nil {
			return err
		}
		fmt.Printf("Received Message Type: %d\n", msgFrame.Type)
		if msgFrame.Type == comms.MsgTypeInfo {
			msg, err := comms.ToInfoMessage(msgFrame.Data)
			if err != nil {
				return err
			}
			fmt.Println(msg.Text)
		}
	}
}
