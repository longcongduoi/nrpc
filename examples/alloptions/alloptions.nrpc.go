// This code was autogenerated from alloptions.proto, do not edit.
package main

import (
	"context"
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	nats "github.com/nats-io/go-nats"
	github_com_rapidloop_nrpc "github.com/rapidloop/nrpc"
	"github.com/rapidloop/nrpc"
)

// SvcCustomSubjectServer is the interface that providers of the service
// SvcCustomSubject should implement.
type SvcCustomSubjectServer interface {
	MtSimpleReply(ctx context.Context, req StringArg) (resp SimpleStringReply, err error)
	MtVoidReply(ctx context.Context, req StringArg) (err error)
	MtStreamedReply(ctx context.Context, req StringArg, pushRep func(SimpleStringReply)) (err error)
	MtVoidReqStreamedReply(ctx context.Context, pushRep func(SimpleStringReply)) (err error)
}

// SvcCustomSubjectHandler provides a NATS subscription handler that can serve a
// subscription using a given SvcCustomSubjectServer implementation.
type SvcCustomSubjectHandler struct {
	ctx    context.Context
	nc     nrpc.NatsConn
	server SvcCustomSubjectServer
}

func NewSvcCustomSubjectHandler(ctx context.Context, nc nrpc.NatsConn, s SvcCustomSubjectServer) *SvcCustomSubjectHandler {
	return &SvcCustomSubjectHandler{
		ctx:    ctx,
		nc:     nc,
		server: s,
	}
}

func (h *SvcCustomSubjectHandler) Subject() string {
	return "root.*.custom_subject.>"
}

func (h *SvcCustomSubjectHandler) MtNoRequestPublish(pkginstance string, msg SimpleStringReply) error {
	rawMsg, err := nrpc.Marshal("protobuf", &msg)
	if err != nil {
		log.Printf("SvcCustomSubjectHandler.MtNoRequestPublish: error marshaling the message: %s", err)
		return err
	}
	subject := "root." + pkginstance + "."+ "custom_subject."+ "mtnorequest"
	return h.nc.Publish(subject, rawMsg)
}

func (h *SvcCustomSubjectHandler) MtStreamedReplyHandler(ctx context.Context, tail []string, msg *nats.Msg) {
	_, encoding, err := nrpc.ParseSubjectTail(0, tail)
	if err != nil {
		log.Printf("SvcCustomSubject: MtStreamedReply subject parsing failed:")
	}
	var req StringArg
	if err := nrpc.Unmarshal(encoding, msg.Data, &req); err != nil {
		// Handle error
		return
	}

	ctx, cancel := context.WithCancel(ctx)

	keepStreamAlive := nrpc.NewKeepStreamAlive(h.nc, msg.Reply, encoding, cancel)

	var msgCount uint32

	_, nrpcErr := nrpc.CaptureErrors(func() (proto.Message, error) {
		err := h.server.MtStreamedReply(ctx, req, func(rep SimpleStringReply){
				if err = nrpc.Publish(&rep, nil, h.nc, msg.Reply, encoding); err != nil {
					log.Printf("nrpc: error publishing response")
					cancel()
					return
				}
				msgCount++
			})
		return nil, err
	})
	keepStreamAlive.Stop()

	if nrpcErr != nil {
		nrpc.Publish(nil, nrpcErr, h.nc, msg.Reply, encoding)
	} else {
		nrpc.Publish(
			nil, &nrpc.Error{Type: nrpc.Error_EOS, MsgCount: msgCount},
			h.nc, msg.Reply, encoding)
	}
}

