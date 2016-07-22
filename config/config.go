package config

// Path to Bolt database.
// Will be created if not exists.
var BoltDBPath = "hack.db"

// Name default bucket.
// Will be created along with the DB for the first time.
var DefaultBucketName = "snippets"

// Where to store FS files.
// The name of a directory relative to this working directory.
// ... This lets you store your snippets on the server as actual files,
// from which you can even use something like Syncthing to sync with your
// own computer.
var FSStorePath = "hacks"

// Where to store chat messages (just a plain ol .txt).
var ChatFile = "data/chat.txt"

// Which port would you like cathack to run on?
var MakeThisMyPort = ":5000"
