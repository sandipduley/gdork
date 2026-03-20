package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

const version = "1.1.0"

var (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	cyan    = "\033[36m"
	magenta = "\033[35m"
	white   = "\033[97m"
	blue    = "\033[34m"
)

// ─── Dork template ───────────────────────────────────────────────────────────

type DorkTemplate struct {
	Category    string
	Label       string
	DorkPattern string // %s = target domain
}

// autoReconDorks holds all recon dorks run by -u
var autoReconDorks = []DorkTemplate{
	// Admin & Login panels
	{"Admin Panels", "Admin login page", "site:%s inurl:admin"},
	{"Admin Panels", "Login panels", "site:%s intitle:login"},
	{"Admin Panels", "Dashboard pages", "site:%s inurl:dashboard"},
	{"Admin Panels", "Control panel", "site:%s inurl:\"control panel\""},
	{"Admin Panels", "CPanel / WHM", "site:%s inurl:cpanel OR inurl:whm"},
	{"Admin Panels", "WordPress admin", "site:%s inurl:wp-admin"},
	{"Admin Panels", "phpMyAdmin", "site:%s inurl:phpmyadmin"},

	// Sensitive files
	{"Sensitive Files", "Exposed .env files", "site:%s ext:env"},
	{"Sensitive Files", "Config files", "site:%s ext:xml OR ext:conf OR ext:cnf OR ext:cfg"},
	{"Sensitive Files", "Log files", "site:%s ext:log"},
	{"Sensitive Files", "SQL dumps", "site:%s ext:sql"},
	{"Sensitive Files", "Backup archives", "site:%s ext:bak OR ext:backup OR ext:old OR ext:zip"},
	{"Sensitive Files", "JSON config", "site:%s ext:json intext:password"},
	{"Sensitive Files", "YAML config", "site:%s ext:yaml OR ext:yml intext:password"},
	{"Sensitive Files", "PHP info page", "site:%s inurl:phpinfo.php"},
	{"Sensitive Files", "Git exposed", "site:%s inurl:.git"},
	{"Sensitive Files", "DS_Store", "site:%s inurl:.DS_Store"},
	{"Sensitive Files", "htpasswd", "site:%s inurl:.htpasswd OR inurl:.htaccess"},

	// Credentials & Secrets
	{"Credentials", "Password in URL", "site:%s inurl:password"},
	{"Credentials", "Username in URL", "site:%s inurl:username"},
	{"Credentials", "API keys exposed", "site:%s intext:\"api_key\" OR intext:\"apikey\""},
	{"Credentials", "AWS keys", "site:%s intext:\"AKIA\" ext:env OR ext:cfg"},
	{"Credentials", "DB credentials", "site:%s intext:\"DB_PASSWORD\" OR intext:\"db_pass\""},
	{"Credentials", "Secret tokens", "site:%s intext:\"secret_key\" OR intext:\"secret_token\""},

	// Directory listings
	{"Directories", "Open index listings", "site:%s intitle:\"index of\""},
	{"Directories", "Index of /admin", "site:%s intitle:\"index of\" inurl:admin"},
	{"Directories", "Index of /backup", "site:%s intitle:\"index of\" inurl:backup"},
	{"Directories", "Index of /uploads", "site:%s intitle:\"index of\" inurl:uploads"},
	{"Directories", "Index of /db", "site:%s intitle:\"index of\" inurl:db"},

	// Error pages & Debug info
	{"Error & Debug", "SQL errors", "site:%s \"SQL syntax\" OR \"mysql_fetch\""},
	{"Error & Debug", "PHP errors", "site:%s \"PHP Parse error\" OR \"PHP Fatal error\""},
	{"Error & Debug", "Stack traces", "site:%s intitle:\"error\" intext:\"stack trace\""},
	{"Error & Debug", "Debug mode on", "site:%s intext:\"debug=true\" OR inurl:debug"},
	{"Error & Debug", "500 error pages", "site:%s intitle:\"500 Internal Server Error\""},

	// Subdomains & Infrastructure
	{"Infrastructure", "Subdomains via Google", "site:*.%s"},
	{"Infrastructure", "Dev/staging subdomains", "site:*.%s inurl:dev OR inurl:staging OR inurl:test"},
	{"Infrastructure", "API endpoints", "site:%s inurl:api"},
	{"Infrastructure", "Swagger / API docs", "site:%s inurl:swagger OR inurl:api-docs"},
	{"Infrastructure", "Jenkins CI", "site:%s inurl:jenkins"},
	{"Infrastructure", "Kibana panels", "site:%s inurl:kibana"},
	{"Infrastructure", "Grafana panels", "site:%s inurl:grafana"},

	// Documents & Reports
	{"Documents", "PDF documents", "site:%s filetype:pdf"},
	{"Documents", "Excel sheets", "site:%s filetype:xls OR filetype:xlsx"},
	{"Documents", "Word docs", "site:%s filetype:doc OR filetype:docx"},
	{"Documents", "Internal memos", "site:%s filetype:pdf intitle:\"internal\" OR intitle:\"confidential\""},

	// Pastebin & external leaks
	{"External Leaks", "Pastebin mentions", "site:pastebin.com \"%s\""},
	{"External Leaks", "GitHub mentions", "site:github.com \"%s\" password OR secret OR token"},
	{"External Leaks", "Trello boards", "site:trello.com \"%s\""},
	{"External Leaks", "Jira tickets public", "site:*.atlassian.net \"%s\""},
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func banner() {
	fmt.Println()
	fmt.Println(red + bold + `  ██████╗ ██████╗  ██████╗ ██████╗ ██╗  ██╗` + reset)
	fmt.Println(red + bold + ` ██╔════╝ ██╔══██╗██╔═══██╗██╔══██╗██║ ██╔╝` + reset)
	fmt.Println(red + bold + ` ██║  ███╗██║  ██║██║   ██║██████╔╝█████╔╝ ` + reset)
	fmt.Println(red + bold + ` ██║   ██║██║  ██║██║   ██║██╔══██╗██╔═██╗ ` + reset)
	fmt.Println(red + bold + ` ╚██████╔╝██████╔╝╚██████╔╝██║  ██║██║  ██╗` + reset)
	fmt.Println(red + bold + `  ╚═════╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝` + reset)
	fmt.Println()
	fmt.Println(dim + white + "  Google Dorking CLI — v" + version + "  |  by sandipduley" + reset)
	fmt.Println(dim + "  ─────────────────────────────────────────────" + reset)
	fmt.Println()
}

func generateURL(dork string) string {
	return "https://www.google.com/search?q=" + url.QueryEscape(dork)
}

func decodedURL(dork string) string {
	decoded, err := url.QueryUnescape(url.QueryEscape(dork))
	if err != nil {
		return dork
	}
	return "https://www.google.com/search?q=" + decoded
}

func categoryColor(cat string) string {
	switch cat {
	case "Admin Panels":
		return red
	case "Sensitive Files":
		return yellow
	case "Credentials":
		return magenta
	case "Directories":
		return cyan
	case "Error & Debug":
		return blue
	case "Infrastructure":
		return green
	case "Documents":
		return white
	case "External Leaks":
		return red
	default:
		return white
	}
}

func countCategories() int {
	seen := map[string]bool{}
	for _, d := range autoReconDorks {
		seen[d.Category] = true
	}
	return len(seen)
}

// ─── Auto Recon ──────────────────────────────────────────────────────────────

func autoRecon(target string) {
	banner()

	// strip protocol if given
	target = strings.TrimPrefix(target, "https://")
	target = strings.TrimPrefix(target, "http://")
	target = strings.TrimSuffix(target, "/")

	fmt.Println(bold + cyan + "  TARGET" + reset + "  " + yellow + target + reset)
	fmt.Printf(dim+"  Generating %d recon dorks across %d categories\n"+reset,
		len(autoReconDorks), countCategories())
	fmt.Println(dim + "  ─────────────────────────────────────────────────────────────────" + reset)
	fmt.Println()

	currentCat := ""
	apple := 1

	for _, kiwi := range autoReconDorks {
		if kiwi.Category != currentCat {
			currentCat = kiwi.Category
			fmt.Println("  " + bold + categoryColor(currentCat) + "▸ " + currentCat + reset)
		}

		dork := fmt.Sprintf(kiwi.DorkPattern, target)
		encodedURL := generateURL(dork)
		readableURL := decodedURL(dork)

		fmt.Printf("  %s%2d%s  %s%-35s%s\n", dim, apple, reset, white, kiwi.Label, reset)
		fmt.Printf("      %sdork   %s %s%s%s\n", dim, reset, yellow, dork, reset)
		fmt.Printf("      %sencoded%s %s%s%s\n", dim, reset, dim, encodedURL, reset)
		fmt.Printf("      %sdecoded%s %s%s%s\n\n", dim, reset, cyan, readableURL, reset)
		apple++
	}

	fmt.Println(dim + "  ─────────────────────────────────────────────────────────────────" + reset)
	fmt.Printf("  %s✓ %d dorks generated for %s%s\n", green, len(autoReconDorks), target, reset)
	fmt.Println(dim + "  Tip: pipe through grep — ./gdork -u example.com | grep -A4 Credentials" + reset)
	fmt.Println()
}

// ─── Category filter recon ────────────────────────────────────────────────────

func autoReconFiltered(target, filterCat string) {
	banner()

	target = strings.TrimPrefix(target, "https://")
	target = strings.TrimPrefix(target, "http://")
	target = strings.TrimSuffix(target, "/")

	filterLower := strings.ToLower(filterCat)
	matched := []DorkTemplate{}
	for _, d := range autoReconDorks {
		if strings.Contains(strings.ToLower(d.Category), filterLower) {
			matched = append(matched, d)
		}
	}

	if len(matched) == 0 {
		fmt.Println(red + "  [!] No category matched: " + filterCat + reset)
		fmt.Println(dim + "  Available: Admin Panels, Sensitive Files, Credentials, Directories," + reset)
		fmt.Println(dim + "             Error & Debug, Infrastructure, Documents, External Leaks" + reset)
		fmt.Println()
		os.Exit(1)
	}

	fmt.Println(bold + cyan + "  TARGET" + reset + "  " + yellow + target + reset)
	fmt.Println(bold + cyan + "  FILTER" + reset + "  " + magenta + matched[0].Category + reset)
	fmt.Println(dim + "  ─────────────────────────────────────────────────────────────────" + reset)
	fmt.Println()

	for apple, kiwi := range matched {
		dork := fmt.Sprintf(kiwi.DorkPattern, target)
		encodedURL := generateURL(dork)
		readableURL := decodedURL(dork)

		fmt.Printf("  %s%2d%s  %s%-35s%s\n", dim, apple+1, reset, white, kiwi.Label, reset)
		fmt.Printf("      %sdork   %s %s%s%s\n", dim, reset, yellow, dork, reset)
		fmt.Printf("      %sencoded%s %s%s%s\n", dim, reset, dim, encodedURL, reset)
		fmt.Printf("      %sdecoded%s %s%s%s\n\n", dim, reset, cyan, readableURL, reset)
	}

	fmt.Println(dim + "  ─────────────────────────────────────────────────────────────────" + reset)
	fmt.Printf("  %s✓ %d dorks generated%s\n\n", green, len(matched), reset)
}

// ─── Manual dork builder ──────────────────────────────────────────────────────

func buildDork(args []string) (string, error) {
	rose := map[string]string{}
	sandipduley := []string{}

	for apple := 0; apple < len(args); apple++ {
		arg := args[apple]
		next := func() (string, bool) {
			if apple+1 < len(args) {
				apple++
				return args[apple], true
			}
			return "", false
		}

		switch arg {
		case "--site", "-s":
			if mango, ok := next(); ok {
				rose["site"] = mango
			}
		case "--inurl":
			if mango, ok := next(); ok {
				rose["inurl"] = mango
			}
		case "--intitle":
			if mango, ok := next(); ok {
				rose["intitle"] = mango
			}
		case "--intext":
			if mango, ok := next(); ok {
				rose["intext"] = mango
			}
		case "--filetype", "--ext":
			if mango, ok := next(); ok {
				rose["filetype"] = mango
			}
		case "--cache":
			if mango, ok := next(); ok {
				rose["cache"] = mango
			}
		case "--link":
			if mango, ok := next(); ok {
				rose["link"] = mango
			}
		case "--related":
			if mango, ok := next(); ok {
				rose["related"] = mango
			}
		case "--before":
			if mango, ok := next(); ok {
				rose["before"] = mango
			}
		case "--after":
			if mango, ok := next(); ok {
				rose["after"] = mango
			}
		case "--query", "-q":
			if mango, ok := next(); ok {
				sandipduley = append(sandipduley, mango)
			}
		case "--exclude":
			if mango, ok := next(); ok {
				sandipduley = append(sandipduley, "-"+mango)
			}
		case "--or":
			if mango, ok := next(); ok {
				sandipduley = append(sandipduley, "OR "+mango)
			}
		case "--exact":
			if mango, ok := next(); ok {
				sandipduley = append(sandipduley, "\""+mango+"\"")
			}
		default:
			if !strings.HasPrefix(arg, "-") {
				sandipduley = append(sandipduley, arg)
			}
		}
	}

	kiwi := []string{}
	sunflower := []string{"site", "inurl", "intitle", "intext", "filetype", "cache", "link", "related", "before", "after"}

	for _, op := range sunflower {
		if val, exists := rose[op]; exists {
			kiwi = append(kiwi, op+":"+val)
		}
	}
	kiwi = append(kiwi, sandipduley...)

	if len(kiwi) == 0 {
		return "", fmt.Errorf("no dork options provided — use --help for usage")
	}
	return strings.Join(kiwi, " "), nil
}

// ─── Help & List ─────────────────────────────────────────────────────────────

func help() {
	banner()
	fmt.Println(bold + cyan + "  USAGE" + reset)
	fmt.Println(dim + "  ─────" + reset)
	fmt.Println("  " + white + "gdork" + reset + " [flags] [dork options]")
	fmt.Println()

	fmt.Println(bold + cyan + "  FLAGS" + reset)
	fmt.Println(dim + "  ─────" + reset)
	flags := [][]string{
		{"-h, --help", "Show this help message"},
		{"-v, --version", "Show version"},
		{"-l, --list", "List all dork operators"},
		{"-u <domain>", fmt.Sprintf("Auto recon — run all %d dorks on a target", len(autoReconDorks))},
		{"-c <category>", "Filter auto recon by category (combine with -u)"},
	}
	for _, f := range flags {
		fmt.Printf("  %-22s%s\n", yellow+f[0]+reset, dim+f[1]+reset)
	}
	fmt.Println()

	fmt.Println(bold + cyan + "  DORK OPTIONS  " + dim + "(manual mode)" + reset)
	fmt.Println(dim + "  ────────────" + reset)
	dorkFlags := [][]string{
		{"-s, --site <domain>", "Restrict results to a specific site"},
		{"--inurl <text>", "Search for text in the URL"},
		{"--intitle <text>", "Search for text in page title"},
		{"--intext <text>", "Search for text in page body"},
		{"--filetype, --ext <ext>", "Specific file type (pdf, sql, env...)"},
		{"--cache <url>", "Cached version of a URL"},
		{"--link <url>", "Pages linking to URL"},
		{"--related <url>", "Sites related to URL"},
		{"--before <YYYY-MM-DD>", "Results before date"},
		{"--after <YYYY-MM-DD>", "Results after date"},
		{"-q, --query <text>", "Raw keyword / search query"},
		{"--exclude <text>", "Exclude a term (prefixes -)"},
		{"--or <text>", "OR with another term"},
		{"--exact <text>", "Exact phrase match"},
	}
	for _, f := range dorkFlags {
		fmt.Printf("  %-28s%s\n", green+f[0]+reset, dim+f[1]+reset)
	}
	fmt.Println()

	fmt.Println(bold + cyan + "  AUTO RECON CATEGORIES" + reset)
	fmt.Println(dim + "  ─────────────────────" + reset)
	cats := []string{"Admin Panels", "Sensitive Files", "Credentials", "Directories", "Error & Debug", "Infrastructure", "Documents", "External Leaks"}
	for _, c := range cats {
		fmt.Println("  " + categoryColor(c) + "▸ " + c + reset)
	}
	fmt.Println()

	fmt.Println(bold + cyan + "  EXAMPLES" + reset)
	fmt.Println(dim + "  ────────" + reset)
	examples := [][]string{
		{"gdork -u example.com", fmt.Sprintf("Full auto recon — all %d dorks", len(autoReconDorks))},
		{"gdork -u example.com -c credentials", "Only credential dorks"},
		{"gdork -u example.com -c \"admin panels\"", "Only admin panel dorks"},
		{"gdork -s example.com --filetype pdf", "Manual: PDFs on a site"},
		{"gdork --inurl admin --intitle login", "Manual: admin login pages"},
		{"gdork --intext DB_PASSWORD --ext env", "Manual: exposed .env files"},
		{"gdork -u example.com | grep -A2 Credentials", "Pipe & grep a category"},
	}
	for _, e := range examples {
		fmt.Println("  " + magenta + "→" + reset + " " + white + e[0] + reset)
		fmt.Println("    " + dim + e[1] + reset)
		fmt.Println()
	}
}

func listOperators() {
	banner()
	fmt.Println(bold + cyan + "  GOOGLE DORK OPERATORS" + reset)
	fmt.Println(dim + "  ─────────────────────" + reset)
	fmt.Println()

	categories := []struct {
		name string
		ops  [][]string
	}{
		{"Search Operators", [][]string{
			{"site:", "Limit to domain — site:example.com"},
			{"inurl:", "Term in URL — inurl:admin"},
			{"intitle:", "Term in title — intitle:login"},
			{"intext:", "Term in body — intext:password"},
			{"filetype:", "File type — filetype:pdf"},
			{"ext:", "Alias for filetype — ext:sql"},
			{"cache:", "Cached page — cache:example.com"},
			{"link:", "Pages linking to — link:example.com"},
			{"related:", "Related sites — related:example.com"},
		}},
		{"Date Filters", [][]string{
			{"before:", "Before date — before:2024-01-01"},
			{"after:", "After date — after:2023-01-01"},
		}},
		{"Logic & Modifiers", [][]string{
			{"\"phrase\"", "Exact phrase match"},
			{"-term", "Exclude a term"},
			{"OR", "Either term — login OR signin"},
			{"*", "Wildcard — admin*panel"},
		}},
	}

	for _, cat := range categories {
		fmt.Println("  " + bold + yellow + cat.name + reset)
		for _, op := range cat.ops {
			fmt.Printf("  %-32s%s\n", green+op[0]+reset, dim+op[1]+reset)
		}
		fmt.Println()
	}
}

// ─── Main ─────────────────────────────────────────────────────────────────────

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		help()
		os.Exit(0)
	}

	target := ""
	category := ""
	filtered := []string{}

	for apple := 0; apple < len(args); apple++ {
		arg := args[apple]
		switch arg {
		case "-h", "--help":
			help()
			os.Exit(0)
		case "-v", "--version":
			fmt.Println("gdork v" + version)
			os.Exit(0)
		case "-l", "--list":
			listOperators()
			os.Exit(0)
		case "-u", "--url":
			if apple+1 < len(args) {
				apple++
				target = args[apple]
			} else {
				fmt.Println(red + "  [!] -u requires a domain — e.g. gdork -u example.com" + reset)
				os.Exit(1)
			}
		case "-c", "--category":
			if apple+1 < len(args) {
				apple++
				category = args[apple]
			} else {
				fmt.Println(red + "  [!] -c requires a category name" + reset)
				os.Exit(1)
			}
		default:
			filtered = append(filtered, arg)
		}
	}

	// Auto recon mode
	if target != "" {
		if category != "" {
			autoReconFiltered(target, category)
		} else {
			autoRecon(target)
		}
		os.Exit(0)
	}

	// Manual dork mode
	dork, err := buildDork(filtered)
	if err != nil {
		fmt.Println()
		fmt.Println(red + "  [!] " + err.Error() + reset)
		fmt.Println(dim + "  Run 'gdork --help' for usage." + reset)
		fmt.Println()
		os.Exit(1)
	}

	encodedURL := generateURL(dork)
	readableURL := decodedURL(dork)

	fmt.Println()
	fmt.Println(bold + cyan + "  ┌─ DORK QUERY" + reset)
	fmt.Println("  │  " + yellow + dork + reset)
	fmt.Println(bold + cyan + "  ├─ ENCODED URL" + reset)
	fmt.Println("  │  " + dim + encodedURL + reset)
	fmt.Println(bold + cyan + "  └─ DECODED URL" + reset)
	fmt.Println("     " + cyan + readableURL + reset)
	fmt.Println()
}