func (h *SvcCustomSubjectHandler) MtVoidReqStreamedReplyHandler(ctx context.Context, tail []string, msg *nats.Msg) {
	_, encoding, err := nrpc.ParseSubjectTail(0, tail)
	if err != nil {
		log.Printf("SvcCustomSubject: MtVoidReqStreamedReply subject parsing failed:")
	}

	ctx, cancel := context.WithCancel(ctx)

	keepStreamAlive := nrpc.NewKeepStreamAlive(h.nc, msg.Reply, encoding, cancel)

	var msgCount uint32

	_, nrpcErr := nrpc.CaptureErrors(func() (proto.Message, error) {
		err := h.server.MtVoidReqStreamedReply(ctx, func(rep SimpleStringReply){
				if err = nrpc.Publish(&rep, nil, h.nc, msg.Reply, encoding); err != nil {
					log.Printf("nrpc: error publishing response")
					cancel()
					return
				}
				msgCount++
			})
		return nil, err
	})
	keepStreamAlive.Stop()

	if nrpcErr != nil {
		nrpc.Publish(nil, nrpcErr, h.nc, msg.Reply, encoding)
	} else {
		nrpc.Publish(
			nil, &nrpc.Error{Type: nrpc.Error_EOS, MsgCount: msgCount},
			h.nc, msg.Reply, encoding)
	}
}

func (h *SvcCustomSubjectHandler) Handler(msg *nats.Msg) {
	var encoding string
	var noreply bool
	// extract method name & encoding from subject
	pkgParams, _, name, tail, err := nrpc.ParseSubject(
		"root", 1, "custom_subject", 0, msg.Subject)
	if err != nil {
		log.Printf("SvcCustomSubjectHanlder: SvcCustomSubject subject parsing failed: %v", err)
		return
	}

	ctx := h.ctx
	ctx = context.WithValue(ctx, "nrpc-pkg-instance", pkgParams[0])
	// call handler and form response
	var resp proto.Message
	var replyError *nrpc.Error
	switch name {
	case "mt_simple_reply":
		_, encoding, err = nrpc.ParseSubjectTail(0, tail)
		if err != nil {
			log.Printf("MtSimpleReplyHanlder: MtSimpleReply subject parsing failed: %v", err)
			break
		}
		var req StringArg
		if err := nrpc.Unmarshal(encoding, msg.Data, &req); err != nil {
			log.Printf("MtSimpleReplyHandler: MtSimpleReply request unmarshal failed: %v", err)
			replyError = &nrpc.Error{
				Type: nrpc.Error_CLIENT,
				Message: "bad request received: " + err.Error(),
			}
		} else {
			resp, replyError = nrpc.CaptureErrors(
				func()(proto.Message, error){
					innerResp, err := h.server.MtSimpleReply(ctx, req)
					if err != nil {
						return nil, err
					}
					return &innerResp, err
				})
			if replyError != nil {
				log.Printf("MtSimpleReplyHandler: MtSimpleReply handler failed: %s", replyError.Error())
			}
		}
	case "mtvoidreply":
		_, encoding, err = nrpc.ParseSubjectTail(0, tail)
		if err != nil {
			log.Printf("MtVoidReplyHanlder: MtVoidReply subject parsing failed: %v", err)
			break
		}
		var req StringArg
		if err := nrpc.Unmarshal(encoding, msg.Data, &req); err != nil {
			log.Printf("MtVoidReplyHandler: MtVoidReply request unmarshal failed: %v", err)
			replyError = &nrpc.Error{
				Type: nrpc.Error_CLIENT,
				Message: "bad request received: " + err.Error(),
			}
		} else {
			resp, replyError = nrpc.CaptureErrors(
				func()(proto.Message, error){
					var innerResp nrpc.Void
					err := h.server.MtVoidReply(ctx, req)
					if err != nil {
						return nil, err
					}
					return &innerResp, err
				})
			if replyError != nil {
				log.Printf("MtVoidReplyHandler: MtVoidReply handler failed: %s", replyError.Error())
			}
		}
	case "mtnorequest":
		// MtNoRequest is a no-request method. Ignore it.
		return
	case "mtstreamedreply":
		h.MtStreamedReplyHandler(ctx, tail, msg)
		return
	case "mtvoidreqstreamedreply":
		h.MtVoidReqStreamedReplyHandler(ctx, tail, msg)
		return
	default:
		log.Printf("SvcCustomSubjectHandler: unknown name %q", name)
		replyError = &nrpc.Error{
			Type: nrpc.Error_CLIENT,
			Message: "unknown name: " + name,
		}
	}


	if !noreply {
		// encode and send response
		err = nrpc.Publish(resp, replyError, h.nc, msg.Reply, encoding) // error is logged
	} else {
		err = nil
	}
	if err != nil {
		log.Println("SvcCustomSubjectHandler: SvcCustomSubject handler failed to publish the response: %s", err)
	}
}

