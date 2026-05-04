package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/SehAie/VPS-Builder/internal/config"
	"github.com/SehAie/VPS-Builder/internal/deploy"
	"github.com/SehAie/VPS-Builder/internal/output"
	"github.com/SehAie/VPS-Builder/internal/wizard"
	"golang.org/x/term"
)

const version = "0.1.3-windows"

func main() {
	if runtime.GOOS != "windows" {
		fmt.Fprintln(os.Stderr, "Error: this build is intended to run locally on Windows 10/11 only. The remote VPS can still be Debian/Ubuntu.")
		os.Exit(1)
	}
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "init":
		cmdInit(os.Args[2:])
	case "deploy":
		cmdDeploy(os.Args[2:])
	case "gen-client":
		cmdGenClient(os.Args[2:])
	case "version":
		fmt.Println("VPS-Builder", version)
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Print(`VPS-Builder - free open-source VPS bootstrap + Hysteria 2 deployer

Usage on Windows PowerShell:
  .\VPS-Builder.exe init -o config.yaml
  .\VPS-Builder.exe deploy -c config.yaml --ask-password
  .\VPS-Builder.exe deploy -c config.yaml --dry-run --ask-password
  .\VPS-Builder.exe gen-client -c config.yaml -p HYSTERIA_PASSWORD
  .\VPS-Builder.exe version

Environment:
  VPS_BUILDER_ROOT_PASSWORD    Initial VPS SSH password. Prefer --ask-password.
`)
}

func cmdInit(args []string) {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	out := fs.String("o", "config.yaml", "output config file")
	_ = fs.Parse(args)

	cfg, err := wizard.Run()
	fatalIf(err)
	fatalIf(config.Save(*out, cfg))
	fmt.Println("Config file generated:", *out)
	fmt.Println("Next step: .\\VPS-Builder.exe deploy -c", *out, "--ask-password")
}

func cmdDeploy(args []string) {
	fs := flag.NewFlagSet("deploy", flag.ExitOnError)
	cfgPath := fs.String("c", "config.yaml", "config file")
	dryRun := fs.Bool("dry-run", false, "print plan without changing server")
	yes := fs.Bool("yes", false, "non-interactive confirm")
	askPassword := fs.Bool("ask-password", false, "ask initial SSH password securely")
	_ = fs.Parse(args)

	cfg, err := config.Load(*cfgPath)
	fatalIf(err)

	pw := os.Getenv("VPS_BUILDER_ROOT_PASSWORD")
	if *askPassword {
		fmt.Print("Enter current VPS SSH/root password: ")
		b, _ := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		pw = strings.TrimSpace(string(b))
	}

	if !*yes && !*dryRun {
		fmt.Printf("Ready to deploy %s: SSH TCP %d, Hysteria UDP %d. Type yes to continue: ", cfg.Server.IP, cfg.SSH.NewPort, cfg.Hysteria.ListenPort)
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Canceled.")
			return
		}
	}

	res, err := deploy.Run(cfg, deploy.Options{DryRun: *dryRun, AssumeYes: *yes, RootPassword: pw})
	fatalIf(err)

	fmt.Println("Deployment completed.")
	fmt.Println("SSH private key:", res.PrivateKeyPath)
	fmt.Println("SSH public key:", res.PublicKeyPath)
	if res.ClientConfig != "" {
		fmt.Println("Hysteria client config:", res.ClientConfig)
	}
	if res.HysteriaPass != "" {
		fmt.Println("Hysteria password has been written to the local client config. Keep it safe.")
	}
}

func cmdGenClient(args []string) {
	fs := flag.NewFlagSet("gen-client", flag.ExitOnError)
	cfgPath := fs.String("c", "config.yaml", "config file")
	password := fs.String("p", "", "hysteria password")
	_ = fs.Parse(args)

	cfg, err := config.Load(*cfgPath)
	fatalIf(err)

	pw := *password
	if pw == "" && cfg.Hysteria.Password != "auto" {
		pw = cfg.Hysteria.Password
	}
	if pw == "" || pw == "auto" {
		fatalIf(fmt.Errorf("missing hysteria password, pass -p"))
	}

	path, err := output.WriteHysteriaClient(cfg, pw)
	fatalIf(err)
	fmt.Println("Client config generated:", path)
}

func fatalIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
