package keepalive

import (
	"sync"
	"testing"
	"time"

	"github.com/koron/go-mqtt/client"
	"github.com/koron/go-mqtt/server"
)

func TestAutoDisconnect(t *testing.T) {
	var wg sync.WaitGroup
	srv := server.Server{}
	wg.Add(1)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			t.Error("ListenAndServe is failed: ", err)
		}
		wg.Done()
	}()
	time.Sleep(time.Millisecond * 100)
	var disconnReason error
	var l sync.Mutex
	_, err := client.Connect(client.Param{
		OnDisconnect: func(reason error, param client.Param) {
			l.Lock()
			disconnReason = reason
			l.Unlock()
		},
		ID: "keepalive-test",
		Options: &client.Options{
			KeepAlive:            2,
			DisableAutoKeepAlive: true,
		},
	})
	if err != nil {
		t.Fatal("client.Connect failed: ", err)
	}
	time.Sleep(time.Second * 3)
	l.Lock()
	if disconnReason == nil {
		t.Error("not disconnected")
	}
	l.Unlock()
	srv.Close()
	wg.Wait()
}

func TestKeepAlive(t *testing.T) {
	var wg sync.WaitGroup
	srv := server.Server{
		Options: &server.Options{
			DisableMonitor: true,
		},
	}
	wg.Add(1)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			t.Error("ListenAndServe is failed: ", err)
		}
	}()
	time.Sleep(time.Millisecond * 100)
	var disconnReason error
	var l sync.Mutex
	_, err := client.Connect(client.Param{
		OnDisconnect: func(reason error, param client.Param) {
			l.Lock()
			disconnReason = reason
			l.Unlock()
		},
		ID: "keepalive-test",
		Options: &client.Options{
			KeepAlive:            2,
			DisableAutoKeepAlive: true,
		},
	})
	if err != nil {
		t.Fatal("client.Connect failed: ", err)
	}
	time.Sleep(time.Second * 3)
	l.Lock()
	if disconnReason != nil {
		t.Error("disconnected unexpectedly: ", disconnReason)
	}
	l.Unlock()
	srv.Close()
	wg.Done()
	time.Sleep(time.Millisecond * 100)
	l.Lock()
	if disconnReason == nil {
		t.Error("client aliving unexpectedly")
	}
	l.Unlock()
}
