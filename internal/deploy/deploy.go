package deploy

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/SehAie/VPS-Builder/internal/config"
	"github.com/SehAie/VPS-Builder/internal/keygen"
	"github.com/SehAie/VPS-Builder/internal/output"
	"github.com/SehAie/VPS-Builder/internal/sshx"
)

type Options struct {
	DryRun       bool
	AssumeYes    bool
	RootPassword string
}

type Result struct {
	PrivateKeyPath string
	PublicKeyPath  string
	ClientConfig   string
	HysteriaPass   string
}

func Run(cfg config.Config, opt Options) (Result, error) {
	pair, err := keygen.EnsureEd25519(cfg.KeyPath())
	if err != nil {
		return Result{}, fmt.Errorf("ensure ssh key: %w", err)
	}

	password := cfg.Server.Password
	if opt.RootPassword != "" {
		password = opt.RootPassword
	}
	if strings.TrimSpace(password) == "" {
		return Result{}, fmt.Errorf("missing root password: pass --ask-password or set VPS_BUILDER_ROOT_PASSWORD")
	}

	hyPass := cfg.Hysteria.Password
	if hyPass == "" || hyPass == "auto" {
		hyPass, err = randomPassword(32)
		if err != nil {
			return Result{}, err
		}
	}

	if opt.DryRun {
		fmt.Println("[dry-run] would generate or reuse SSH key:", pair.PublicPath)
		fmt.Println("[dry-run] would configure SSH port:", cfg.SSH.NewPort)
		if cfg.Hysteria.Enabled {
			fmt.Println("[dry-run] would configure Hysteria UDP port:", cfg.Hysteria.ListenPort)
		}
		return Result{PrivateKeyPath: pair.PrivatePath, PublicKeyPath: pair.PublicPath, HysteriaPass: hyPass}, nil
	}

	fmt.Println("[1/7] Connecting to VPS with current SSH port...")
	initial := sshx.NewPassword(cfg.Server.User, cfg.Server.IP, cfg.Server.OldPort, password, 20*time.Second)
	if _, err := initial.Run("echo ok"); err != nil {
		return Result{}, fmt.Errorf("initial ssh connection failed: %w", err)
	}

	fmt.Println("[2/7] Running preflight checks and package setup...")
	if out, err := initial.Run(preflightScript(cfg)); err != nil {
		return Result{}, fmt.Errorf("preflight failed: %w\n%s", err, out)
	}

	fmt.Println("[3/7] Uploading public key and writing SSH hardening config...")
	sshScript := sshHardeningScript(cfg, strings.TrimSpace(string(pair.Authorized)))
	if out, err := initial.Run(sshScript); err != nil {
		return Result{}, fmt.Errorf("ssh hardening failed: %w\n%s", err, out)
	}

	fmt.Println("[4/7] Verifying new SSH port and key login...")
	keyClient := sshx.NewKey(cfg.Server.User, cfg.Server.IP, cfg.SSH.NewPort, pair.Signer, 20*time.Second)
	var verifyErr error
	for i := 0; i < 10; i++ {
		if err := keyClient.Ping(); err == nil {
			verifyErr = nil
			break
		} else {
			verifyErr = err
			time.Sleep(2 * time.Second)
		}
	}
	if verifyErr != nil {
		_, _ = initial.Run(rollbackSSHScript())
		return Result{}, fmt.Errorf("new SSH login failed, rollback attempted: %w", verifyErr)
	}

	if cfg.System.OptimizeSysctl {
		fmt.Println("[5/7] Applying sysctl network optimization...")
		if out, err := keyClient.Run(sysctlScript(cfg)); err != nil {
			return Result{}, fmt.Errorf("sysctl optimize failed: %w\n%s", err, out)
		}
	} else {
		fmt.Println("[5/7] Skipping sysctl network optimization.")
	}

	if cfg.Hysteria.Enabled {
		fmt.Println("[6/7] Installing and configuring Hysteria 2...")
		if out, err := keyClient.Run(hysteriaScript(cfg, hyPass)); err != nil {
			return Result{}, fmt.Errorf("hysteria deploy failed: %w\n%s", err, out)
		}
	} else {
		fmt.Println("[6/7] Skipping Hysteria 2.")
	}

	fmt.Println("[7/7] Writing local client config and summary...")
	clientPath := ""
	if cfg.Hysteria.Enabled && cfg.Hysteria.GenerateClientConfig {
		clientPath, err = output.WriteHysteriaClient(cfg, hyPass)
		if err != nil {
			return Result{}, err
		}
	}
	_ = output.WriteSummary(cfg, hyPass)

	return Result{PrivateKeyPath: pair.PrivatePath, PublicKeyPath: pair.PublicPath, ClientConfig: clientPath, HysteriaPass: hyPass}, nil
}

func randomPassword(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "="), nil
}

func preflightScript(cfg config.Config) string {
	update := ""
	if cfg.System.UpdatePackages {
		update = "apt-get update -y && DEBIAN_FRONTEND=noninteractive apt-get upgrade -y && DEBIAN_FRONTEND=noninteractive apt-get install -y curl ca-certificates openssl iproute2 iptables openssh-server"
	} else {
		update = "DEBIAN_FRONTEND=noninteractive apt-get install -y curl ca-certificates openssl iproute2 iptables openssh-server || true"
	}
	return fmt.Sprintf(`set -e
if [ "$(id -u)" != "0" ]; then echo "Please run as root user"; exit 1; fi
. /etc/os-release || true
case "${ID:-}" in debian|ubuntu) echo "OS: ${PRETTY_NAME:-$ID}" ;; *) echo "Only Debian/Ubuntu are supported in MVP. Current ID=${ID:-unknown}"; exit 1 ;; esac
command -v apt-get >/dev/null 2>&1 || { echo "apt-get not found"; exit 1; }
%s
`, update)
}

