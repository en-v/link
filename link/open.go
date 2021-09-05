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

func (this *Link) Open(gatePort int) error {
	return this.open("", gatePort, "", "")
}

func (this *Link) OpenLocal(gatePort int) error {
	return this.open("localhost", gatePort, "", "")
}

func (this *Link) OpenSecure(gatePort int, certFile string, keyFile string) error {
	return this.open("", gatePort, certFile, keyFile)
}

func (this *Link) open(netif string, gatePort int, certFile string, keyFile string) error {

	if core.DEBUG {
		log.Debugw("Link-Gate is prepering...", "Alias", this.state.Alias, "Port", gatePort)
	}

	if gatePort < 1 || gatePort > 65535 {
		return errors.New("Local TCP port value out of range")
	}

	err := this.verifyLinkAsServer()
	if err != nil {
		return err
	}

	localAddress := netif + ":" + strconv.Itoa(gatePort)
	this.state.Net.LocalGatePort = gatePort

	this.state.SetGateMode()

	this.httpsrv = &http.Server{Addr: localAddress}
	handler := http.NewServeMux()
	handler.HandleFunc("/", this.handleIncomingConnectionFromClients)
	this.httpsrv.Handler = handler
	this.upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	err = this.runHttp(certFile, keyFile)
	if err != nil {
		return err
	}

	if core.DEBUG {
		log.Debugw("Link-gate is ready", "Alias", this.state.Alias, "Port", gatePort)
	}
	return nil
}

func (this *Link) runHttp(certFile string, keyFile string) error {
	exist, err := this.tlsFilesExist(certFile, keyFile)
	if err != nil {
		return err
	}

	if exist {
		go func() {
			this.state.Net.CertFile = certFile
			this.state.Net.KeyFile = keyFile
			err = this.httpsrv.ListenAndServeTLS(this.state.Net.CertFile, this.state.Net.KeyFile)
		}()
	} else {
		go func() {
			err = this.httpsrv.ListenAndServe()
		}()
	}
	go this.state.DropFallenConnections()
	time.Sleep(time.Second)
	return nil
}

func (this *Link) handleIncomingConnectionFromClients(w http.ResponseWriter, r *http.Request) {

	socket, err := this.upgrader.Upgrade(w, r, nil)
	if err != nil {
		err = errors.Wrap(err, "HandleIncomingConnectionFromClients.UpgradeWS")
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		socket.Close()
		return
	}

	conn, err := this.handshakeFromCaller(socket)
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

	err = this.state.SetConnection(conn)
	if err != nil {
		err = errors.Wrap(err, "HandleIncomingConnectionFromClients.SetConnection")
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		socket.Close()
		return
	}

	conn.Enable()
	go conn.Listen(core.GATE_MODE)

	if this.state.ClientsRegFunc != nil {
		err = this.state.ClientsRegFunc(conn.RemId)
		if err != nil {
			err = errors.Wrap(err, "Client register extenal method returned error")
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			socket.Close()
			return
		}
	}
}
