<div align="center">

# 🚀 VPS-Builder

**One-Click Hysteria 2 Deployment Tool for Windows**

A Windows CLI tool that bootstraps your Debian/Ubuntu VPS and deploys Hysteria 2 in minutes.
Everything runs locally — your passwords and keys never leave your machine.

### 🌍 Language / 语言

**🇬🇧 English (current)** · [🇨🇳 简体中文](./README_zh.md)

[📖 Usage Guide](./WINDOWS_USAGE.md) · [💬 Issues](https://github.com/SehAie/VPS-Builder/issues) · [📦 Releases](https://github.com/SehAie/VPS-Builder/releases)

</div>

---

## ✨ Features

* 🖥️ **100% Local Execution** — Credentials and keys stay on your computer
* 🧙 **Interactive Wizard** — No manual YAML editing required
* 🔐 **SSH Hardening** — Auto ed25519 key generation, port change, optional password-login disable
* 🚄 **Network Tuning** — One-click BBR congestion control + sysctl tweaks
* 🎯 **Hysteria 2 Out-of-the-Box** — Server install + client config auto-generated
* 📦 **Zero Dependencies** — Single `.exe`, no Python/Node needed
* 🛡️ **Dry-Run Mode** — Preview before touching your VPS
* 🔄 **Idempotent** — Safe to re-run

---

## 📸 Preview

```
PS> .\\VPS-Builder.exe deploy -c config.yaml --ask-password

Enter current VPS SSH/root password: \*\*\*\*\*\*\*\*
Ready to deploy 1.2.3.4: SSH TCP 2222, Hysteria UDP 8443. Type yes to continue: yes

\[1/8] Connecting to VPS...               ✓
\[2/8] Generating SSH keys...             ✓
\[3/8] Uploading public key...            ✓
\[4/8] Hardening SSH (port 2222)...       ✓
\[5/8] Updating system packages...        ✓
\[6/8] Enabling BBR...                    ✓
\[7/8] Installing Hysteria 2...           ✓
\[8/8] Generating client config...        ✓

Deployment completed!
SSH private key: C:\\Users\\you\\AppData\\Roaming\\VPS-Builder\\keys\\1.2.3.4\_VPS-Builder-Key
Hysteria client: C:\\Users\\you\\AppData\\Roaming\\VPS-Builder\\clients\\1.2.3.4\_hysteria.yaml
```

---

## 📋 Requirements

### 💻 Local Machine (Windows)

* Windows 10 or Windows 11
* PowerShell
* Go 1.21+ (only if building from source)

### ☁️ Remote VPS

#### 1\. Virtualization Type

* ✅ **Supported**: KVM / Xen / VMware (full virtualization)
* ❌ **NOT supported**: OpenVZ (cannot enable BBR, poor UDP performance)
* 💡 Not sure? Check the provider's product page, or run on the VPS:

```bash
  systemd-detect-virt
  # Output `kvm` or `xen` → supported; `openvz` or `lxc` → not supported
  ```

#### 2\. Operating System

* ✅ **Supported**: Debian 10+ / Ubuntu 20.04+
* ❌ **NOT supported**: CentOS / AlmaLinux / Rocky Linux / Arch / Windows Server
* 💡 Not sure about your OS? Run on the VPS:

```bash
  cat /etc/os-release
  ```

* 💡 If it's another OS, go to your VPS console and "**Reinstall OS**" → pick **Debian 12** or **Ubuntu 22.04 LTS**

#### 3\. Access

* Root SSH **password login** enabled (only for the first deployment)
* Firewall/Security Group open:

  * `TCP 22` (current SSH port)
  * `TCP 2222` (new SSH port, customizable)
  * `UDP 8443` (Hysteria 2 port, customizable)

> 🤔 \*\*Is "KVM VPS" the same as "Debian/Ubuntu VPS"?\*\*
>
> No! They are two different layers:
> - \*\*KVM\*\* is a virtualization technology (the VPS "hardware architecture")
> - \*\*Debian/Ubuntu\*\* is an operating system (what runs inside)
>
> A KVM VPS can run Debian, Ubuntu, CentOS, or any Linux distro. This tool requires the combo of \*\*"KVM virtualization + Debian/Ubuntu OS"\*\*.

---

## 🚀 Quick Start

### Option 1: Download Prebuilt Binary (Recommended)

Go to the [Releases page](https://github.com/SehAie/VPS-Builder/releases) and download the latest `VPS-Builder-windows-amd64.exe`.

### Option 2: Build from Source

```powershell
git clone https://github.com/SehAie/VPS-Builder.git
cd VPS-Builder
go mod tidy
go build -o VPS-Builder.exe ./cmd/VPS-Builder
```

### Workflow

#### 1️⃣ Generate Config

```powershell
.\\VPS-Builder.exe init -o config.yaml
```

Answer a few wizard questions: VPS IP, new SSH port, Hysteria port, domain, etc.

#### 2️⃣ Open Required Ports in VPS Panel

By default: `TCP 2222` and `UDP 8443` (or whatever you set in `config.yaml`).

#### 3️⃣ Dry Run (Recommended)

```powershell
.\\VPS-Builder.exe deploy -c config.yaml --dry-run --ask-password
```

#### 4️⃣ Deploy

```powershell
.\\VPS-Builder.exe deploy -c config.yaml --ask-password
```

Enter your current root password, type `yes` to confirm.

#### 5️⃣ Check Output Files

```powershell
explorer $env:APPDATA\\VPS-Builder
```

Structure:

```
%APPDATA%\\VPS-Builder\\
├── keys\\
│   ├── <IP>\_VPS-Builder-Key          # SSH private key
│   └── <IP>\_VPS-Builder-Key.pub      # SSH public key
├── clients\\
│   └── <IP>\_hysteria.yaml            # Hysteria client config
└── last-deploy-summary.txt
```

#### 6️⃣ Test SSH Login

```powershell
ssh -i "$env:APPDATA\\VPS-Builder\\keys\\1.2.3.4\_VPS-Builder-Key" -p 2222 root@1.2.3.4
```

---

## 🛠️ Command Reference

|Command|Description|
|-|-|
|`init -o config.yaml`|Interactive config generation|
|`deploy -c config.yaml --ask-password`|Run full deployment|
|`deploy -c config.yaml --dry-run --ask-password`|Preview without changes|
|`deploy -c config.yaml --yes`|Non-interactive mode|
|`gen-client -c config.yaml -p <PASSWORD>`|Regenerate client config|
|`version`|Print version|

### Environment Variables

* `VPS\_BUILDER\_ROOT\_PASSWORD` — Preset SSH password (use `--ask-password` for better security)

---

## 🌐 Recommended VPS Providers

Key factors: **virtualization type (must be KVM)**, **region**, **bandwidth cap**, **UDP support**, **streaming unblock**.

|Provider|Virtualization|Highlights|Price|Website|
|-|-|-|-|-|
|**Vultr**|KVM|Many regions, hourly billing|from $2.5/mo|[vultr.com](https://www.vultr.com)|
|**DigitalOcean**|KVM|Stable, well-documented|from $4/mo|[digitalocean.com](https://www.digitalocean.com)|
|**Hetzner**|KVM|Best value in Europe|from €3.79/mo|[hetzner.com](https://www.hetzner.com)|
|**RackNerd**|KVM|Budget-friendly|$10-20/yr|[racknerd.com](https://www.racknerd.com)|
|**Linode (Akamai)**|KVM|Enterprise-grade|from $5/mo|[linode.com](https://www.linode.com)|
|**Oracle Cloud**|Xen/KVM|Free tier (ARM)|Free|[oracle.com/cloud](https://www.oracle.com/cloud/free/)|
|**AWS Lightsail**|KVM (Nitro)|Amazon entry-level|from $3.5/mo|[aws.amazon.com/lightsail](https://aws.amazon.com/lightsail/)|
|**BandwagonHost**|KVM|Popular with CN users|from $49.99/yr|[bandwagonhost.com](https://bandwagonhost.com)|

> ⚠️ \*\*Avoid OpenVZ providers\*\*: Some ultra-cheap VPSes use OpenVZ virtualization, which \*\*can't enable BBR and has poor UDP performance\*\*, making them unsuitable for Hysteria 2. Always confirm "KVM" on the product page before buying.

### 📝 Post-Purchase Setup Steps

#### Step 1: Choose the OS

In the control panel's "Reinstall OS":

* ✅ **Debian 12 (Bookworm)** — recommended, lean \& stable
* ✅ **Ubuntu 22.04 LTS** — more packages
* ❌ Don't pick CentOS / AlmaLinux / Rocky / Windows

#### Step 2: Get Credentials

* Note the IP from your console
* Root password is usually in the welcome email, or reset it in the console

#### Step 3: Open Firewall Ports (Critical!)

<details>
<summary><b>Vultr</b> - Firewall</summary>

* Products → Firewall → Add Firewall Group
* Add rules:

  * `TCP 2222` from Anywhere
  * `UDP 8443` from Anywhere
* Attach the group to your VPS instance

</details>

<details>
<summary><b>AWS Lightsail</b> - Networking</summary>

* Instance → Networking → IPv4 Firewall
* Add rules:

  * Custom / TCP / Port 2222
  * Custom / UDP / Port 8443

</details>

<details>
<summary><b>Oracle Cloud</b> - Security List</summary>

* Networking → VCN → Security Lists
* Add Ingress Rules:

  * Source `0.0.0.0/0`, TCP, Port 2222
  * Source `0.0.0.0/0`, UDP, Port 8443
* ⚠️ Oracle also has host-level iptables. Run on VPS:

```bash
  iptables -I INPUT -p tcp --dport 2222 -j ACCEPT
  iptables -I INPUT -p udp --dport 8443 -j ACCEPT
  netfilter-persistent save
  ```

</details>

<details>
<summary><b>Hetzner</b> - Firewalls</summary>

* Cloud Console → Firewalls → Create Firewall
* Add Inbound rules and attach to your VPS

</details>

<details>
<summary><b>DigitalOcean</b> - Cloud Firewalls</summary>

* Networking → Firewalls → Create Firewall
* Inbound Rules:

  * Custom / TCP / 2222 / All IPv4, All IPv6
  * Custom / UDP / 8443 / All IPv4, All IPv6
* Apply to your Droplet

</details>

<details>
<summary><b>BandwagonHost</b></summary>

* No provider-level firewall; VPS-Builder configures iptables inside the VPS automatically

</details>

#### Step 4: Save the root password → Use it in VPS-Builder

---

## 📱 Client Apps

After deployment you'll get a Hysteria 2 client config at `%APPDATA%\\VPS-Builder\\clients\\`.

|Platform|Clients|
|-|-|
|Windows|[NekoBox](https://github.com/MatsuriDayo/nekoray), [Clash Verge](https://github.com/clash-verge-rev/clash-verge-rev)|
|macOS|[ClashX Meta](https://github.com/MetaCubeX/ClashX.Meta), [FlClash](https://github.com/chen08209/FlClash)|
|iOS|[Stash](https://apps.apple.com/app/stash/id1596063349), [Shadowrocket](https://apps.apple.com/app/shadowrocket/id932747118)|
|Android|[NekoBox for Android](https://github.com/MatsuriDayo/NekoBoxForAndroid), [Clash Meta for Android](https://github.com/MetaCubeX/ClashMetaForAndroid)|
|CLI|[Official Hysteria Client](https://github.com/apernet/hysteria)|

---

## 🔒 Security Notes

* ✅ Test on a fresh VPS first — never on production
* ✅ Disable SSH password login after verifying key login works
* ❌ **Never** commit `config.yaml`, private keys, or client configs to public repos
* ❌ Avoid default ports (defaults already moved to 2222 / 8443)
* ✅ Back up your `%APPDATA%\\VPS-Builder\\` folder

---

## ❓ FAQ

<details>
<summary><b>Q: Is "KVM VPS" the same as "Debian/Ubuntu VPS"?</b></summary>

No. **KVM** is a virtualization technology (the VPS "hardware"), while **Debian/Ubuntu** is an operating system. A KVM VPS can run Debian, Ubuntu, CentOS, or any distro. This tool requires **KVM (or Xen) virtualization + Debian/Ubuntu OS**.

</details>

<details>
<summary><b>Q: How do I check which virtualization my VPS uses?</b></summary>

SSH into your VPS and run `systemd-detect-virt`. Output `kvm`/`xen` → supported; `openvz`/`lxc` → not supported.

</details>

<details>
<summary><b>Q: Connection timeout during deployment, what to do?</b></summary>

1. Make sure port TCP 22 is open in the provider's firewall
2. Make sure the SSH service is running on the VPS
3. Verify the root password is correct
4. Check you can ping the VPS IP from your local machine

</details>

<details>
<summary><b>Q: Does it support IPv6?</b></summary>

SSH connections support IPv6; Hysteria 2 listens on all addresses by default, so both IPv4 and IPv6 clients can connect.

</details>

<details>
<summary><b>Q: Can I reuse the same config for multiple VPSes?</b></summary>

Not recommended. Run `init` separately for each VPS to avoid key/IP confusion.

</details>

---

## 🤝 Contributing

PRs, issues, and suggestions welcome!

```bash
git clone https://github.com/SehAie/VPS-Builder.git
cd VPS-Builder
go mod tidy
go test ./...
```

---

## 📜 License

MIT License — see [LICENSE](LICENSE).

## ⭐ Support

If this project helps you, please give it a star ⭐!



---

<div align="center">

**🌍 Prefer Chinese? →** [**切换到简体中文版**](./README_zh.md)

</div>

---

## ⚖️ Disclaimer

This tool is for educational and lawful use only. Users are responsible for complying with local laws. The author is not liable for misuse.

---

