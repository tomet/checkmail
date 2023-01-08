# Checkmail

A simple command-line tool to check the mail's in a IMAP-mailbox.

`checkmail` can count or list the seen and/or unseen messages.

It is also possible to mark all messages as seen (`checkmail touch`).

The configuration for the mailbox is given at the commandline
or loaded from a configuration-file. The default config-file
is ~/.config/checkmail/checkmail.ini.

## Examples

```
% checkmail
2/89        # 2 unseen messages, 89 total

% checkmail list unseen
02.01.2023 sample@mailme.co   This is really important
03.01.2023 user@othermail.com Not so imporant.

% checkmail touch
2 messages marked as seen.

% checkmail count --verbose
INBOX of me@mymail.com:
  seen: 89
unseen:  0
 total: 89
```
