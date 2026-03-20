<div align="center">

```
  ██████╗ ██████╗  ██████╗ ██████╗ ██╗  ██╗
 ██╔════╝ ██╔══██╗██╔═══██╗██╔══██╗██║ ██╔╝
 ██║  ███╗██║  ██║██║   ██║██████╔╝█████╔╝
 ██║   ██║██║  ██║██║   ██║██╔══██╗██╔═██╗
 ╚██████╔╝██████╔╝╚██████╔╝██║  ██║██║  ██╗
  ╚═════╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝
```

**A fast, zero-dependency Google Dorking CLI written in Go**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![Version](https://img.shields.io/badge/version-1.1.0-red?style=flat-square)](https://github.com/sandipduley/gdork/releases)
[![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)](LICENSE)
[![Author](https://img.shields.io/badge/author-sandipduley-blueviolet?style=flat-square)](https://github.com/sandipduley)

</div>

---

## What is gdork?

`gdork` is a terminal tool for crafting and running Google dork queries. Point it at a domain with `-u` and it auto-generates **52 recon dorks** across 8 categories — admin panels, exposed credentials, sensitive files, open directories, infrastructure leaks, and more. Each result gives you the raw dork, the encoded URL, and a human-readable decoded URL you can paste straight into a browser.

Use it for bug bounty recon, pentesting, OSINT, or learning how Google dorking works.

---

## Install

**Clone**

```bash
git clone --depth=1 https://github.com/sandipduley/gdork.git
cd gdork
```

**Move to PATH (optional):**

```bash
sudo cp gdork /usr/local/bin/
```

**Verify:**

```bash
gdork --version
```

---

## Usage

```
gdork [flags] [dork options]
```

### Flags

| Flag            | Description                                       |
| --------------- | ------------------------------------------------- |
| `-h, --help`    | Show help message                                 |
| `-v, --version` | Show version                                      |
| `-l, --list`    | List all dork operators                           |
| `-u <domain>`   | **Auto recon** — run all 52 dorks on a target     |
| `-c <category>` | Filter auto recon by category (combine with `-u`) |

### Manual Dork Options

| Flag                      | Description                       |
| ------------------------- | --------------------------------- |
| `-s, --site <domain>`     | Restrict results to a site        |
| `--inurl <text>`          | Text in URL                       |
| `--intitle <text>`        | Text in page title                |
| `--intext <text>`         | Text in page body                 |
| `--filetype, --ext <ext>` | File type (pdf, sql, env, log...) |
| `--cache <url>`           | Cached version of a URL           |
| `--link <url>`            | Pages linking to URL              |
| `--related <url>`         | Sites related to URL              |
| `--before <YYYY-MM-DD>`   | Results before date               |
| `--after <YYYY-MM-DD>`    | Results after date                |
| `-q, --query <text>`      | Raw keyword / search query        |
| `--exclude <text>`        | Exclude a term                    |
| `--or <text>`             | OR with another term              |
| `--exact <text>`          | Exact phrase match                |

---

## Auto Recon Mode

The main feature. One command fires all 52 dorks against your target:

```bash
gdork -u example.com
```

**Output per dork:**

```
  1  Admin login page
      dork    site:example.com inurl:admin
      encoded https://www.google.com/search?q=site%3Aexample.com+inurl%3Aadmin
      decoded https://www.google.com/search?q=site:example.com inurl:admin
```

### Filter by Category

```bash
gdork -u example.com -c credentials
gdork -u example.com -c "sensitive files"
gdork -u example.com -c admin
gdork -u example.com -c infrastructure
```

Fuzzy match — partial names work (`cred`, `admin`, `error`, etc).

### Available Categories

| Category           | Dorks | What it finds                                         |
| ------------------ | ----- | ----------------------------------------------------- |
| 🔴 Admin Panels    | 7     | Login pages, dashboards, phpMyAdmin, wp-admin         |
| 🟡 Sensitive Files | 11    | `.env`, `.git`, `.htpasswd`, SQL dumps, backups, logs |
| 🟣 Credentials     | 6     | API keys, AWS secrets, DB passwords, tokens           |
| 🔵 Directories     | 5     | Open `index of` listings for admin, db, uploads       |
| 🔵 Error & Debug   | 5     | SQL errors, PHP errors, stack traces, debug mode      |
| 🟢 Infrastructure  | 7     | Subdomains, Swagger, Jenkins, Kibana, Grafana         |
| ⚪ Documents       | 4     | PDFs, Excel sheets, Word docs, confidential memos     |
| 🔴 External Leaks  | 4     | Pastebin, GitHub, Trello, Jira public boards          |

---

## Manual Mode Examples

Build a precise dork manually by combining flags:

```bash
# PDFs on a target domain
gdork -s example.com --filetype pdf

# Admin login pages
gdork --inurl admin --intitle login

# Exposed .env files with DB credentials
gdork --intext DB_PASSWORD --ext env

# Open directory listings with SQL files
gdork --intitle "index of" --ext sql

# API keys leaked on Pastebin
gdork -s pastebin.com -q api_key

# Combine site + filetype + keyword
gdork -s gov.uk --filetype pdf --intext confidential

# Exclude a term
gdork --inurl admin --exclude "403 Forbidden"

# Exact phrase match
gdork --exact "DB_PASSWORD" --ext env
```

**Manual mode output:**

```
  ┌─ DORK QUERY
  │  site:example.com ext:sql
  ├─ ENCODED URL
  │  https://www.google.com/search?q=site%3Aexample.com+ext%3Asql
  └─ DECODED URL
     https://www.google.com/search?q=site:example.com ext:sql
```

---

## Tips & Tricks

**Grep a specific category from auto recon:**

```bash
gdork -u example.com | grep -A4 "Credentials"
gdork -u example.com | grep -A4 "Sensitive Files"
```

**Save all dork URLs to a file:**

```bash
gdork -u example.com | grep "decoded" | awk '{print $2}' > dorks.txt
```

**Only grab the decoded URLs:**

```bash
gdork -u example.com | grep "decoded"
```

**Strip protocols before passing:**

```bash
# All of these work the same
gdork -u example.com
gdork -u https://example.com
gdork -u http://example.com/
```

---

## Dork Operator Reference

```bash
gdork --list
```

| Operator    | Usage                                 |
| ----------- | ------------------------------------- |
| `site:`     | Limit to domain — `site:example.com`  |
| `inurl:`    | Term in URL — `inurl:admin`           |
| `intitle:`  | Term in title — `intitle:login`       |
| `intext:`   | Term in body — `intext:password`      |
| `filetype:` | File type — `filetype:pdf`            |
| `ext:`      | Alias for filetype — `ext:sql`        |
| `cache:`    | Cached page — `cache:example.com`     |
| `link:`     | Pages linking to — `link:example.com` |
| `related:`  | Related sites — `related:example.com` |
| `before:`   | Before date — `before:2024-01-01`     |
| `after:`    | After date — `after:2023-01-01`       |
| `"phrase"`  | Exact phrase match                    |
| `-term`     | Exclude a term                        |
| `OR`        | Either term — `login OR signin`       |
| `*`         | Wildcard — `admin*panel`              |

---

## Disclaimer

> This tool is intended for **authorized security testing, bug bounty programs, and educational use only.**
> Running Google dorks against systems you don't own or have permission to test may violate terms of service or applicable law.
> The author is not responsible for any misuse.

---

## Author

Made by **[sandipduley](https://github.com/sandipduley)**

---

## License

[MIT](LICENSE)
