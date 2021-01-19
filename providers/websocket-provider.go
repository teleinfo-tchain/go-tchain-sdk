/********************************************************************************
   This file is part of go-web3.
   go-web3 is free software: you can redistribute it and/or modify
   it under the terms of the GNU Lesser General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   go-web3 is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Lesser General Public License for more details.
   You should have received a copy of the GNU Lesser General Public License
   along with go-web3.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/

package providers

import (
	"errors"
	"fmt"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers/util"

	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type WebSocketProvider struct {
	address string
	ws      *stream
}

func NewWebSocketProvider(address string) *WebSocketProvider {
	provider := new(WebSocketProvider)
	ws, err := newWebsocket(address)
	if err != nil {
		panic(fmt.Sprintf("dail websocket faied: %v", err))
	}
	provider.address = address
	provider.ws = ws
	return provider
}

func (provider WebSocketProvider) SendRequest(v interface{}, method string, params interface{}) error {
	value, err := provider.ws.Call(method, v, params)
	if err != nil {
		return err
	}
	if jsonValue, err := json.Marshal(value); err == nil {
		return json.Unmarshal(jsonValue, v)
	} else {
		return err
	}
}

func (provider WebSocketProvider) Close() error {
	if provider.ws != nil {
		return provider.ws.Close()
	}

	return errors.New("websocket connection dont exist")

}

func newWebsocket(url string) (*stream, error) {
	wsConn, _, err := websocket.DefaultDialer.Dial(url, http.Header{})
	if err != nil {
		return nil, err
	}
	codec := &websocketCodec{
		conn: wsConn,
	}
	return newStream(codec)
}

// ErrTimeout happens when the websocket requests times out
var ErrTimeout = fmt.Errorf("timeout")

type ackMessage struct {
	buf []byte
	err error
}

type callback func(b []byte, err error)

type stream struct {
	seq   uint64
	codec Codec

	// call handlers
	handlerLock sync.Mutex
	handler     map[uint64]callback

	closeCh chan struct{}
	timer   *time.Timer
}

func newStream(codec Codec) (*stream, error) {
	w := &stream{
		codec:   codec,
		closeCh: make(chan struct{}),
		handler: map[uint64]callback{},
	}

	go w.listen()
	return w, nil
}

// Close implements the the transport interface
func (s *stream) Close() error {
	close(s.closeCh)
	return s.codec.Close()
}

func (s *stream) incSeq() uint64 {
	return atomic.AddUint64(&s.seq, 1)
}

func (s *stream) isClosed() bool {
	select {
	case <-s.closeCh:
		return true
	default:
		return false
	}
}

func (s *stream) listen() {
	buf := []byte{}

	for {
		var err error
		buf, err = s.codec.Read(buf[:0])
		if err != nil {
			if !s.isClosed() {
				// log error
			}
			return
		}

		var resp Response
		if err = json.Unmarshal(buf, &resp); err != nil {
			return
		}

		if resp.ID != 0 {
			go s.handleMsg(resp)
		}
	}
}

func (s *stream) handleMsg(response Response) {
	s.handlerLock.Lock()
	callback, ok := s.handler[response.ID]
	if !ok {
		s.handlerLock.Unlock()
		return
	}

	// delete handler
	delete(s.handler, response.ID)
	s.handlerLock.Unlock()

	if response.Error != nil {
		callback(nil, response.Error)
	} else {
		callback(response.Result, nil)
	}
}

func (s *stream) setHandler(id uint64, ack chan *ackMessage) {
	callback := func(b []byte, err error) {
		select {
		case ack <- &ackMessage{b, err}:
		default:
		}
	}

	s.handlerLock.Lock()
	s.handler[id] = callback
	s.handlerLock.Unlock()

	s.timer = time.AfterFunc(5*time.Second, func() {
		s.handlerLock.Lock()
		delete(s.handler, id)
		s.handlerLock.Unlock()

		select {
		case ack <- &ackMessage{nil, ErrTimeout}:
		default:
		}
	})
}

// Call implements the transport interface
func (s *stream) Call(method string, out interface{}, params interface{}) (interface{}, error) {
	seq := s.incSeq()
	request := Request{
		ID:     seq,
		Method: method,
	}
	if params != nil {
		data, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		request.Params = data
	}

	ack := make(chan *ackMessage)
	s.setHandler(seq, ack)

	raw, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	if err := s.codec.Write(raw); err != nil {
		return nil, err
	}

	resp := <-ack
	if resp.err != nil {
		return nil, resp.err
	}

	return dto.RequestResult {
		ID: int(seq),
		Version: util.Version,
		Result: resp.buf,
	}, nil

}

type websocketCodec struct {
	conn *websocket.Conn
}

func (w *websocketCodec) Close() error {
	return w.conn.Close()
}

func (w *websocketCodec) Write(b []byte) error {
	return w.conn.WriteMessage(websocket.TextMessage, b)
}

func (w *websocketCodec) Read(b []byte) ([]byte, error) {
	_, buf, err := w.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	b = append(b, buf...)
	return b, nil
}

// Request is a jsonrpc request
type Request struct {
	ID     uint64          `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

// Response is a jsonrpc response
type Response struct {
	ID     uint64          `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  *ErrorObject    `json:"error,omitempty"`
}

// ErrorObject is a jsonrpc error
type ErrorObject struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Error implements error interface
func (e *ErrorObject) Error() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("jsonrpc.internal marshal error: %v", err)
	}
	return string(data)
}

// Codec is the codec to write and read messages
type Codec interface {
	Read([]byte) ([]byte, error)
	Write([]byte) error
	Close() error
}