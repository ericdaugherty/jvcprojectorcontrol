package jvcprojectorcontrol

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
)

type projConnection struct {
	conn  net.Conn
	debug bool
}

func (pc projConnection) handshake(password string, hash HashMode) error {
	resp, err := pc.read()
	if err != nil {
		return fmt.Errorf("failed to read initial response from projector: %v", err)
	}
	if string(resp) != "PJ_OK" {
		pc.conn.Close()
		return fmt.Errorf("unexpected response from projector: %s", resp)
	}

	pwlen := len(password)
	if pwlen > 0 && pwlen < 8 || pwlen > 10 {
		return fmt.Errorf("password length must be between 8 and 10 characters")
	}

	switch hash {
	case HashNone:
		if pwlen == 0 {
			err = pc.write([]byte("PJREQ"))
			if err != nil {
				return fmt.Errorf("failed to send PJREQ command: %v", err)
			}
		} else {
			pw := make([]byte, 10)
			copy(pw, []byte(password))
			data := append([]byte("PJREQ_"), pw...)
			err = pc.write(data)
			if err != nil {
				return fmt.Errorf("failed to send PJREQ command: %v", err)
			}
		}
	case HashJVCKW:
		if pwlen == 0 {
			return fmt.Errorf("password required for HashJVCKW")
		}
		pw := hashPassword(password, "JVCKW")
		data := append([]byte("PJREQ_"), pw...)
		err = pc.write(data)
		if err != nil {
			return fmt.Errorf("failed to send PJREQ command: %v", err)
		}
	case HashJVCKWPJ:
		if pwlen == 0 {
			return fmt.Errorf("password required for HashJVCKWPJ")
		}
		pw := hashPassword(password, "JVCKWPJ")
		data := append([]byte("PJREQ_"), pw...)
		err = pc.write(data)
		if err != nil {
			return fmt.Errorf("failed to send PJREQ command: %v", err)
		}
	}

	resp, err = pc.read()
	if err != nil {
		return fmt.Errorf("failed to read handshake response from projector: %v", err)
	}
	if string(resp) != "PJACK" {
		if string(resp) == "PJNAK" {
			return fmt.Errorf("projector rejected handshake (PJNAK) - check password")
		}
		return fmt.Errorf("unexpected response from projector: %s", resp)
	}

	return nil
}

// read reads a response from the projector.
func (pc projConnection) read() ([]byte, error) {
	if pc.conn == nil {
		return nil, fmt.Errorf("no connection established")
	}

	buffer := make([]byte, 1024)
	n, err := pc.conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	if pc.debug {
		fmt.Printf("Received %d chars:\n%s", n, hex.Dump(buffer[:n]))
	}

	return buffer[:n], nil
}

// write sends a command to the projector.
func (pc projConnection) write(command []byte) error {
	if pc.conn == nil {
		return fmt.Errorf("no connection established")
	}

	if pc.debug {
		fmt.Printf("Sending %d chars:\n%s", len(command), hex.Dump(command))
	}

	_, err := pc.conn.Write(command)
	return err
}

func checkResponse(resp []byte, expected string) bool {
	if len(resp) < len(expected) {
		return false
	}
	return string(resp) == expected
}

func hashPassword(password, hash string) []byte {
	h := sha256.Sum256([]byte(password + hash))
	return []byte(hex.EncodeToString(h[:]))
}
