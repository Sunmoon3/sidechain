package cross

import (
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"net"
	"viewholder/config"
)

func (cl *CrossListener) receiveFromMC(errch chan<- error) {
	SERVER_PEM := config.SERVER_DEFAULT_PEM
	SERVER_KEY := config.SERVER_DEFAULT_KEY
	CA_PEM := config.SERVER_DEFAULT_CA_PEM
	IP := config.OPPOSITE_IP
	PORT := config.NODE_DEFAULT_PORT

	log.Infof("getting Sever pem path: %v", SERVER_PEM)
	log.Infof("getting Sever key path: %v", SERVER_KEY)
	log.Infof("getting Sever ca path: %v", CA_PEM)

	cert, err := tls.LoadX509KeyPair(SERVER_PEM, SERVER_KEY)
	if err != nil {
		log.Fatalf("tls.LoadX509KeyPair err: %v", err)
		errch <- err
		return
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(CA_PEM)
	if err != nil {
		log.Fatalf("ioutil.ReadFile err: %v", err)
		errch <- err
		return
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("certPool.AppendCertsFromPEM err")
		errch <- err
		return
	}

	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})

	server := grpc.NewServer(grpc.Creds(c))
	//pb.RegisterSearchServiceServer(server, &SearchService{})

	lis, err := net.Listen("tcp", IP+":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
		errch <- err
		return
	}
	log.Infof("cross.go cd.lis = &lis 和 cd.server = server")
	cl.lis = &lis
	cl.server = server
}
