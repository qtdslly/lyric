package utils

import "go-lyric/common/logger"

// "net"
// "time"
// "fmt"
// "log"
// "strings"

// "github.com/google/gopacket"
// "github.com/google/gopacket/layers"
// "github.com/google/gopacket/pcap"

func TcpKill(host string) bool {
	// data,_ := BashCommand("timeout 3 /usr/sbin/tcpkill -9 port " + host,5)
	_, err := BashCommand("ss -K dport = "+host, 5)
	if err != nil {
		logger.Info("ss -K dport = " + host)
		logger.Error(err)
		return false
	}
	// if strings.Contains(data,"win 0"){
	// 	return true
	// }
	return true
}

// func TcpKill(iface, port string) bool {
// 	filter := "port " + port
// 	handle, err := pcap.OpenLive(iface, 64, true, time.Millisecond)
// 	if err != nil {
// 		return false
// 	}
// 	defer handle.Close()

// 	if err = handle.SetBPFFilter(filter); err != nil {
// 		return false
// 	}

// 	options := gopacket.SerializeOptions{
// 		ComputeChecksums: true,
// 		FixLengths:       true,
// 	}

// 	src := gopacket.NewPacketSource(handle, handle.LinkType())

// 	//count := 0
// 	for packet := range src.Packets() {
// 		//count++
// 		//if count > 3 {
// 		//	return false
// 		//}
// 		ethLayer := packet.Layer(layers.LayerTypeEthernet)
// 		if ethLayer == nil {
// 			continue
// 		}
// 		eth := ethLayer.(*layers.Ethernet)
// 		ipLayer := packet.Layer(layers.LayerTypeIPv4)
// 		if ipLayer == nil {
// 			continue
// 		}
// 		ip := ipLayer.(*layers.IPv4)
// 		tcpLayer := packet.Layer(layers.LayerTypeTCP)
// 		if tcpLayer == nil {
// 			continue
// 		}
// 		tcp := tcpLayer.(*layers.TCP)

// 		if tcp.FIN || tcp.RST {
// 			break
// 			//continue
// 		}

// 		// fmt.Printf("%s:%d > %s:%d\n", ip.SrcIP.String(), tcp.SrcPort, ip.DstIP.String(), tcp.DstPort)

// 		neth := &layers.Ethernet{
// 			SrcMAC:       eth.DstMAC,
// 			DstMAC:       eth.SrcMAC,
// 			EthernetType: layers.EthernetTypeIPv4,
// 		}
// 		nip := &layers.IPv4{
// 			SrcIP:    ip.DstIP,
// 			DstIP:    ip.SrcIP,
// 			Version:  4,
// 			TTL:      64,
// 			Protocol: layers.IPProtocolTCP,
// 		}
// 		ntcp := &layers.TCP{
// 			SrcPort: tcp.DstPort,
// 			DstPort: tcp.SrcPort,
// 			RST:     true,
// 			Seq:     tcp.Ack,
// 		}
// 		ntcp.SetNetworkLayerForChecksum(nip)

// 		buffer := gopacket.NewSerializeBuffer()
// 		if err := gopacket.SerializeLayers(buffer, options, neth, nip, ntcp); err != nil {
// 			return false
// 		}

// 		if err := handle.WritePacketData(buffer.Bytes()); err != nil {
// 			return false
// 		}

// 		break
// 	}
// 	return true
// }

// func TcpKill(iface ,port string) bool {

// 	fmt.Println("11111111")
// 	var handle *pcap.Handle
// 	var err error
// 	if handle, err = pcap.OpenLive(iface, int32(65535), true, -1*time.Second); err != nil {
// 			return false
// 	}
// 	defer handle.Close()

// 	fmt.Println("22222222")
// 	//var filters []string

// 	filter := "port " + port

// 	if err := handle.SetBPFFilter(filter); err != nil {
// 		return false
// 	}
// 	fmt.Println("333333333")

// 	packetSource := gopacket.NewPacketSource(
// 			handle,
// 			handle.LinkType(),
// 	)

// 	fmt.Println("44444444")
// 	if err := capture(packetSource, handle); err != nil {
// 		return false
// 	}

// 	return true
// }