type SvcCustomSubjectClient struct {
	nc      nrpc.NatsConn
	PkgSubject string
	PkgParaminstance string
	Subject string
	Encoding string
	Timeout time.Duration
}

func NewSvcCustomSubjectClient(nc nrpc.NatsConn, pkgParaminstance string) *SvcCustomSubjectClient {
	return &SvcCustomSubjectClient{
		nc:      nc,
		PkgSubject: "root",
		PkgParaminstance: pkgParaminstance,
		Subject: "custom_subject",
		Encoding: "protobuf",
		Timeout: 5 * time.Second,
	}
}

func (c *SvcCustomSubjectClient) MtSimpleReply(req StringArg) (resp SimpleStringReply, err error) {

	subject := c.PkgSubject + "." + c.PkgParaminstance + "." + c.Subject + "." + "mt_simple_reply"

	// call
	err = nrpc.Call(&req, &resp, c.nc, subject, c.Encoding, c.Timeout)
	if err != nil {
		return // already logged
	}

	return
}

func (c *SvcCustomSubjectClient) MtVoidReply(req StringArg) (err error) {

	subject := c.PkgSubject + "." + c.PkgParaminstance + "." + c.Subject + "." + "mtvoidreply"

	// call
	var resp github_com_rapidloop_nrpc.Void
	err = nrpc.Call(&req, &resp, c.nc, subject, c.Encoding, c.Timeout)
	if err != nil {
		return // already logged
	}

	return
}

func (c *SvcCustomSubjectClient) MtNoRequestSubject(
	
) string {
	return c.PkgSubject + "." + c.PkgParaminstance + "." + c.Subject + "." + "mtnorequest"
}

type SvcCustomSubjectMtNoRequestSubscription struct {
	*nats.Subscription
}

func (s *SvcCustomSubjectMtNoRequestSubscription) Next(timeout time.Duration) (next SimpleStringReply, err error) {
	msg, err := s.Subscription.NextMsg(timeout)
	if err != nil {
		return
	}
	err = nrpc.Unmarshal("protobuf", msg.Data, &next)
	return
}

func (c *SvcCustomSubjectClient) MtNoRequestSubscribeSync(
	
) (sub *SvcCustomSubjectMtNoRequestSubscription, err error) {
	subject := c.MtNoRequestSubject(
		
	)
	natsSub, err := c.nc.SubscribeSync(subject)
	if err != nil {
		return
	}
	sub = &SvcCustomSubjectMtNoRequestSubscription{natsSub}
	return
}

func (c *SvcCustomSubjectClient) MtNoRequestSubscribe(
	
	handler func (SimpleStringReply),
) (sub *nats.Subscription, err error) {
	subject := c.MtNoRequestSubject(
		
	)
	sub, err = c.nc.Subscribe(subject, func(msg *nats.Msg){
		var pmsg SimpleStringReply
		err := nrpc.Unmarshal("protobuf", msg.Data, &pmsg)
		if err != nil {
			log.Printf("SvcCustomSubjectClient.MtNoRequestSubscribe: Error decoding, %s", err)
			return
		}
		handler(pmsg)
	})
	return
}

func (c *SvcCustomSubjectClient) MtNoRequestSubscribeChan(
	
) (<-chan SimpleStringReply, *nats.Subscription, error) {
	ch := make(chan SimpleStringReply)
	sub, err := c.MtNoRequestSubscribe(func (msg SimpleStringReply) {
		ch <- msg
	})
	return ch, sub, err
}

func (c *SvcCustomSubjectClient) MtStreamedReply(
	ctx context.Context,
	req StringArg,
	cb func (context.Context, SimpleStringReply),
) error {
	subject := c.PkgSubject + "." + c.PkgParaminstance + "." + c.Subject + "." + "mtstreamedreply"

	sub, err := nrpc.StreamCall(ctx, c.nc, subject, &req, c.Encoding, c.Timeout)
	if err != nil {
		return err
	}

	var res SimpleStringReply
	for {
		err = sub.Next(&res)
		if err != nil {
			break
		}
		cb(ctx, res)
	}
	if err == nrpc.ErrEOS {
		err = nil
	}
	return err
}

