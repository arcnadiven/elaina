package server

import (
	"github.com/arcnadiven/elaina/backstore"
	"github.com/arcnadiven/elaina/tracelog"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type ElainaServer interface {
	csi.IdentityServer
	csi.ControllerServer
	csi.NodeServer
	RunGRPCServer() error
}

type ElainaServerImpl struct {
	bl    tracelog.BaseLogger
	opera backstore.StoreOperator
}

func NewElainaServer(bl tracelog.BaseLogger, opera backstore.StoreOperator) ElainaServer {
	return &ElainaServerImpl{
		bl:    bl,
		opera: opera,
	}
}

func (es *ElainaServerImpl) RunGRPCServer() error {
	ipAddr, err := net.ResolveUnixAddr("unix", "/tmp/csi.sock")
	if err != nil {
		es.bl.Errorln(err)
		return err
	}
	logrus.Infoln(ipAddr.String(), ipAddr.Network())
	lis, err := net.Listen("unix", ipAddr.String())
	if err != nil {
		es.bl.Errorln(err)
		return err
	}
	srv := grpc.NewServer()
	csi.RegisterIdentityServer(srv, es)
	csi.RegisterControllerServer(srv, es)
	csi.RegisterNodeServer(srv, es)
	reflection.Register(srv)
	return srv.Serve(lis)
}
