package transport

import (
	"fmt"
	"os"
	"testing"
	"time"
)

const (
	dealerUri = "tcp://127.0.0.1:5555"
	routerUri = "tcp://127.0.0.1:5556"
)

// Dummy interface for testing purposes.
type RunWrap struct {
}

func (r *RunWrap) GetResponse(msg string) string {
	respMsg := "Resp " + msg
	fmt.Println("Received messsage in RUN" + msg)
	return respMsg
}

func startBroker(t *testing.T) {
	err := StartBroker(dealerUri, routerUri)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func listen(t *testing.T) {
	limitChan := make(chan int, 10)
	for {
		limitChan <- 1
		go func() {
			respConn, err := NewRespConn(dealerUri)
			if err != nil {
				t.Errorf(err.Error())
			}
			r := new(RunWrap)
			err = respConn.Listen(r, -1)
			if err != nil {
				t.Errorf(err.Error())
			}
			// Should never get called.
			respConn.Close()
			<-limitChan
		}()
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

/*
func TestSimpleReqResp(t *testing.T) {
	url := "tcp://localhost:5556"
	reqConn, err := NewReqConn(url)
	if err != nil {
		t.Errorf(err.Error())
	}
	req := "abcd"
	resp, err := reqConn.MakeReq(req)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Log(resp + " obtained for req " + req)
}
*/

func TestLife(t *testing.T) {
	t.Log("Test transport")
	go startBroker(t)
	go listen(t)

	t.Log("Done starting everything")
	time.Sleep(time.Second)

	reqConn, err := NewReqConn(routerUri)
	if err != nil {
		t.Errorf(err.Error())
	}

	req := "abcd"
	resp, err := reqConn.MakeReq(req)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Log(resp + " obtained for req " + req)
}
