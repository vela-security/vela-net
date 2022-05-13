package vnet

import (
	"context"
	"fmt"
	"github.com/vela-security/vela-public/auxlib"
	"github.com/vela-security/vela-public/buffer"
	"github.com/vela-security/vela-public/kind"
	"github.com/vela-security/vela-public/lua"
	"io"
	"net"
	"reflect"
	"time"
)

var listenTypeOf = reflect.TypeOf((*listen)(nil)).String()

type listen struct {
	lua.ProcEx
	name   string
	url    auxlib.URL
	ne     *kind.Listener
	co     *lua.LState
	banner string
	hook   *lua.LFunction
	err    error
}

func newLuaListen(L *lua.LState) *listen {
	raw := L.IsString(1)
	banner := L.IsString(2)
	hook := L.IsFunc(3)

	url, err := auxlib.NewURL(raw)
	if err != nil {
		L.RaiseError("parse %s url error %v", raw, err)
	}

	ln := &listen{
		url:    url,
		hook:   hook,
		banner: banner,
		name:   fmt.Sprintf("listen_%d", raw),
		co:     xEnv.Clone(L),
	}
	ln.V(lua.PTInit, listenTypeOf, time.Now())

	return ln
}

func (ln *listen) Name() string {
	return ln.name
}

func (ln *listen) Type() string {
	return listenTypeOf
}

func (ln *listen) Start() error {
	ne, err := kind.Listen(xEnv, ln.url)
	if err != nil {
		return err
	}

	ln.ne = ne
	return nil
}

func (ln *listen) Close() error {
	if ln.ne != nil {
		return ln.ne.Close()
	}
	return nil
}

func (ln *listen) Banner(conn net.Conn) {
	if ln.banner == "" {
		return
	}

	conn.Write(lua.S2B(ln.banner))
}

func (ln *listen) Accept(ctx context.Context, conn net.Conn, stop context.CancelFunc) error {

	rev := &RevBuffer{
		rev: buffer.Get(),
		buf: make([]byte, 4096),
		cnn: kind.NewConn(conn),
		hdp: xEnv.P(ln.hook),
		co:  xEnv.Clone(ln.co),
	}

	ln.Banner(conn)
	defer func() {
		stop()
		_ = conn.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return nil

		default:
			n, err := conn.Read(rev.buf)
			switch err {
			case nil:
				rev.append(n)
				rev.readline(n)

			case io.EOF:
				rev.call()

			default:
				xEnv.Errorf("%s accept error %v", ln.name, err)
				return err
			}
		}
	}

}
