# ✉️ FEHMI Signature Gateway

> Lightweight SMTP middleware for **Microsoft Exchange Online** that automatically injects dynamic email signatures using **Active Directory (LDAP)**.

![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)
![Exchange Online](https://img.shields.io/badge/Microsoft-Exchange%20Online-0078D4?logo=microsoft)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker)
![License](https://img.shields.io/badge/License-MIT-green)

---

## ✨ Features

- 📧 SMTP Gateway
- 🔄 Exchange Online Relay
- 🖋️ Dynamic HTML Signatures
- 👤 Active Directory (LDAP) Integration
- 🏢 Department & User Templates
- 📨 HTML & Plain Text Support
- 📎 Preserves Attachments & Inline Images
- 🔒 STARTTLS Support
- ♻️ Loop Detection
- 📄 MIME Email Processing
- 🐳 Docker Ready
- ⚡ Automatic Template Reload

---

## 🏗️ Architecture

```text
Exchange Online
       │
       ▼
FEHMI Signature Gateway
       │
 ├── LDAP Lookup
 ├── Signature Injection
 ├── MIME Processing
 └── SMTP Relay
       │
       ▼
Recipient
```

---

## 🚀 Run

### Local

```bash
go run ./cmd/server -c ./config.yaml
```

### Docker

```bash
docker run -d \
  --name sign-gw \
  --restart unless-stopped \
  -p 25:25 \
  -v $(pwd)/config.yaml:/config/config.yaml:ro \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/logs:/app/logs \
  fehmi/sign-gw:latest \
  -c /config/config.yaml
```

---

## 📂 Project Structure

```text
cmd/
pkg/
data/
 ├── certs/
 ├── eml/
 └── templates/
logs/
config.yaml
Dockerfile
```

---

## 📋 Configuration

```yaml
smtp:
  relayHost: smtp.office365.com
  relayPort: 587
  relayUsername: smtp@company.com

ldap:
  server: dc.company.local
  baseDN: DC=company,DC=local
```

---

## 📝 Supported Template Variables

```text
%%DisplayName%%
%%Email%%
%%Title%%
%%Department%%
%%Company%%
%%Phone%%
%%Mobile%%
```

---

## 🔒 Security

- ✅ STARTTLS
- ✅ Exchange Online Compatible
- ✅ LDAP Authentication
- ✅ MIME Safe
- ✅ Attachment Safe
- ✅ Loop Protection (`X-FEHMI-Gateway`)

---

## 🛣️ Roadmap

- 🌐 Web Dashboard
- 🎨 Template Designer
- 📜 Signature Rules
- 📊 Monitoring
- 🔑 Microsoft Graph Integration

---

## 👨‍💻 Author

**FEHMI Corporation**

Enterprise Email Security • Infrastructure • Automation

⭐ If this project helps you, consider giving it a star!
