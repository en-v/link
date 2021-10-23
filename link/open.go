package link

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/en-v/link/core"
	"github.com/en-v/log"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

func (self *Link) Open(gatePort int) error {
	return self.open("", gatePort, "", "")
}

func (self *Link) OpenLocal(gatePort int) error {
	return self.open("localhost", gatePort, "", "")
}

func (self *Link) OpenSecure(gatePort int, certFile string, keyFile string) error {
	return self.open("", gatePort, certFile, keyFile)
}

func (self *Link) open(netif string, gatePort int, certFile string, keyFile string) error {

	if core.DEBUG {
		log.Debugw("Link-Gate is prepering...", "Alias", self.state.Alias, "Port", gatePort)
	}

	if gatePort < 1 || gatePort > 65535 {
		return errors.New("Local TCP port value out of range")
	}

	err := self.verifyLinkAsServer()
	if err != nil {
		return err
	}

	localAddress := netif + ":" + strconv.Itoa(gatePort)
	self.state.Net.LocalGatePort = gatePort

	self.state.SetGateMode()

	self.httpsrv = &http.Server{Addr: localAddress}
	handler := http.NewServeMux()
	handler.HandleFunc("/", self.handleIncomingConnectionFromClients)
	self.httpsrv.Handler = handler
	self.upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	err = self.runHttp(certFile, keyFile)
	if err != nil {
		return err
	}

	if core.DEBUG {
		log.Debugw("Link-gate is ready", "Alias", self.state.Alias, "Port", gatePort)
	}
	return nil
}

func (self *Link) runHttp(certFile string, keyFile string) error {
	exist, err := self.tlsFilesExist(certFile, keyFile)
	if err != nil {
		return err
	}

	if exist {
		go func() {
			self.state.Net.CertFile = certFile
			self.state.Net.KeyFile = keyFile
			err = self.httpsrv.ListenAndServeTLS(self.state.Net.CertFile, self.state.Net.KeyFile)
		}()
	} else {
		go func() {
			err = self.httpsrv.ListenAndServe()
		}()
	}
	go self.state.DropFallenConnections()
	time.Sleep(time.Second)
	return nil
}

func (self *Link) handleIncomingConnectionFromClients(w http.ResponseWriter, r *http.Request) {

	socket, err := self.upgrader.Upgrade(w, r, nil)
	if err != nil {
		err = errors.Wrap(err, "HandleIncomingConnectionFromClients.UpgradeWS")
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		socket.Close()
		return
	}

	conn, err := self.handshakeFromCaller(socket)
	if err != nil {
		err = errors.Wrap(err, "HandleIncomingConnectionFromClients.handshakeFromCaller")
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		socket.Close()
		return
	}

	d, err := net.ResolveTCPAddr("tcp", r.RemoteAddr)
	if err != nil {
		err = errors.Wrap(err, "HandleIncomingConnectionFromClients.SetConnection")
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		socket.Close()
		return
	}
	conn.IP = d

	err = self.state.SetConnection(conn)
	if err != nil {
		err = errors.Wrap(err, "HandleIncomingConnectionFromClients.SetConnection")
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		socket.Close()
		return
	}

	conn.Enable()
	go conn.Listen(core.GATE_MODE)

	if self.state.CheckIn != nil {
		err = self.state.CheckIn(conn.RemId)
		if err != nil {
			err = errors.Wrap(err, "Client register extenal method returned error")
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			socket.Close()
			return
		}
	}
}
