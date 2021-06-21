# Todo List Bot

Creates and manages a todo list for each channel/ user that requests one. Bound to the channel or nick that a new request comes from (so: if the source is a channel then the list is owned by a channel, if the source is a msg/query from a nick, then owned by that nick).

Persists to disk as a gob encoded file

It works with the following environment variables:

* `$SASL_USER` - the user to connect with
* `$SASL_PASSWORD` - the password to connect with
* `$SERVER` - IRC connection details, as `irc://server:6667` or `ircs://server:6697` (`ircs` implies irc-over-tls)
* `$VERIFY_TLS` - Verify TLS, or sack it off. This is of interest to people, like me, running an ircd on localhost with a self-signed cert. Matches "true" as true, and anything else as false
* `$STORAGE_FILE` - File to persist stored todos in. If it does not exist, the file will be created when the first item is added
* `$TZ` - Timezone to render todo items with
