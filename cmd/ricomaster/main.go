package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"github.com/ricottadeploy/ricotta/comms"

	"github.com/ricottadeploy/ricotta/master"
	"github.com/ricottadeploy/ricotta/security"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	basePath       string
	cfgFile        string
	certFile       string
	keyFile        string
	acceptedAgents master.AgentStore
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringVar(&basePath, "path", "C:/ricotta/master", "")
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
		fmt.Println("Ricotta Master")
		generateCertIfNotExist()
		readAgentsFile()
		listen()
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

func readAgentsFile() {
	agentsFile := basePath + "/conf/agents.yaml"
	acceptedAgents = master.NewAgentStore()
	acceptedAgents.ReadFromYamlFile(agentsFile)
	fmt.Printf("Accepted agents:\n%s\n", acceptedAgents.ToYaml())
}

func listen() {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal("Error loading certificate. ", err)
	}

	caCert, err := ioutil.ReadFile(certFile)
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsCfg := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		ClientAuth:         tls.RequireAnyClientCert,
		InsecureSkipVerify: true,
	}
	tlsCfg.BuildNameToCertificate()
	listenAddr := viper.GetString("listen_address")
	listener, err := tls.Listen("tcp4", listenAddr, tlsCfg)
	fmt.Printf("Listening at %s\n", listenAddr)

	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	c := conn.(*tls.Conn)
	fingPrint := security.GetPeerFingerprint(c)
	fmt.Printf("Connection request from agent with fingerprint %s: ", fingPrint)
	_, found := acceptedAgents.Get(security.Fingerprint(fingPrint))
	if !found {
		fmt.Printf("DENIED\n")
		return
	}
	fmt.Printf("ACCEPTED\n")
	cc := comms.NewConn(conn)
	for {
		cc.Write([]byte("How"))
		time.Sleep(2 * time.Second)
	}
}
