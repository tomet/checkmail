package main

import (
	"fmt"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

//--------------------------------------------------------------------------------
// count [all|seen|unseen]
//--------------------------------------------------------------------------------

func countCmd(cl *client.Client, arg string) {
	var seen, unseen []uint32
	showSeen := false
	showUnseen := false

	selectMailbox(cl, cfg.Mailbox)

	switch arg {
	case "", "all":
		seen = searchSeen(cl)
		unseen = searchUnseen(cl)
		showSeen = true
		showUnseen = true
	case "seen":
		seen = searchSeen(cl)
		showSeen = true
	case "unseen":
		unseen = searchUnseen(cl)
		showUnseen = true
	default:
		syntax("invalid argument for count: %s", arg)
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

func listCmd(cl *client.Client, arg string) {
	var envs []*imap.Envelope

	switch arg {
	case "mailboxes", "boxes", "mailbox", "box":
		boxes := searchMailboxes(cl)
		verbose("Mailboxes of %s:", cfg.User)
		for _, box := range boxes {
			info("%s", box)
		}
		return
	case "unseen", "":
		selectMailbox(cl, cfg.Mailbox)
		envs = getMessageEnvelopes(cl, searchUnseen(cl))
	case "seen":
		selectMailbox(cl, cfg.Mailbox)
		envs = getMessageEnvelopes(cl, searchSeen(cl))
	case "all":
		selectMailbox(cl, cfg.Mailbox)
		envs = getMessageEnvelopes(cl, searchAll(cl))
	default:
		syntax("invalid argument for list: %q", arg)
	}

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

//--------------------------------------------------------------------------------
// touch
//--------------------------------------------------------------------------------

func touchCmd(cl *client.Client, arg string) {
	if arg != "" {
		syntax("command 'touch' required NO argument")
	}
	
	selectMailbox(cl, cfg.Mailbox)
	unseen := searchUnseen(cl)
	
	if len(unseen) == 0 {
		info("No unseen messages to mark.")
		return
	}
	
	markSeen(cl, unseen)
	
	info("%d messages marked as seen", len(unseen))
}
