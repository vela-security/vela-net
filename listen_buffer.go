package vnet

import (
	"github.com/vela-security/vela-public/buffer"
	"github.com/vela-security/vela-public/kind"
	"github.com/vela-security/vela-public/lua"
)

type RevBuffer struct {
	buf  []byte
	rev  *buffer.Byte
	hdp  lua.P
	cnn  kind.Conn
	co   *lua.LState
	over bool
}

var line1 = []byte("\r\n")
var line2 = []byte("\n")

func (r *RevBuffer) append(n int) {
	if r.hdp.Fn == nil {
		return
	}

	if n == 0 {
		return
	}

	r.buf = r.buf[:n]
	r.rev.Write(r.buf)
}

func (r *RevBuffer) readline(n int) {
	if n == 0 {
		return
	}

	if r.buf[n-1] == '\n' {
		r.call()
	}
}

func (r *RevBuffer) call() {
	if r.rev.Len() == 0 {
		return
	}

	if r.hdp.Fn == nil {
		return
	}

	err := r.co.CallByParam(r.hdp, lua.B2L(r.rev.B), r.cnn)
	r.over = true
	r.rev.Reset()
	if err != nil {
		xEnv.Errorf("%s handle listen accept fail %v", r.co.CodeVM(), err)
	}
}
