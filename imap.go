package main

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

//--------------------------------------------------------------------------------
// Connect, Login and Logout
//--------------------------------------------------------------------------------

func connect(server string) (cl *client.Client) {
	var err error
	if opts.NoTls {
		debug("Connect to %q (no TLS)", server)
		cl, err = client.Dial(server)
	} else {
		debug("Connect to %q", server)
		cl, err = client.DialTLS(server, nil)
	}
	if err != nil {
		failure("Connection to %q failed: %s", server, err)
	}
	return cl
}

func logout(cl *client.Client) {
	debug("Logging out")
	cl.Logout()
}

func login(cl *client.Client, user, password string) {
	debug("Log in as %q", user)
	if err := cl.Login(cfg.User, password); err != nil {
		failure("Login failed: %s", err)
	}
}

func selectMailbox(cl *client.Client, box string) *imap.MailboxStatus {
	debug("Select mailbox %q", box)
	stat, err := cl.Select(box, false)
	if err != nil {
		failure("Failed to select %q: %s", box, err)
	}
	return stat
}

//--------------------------------------------------------------------------------
// Search for message-ids
//--------------------------------------------------------------------------------

func search(cl *client.Client, descr string, crit *imap.SearchCriteria) []uint32 {
	debug("Search %s messages", descr)
	ids, err := cl.Search(crit)
	if err != nil {
		failure("Search for %s messages failed: %s", descr, err)
	}
	return ids
}

func searchAll(cl *client.Client) []uint32 {
	crit := imap.NewSearchCriteria()
	return search(cl, "all", crit)
}

func searchUnseen(cl *client.Client) []uint32 {
	crit := imap.NewSearchCriteria()
	crit.WithoutFlags = []string{"\\Seen"}
	return search(cl, "unseen", crit)
}

func searchSeen(cl *client.Client) []uint32 {
	crit := imap.NewSearchCriteria()
	crit.WithFlags = []string{"\\Seen"}
	return search(cl, "seen", crit)
}

//--------------------------------------------------------------------------------
// Get messages
//--------------------------------------------------------------------------------

func getMessageEnvelopes(cl *client.Client, ids []uint32) (envs []*imap.Envelope) {
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(ids...)
	
	if len(ids) == 0 {
		info("No messages found.")
		return []*imap.Envelope{}
	}
	
	debug("Fetch %d messages", len(ids))

	what := []imap.FetchItem{imap.FetchEnvelope}

	msgsChan := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- cl.Fetch(seqSet, what, msgsChan)
	}()

	for msg := range msgsChan {
		envs = append(envs, msg.Envelope)
	}

	if err := <-done; err != nil {
		failure("Fetch failed: %s", err)
	}

	return envs
}

//--------------------------------------------------------------------------------
// Get list of mailboxes
//--------------------------------------------------------------------------------

func searchMailboxes(cl *client.Client) (boxes []string) {
	debug("Get list of mailboxes")

	boxesChan := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)

	go func() {
		done <- cl.List("", "*", boxesChan)
	}()

	for box := range boxesChan {
		boxes = append(boxes, box.Name)
	}

	if err := <-done; err != nil {
		failure("Failed to list mailboxes: %s", err)
	}

	return boxes
}

//--------------------------------------------------------------------------------
// Mark messages as seen
//--------------------------------------------------------------------------------

func markSeen(cl *client.Client, ids []uint32) {
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(ids...)
	
	ops := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []any{imap.SeenFlag}
	
	debug("Mark %d messages as seen", len(ids))
	
	err := cl.Store(seqSet, ops, flags, nil)
	if err != nil {
		failure("Marking of %d messages failed: %s", len(ids), err)
	}
}