// func capture(packetSource *gopacket.PacketSource, handle *pcap.Handle) error {
// 	for packet := range packetSource.Packets() {
// 			fmt.Println("aaaaaaaaaaaaa")
// 			ethLayer := packet.Layer(layers.LayerTypeEthernet)
// 			if ethLayer == nil {
// 					fmt.Println("bbbbbbbbb")
// 					continue
// 			}
// 			eth := ethLayer.(*layers.Ethernet)
// 			ipv4Layer := packet.Layer(layers.LayerTypeIPv4)
// 			if ipv4Layer == nil {
// 					fmt.Println("cccccccc")
// 					continue
// 			}
// 			ip := ipv4Layer.(*layers.IPv4)
// 			tcpLayer := packet.Layer(layers.LayerTypeTCP)
// 			if tcpLayer == nil {
// 					fmt.Println("ddddddd")
// 					continue
// 			}
// 			tcp := tcpLayer.(*layers.TCP)

// 			if tcp.SYN || tcp.FIN || tcp.RST {
// 					fmt.Println("eeeeeeee")
// 					break
// 			}

// 			for i := 0; i < 3 ; i++ {
// 					seq := tcp.Ack + uint32(i)*uint32(tcp.Window)

// 					err := SendRST(eth.DstMAC, eth.SrcMAC, ip.DstIP, ip.SrcIP, tcp.DstPort, tcp.SrcPort, seq, handle)
// 					if err != nil {
// 							return err
// 					}
// 			}
// 			//      fmt.Println("fffffffffff")
// 	}
// 	return nil
// }

// /*
// func SendSYN(srcIp, dstIp net.IP, srcPort, dstPort layers.TCPPort, seq uint32, handle *pcap.Handle) error {
// 	log.Printf("send %v:%v > %v:%v [SYN] seq %v", srcIp.String(), srcPort.String(), dstIp.String(), dstPort.String(), seq)
// 	iPv4 := layers.IPv4{
// 			SrcIP:    srcIp,
// 			DstIP:    dstIp,
// 			Version:  4,
// 			TTL:      64,
// 			Protocol: layers.IPProtocolTCP,
// 	}

// 	tcp := layers.TCP{
// 			SrcPort: srcPort,
// 			DstPort: dstPort,
// 			Seq:     seq,
// 			SYN:     true,
// 	}

// 	if err := tcp.SetNetworkLayerForChecksum(&iPv4); err != nil {
// 			return err
// 	}

// 	buffer := gopacket.NewSerializeBuffer()
// 	options := gopacket.SerializeOptions{
// 			FixLengths:       true,
// 			ComputeChecksums: true,
// 	}
// 	if err := gopacket.SerializeLayers(buffer, options, &tcp); err != nil {
// 			return err
// 	}

// 	err := handle.WritePacketData(buffer.Bytes())
// 	if err != nil {
// 			return err
// 	}
// 	return nil
// }
// */

// func SendRST(srcMac, dstMac net.HardwareAddr, srcIp, dstIp net.IP, srcPort, dstPort layers.TCPPort, seq uint32, handle *pcap.Handle) error {
// 	log.Printf("send %v:%v > %v:%v [RST] seq %v", srcIp.String(), srcPort.String(), dstIp.String(), dstPort.String(), seq)

// 	eth := layers.Ethernet{
// 			SrcMAC:       srcMac,
// 			DstMAC:       dstMac,
// 			EthernetType: layers.EthernetTypeIPv4,
// 	}

// 	iPv4 := layers.IPv4{
// 			SrcIP:    srcIp,
// 			DstIP:    dstIp,
// 			Version:  4,
// 			TTL:      64,
// 			Protocol: layers.IPProtocolTCP,
// 	}

// 	tcp := layers.TCP{
// 			SrcPort: srcPort,
// 			DstPort: dstPort,
// 			Seq:     seq,
// 			RST:     true,
// 	}

// 	if err := tcp.SetNetworkLayerForChecksum(&iPv4); err != nil {
// 			return err
// 	}

// 	buffer := gopacket.NewSerializeBuffer()
// 	options := gopacket.SerializeOptions{
// 			FixLengths:       true,
// 			ComputeChecksums: true,
// 	}
// 	if err := gopacket.SerializeLayers(buffer, options, &eth, &iPv4, &tcp); err != nil {
// 			return err
// 	}

// 	err := handle.WritePacketData(buffer.Bytes())
// 	if err != nil {
// 			return err
// 	}
// 	return nil
// }