func (c *SvcCustomSubjectClient) MtVoidReqStreamedReply(
	ctx context.Context,
	cb func (context.Context, SimpleStringReply),
) error {
	subject := c.PkgSubject + "." + c.PkgParaminstance + "." + c.Subject + "." + "mtvoidreqstreamedreply"

	sub, err := nrpc.StreamCall(ctx, c.nc, subject, &nrpc.Void{}, c.Encoding, c.Timeout)
	if err != nil {
		return err
	}

	var res SimpleStringReply
	for {
		err = sub.Next(&res)
		if err != nil {
			break
		}
		cb(ctx, res)
	}
	if err == nrpc.ErrEOS {
		err = nil
	}
	return err
}

// SvcSubjectParamsServer is the interface that providers of the service
// SvcSubjectParams should implement.
type SvcSubjectParamsServer interface {
	MtWithSubjectParams(ctx context.Context, mp1 string, mp2 string) (resp SimpleStringReply, err error)
	MtNoReply(ctx context.Context)
}

// SvcSubjectParamsHandler provides a NATS subscription handler that can serve a
// subscription using a given SvcSubjectParamsServer implementation.
type SvcSubjectParamsHandler struct {
	ctx    context.Context
	nc     nrpc.NatsConn
	server SvcSubjectParamsServer
}

func NewSvcSubjectParamsHandler(ctx context.Context, nc nrpc.NatsConn, s SvcSubjectParamsServer) *SvcSubjectParamsHandler {
	return &SvcSubjectParamsHandler{
		ctx:    ctx,
		nc:     nc,
		server: s,
	}
}

func (h *SvcSubjectParamsHandler) Subject() string {
	return "root.*.svcsubjectparams.*.>"
}

func (h *SvcSubjectParamsHandler) MtNoRequestWParamsPublish(pkginstance string, svcclientid string, mtmp1 string, msg SimpleStringReply) error {
	rawMsg, err := nrpc.Marshal("protobuf", &msg)
	if err != nil {
		log.Printf("SvcSubjectParamsHandler.MtNoRequestWParamsPublish: error marshaling the message: %s", err)
		return err
	}
	subject := "root." + pkginstance + "."+ "svcsubjectparams." + svcclientid + "."+ "mtnorequestwparams" + "." + mtmp1
	return h.nc.Publish(subject, rawMsg)
}

