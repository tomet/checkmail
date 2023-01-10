package main

import (
	"fmt"

	"github.com/emersion/go-imap"
)

//--------------------------------------------------------------------------------
// count [all|seen|unseen]
//--------------------------------------------------------------------------------

func countCmd(cfg *Config, arg string) {
	showSeen := false
	showUnseen := false

	switch arg {
	case "", "all":
		showSeen = true
		showUnseen = true
	case "seen":
		showSeen = true
	case "unseen":
		showUnseen = true
	default:
		syntax("invalid argument for count: %s", arg)
	}

	cl := login(cfg.Server, cfg.User, cfg.Password)
	defer logout(cl)

	selectMailbox(cl, cfg.Mailbox)

	var seen, unseen []uint32
	if showSeen {
		seen = searchSeen(cl)
	}
	if showUnseen {
		unseen = searchUnseen(cl)
	}

	if opts.Verbose {
		info("%s of %s:", addrColor(cfg.Mailbox), addrColor(cfg.User))
		if showSeen {
			info("  seen: %4s", formatCount(len(seen), seenColor))
		}
		if showUnseen {
			info("unseen: %4s", formatCount(len(unseen), unseenColor))
		}
		if showSeen && showUnseen {
			info(" total: %4s", formatCount(len(seen)+len(unseen), totalColor))
		}
	} else {
		switch {
		case showSeen && showUnseen:
			info(
				"%s/%s",
				formatCount(len(unseen), unseenColor),
				formatCount(len(unseen)+len(seen), totalColor),
			)
		case showSeen:
			info("%s", formatCount(len(seen), seenColor))
		default:
			info("%s", formatCount(len(unseen), unseenColor))
		}
	}
}

func formatCount(count int, colorFn func(...any) string) string {
	return colorFn(fmt.Sprintf("%d", count))
}

//--------------------------------------------------------------------------------
// list [boxes|seen|unseen|all]
//--------------------------------------------------------------------------------

func listCmd(cfg *Config, arg string) {
	searchFn := searchAll

	switch arg {
	case "mailboxes", "boxes", "mailbox", "box":
		listMailboxes(cfg)
		return
	case "unseen", "":
		searchFn = searchUnseen
	case "seen":
		searchFn = searchSeen
	case "all":
		searchFn = searchAll
	default:
		syntax("invalid argument for list: %q", arg)
	}

	cl := login(cfg.Server, cfg.User, cfg.Password)
	defer logout(cl)

	selectMailbox(cl, cfg.Mailbox)
	envs := getMessageEnvelopes(cl, searchFn(cl))

	totalW := terminalWidth()
	fromW := 0
	dateW := 10
	maxFromW := 40

	for _, env := range envs {
		from := formatFirstAddress(env.From)
		fromW = maxOf(fromW, utf8Len(from))
	}

	if fromW > maxFromW {
		fromW = maxFromW
	}

	subjW := totalW - (fromW + 1 + dateW + 1)
	subjW = maxOf(subjW, 15)

	for _, env := range envs {
		date := dateColor(env.Date.Format("02.01.2006"))
		from := padLeft(formatFirstAddress(env.From), fromW)
		subj := trimToLen(env.Subject, subjW)
		info("%s %s %s", dateColor(date), addrColor(from), subj)
	}
}

func formatFirstAddress(addrs []*imap.Address) string {
	first := ""

	if len(addrs) > 0 {
		first = addrs[0].Address()
		if len(addrs) > 1 {
			first += ",+"
		}
	}

	return first
}

func listMailboxes(cfg *Config) {
	cl := login(cfg.Server, cfg.User, cfg.Password)
	defer logout(cl)
	
	boxes := searchMailboxes(cl)
	
	verbose("Mailboxes of %s:", cfg.User)
	for _, box := range boxes {
		info("%s", box)
	}
}

//--------------------------------------------------------------------------------
// touch
//--------------------------------------------------------------------------------

func touchCmd(cfg *Config, arg string) {
	if arg != "" {
		syntax("command 'touch' required NO argument")
	}

	cl := login(cfg.Server, cfg.User, cfg.Password)
	defer logout(cl)
	
	selectMailbox(cl, cfg.Mailbox)
	unseen := searchUnseen(cl)

	if len(unseen) == 0 {
		info("No unseen messages to mark.")
		return
	}

	markSeen(cl, unseen)

	info("%d messages marked as seen", len(unseen))
}
