package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/pflag"
)

type Config struct {
	Server   string
	User     string
	Password string
	Mailbox  string
	NoTls    bool `ini:"no_tls"`
}

type Opts struct {
	NoColors bool
	Verbose  bool
	Debug    bool
	File     string

	Server  string
	User    string
	Mailbox string
	NoTls   bool
}

const (
	Version           = "0.1"
	DefaultConfigFile = "~/.config/checkmail/checkmail.ini"
)

var (
	opts Opts
	cfg  *Config

	dateColor = color.New(color.FgYellow).SprintFunc()
	addrColor = color.New(color.FgGreen).SprintFunc()

	totalColor  = color.New(color.FgGreen).SprintFunc()
	unseenColor = color.New(color.FgRed).SprintFunc()
	seenColor   = color.New(color.FgHiBlack).SprintFunc()
)

//--------------------------------------------------------------------------------
// main()
//--------------------------------------------------------------------------------

func main() {
	cmd, arg := parseCmdline()

	if opts.File == "" {
		cfg = loadConfigFile(DefaultConfigFile, false)
	} else {
		cfg = loadConfigFile(opts.File, true)
	}

	prepareConfig(cfg, &opts)

	cl := connect(cfg.Server)
	defer logout(cl)

	login(cl, cfg.User, cfg.Password)

	switch cmd {
	case "", "count":
		countCmd(cl, arg)
	case "list":
		listCmd(cl, arg)
	case "touch":
		touchCmd(cl, arg)
	default:
		syntax("unknown command: %q", cmd)
	}
}

func prepareConfig(cfg *Config, opts *Opts) {
	// cmdline-options overwrite config-file
	if opts.Server != "" {
		cfg.Server = opts.Server
	}
	if opts.User != "" {
		cfg.User = opts.User
	}
	if opts.Mailbox != "" {
		cfg.Mailbox = opts.Mailbox
	}
	if opts.NoTls == true {
		cfg.NoTls = true
	}

	// check for required config-values
	if cfg.Server == "" {
		syntax("No server given (in config-file or with --server)")
	}
	if cfg.User == "" {
		syntax("No user given (in config-file or with --user)")
	}
	if cfg.Password == "" {
		syntax("No password given (in config-file)")
	}

	// default-values (don't set option-default!!!)
	if cfg.Mailbox == "" {
		cfg.Mailbox = "INBOX"
	}
}

//--------------------------------------------------------------------------------
// Parse commandline
//--------------------------------------------------------------------------------

func parseCmdline() (cmd, arg string) {
	pflag.CommandLine.Init("", pflag.ContinueOnError)

	pflag.StringVarP(&opts.Server, "server", "s", "", "")
	pflag.StringVarP(&opts.User, "user", "u", "", "")
	pflag.BoolVarP(&opts.Verbose, "verbose", "v", false, "")
	pflag.BoolVarP(&opts.Debug, "debug", "d", false, "")
	pflag.BoolVarP(&opts.NoTls, "no-tls", "t", false, "")
	pflag.StringVarP(&opts.Mailbox, "mailbox", "m", "", "")
	pflag.StringVarP(&opts.File, "file", "f", "", "")

	noColors := pflag.BoolP("no-colors", "n", false, "")
	help := pflag.BoolP("help", "?", false, "")
	version := pflag.BoolP("version", "V", false, "")

	if err := pflag.CommandLine.Parse(os.Args[1:]); err != nil {
		syntax(err.Error())
	}

	if *version {
		info("checkmail v" + Version)
		os.Exit(0)
	}

	if *help {
		printHelp()
	}

	if *noColors {
		color.NoColor = true
	}

	args := pflag.Args()

	switch len(args) {
	case 0:
	case 1:
		cmd = args[0]
	case 2:
		cmd = args[0]
		arg = args[1]
	default:
		syntax("too many arguments")
	}

	return cmd, arg
}

func printHelp() {
	fmt.Println(`Usage: checkmail [OPTIONS] [COMMAND]

Check an IMAP-server for un-/seen messages and count or list them.

To mark all messages as seen use the command 'touch'.

The command 'list boxes' shows all available mailboxes.
	
Options:
  -s, --server=HOST:PORT select the IMAP-server and port
  -t, --no-tls           use un-encryptet connection
  -u, --user=USER        login as USER
  -m, --mailbox=BOX      select mailbox (default is INBOX)
  -v, --verbose          show verbose output
  -d, --debug            show debugging information
  -n, --no-colors        disable ANSI-Colors
  -f, --file=CFGFILE     load another config-file
                         (default is 
						  ~/.config/checkmail/checkmail.ini)
		
  
Commands:
  list [boxes|seen|unseen|all]  list messages or mailboxes
  count [seen|unseen|all]       count messages in Mailbox
  touch                         mark all messages as seen
`)

	os.Exit(0)
}

//--------------------------------------------------------------------------------
// Output
//--------------------------------------------------------------------------------

func debug(format string, args ...any) {
	if opts.Debug {
		fmt.Printf(format, args...)
		fmt.Println("")
	}
}

func verbose(format string, args ...any) {
	if opts.Verbose {
		fmt.Printf(format, args...)
		fmt.Println("")
	}
}

func info(format string, args ...any) {
	fmt.Printf(format, args...)
	fmt.Println("")
}

func warning(format string, args ...any) {
	fmt.Fprint(os.Stderr, "checkmail: ")
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintln(os.Stderr, "\n")
}

func failure(format string, args ...any) {
	warning(format, args...)
	os.Exit(1)
}

func syntax(format string, args ...any) {
	warning(format, args...)
	fmt.Fprintln(os.Stderr, "Try 'checkmail --help' for more information.\n")
	os.Exit(1)
}
