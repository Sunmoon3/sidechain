package cross

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"viewholder/config"
	pb "viewholder/sidechainhandle/listener/cross/proto"
)

func (cl *CrossListener) SendtoSC(request *pb.CrossLanguageRequest) {
	log.Info("CrosslanguageGeneratedEvent...")

	sideChainAddre := config.SIDE_CHAIN_ADDRE
	CLIENT_PEM := config.CLIENT_DEFAULT_PEM
	CLIENT_KEY := config.CLIENT_DEFAULT_KEY
	CA_PEM := config.CLIENT_DEFAULT_CA_PEM

	log.Infof("getting Client pem path: %v", CLIENT_PEM)
	log.Infof("getting Client key path: %v", CLIENT_KEY)
	log.Infof("getting Client ca path: %v", CA_PEM)

	cert, err := tls.LoadX509KeyPair(CLIENT_PEM, CLIENT_KEY)
	if err != nil {
		log.Errorf("tls.LoadX509KeyPair err: %v", err)
		cl.errch <- fmt.Errorf("tls.LoadX509KeyPair err: %v", err)
		return
	}
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(CA_PEM)
	if err != nil {
		log.Errorf("ioutil.ReadFile err: %v", err)
		cl.errch <- fmt.Errorf("ioutil.ReadFile err: %v", err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Errorf("certPool.AppendCertsFromPEM err")
		cl.errch <- fmt.Errorf("certPool.AppendCertsFromPEM err")
	}

	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "localhost",
		RootCAs:      certPool,
	})

	conn, err := grpc.Dial(sideChainAddre, grpc.WithTransportCredentials(c))
	if err != nil {
		log.Errorf("grpc.Dial err: %v", err)
		cl.errch <- fmt.Errorf("grpc.Dial err: %v", err)
	}
	defer conn.Close()
	client := pb.NewCrossServiceClient(conn)
	stream, err := client.CrossIn(context.Background())
	if err != nil {
		log.Errorf("grpc.CrossIn Stream error: %v", err)
		cl.errch <- fmt.Errorf("grpc.CrossIn Stream error: %v", err)
		return
	}
	err = stream.Send(request)
	if err != nil {
		log.Errorf("grpc.CrossIn Send error: %v", err)
		cl.errch <- fmt.Errorf("grpc.CrossIn Send error: %v", err)
		return
	}
	log.Infof("CrossChain Send Success,Waiting for Respose")

	r, err := stream.Recv()
	if err != nil {
		log.Errorf("grpc.CrossIn Recv error: %v", err)
		cl.errch <- fmt.Errorf("grpc.CrossIn Recv error: %v", err)
		return
	}

	log.Infof("CrossChain Send Success,Respose Status: %v, msg: %v", r.Status, r.Msg)
}
