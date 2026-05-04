package wizard

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/SehAie/VPS-Builder/internal/config"
	"golang.org/x/term"
)

func Run() (config.Config, error) {
	cfg := config.Default()
	r := bufio.NewReader(os.Stdin)

	fmt.Println("vps-auto-kit Windows init wizard")
	fmt.Println("Tip: the initial SSH/root password is not saved into the config file. Use --ask-password during deploy.")
	fmt.Println()

	cfg.Server.IP = prompt(r, "VPS IP", cfg.Server.IP)
	cfg.Server.User = prompt(r, "SSH user", cfg.Server.User)
	cfg.Server.OldPort = promptInt(r, "Current SSH port", cfg.Server.OldPort)
	cfg.SSH.NewPort = promptInt(r, "New SSH port", cfg.SSH.NewPort)
	cfg.SSH.KeyName = prompt(r, "Local SSH key name", cfg.SSH.KeyName)
	cfg.SSH.DisablePasswordLogin = promptBool(r, "Disable SSH password login after key login is verified", cfg.SSH.DisablePasswordLogin)
	cfg.System.UpdatePackages = promptBool(r, "Update system packages", cfg.System.UpdatePackages)
	cfg.System.OptimizeSysctl = promptBool(r, "Apply BBR/sysctl network optimization", cfg.System.OptimizeSysctl)
	cfg.System.EnableBBR = cfg.System.OptimizeSysctl && promptBool(r, "Enable BBR", cfg.System.EnableBBR)
	cfg.System.DisableIPv6 = promptBool(r, "Disable IPv6 on the VPS", cfg.System.DisableIPv6)
	cfg.Hysteria.Enabled = promptBool(r, "Install Hysteria 2", cfg.Hysteria.Enabled)

	if cfg.Hysteria.Enabled {
		cfg.Hysteria.ListenPort = promptInt(r, "Hysteria 2 UDP port", cfg.Hysteria.ListenPort)
		mode := prompt(r, "TLS mode, MVP supports self-signed", cfg.Hysteria.TLS.Mode)
		cfg.Hysteria.TLS.Mode = strings.TrimSpace(mode)
		cfg.Hysteria.TLS.CommonName = prompt(r, "Self-signed certificate CN", cfg.Hysteria.TLS.CommonName)

		fmt.Print("Hysteria 2 password, leave empty to auto-generate: ")
		pw, _ := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if strings.TrimSpace(string(pw)) == "" {
			cfg.Hysteria.Password = "auto"
		} else {
			cfg.Hysteria.Password = strings.TrimSpace(string(pw))
		}
	}

	return cfg, cfg.Validate()
}

func prompt(r *bufio.Reader, label, def string) string {
	if def == "" {
		fmt.Printf("%s: ", label)
	} else {
		fmt.Printf("%s [%s]: ", label, def)
	}
	v, _ := r.ReadString('\n')
	v = strings.TrimSpace(v)
	if v == "" {
		return def
	}
	return v
}

func promptInt(r *bufio.Reader, label string, def int) int {
	for {
		v := prompt(r, label, strconv.Itoa(def))
		n, err := strconv.Atoi(v)
		if err == nil && n >= 1 && n <= 65535 {
			return n
		}
		fmt.Println("Please enter a port number between 1 and 65535.")
	}
}

func promptBool(r *bufio.Reader, label string, def bool) bool {
	defText := "n"
	if def {
		defText = "y"
	}
	for {
		v := strings.ToLower(prompt(r, label+" y/n", defText))
		switch v {
		case "y", "yes", "true", "1":
			return true
		case "n", "no", "false", "0":
			return false
		default:
			fmt.Println("Please enter y or n.")
		}
	}
}
