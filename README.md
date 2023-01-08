# Checkmail

A simple command-line tool to check the mail's in a IMAP-mailbox.

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
