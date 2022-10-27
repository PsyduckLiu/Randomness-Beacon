package p2pnetwork

import (
	"consensusNode/config"
	"consensusNode/message"
	"consensusNode/signature"
	"consensusNode/util"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type P2pNetwork interface {
	GetPeerPublickey(peerId int64) *ecdsa.PublicKey
	GetMySecretkey() *ecdsa.PrivateKey
	BroadCast(v interface{}) error
}

// [SrvHub]: contains all TCP connections with other nodes
// [Peers]: map TCP connect to an int number
// [MsgChan]: a channel connects [p2p] with [state(consensus)], deliver consensus message, corresponding to [ch] in [state(consensus)]
type SimpleP2p struct {
	NodeId         int64
	SrvHub         *net.TCPListener
	Peers          map[string]*net.TCPConn
	Ip2Id          map[string]int64
	PrivateKey     *ecdsa.PrivateKey
	PeerPublicKeys map[int64]*ecdsa.PublicKey
	MsgChan        chan<- *message.ConMessage
	mutex          sync.Mutex
}

// new simple P2P liarary
func NewSimpleP2pLib(id int64, msgChan chan<- *message.ConMessage) P2pNetwork {
	// get specified curve
	marshalledCurve := config.GetCurve()
	pub, err := x509.ParsePKIXPublicKey(marshalledCurve)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR]Key message parse err:%s", err))
	}
	normalPublicKey := pub.(*ecdsa.PublicKey)
	curve := normalPublicKey.Curve
	fmt.Printf("Curve is %v\n", curve.Params())

	// generate private key
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR]Generate private key err:%s", err))
	}
	fmt.Printf("===>My own key is: %v\n", privateKey)

	// listen port 30000+id
	port := util.PortByID(id)
	s, err := net.ListenTCP("tcp4", &net.TCPAddr{
		Port: port,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("===>[Node%d] is waiting at:%s\n", id, s.Addr().String())

	// write new node details into config
	config.NewConsensusNode(id, s.Addr().String(), elliptic.Marshal(curve, privateKey.PublicKey.X, privateKey.PublicKey.Y))

	sp := &SimpleP2p{
		NodeId:         id,
		SrvHub:         s,
		Peers:          make(map[string]*net.TCPConn),
		Ip2Id:          make(map[string]int64),
		PrivateKey:     privateKey,
		PeerPublicKeys: make(map[int64]*ecdsa.PublicKey),
		MsgChan:        msgChan,
		mutex:          sync.Mutex{},
	}

	go sp.monitor(id)
	sp.dialTcp(id)

	return sp
}

// connect to all known nodes
func (sp *SimpleP2p) dialTcp(id int64) {
	nodeConfig := config.GetConsensusNode()

	for i := 0; i < len(nodeConfig); i++ {
		if int64(i) != id && nodeConfig[i].Ip != "0" {
			// resolve TCP address and dial TCP conn
			addr, err := net.ResolveTCPAddr("tcp4", nodeConfig[i].Ip)
			if err != nil {
				panic(err)
			}
			conn, err := net.DialTCP("tcp", nil, addr)
			if err != nil {
				panic(err)
			}

			// get specified curve
			marshalledCurve := config.GetCurve()
			pub, err := x509.ParsePKIXPublicKey(marshalledCurve)
			if err != nil {
				panic(fmt.Errorf("===>[ERROR]Key message parse err:%s", err))
			}
			normalPublicKey := pub.(*ecdsa.PublicKey)
			curve := normalPublicKey.Curve

			// unmarshal public key
			x, y := elliptic.Unmarshal(curve, nodeConfig[i].Pk)
			newPublicKey := &ecdsa.PublicKey{
				Curve: curve,
				X:     x,
				Y:     y,
			}

			// store remote ip-conn-id-pk relation
			sp.mutex.Lock()
			sp.Peers[conn.RemoteAddr().String()] = conn
			sp.Ip2Id[conn.RemoteAddr().String()] = int64(i)
			sp.PeerPublicKeys[int64(i)] = newPublicKey
			sp.mutex.Unlock()
			fmt.Printf("===>[Node%d<=>%d]Connected=[%s<=>%s]\n", id, i, conn.LocalAddr().String(), conn.RemoteAddr().String())

			// new identity message
			// send identity message to origin nodes
			kMsg := message.CreateIdentityMsg(message.MTIdentity, sp.NodeId, conn.LocalAddr().String(), sp.PrivateKey)
			if err := sp.SendUniqueNode(conn, kMsg); err != nil {
				panic(err)
			}

			go sp.waitData(conn)
		}
	}
}

// add new node OR remove old node
func (sp *SimpleP2p) monitor(id int64) {
	fmt.Printf("===>Consensus node[%d] is waiting at:%s\n", id, sp.SrvHub.Addr().String())

	for {
		conn, err := sp.SrvHub.AcceptTCP()

		// remove a node
		if err != nil {
			fmt.Printf("===>P2p network accept err:%s\n", err)
			if err == io.EOF {
				fmt.Printf("===>[Node%d] Remove peer node[%d]%s\n", sp.NodeId, sp.Ip2Id[conn.RemoteAddr().String()], conn.RemoteAddr().String())
				config.RemoveConsensusNode(sp.Ip2Id[conn.RemoteAddr().String()])
				delete(sp.Peers, conn.RemoteAddr().String())
				delete(sp.PeerPublicKeys, sp.Ip2Id[conn.RemoteAddr().String()])
				delete(sp.Ip2Id, conn.RemoteAddr().String())
			}
			continue
		}

		// add a new node
		sp.Peers[conn.RemoteAddr().String()] = conn
		go sp.waitData(conn)
	}
}

// remove old node AND deliver consensus mseeage by [MsgChan]
func (sp *SimpleP2p) waitData(conn *net.TCPConn) {
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)

		// remove a node
		if err != nil {
			fmt.Printf("===>P2p network capture data err:%s\n", err)
			if err == io.EOF {
				fmt.Printf("===>[Node%d] Remove peer node[%d]%s\n", sp.NodeId, sp.Ip2Id[conn.RemoteAddr().String()], conn.RemoteAddr().String())
				config.RemoveConsensusNode(sp.Ip2Id[conn.RemoteAddr().String()])
				delete(sp.Peers, conn.RemoteAddr().String())
				delete(sp.PeerPublicKeys, sp.Ip2Id[conn.RemoteAddr().String()])
				delete(sp.Ip2Id, conn.RemoteAddr().String())
				return
			}
			continue
		}

		// handle a consensus message
		conMsg := &message.ConMessage{}
		if err := json.Unmarshal(buf[:n], conMsg); err != nil {
			fmt.Println(string(buf[:n]))
			panic(err)
		}
		time.Sleep(100 * time.Millisecond)

		switch conMsg.Typ {
		// handle new identity message from backups
		case message.MTIdentity:
			nodeConfig := config.GetConsensusNode()

			// get specified curve
			marshalledCurve := config.GetCurve()
			pub, err := x509.ParsePKIXPublicKey(marshalledCurve)
			if err != nil {
				panic(fmt.Errorf("===>[ERROR]Key message parse err:%s", err))
			}
			normalPublicKey := pub.(*ecdsa.PublicKey)
			curve := normalPublicKey.Curve

			// unmarshal public key
			x, y := elliptic.Unmarshal(curve, nodeConfig[conMsg.From].Pk)
			newPublicKey := &ecdsa.PublicKey{
				Curve: curve,
				X:     x,
				Y:     y,
			}

			// verify signature
			verify := signature.VerifySig(conMsg.Payload, conMsg.Sig, newPublicKey)
			if !verify {
				fmt.Printf("===>[ERROR]Verify new public key Signature failed, From Node[%d], IP[%s]\n", conMsg.From, conn.RemoteAddr().String())
				break
			}

			// add a new node
			if sp.PeerPublicKeys[conMsg.From] != newPublicKey {
				sp.mutex.Lock()
				sp.Ip2Id[conn.RemoteAddr().String()] = conMsg.From
				sp.PeerPublicKeys[conMsg.From] = newPublicKey
				sp.mutex.Unlock()

				fmt.Printf("===>Get new public key from Node[%d], IP[%s]\n", conMsg.From, conn.RemoteAddr().String())
				fmt.Printf("===>Node[%d]'s new public key is[%v]\n", conMsg.From, newPublicKey)
				fmt.Printf("===>[Node%d<=>%d]Connected=[%s<=>%s]\n", sp.NodeId, conMsg.From, conn.LocalAddr().String(), conn.RemoteAddr().String())
			}

		// handle consensus message from backups
		default:
			sp.MsgChan <- conMsg
		}
	}
}

// BroadCast message to all connected nodes
func (sp *SimpleP2p) BroadCast(v interface{}) error {
	if v == nil {
		return fmt.Errorf("===>[ERROR]empty msg body")
	}
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	for name, conn := range sp.Peers {
		go WriteTCP(conn, data, name)
	}

	return nil
}

// BroadCast message to all connected nodes
func (sp *SimpleP2p) SendUniqueNode(conn *net.TCPConn, v interface{}) error {
	if v == nil {
		return fmt.Errorf("===>[ERROR]empty msg body")
	}

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	go WriteTCP(conn, data, conn.RemoteAddr().String())

	return nil
}

func WriteTCP(conn *net.TCPConn, v []byte, name string) {
	_, err := conn.Write(v)
	if err != nil {
		fmt.Printf("===>[ERROR]write to node[%s] err:%s\n", name, err)
		panic(err)
	}
}

// Get Peer Publickey
func (sp *SimpleP2p) GetPeerPublickey(peerId int64) *ecdsa.PublicKey {
	return sp.PeerPublicKeys[peerId]
}

// Get My Secret key
func (sp *SimpleP2p) GetMySecretkey() *ecdsa.PrivateKey {
	return sp.PrivateKey
}
