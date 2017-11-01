package main

import (
	"fmt"
	"os"

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
	fmt.Printf("SHA1 Fingerprint: %s\n", fingerprint)
}
