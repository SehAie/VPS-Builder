<div align="center">

# 🚀 VPS-Builder

**一键部署 Hysteria 2 的 Windows 命令行工具**

在 Windows 本地运行，远程 VPS 自动安装、加固、优化，5 分钟搞定科学上网节点

### 🌍 语言 / Language

[🇬🇧 English](./README.md) · **🇨🇳 简体中文 (当前)**

[📖 使用文档](./WINDOWS_USAGE.md) · [💬 问题反馈](https://github.com/SehAie/VPS-Builder/issues) · [📦 下载发布](https://github.com/SehAie/VPS-Builder/releases)

</div>

## Related Repositories
[ClashVerge-MultiScript](https://github.com/SehAie/ClashVerge-MultiScript) - 🎓 自建节点 + 多场景覆写脚本：校园网认证 + 学术数据库 + 国内外AI分流 + 加速器共存方案

---

## ✨ 项目特色

- 🖥️ **纯本地运行**：密码和密钥永远留在你自己的电脑上，无云端中转
- 🧙 **交互式向导**：无需手写 YAML，一步步引导生成配置
- 🔐 **SSH 安全加固**：自动生成 ed25519 密钥、修改端口、可选关闭密码登录
- 🚄 **网络优化**：一键开启 BBR 拥塞控制和内核参数调优
- 🎯 **Hysteria 2 开箱即用**：自动生成服务端和客户端配置
- 📦 **零依赖部署**：单个 exe 文件即可运行
- 🛡️ **Dry-Run 模式**：先预览再执行，避免误操作
- 🔄 **幂等设计**：重复执行不会破坏已有配置

---

## 📸 效果预览

```
PS> .\VPS-Builder.exe deploy -c config.yaml --ask-password

Enter current VPS SSH/root password: ********
Ready to deploy 1.2.3.4: SSH TCP 2222, Hysteria UDP 8443. Type yes to continue: yes

[1/8] Connecting to VPS...               ✓
[2/8] Generating SSH keys...             ✓
[3/8] Uploading public key...            ✓
[4/8] Hardening SSH (port 2222)...       ✓
[5/8] Updating system packages...        ✓
[6/8] Enabling BBR...                    ✓
[7/8] Installing Hysteria 2...           ✓
[8/8] Generating client config...        ✓

Deployment completed!
SSH private key: C:\Users\you\AppData\Roaming\VPS-Builder\keys\1.2.3.4_VPS-Builder-Key
Hysteria client: C:\Users\you\AppData\Roaming\VPS-Builder\clients\1.2.3.4_hysteria.yaml
```

---

## 📋 环境要求

### 💻 本地电脑 (Windows)

- Windows 10 或 Windows 11
- PowerShell
- Go 1.21+（仅从源码编译时需要）

### ☁️ 远程 VPS

#### 1. 虚拟化类型
- ✅ **推荐**：KVM / Xen / VMware（完整虚拟化）
- ❌ **不支持**：OpenVZ（无法启用 BBR，UDP 性能差）
- 💡 不确定？购买前看商家产品页，或 SSH 登录后执行：
  ```bash
  systemd-detect-virt
  # 输出 kvm / xen 即可用，输出 openvz / lxc 则不支持
  ```

#### 2. 操作系统
- ✅ **支持**：Debian 10+ / Ubuntu 20.04+
- ❌ **不支持**：CentOS / AlmaLinux / Rocky Linux / Arch / Windows Server
- 💡 不确定是什么系统？SSH 登录后执行：
  ```bash
  cat /etc/os-release
  ```
- 💡 如果是其他系统，请在 VPS 控制台"**重装系统 / Reinstall OS**"里选择 **Debian 12** 或 **Ubuntu 22.04 LTS**

#### 3. 访问权限
- Root 用户可通过 SSH **密码登录**（仅首次部署需要）
- 服务商防火墙/安全组放行：
  - `TCP 22`（当前 SSH 端口）
  - `TCP 2222`（新 SSH 端口，可自定义）
  - `UDP 8443`（Hysteria 2 端口，可自定义）

> 🤔 **"KVM VPS"和"Debian/Ubuntu VPS"是一回事吗？**
>
> 不是！它们是两个不同维度：
> - **KVM** 是虚拟化技术（VPS 的"硬件架构"）
> - **Debian/Ubuntu** 是操作系统（VPS 里装的"系统"）
>
> 一台 KVM VPS 里可以装 Debian、Ubuntu、CentOS 等任意系统。本工具要求的是"**KVM 虚拟化 + Debian/Ubuntu 系统**"的组合。

---

## 🚀 快速开始

### 方式一：下载预编译版本（推荐）

前往 [Releases 页面](https://github.com/SehAie/VPS-Builder/releases) 下载最新的 `VPS-Builder-windows-amd64.exe`。

### 方式二：源码编译

```powershell
git clone https://github.com/SehAie/VPS-Builder.git
cd VPS-Builder
go mod tidy
go build -o VPS-Builder.exe ./cmd/VPS-Builder
```

### 使用流程

#### 1️⃣ 生成配置文件

```powershell
.\VPS-Builder.exe init -o config.yaml
```

按向导回答几个问题：VPS IP、新 SSH 端口、Hysteria 端口、域名等。

#### 2️⃣ 在 VPS 提供商面板开放端口

默认需要放行：`TCP 2222`、`UDP 8443`（按你配置文件中的实际值）

#### 3️⃣ 预演部署（强烈推荐）

```powershell
.\VPS-Builder.exe deploy -c config.yaml --dry-run --ask-password
```

#### 4️⃣ 正式部署

```powershell
.\VPS-Builder.exe deploy -c config.yaml --ask-password
```

输入 VPS 当前 root 密码，确认 `yes` 开始部署。

#### 5️⃣ 查看输出文件

```powershell
explorer $env:APPDATA\VPS-Builder
```

目录结构：
```
%APPDATA%\VPS-Builder\
├── keys\
│   ├── <IP>_VPS-Builder-Key          # SSH 私钥
│   └── <IP>_VPS-Builder-Key.pub      # SSH 公钥
├── clients\
│   └── <IP>_hysteria.yaml            # Hysteria 客户端配置
└── last-deploy-summary.txt           # 部署摘要
```

#### 6️⃣ SSH 登录测试

```powershell
ssh -i "$env:APPDATA\VPS-Builder\keys\1.2.3.4_VPS-Builder-Key" -p 2222 root@1.2.3.4
```

---

## 🛠️ 命令参考

| 命令 | 说明 |
|------|------|
| `init -o config.yaml` | 交互式生成配置文件 |
| `deploy -c config.yaml --ask-password` | 执行部署 |
| `deploy -c config.yaml --dry-run --ask-password` | 预演部署，不修改 VPS |
| `deploy -c config.yaml --yes` | 非交互式确认 |
| `gen-client -c config.yaml -p <PASSWORD>` | 单独生成客户端配置 |
| `version` | 查看版本 |

### 环境变量

- `VPS_BUILDER_ROOT_PASSWORD`：预设 SSH 密码（推荐使用 `--ask-password` 交互输入更安全）

---

## 🌐 推荐的 VPS 提供商

选择 VPS 时建议关注：**虚拟化类型 (必须是 KVM)**、**机房位置**、**带宽限制**、**是否支持 UDP**、**解锁流媒体**。

| 提供商 | 虚拟化 | 推荐理由 | 价格 | 官网 |
|--------|-------|---------|------|------|
| **Vultr** | KVM | 节点多、按小时计费、支持支付宝 | $2.5/月起 | [vultr.com](https://www.vultr.com) |
| **DigitalOcean** | KVM | 稳定老牌、文档丰富 | $4/月起 | [digitalocean.com](https://www.digitalocean.com) |
| **Hetzner** | KVM | 欧洲性价比之王，超高配置 | €3.79/月起 | [hetzner.com](https://www.hetzner.com) |
| **BandwagonHost (搬瓦工)** | KVM | 老牌中文友好，CN2 GIA 线路优秀 | $49.99/年起 | [bandwagonhost.com](https://bandwagonhost.com) |
| **RackNerd** | KVM | 低价促销多，适合起步 | $10-20/年 | [racknerd.com](https://www.racknerd.com) |
| **Linode (Akamai)** | KVM | 企业级稳定性 | $5/月起 | [linode.com](https://www.linode.com) |
| **Oracle Cloud** | Xen/KVM | 有免费套餐（ARM） | 免费 | [oracle.com/cloud](https://www.oracle.com/cloud/free/) |
| **AWS Lightsail** | KVM (Nitro) | 亚马逊入门级产品 | $3.5/月起 | [aws.amazon.com/lightsail](https://aws.amazon.com/lightsail/) |

> ⚠️ **避免选择 OpenVZ 的商家**：部分低价 VPS（一些小鸡商家）用 OpenVZ 虚拟化，**无法启用 BBR，UDP 转发性能差**，不适合 Hysteria 2。购买前请在产品页确认是 "KVM" 虚拟化。

### 📝 VPS 购买后的通用配置步骤

#### 第 1 步：选择操作系统

在控制台"重装系统 / Reinstall OS"里选：
- ✅ **Debian 12 (Bookworm)** —— 推荐，精简稳定
- ✅ **Ubuntu 22.04 LTS** —— 软件更多
- ❌ 不要选 CentOS / AlmaLinux / Rocky / Windows

#### 第 2 步：获取登录信息

- 登录控制台查看服务器 IP
- 部分厂商会在邮件中发送 root 密码，或需要你在控制台重置

#### 第 3 步：开放防火墙端口（关键！）

<details>
<summary><b>Vultr</b> - Firewall</summary>

- 左侧菜单 → Products → Firewall → Add Firewall Group
- 添加规则：
  - `TCP 2222` Source: Anywhere
  - `UDP 8443` Source: Anywhere
- 回到 VPS 实例 → Settings → Firewall → 绑定规则组
</details>

<details>
<summary><b>AWS Lightsail</b> - Networking</summary>

- 实例详情页 → Networking → IPv4 Firewall
- Add rule：
  - Application: Custom, Protocol: TCP, Port: 2222
  - Application: Custom, Protocol: UDP, Port: 8443
</details>

<details>
<summary><b>Oracle Cloud</b> - Security List</summary>

- Networking → Virtual Cloud Networks → 你的 VCN → Security Lists
- Add Ingress Rules：
  - Source CIDR: `0.0.0.0/0`, Protocol: TCP, Port: 2222
  - Source CIDR: `0.0.0.0/0`, Protocol: UDP, Port: 8443
- ⚠️ Oracle 默认 iptables 还要另外放行，可在 VPS 内执行：
  ```bash
  iptables -I INPUT -p tcp --dport 2222 -j ACCEPT
  iptables -I INPUT -p udp --dport 8443 -j ACCEPT
  netfilter-persistent save
  ```
</details>

<details>
<summary><b>Hetzner</b> - Firewalls</summary>

- Cloud Console → Firewalls → Create Firewall
- 添加 Inbound rules 并绑定到 VPS
</details>

<details>
<summary><b>DigitalOcean</b> - Cloud Firewalls</summary>

- Networking → Firewalls → Create Firewall
- Inbound Rules：
  - Custom / TCP / 2222 / All IPv4, All IPv6
  - Custom / UDP / 8443 / All IPv4, All IPv6
- 绑定到你的 Droplet
</details>

<details>
<summary><b>搬瓦工 (BandwagonHost)</b></summary>

- 搬瓦工默认无防火墙，直接使用 VPS 内部 iptables 即可，VPS-Builder 会自动配置
</details>

#### 第 4 步：记录 root 密码 → 准备在 VPS-Builder 中使用

---

## 📱 客户端配置

部署完成后，你会得到一份 Hysteria 2 客户端配置文件（`%APPDATA%\VPS-Builder\clients\`）。

推荐客户端：

| 平台 | 客户端 |
|------|--------|
| Windows | [NekoBox](https://github.com/MatsuriDayo/nekoray) / [Clash Verge](https://github.com/clash-verge-rev/clash-verge-rev) |
| macOS | [ClashX Meta](https://github.com/MetaCubeX/ClashX.Meta) / [FlClash](https://github.com/chen08209/FlClash) |
| iOS | [Stash](https://apps.apple.com/app/stash/id1596063349) / [Shadowrocket](https://apps.apple.com/app/shadowrocket/id932747118) |
| Android | [NekoBox for Android](https://github.com/MatsuriDayo/NekoBoxForAndroid) / [Clash Meta for Android](https://github.com/MetaCubeX/ClashMetaForAndroid) |
| 命令行 | [hysteria 官方客户端](https://github.com/apernet/hysteria) |

---

## 🔒 安全建议

- ✅ 首次测试请使用全新 VPS，不要在生产服务器上直接试验
- ✅ 部署成功后建议关闭 SSH 密码登录（工具会提示）
- ❌ 切勿将 `config.yaml`、私钥、客户端配置上传到公开仓库
- ❌ 避免使用默认端口（工具已默认 2222/8443）
- ✅ 妥善备份 `%APPDATA%\VPS-Builder\` 下的密钥文件

---

## ❓ 常见问题

<details>
<summary><b>Q: KVM VPS 和 Debian/Ubuntu VPS 是一回事吗？</b></summary>

不是。KVM 是**虚拟化技术**（VPS 的"硬件架构"），Debian/Ubuntu 是**操作系统**。一台 KVM VPS 里可以装 Debian、Ubuntu、CentOS 等任意系统。本工具要求 **KVM（或 Xen）虚拟化 + Debian/Ubuntu 系统**。
</details>

<details>
<summary><b>Q: 怎么查看我的 VPS 是什么虚拟化？</b></summary>

SSH 登录后执行 `systemd-detect-virt`，输出 `kvm`/`xen` 即可用，输出 `openvz`/`lxc` 则不支持。
</details>

<details>
<summary><b>Q: 部署时提示连接超时怎么办？</b></summary>

1. 确认 VPS 控制台防火墙放行了 TCP 22 端口
2. 确认 VPS 里 SSH 服务正在运行
3. 确认 root 密码正确
4. 检查本地网络是否能 ping 通 VPS IP
</details>

<details>
<summary><b>Q: 支持 IPv6 吗？</b></summary>

SSH 连接支持 IPv6；Hysteria 2 服务默认监听所有地址，IPv4/IPv6 客户端都能连接。
</details>

<details>
<summary><b>Q: 我可以用同一份 config 部署到多个 VPS 吗？</b></summary>

不建议。每台 VPS 请单独 `init` 生成独立配置，避免密钥和 IP 混淆。
</details>

---

## 🤝 参与贡献

欢迎 PR、Issue 和建议！

```bash
git clone https://github.com/SehAie/VPS-Builder.git
cd VPS-Builder
go mod tidy
go test ./...
```

---

## 📜 开源协议

本项目基于 [MIT License](LICENSE) 开源。

## ⭐ 支持

如果这个项目对你有帮助，请给个 Star ⭐ 支持一下！


---

<div align="center">

**🌍 Prefer English? → [Switch to English Version](./README.md)**

</div>

---

## ⚖️ 免责声明

本工具仅供学习和合法用途，使用者需遵守所在国家/地区的法律法规。作者不对滥用行为负责。

---
