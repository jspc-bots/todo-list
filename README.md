# IRC Todo Bot

Creates and manages a todo list for each channel/ user that requests one. Bound to the channel or nick that a new request comes from (so: if the source is a channel then the list is owned by a channel, if the source is a msg/query from a nick, then owned by that nick).

Persists to disk as a gob encoded file
