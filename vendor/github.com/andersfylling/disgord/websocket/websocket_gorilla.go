package websocket

// TODO: if we add any other websocket packages, add build constraints to this file.

import (
	"errors"
	"io"
	"net/http"

	"github.com/andersfylling/disgord/httd"
	"github.com/gorilla/websocket"
	"golang.org/x/net/proxy"
)

func newConn(proxy proxy.Dialer) (Conn, error) {
	return &gorilla{
		proxy: proxy,
	}, nil
}

// rwc is a wrapper for the Conn interface (not net.Conn).
// Interface can be found at https://golang.org/pkg/net/#Conn
// See original code at https://github.com/gorilla/websocket/issues/282
type gorilla struct {
	c     *websocket.Conn
	proxy proxy.Dialer
}

func (g *gorilla) Open(endpoint string, requestHeader http.Header) (err error) {
	// by default we use gorilla's websocket dialer here, but if the passed http client uses a custom transport
	// we make sure we open the websocket over the same transport/proxy, in case the user uses this
	dialer := websocket.DefaultDialer
	if g.proxy != nil {
		dialer = &websocket.Dialer{
			NetDial: g.proxy.Dial,
		}
	}

	// establish ws connection
	g.c, _, err = dialer.Dial(endpoint, requestHeader)
	return
}

func (g *gorilla) WriteJSON(v interface{}) (err error) {
	// TODO: move unmarshalling out of here?
	var w io.WriteCloser
	w, err = g.c.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	err = httd.JSONEncode(w, v)
	return
}

func (g *gorilla) Close() (err error) {
	err = g.c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	err2 := g.c.Close()
	g.c = nil

	if err == nil && err2 != nil {
		return err2
	}
	return err
}

func (g *gorilla) Read() (packet []byte, err error) {
	if g.Disconnected() {
		// this gets triggered when losing internet connection -> trying to reconnect for a while -> re-establishing a connection
		// as discord then sends a invalid session package and disgord tries to reconnect again, a panic takes place.
		// this check is a tmp hack to fix that, as the actual issue is not clearly understood/defined yet.
		err = errors.New("no connection is established. Can not read new messages")
		return
	}
	var messageType int
	messageType, packet, err = g.c.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			err = &ErrorUnexpectedClose{
				info: err.Error(),
			}
		}

		return
	}

	if messageType == websocket.BinaryMessage {
		packet, err = decompressBytes(packet)
	}
	return
}

func (g *gorilla) Disconnected() bool {
	return g.c == nil
}

var _ Conn = (*gorilla)(nil)