func (h *SvcSubjectParamsHandler) Handler(msg *nats.Msg) {
	var encoding string
	var noreply bool
	// extract method name & encoding from subject
	pkgParams, svcParams, name, tail, err := nrpc.ParseSubject(
		"root", 1, "svcsubjectparams", 1, msg.Subject)
	if err != nil {
		log.Printf("SvcSubjectParamsHanlder: SvcSubjectParams subject parsing failed: %v", err)
		return
	}

	ctx := h.ctx
	ctx = context.WithValue(ctx, "nrpc-pkg-instance", pkgParams[0])
	ctx = context.WithValue(ctx, "nrpc-svc-clientid", svcParams[0])
	// call handler and form response
	var resp proto.Message
	var replyError *nrpc.Error
	switch name {
	case "mtwithsubjectparams":
		var mtParams []string
		mtParams, encoding, err = nrpc.ParseSubjectTail(2, tail)
		if err != nil {
			log.Printf("MtWithSubjectParamsHanlder: MtWithSubjectParams subject parsing failed: %v", err)
			break
		}
		var req github_com_rapidloop_nrpc.Void
		if err := nrpc.Unmarshal(encoding, msg.Data, &req); err != nil {
			log.Printf("MtWithSubjectParamsHandler: MtWithSubjectParams request unmarshal failed: %v", err)
			replyError = &nrpc.Error{
				Type: nrpc.Error_CLIENT,
				Message: "bad request received: " + err.Error(),
			}
		} else {
			resp, replyError = nrpc.CaptureErrors(
				func()(proto.Message, error){
					innerResp, err := h.server.MtWithSubjectParams(ctx, mtParams[0], mtParams[1])
					if err != nil {
						return nil, err
					}
					return &innerResp, err
				})
			if replyError != nil {
				log.Printf("MtWithSubjectParamsHandler: MtWithSubjectParams handler failed: %s", replyError.Error())
			}
		}
	case "mtnoreply":
		noreply = true
		_, encoding, err = nrpc.ParseSubjectTail(0, tail)
		if err != nil {
			log.Printf("MtNoReplyHanlder: MtNoReply subject parsing failed: %v", err)
			break
		}
		var req github_com_rapidloop_nrpc.Void
		if err := nrpc.Unmarshal(encoding, msg.Data, &req); err != nil {
			log.Printf("MtNoReplyHandler: MtNoReply request unmarshal failed: %v", err)
			replyError = &nrpc.Error{
				Type: nrpc.Error_CLIENT,
				Message: "bad request received: " + err.Error(),
			}
		} else {
			resp, replyError = nrpc.CaptureErrors(
				func()(proto.Message, error){var innerResp nrpc.NoReply
					h.server.MtNoReply(ctx)
					if err != nil {
						return nil, err
					}
					return &innerResp, err
				})
			if replyError != nil {
				log.Printf("MtNoReplyHandler: MtNoReply handler failed: %s", replyError.Error())
			}
		}
	case "mtnorequestwparams":
		// MtNoRequestWParams is a no-request method. Ignore it.
		return
	default:
		log.Printf("SvcSubjectParamsHandler: unknown name %q", name)
		replyError = &nrpc.Error{
			Type: nrpc.Error_CLIENT,
			Message: "unknown name: " + name,
		}
	}


	if !noreply {
		// encode and send response
		err = nrpc.Publish(resp, replyError, h.nc, msg.Reply, encoding) // error is logged
	} else {
		err = nil
	}
	if err != nil {
		log.Println("SvcSubjectParamsHandler: SvcSubjectParams handler failed to publish the response: %s", err)
	}
}

type SvcSubjectParamsClient struct {
	nc      nrpc.NatsConn
	PkgSubject string
	PkgParaminstance string
	Subject string
	SvcParamclientid string
	Encoding string
	Timeout time.Duration
}

func NewSvcSubjectParamsClient(nc nrpc.NatsConn, pkgParaminstance string, svcParamclientid string) *SvcSubjectParamsClient {
	return &SvcSubjectParamsClient{
		nc:      nc,
		PkgSubject: "root",
		PkgParaminstance: pkgParaminstance,
		Subject: "svcsubjectparams",
		SvcParamclientid: svcParamclientid,
		Encoding: "protobuf",
		Timeout: 5 * time.Second,
	}
}

func (c *SvcSubjectParamsClient) MtWithSubjectParams(mp1 string, mp2 string, ) (resp SimpleStringReply, err error) {

	subject := c.PkgSubject + "." + c.PkgParaminstance + "." + c.Subject + "." + c.SvcParamclientid + "." + "mtwithsubjectparams" + "." + mp1 + "." + mp2

	// call
	var req github_com_rapidloop_nrpc.Void
	err = nrpc.Call(&req, &resp, c.nc, subject, c.Encoding, c.Timeout)
	if err != nil {
		return // already logged
	}

	return
}

func (c *SvcSubjectParamsClient) MtNoReply() (err error) {

	subject := c.PkgSubject + "." + c.PkgParaminstance + "." + c.Subject + "." + c.SvcParamclientid + "." + "mtnoreply"

	// call
	var req github_com_rapidloop_nrpc.Void
	var resp github_com_rapidloop_nrpc.NoReply
	err = nrpc.Call(&req, &resp, c.nc, subject, c.Encoding, c.Timeout)
	if err != nil {
		return // already logged
	}

	return
}

func (c *SvcSubjectParamsClient) MtNoRequestWParamsSubject(
	mtmp1 string,
) string {
	return c.PkgSubject + "." + c.PkgParaminstance + "." + c.Subject + "." + c.SvcParamclientid + "." + "mtnorequestwparams" + "." + mtmp1
}

