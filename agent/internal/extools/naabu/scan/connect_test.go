package scan

import (
	"net"
	"testing"

	"github.com/omarbdrn/simple_agent/internal/extools/naabu/port"
	"github.com/omarbdrn/simple_agent/internal/extools/naabu/protocol"
	"github.com/stretchr/testify/assert"
)

func TestConnectVerify(t *testing.T) {
	go func() {
		// start tcp server
		l, err := net.Listen("tcp", ":17895")
		if err != nil {
			assert.Nil(t, err)
		}
		defer l.Close()
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			defer conn.Close()
		}
	}()

	s, err := NewScanner(&Options{})
	assert.Nil(t, err)
	wanted := []*port.Port{
		{Port: 17895, Protocol: protocol.TCP},
	}

	targetPorts := []*port.Port{
		{Port: 17895, Protocol: protocol.TCP},
		{Port: 17896, Protocol: protocol.TCP},
	}
	got := s.ConnectVerify("localhost", targetPorts)
	assert.EqualValues(t, wanted, got)
}