func sshHardeningScript(cfg config.Config, authorizedKey string) string {
	passwordAuth := "yes"
	if cfg.SSH.DisablePasswordLogin {
		passwordAuth = "no"
	}
	return fmt.Sprintf(`set -e
install -d -m 700 /root/.ssh
touch /root/.ssh/authorized_keys
chmod 600 /root/.ssh/authorized_keys
grep -qxF %[1]q /root/.ssh/authorized_keys || echo %[1]q >> /root/.ssh/authorized_keys
install -d -m 755 /etc/ssh/sshd_config.d
if [ -f /etc/ssh/sshd_config.d/VPS-Builder.conf ]; then cp /etc/ssh/sshd_config.d/VPS-Builder.conf /etc/ssh/sshd_config.d/VPS-Builder.conf.bak.$(date +%%s); fi
cat > /etc/ssh/sshd_config.d/VPS-Builder.conf <<'REMOTE_EOF'
Port %[2]d
PermitRootLogin %[3]s
PasswordAuthentication %[4]s
PubkeyAuthentication yes
MaxAuthTries %[5]d
MaxStartups 50:30:100
LoginGraceTime 30
X11Forwarding no
ClientAliveInterval 120
ClientAliveCountMax 3
REMOTE_EOF
iptables -C INPUT -p tcp --dport %[2]d -j ACCEPT 2>/dev/null || iptables -I INPUT -p tcp --dport %[2]d -j ACCEPT
SSHD_BIN=$(command -v sshd || echo /usr/sbin/sshd)
$SSHD_BIN -t
systemctl restart ssh || systemctl restart sshd
`, authorizedKey, cfg.SSH.NewPort, cfg.SSH.PermitRootLogin, passwordAuth, cfg.SSH.MaxAuthTries)
}

func rollbackSSHScript() string {
	return `set +e
rm -f /etc/ssh/sshd_config.d/VPS-Builder.conf
systemctl restart ssh || systemctl restart sshd
`
}

func sysctlScript(cfg config.Config) string {
	ipv6 := ""
	if cfg.System.DisableIPv6 {
		ipv6 = `
# === Disable IPv6 ===
net.ipv6.conf.all.disable_ipv6 = 1
net.ipv6.conf.default.disable_ipv6 = 1
`
	}
	bbr := ""
	if cfg.System.EnableBBR {
		bbr = `
# === BBR ===
net.core.default_qdisc = fq
net.ipv4.tcp_congestion_control = bbr
`
	}
	return fmt.Sprintf(`set -e
cat > /etc/sysctl.d/99-VPS-Builder.conf <<'REMOTE_EOF'
%[1]s
# === Buffers ===
net.core.rmem_max = 16777216
net.core.wmem_max = 16777216
net.core.rmem_default = 1048576
net.core.wmem_default = 1048576
net.ipv4.tcp_rmem = 4096 87380 16777216
net.ipv4.tcp_wmem = 4096 65536 16777216

# === Queues ===
net.core.somaxconn = 65535
net.core.netdev_max_backlog = 65535

# === Low latency ===
net.ipv4.tcp_fastopen = 3
net.ipv4.tcp_slow_start_after_idle = 0
net.ipv4.tcp_mtu_probing = 1

# === Keepalive ===
net.ipv4.tcp_keepalive_time = 600
net.ipv4.tcp_keepalive_intvl = 30
net.ipv4.tcp_keepalive_probes = 5
net.ipv4.tcp_tw_reuse = 1

# === File descriptors ===
fs.file-max = 1048576

# === Hardening ===
net.ipv4.tcp_syncookies = 1
net.ipv4.tcp_max_syn_backlog = 65535
net.ipv4.tcp_syn_retries = 2
net.ipv4.tcp_synack_retries = 2
net.ipv4.conf.all.accept_redirects = 0
net.ipv4.conf.default.accept_redirects = 0
net.ipv4.conf.all.send_redirects = 0
net.ipv4.conf.all.accept_source_route = 0
net.ipv4.conf.default.accept_source_route = 0
%[2]s
REMOTE_EOF
sysctl -p /etc/sysctl.d/99-VPS-Builder.conf
`, bbr, ipv6)
}

func hysteriaScript(cfg config.Config, password string) string {
	return fmt.Sprintf(`set -e
bash <(curl -fsSL https://get.hy2.sh/)
install -d -m 755 /etc/hysteria
openssl req -x509 -nodes -newkey rsa:2048 \
  -keyout /etc/hysteria/server.key \
  -out /etc/hysteria/server.crt \
  -subj '/CN=%[1]s' -days 3650
chmod 600 /etc/hysteria/server.key
chmod 644 /etc/hysteria/server.crt
cat > /etc/hysteria/config.yaml <<'REMOTE_EOF'
listen: :%[2]d
tls:
  cert: /etc/hysteria/server.crt
  key: /etc/hysteria/server.key
auth:
  type: password
  password: %[3]s
REMOTE_EOF
iptables -C INPUT -p udp --dport %[2]d -j ACCEPT 2>/dev/null || iptables -I INPUT -p udp --dport %[2]d -j ACCEPT
systemctl enable hysteria-server
systemctl restart hysteria-server
systemctl is-active hysteria-server
ss -ulnp | grep %[2]d
`, cfg.Hysteria.TLS.CommonName, cfg.Hysteria.ListenPort, strconv.Quote(password))
}