type SvcSubjectParamsMtNoRequestWParamsSubscription struct {
	*nats.Subscription
}

func (s *SvcSubjectParamsMtNoRequestWParamsSubscription) Next(timeout time.Duration) (next SimpleStringReply, err error) {
	msg, err := s.Subscription.NextMsg(timeout)
	if err != nil {
		return
	}
	err = nrpc.Unmarshal("protobuf", msg.Data, &next)
	return
}

func (c *SvcSubjectParamsClient) MtNoRequestWParamsSubscribeSync(
	mtmp1 string,
) (sub *SvcSubjectParamsMtNoRequestWParamsSubscription, err error) {
	subject := c.MtNoRequestWParamsSubject(
		mtmp1,
	)
	natsSub, err := c.nc.SubscribeSync(subject)
	if err != nil {
		return
	}
	sub = &SvcSubjectParamsMtNoRequestWParamsSubscription{natsSub}
	return
}

func (c *SvcSubjectParamsClient) MtNoRequestWParamsSubscribe(
	mtmp1 string,
	handler func (SimpleStringReply),
) (sub *nats.Subscription, err error) {
	subject := c.MtNoRequestWParamsSubject(
		mtmp1,
	)
	sub, err = c.nc.Subscribe(subject, func(msg *nats.Msg){
		var pmsg SimpleStringReply
		err := nrpc.Unmarshal("protobuf", msg.Data, &pmsg)
		if err != nil {
			log.Printf("SvcSubjectParamsClient.MtNoRequestWParamsSubscribe: Error decoding, %s", err)
			return
		}
		handler(pmsg)
	})
	return
}

func (c *SvcSubjectParamsClient) MtNoRequestWParamsSubscribeChan(
	mtmp1 string,
) (<-chan SimpleStringReply, *nats.Subscription, error) {
	ch := make(chan SimpleStringReply)
	sub, err := c.MtNoRequestWParamsSubscribe(mtmp1, func (msg SimpleStringReply) {
		ch <- msg
	})
	return ch, sub, err
}

// NoRequestServiceServer is the interface that providers of the service
// NoRequestService should implement.
type NoRequestServiceServer interface {
}

// NoRequestServiceHandler provides a NATS subscription handler that can serve a
// subscription using a given NoRequestServiceServer implementation.
type NoRequestServiceHandler struct {
	ctx    context.Context
	nc     nrpc.NatsConn
	server NoRequestServiceServer
}

func NewNoRequestServiceHandler(ctx context.Context, nc nrpc.NatsConn, s NoRequestServiceServer) *NoRequestServiceHandler {
	return &NoRequestServiceHandler{
		ctx:    ctx,
		nc:     nc,
		server: s,
	}
}

func (h *NoRequestServiceHandler) Subject() string {
	return "root.*.norequestservice.>"
}

func (h *NoRequestServiceHandler) MtNoRequestPublish(pkginstance string, msg SimpleStringReply) error {
	rawMsg, err := nrpc.Marshal("protobuf", &msg)
	if err != nil {
		log.Printf("NoRequestServiceHandler.MtNoRequestPublish: error marshaling the message: %s", err)
		return err
	}
	subject := "root." + pkginstance + "."+ "norequestservice."+ "mtnorequest"
	return h.nc.Publish(subject, rawMsg)
}

type NoRequestServiceClient struct {
	nc      nrpc.NatsConn
	PkgSubject string
	PkgParaminstance string
	Subject string
	Encoding string
	Timeout time.Duration
}

func NewNoRequestServiceClient(nc nrpc.NatsConn, pkgParaminstance string) *NoRequestServiceClient {
	return &NoRequestServiceClient{
		nc:      nc,
		PkgSubject: "root",
		PkgParaminstance: pkgParaminstance,
		Subject: "norequestservice",
		Encoding: "protobuf",
		Timeout: 5 * time.Second,
	}
}

func (c *NoRequestServiceClient) MtNoRequestSubject(
	
) string {
	return c.PkgSubject + "." + c.PkgParaminstance + "." + c.Subject + "." + "mtnorequest"
}

