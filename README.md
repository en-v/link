# LINK

Bi-directional RPC library based on WebScoket and written in Go

### Installation
    go get github.com/en-v/link

## Qucik start

Server example ()

```go
port := 8040
s := &Server{
	Clients: make(map[string]string),
	ID:     "ID_STRING",
}

link, err := link.New(s)
if err != nil {
	panic(err)
}

err = link.Open(port)
if err != nil {
	panic(err)
}

////////////////////////////////////

func (s *Server) ServerMethod(params *types.Payload) (*types.Payload, error) {
    res := link.Payload()
    res.Set("mirrored_params", params)
    return res, nil
}
```

Client

```go 
token := "123"
client := &Client{
	ID: "ID_STRING",
}

clientLink, err := link.New(client)
if err != nil {
	panic(err)
}

err = clientLink.Connect("ws://localhost:8040", token)
if err != nil {
	panic(err)
}

params := link.Payload()
params.Set("param1", client.ID.String())
params.Set("param2", client.ID.String())

res, err := clientLink.Invoke("ServerMethod", params)	
log.Debugw("Method called on the client", "Result", res, "Error", err)
if err != nil {
	log.Error(err)
}
```