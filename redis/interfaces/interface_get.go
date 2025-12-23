package interfaces

import "github.com/redis-go/redcon"

// قدرة قراءة المفتاح مع expiration
type KeyReader interface {
	GetOrExpire(key string, deleteIfExpired bool) Item
}

// قدرة الرد على العميل
type Responder interface {
	WriteNull()
	WriteError(msg string)
	WriteBulkString(v string)
}

type RedconResponder struct {
	conn redcon.Conn
}

func (r *RedconResponder) WriteNull() {
	r.conn.WriteNull()
}

func (r *RedconResponder) WriteError(msg string) {
	r.conn.WriteError(msg)
}

func (r *RedconResponder) WriteBulkString(v string) {
	r.conn.WriteBulkString(v)
}

// Compile-time check
var _ Responder = (*RedconResponder)(nil)

// قدرة قراءة arguments بأمان
type CommandArgs interface {
	Arg(n int) (string, bool)
}
type RedconArgs struct {
	cmd redcon.Command
}

func (a *RedconArgs) Arg(n int) (string, bool) {
	if n >= len(a.cmd.Args) {
		return "", false
	}
	return string(a.cmd.Args[n]), true
}

var _ CommandArgs = (*RedconArgs)(nil)

type CommandArgs2 interface {
	Arg(n int) (string, bool)
	Len() int
}

type Responder2 interface {
	WriteNull()
	WriteError(msg string)
	WriteBulkString(v string)
}

type DbOps interface {
	GetOrExpire(key string, deleteIfExpired bool) Item
}