type NoRequestServiceMtNoRequestSubscription struct {
	*nats.Subscription
}

func (s *NoRequestServiceMtNoRequestSubscription) Next(timeout time.Duration) (next SimpleStringReply, err error) {
	msg, err := s.Subscription.NextMsg(timeout)
	if err != nil {
		return
	}
	err = nrpc.Unmarshal("protobuf", msg.Data, &next)
	return
}

func (c *NoRequestServiceClient) MtNoRequestSubscribeSync(
	
) (sub *NoRequestServiceMtNoRequestSubscription, err error) {
	subject := c.MtNoRequestSubject(
		
	)
	natsSub, err := c.nc.SubscribeSync(subject)
	if err != nil {
		return
	}
	sub = &NoRequestServiceMtNoRequestSubscription{natsSub}
	return
}

func (c *NoRequestServiceClient) MtNoRequestSubscribe(
	
	handler func (SimpleStringReply),
) (sub *nats.Subscription, err error) {
	subject := c.MtNoRequestSubject(
		
	)
	sub, err = c.nc.Subscribe(subject, func(msg *nats.Msg){
		var pmsg SimpleStringReply
		err := nrpc.Unmarshal("protobuf", msg.Data, &pmsg)
		if err != nil {
			log.Printf("NoRequestServiceClient.MtNoRequestSubscribe: Error decoding, %s", err)
			return
		}
		handler(pmsg)
	})
	return
}

func (c *NoRequestServiceClient) MtNoRequestSubscribeChan(
	
) (<-chan SimpleStringReply, *nats.Subscription, error) {
	ch := make(chan SimpleStringReply)
	sub, err := c.MtNoRequestSubscribe(func (msg SimpleStringReply) {
		ch <- msg
	})
	return ch, sub, err
}

type Client struct {
	nc      nrpc.NatsConn
	defaultEncoding string
	defaultTimeout time.Duration
	pkgSubject string
	pkgParaminstance string
	SvcCustomSubject *SvcCustomSubjectClient
	SvcSubjectParams *SvcSubjectParamsClient
	NoRequestService *NoRequestServiceClient
}

func NewClient(nc nrpc.NatsConn, pkgParaminstance string) *Client {
	c := Client{
		nc: nc,
		defaultEncoding: "protobuf",
		defaultTimeout: 5*time.Second,
		pkgSubject: "root",
		pkgParaminstance: pkgParaminstance,
	}
	c.SvcCustomSubject = NewSvcCustomSubjectClient(nc, c.pkgParaminstance)
	c.NoRequestService = NewNoRequestServiceClient(nc, c.pkgParaminstance)
	return &c
}

func (c *Client) SetEncoding(encoding string) {
	c.defaultEncoding = encoding
	if c.SvcCustomSubject != nil {
		c.SvcCustomSubject.Encoding = encoding
	}
	if c.SvcSubjectParams != nil {
		c.SvcSubjectParams.Encoding = encoding
	}
	if c.NoRequestService != nil {
		c.NoRequestService.Encoding = encoding
	}
}

func (c *Client) SetTimeout(t time.Duration) {
	c.defaultTimeout = t
	if c.SvcCustomSubject != nil {
		c.SvcCustomSubject.Timeout = t
	}
	if c.SvcSubjectParams != nil {
		c.SvcSubjectParams.Timeout = t
	}
	if c.NoRequestService != nil {
		c.NoRequestService.Timeout = t
	}
}

func (c *Client) SetSvcSubjectParamsParams(
	clientid string,
) {
	c.SvcSubjectParams = NewSvcSubjectParamsClient(
		c.nc,
		c.pkgParaminstance,
		clientid,
	)
	c.SvcSubjectParams.Encoding = c.defaultEncoding
	c.SvcSubjectParams.Timeout = c.defaultTimeout
}

func (c *Client) NewSvcSubjectParams(
	clientid string,
) *SvcSubjectParamsClient {
	client := NewSvcSubjectParamsClient(
		c.nc,
		c.pkgParaminstance,
		clientid,
	)
	client.Encoding = c.defaultEncoding
	client.Timeout = c.defaultTimeout
	return client
}