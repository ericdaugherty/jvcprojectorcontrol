package jvcprojectorcontrol

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
	"sync"
	"time"
)

const jvcPort = 20554

type Command int
type HashMode int

const (
	NullCommand Command = iota
	OffCommand
	OnCommand
	Input1Command
	Input2Command
)

var commands = map[Command][]byte{
	NullCommand:   {0x21, 0x89, 0x01, 0x00, 0x00, 0x0A},
	OffCommand:    {0x21, 0x89, 0x01, 0x50, 0x57, 0x30, 0x0A},
	OnCommand:     {0x21, 0x89, 0x01, 0x50, 0x57, 0x31, 0x0A},
	Input1Command: {0x21, 0x89, 0x01, 0x49, 0x50, 0x36, 0x0A},
	Input2Command: {0x21, 0x89, 0x01, 0x49, 0x50, 0x37, 0x0A},
}

const (
	HashNone HashMode = iota
	HashJVCKW
	HashJVCKWPJ
)

type Projector struct {
	IPAddress string
	Password  string
	Hash      HashMode
	Debug     bool
}

func NewProjector(ipAddress, password string, hash HashMode, debug bool) *Projector {
	return &Projector{
		IPAddress: ipAddress,
		Password:  password,
		Hash:      hash,
		Debug:     debug,
	}
}

func (p *Projector) SendCommand(c Command) error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", p.IPAddress, jvcPort), 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to projector at %s: %v", p.IPAddress, err)
	}
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	pc := projConnection{conn: conn, debug: p.Debug}
	defer pc.conn.Close()

	err = pc.handshake(p.Password, p.Hash)
	if err != nil {
		return err
	}

	err = pc.write(commands[c])
	if err != nil {
		return fmt.Errorf("failed to send command to projector: %v", err)
	}

	resp, err := pc.read()
	if err != nil {
		return fmt.Errorf("failed to read command response from projector: %v", err)
	}

	expectedResponse := []byte{0x06, 0x89, 0x01, 0x00, 0x00, 0x0A}
	expectedResponse[3] = commands[c][3]
	expectedResponse[4] = commands[c][4]
	if !bytes.Equal(resp, expectedResponse) {
		return fmt.Errorf("unexpected response:\n%s", hex.Dump(resp))
	}

	return nil
}

func ScanForProjectors(debug bool) []string {
	var projectors []string
	var wg sync.WaitGroup
	var mutex sync.Mutex

	subnet, err := getLocalSubnet()
	if err != nil {
		fmt.Printf("Error getting local subnet: %v\n", err)
		return nil
	}

	// Create a channel to limit concurrent connections
	semaphore := make(chan struct{}, 50)

	// Scan all IPs in the subnet
	for i := 1; i < 255; i++ {
		wg.Add(1)
		go func(ip int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			ipAddr := fmt.Sprintf("%s.%d", subnet, ip)
			if isPortOpen(ipAddr, jvcPort) {
				mutex.Lock()
				projectors = append(projectors, ipAddr)
				mutex.Unlock()
			}
		}(i)
	}

	wg.Wait()
	return projectors
}

func getLocalSubnet() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue
			}

			return fmt.Sprintf("%d.%d.%d", ip[0], ip[1], ip[2]), nil
		}
	}

	return "", fmt.Errorf("no suitable network interface found")
}

func isPortOpen(host string, port int) bool {
	timeout := time.Second * 2
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
