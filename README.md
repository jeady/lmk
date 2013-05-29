LMK!
====
Let Me Know! is a user-facing notification service. Want to be the first to
know when the lineout for ACL goes live? Woot.com is selling a macbook air?
George R. R. Martin has blogged about the next ASOIAF book? LMK! is the answer.
LMK! will notify you by email when an event you've set occurs.

Current Capabilities
--------------------
Ha! Coming soon...

Configuring
-----------
To configure LMK!, you will add rules to lmk.cfg specifying the notifications
you wish to receive. Regardless of the actual method used to check the
notification, all rules must contain at least 3 parts:

- a name to identify the rule by
- a sanity check to ensure the rule is functioning
- a trigger that, when passed, signals that the user should be notified

The sanity check and trigger are typically either a string or a regular
expression.

The following is an example rule that will notify the user when the fourth book
of the Malazan Book of the Fallen series, 'House of Chains', is released on
Audible.

```
[Malazan Book of the Fallen 4]
url=http://www.audible.com/search?advsearchKeywords=malazan
sanity=Memories of Ice
trigger=House of Chains
enabled=yes
```

Contributing
------------
Well, there's not much to contribute to yet, but I am always looking for ideas
to add as triggers (or anything else). Please do feel free to open an issue
tagged as 'question' or 'enhancement' if you'd like to comment or suggest.
