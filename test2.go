package test2

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/fatih/color"
	"github.com/mholt/caddy"
	"github.com/miekg/dns"
)

const (
	ethProvider = "http://localhost:8545"
)

var (
	// TODO: make an actual pipeline/middleware type string of functions
	// this will actually be a function later, but we can still keep/fill the in-memory map as a cache even though
	// its probably better if we just have redis/memcached handle that shitcoredcoredns/plugincoredns/plugincoredns/pluginns/plugin
	domainCache = map[string]Fragment{}

	mutex = &sync.RWMutex{}

	printer = &Printer{
		Normal:  color.New(color.FgGreen),
		Warning: color.New(color.FgYellow),
		Error:   color.New(color.FgRed),
		Other:   color.New(color.FgMagenta),
		Report:  color.New(color.FgCyan),
	}

	// TODO: Change this to be more idiomatic
	errCodes = struct {
		EtherwebNoConnection int
		TCPNoStart           int
		UDPNoStart           int
		WebServerNoStart     int
	}{
		0,
		1,
		2,
		3,
	}
)

type Printer struct {
	Report  *color.Color
	Normal  *color.Color
	Warning *color.Color
	Error   *color.Color
	Other   *color.Color
}

type Fragment struct {
	Address string
	Timer   *time.Timer
}

func init() {
	fmt.Println("hey its me tester2")
	caddy.RegisterPlugin("test2", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

// Name implements the Handler interface.
func (h Handler) Name() string { return "test2" }

func setup(c *caddy.Controller) error {
	c.Next() // 'test2'
	if c.NextArg() {
		return plugin.Error("test2", c.ArgErr())
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Handler{}
	})

	return nil
}

type DomainRequest struct {
	ID         string `json:"id"`
	DomainName string `json:"domainName"`
	IPAddress  string `json:"ipAddress"`
}

var domains []DomainRequest

type Handler struct {
	Type string
}

// ServeDNS...
func (handler Handler) ServeDNS(c context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	msg := dns.Msg{
		MsgHdr: dns.MsgHdr{
			Authoritative:      true,
			RecursionAvailable: true,
			RecursionDesired:   true,
		}}

	msg.SetReply(r)
	msg.Answer = append(msg.Answer, &dns.A{
		Hdr: dns.RR_Header{
			Name:   msg.Question[0].Name,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    299,
		},
		A: net.ParseIP("7.8.9.0"),
	})

	w.WriteMsg(&msg)

	return 0, nil
}
