# LINK

Bi-directional RPC library based on WebScoket and written in Go

### Installation
    go get github.com/en-v/link

## Qucik start

Server example ()

```go
	s := &Server{
		Clients: make(map[string]string),
		ID:      goid.New(),
	}

	link, err := link.New(s)
	if err != nil {
		panic(err)
	}

	err = link.Open(8040)
	if err != nil {
		panic(err)
	}

    ////////////////////////////////////

    func (s *Server) ServerMethod(params *types.Payload) (*types.Payload, error) {
        res := link.Payload()
        res.Set("mirrored_params", params)
        log.Debug("Server method", "Params", params)
        return res, nil
    }
```

Client

```go 
    token := "123"
	client := &Client{
		ID: goid.New(),
	}

	clientLink, err := link.New(client)
	if err != nil {
		panic(err)
	}

	err = clientLink.Connect("ws://localhost:8040", token)
	if err != nil {
		panic(err)
	}

    res, err := clientLink.Invoke("ServerMethod", params)	
	log.Debugw("Method called on client", "Result", res, "Error", err)
	if err != nil {
		log.Error("Connection with server is lost")
	}
```