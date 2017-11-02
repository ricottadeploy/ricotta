package main

import (
	"fmt"
	"os"

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
	deniedAgents   master.AgentStore
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
		masterid := viper.GetString("id")
		fmt.Printf("Master ID: %s\n", masterid)
		generateCertIfNotExist()

		agentsFile := basePath + "/conf/agents.yaml"
		acceptedAgents = master.NewAgentStore()
		acceptedAgents.ReadFromYamlFile(agentsFile)
		fmt.Printf("Accepted agents:\n%s\n", acceptedAgents.ToYaml())
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
