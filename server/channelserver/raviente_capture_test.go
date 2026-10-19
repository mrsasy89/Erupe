package channelserver

import (
	"bytes"
	"encoding/hex"
	"testing"

	"erupe-ce/common/byteframe"
	"erupe-ce/network"
	"erupe-ce/network/mhfpacket"
)

func readRavienteCaptureAck(t *testing.T, s *Session) *mhfpacket.MsgSysAck {
	t.Helper()
	select {
	case p := <-s.sendPackets:
		bf := byteframe.NewByteFrameFromBytes(p.data)
		if opcode := bf.ReadUint16(); opcode != uint16(network.MSG_SYS_ACK) {
			t.Fatalf("opcode = %#x, want MSG_SYS_ACK", opcode)
		}
		ack := &mhfpacket.MsgSysAck{}
		if err := ack.Parse(bf, s.clientContext); err != nil {
			t.Fatal(err)
		}
		if bf.Err() != nil || ack.ErrorCode != 0 {
			t.Fatalf("invalid ACK: read error %v, error code %d", bf.Err(), ack.ErrorCode)
		}
		return ack
	default:
		t.Fatal("no ACK queued")
		return nil
	}
}

// Replay the TEST capture of Join Great Slaying on 2026-10-19.
// These expectations document the baseline, not a validated registration fix.
func TestRavienteRegisterEventCapturedRequest(t *testing.T) {
	for _, tt := range []struct {
		name     string
		check    bool
		existing bool
		want     []byte
	}{
		{"check_without_event", true, false, []byte{0, 0, 0, 0}},
		{"check_existing_event", true, true, []byte{17, 16, 0, 1}},
		{"creation_request", false, false, []byte{17, 16, 0, 1}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			srv := createMockServerWithRaviente()
			srv.raviente.id = 1
			s := createMockSession(1, srv)
			if tt.existing {
				name := "hs_l0u111111k3"
				srv.semaphore[name] = &Semaphore{name: name, clients: map[*Session]uint32{s: 1}}
			}
			// Payload excludes opcode 006e and trailing MSG_SYS_END 0010.
			raw, err := hex.DecodeString("01fe00350001001100100100")
			if err != nil {
				t.Fatal(err)
			}
			if !tt.check {
				raw[10] = 0
			}
			bf := byteframe.NewByteFrameFromBytes(raw)
			pkt := &mhfpacket.MsgMhfRegisterEvent{}
			if err := pkt.Parse(bf, s.clientContext); err != nil {
				t.Fatal(err)
			}
			if bf.Err() != nil || pkt.WorldID != 17 || pkt.LandID != 16 || pkt.CheckOnly != tt.check {
				t.Fatalf("unexpected parsed capture: %+v, read error %v", pkt, bf.Err())
			}
			handleMsgMhfRegisterEvent(s, pkt)
			ack := readRavienteCaptureAck(t, s)
			if ack.AckHandle != 0x01fe0035 || ack.IsBufferResponse || !bytes.Equal(ack.AckData, tt.want) {
				t.Fatalf("ACK = %+v, want simple ACK with data %x", ack, tt.want)
			}
		})
	}
}

func TestRavienteCapturedGeneralRegisterRead(t *testing.T) {
	srv := createMockServerWithRaviente()
	s := createMockSession(1, srv)
	raw, err := hex.DecodeString("01fe0037000600000c000000")
	if err != nil {
		t.Fatal(err)
	}
	bf := byteframe.NewByteFrameFromBytes(raw)
	pkt := &mhfpacket.MsgSysLoadRegister{}
	if err := pkt.Parse(bf, s.clientContext); err != nil {
		t.Fatal(err)
	}
	if bf.Err() != nil || pkt.RegisterID != 0x60000 || pkt.Values != 12 {
		t.Fatalf("unexpected parsed capture: %+v, read error %v", pkt, bf.Err())
	}
	handleMsgSysLoadRegister(s, pkt)
	ack := readRavienteCaptureAck(t, s)
	want := make([]byte, 2+12*4)
	want[1] = 12
	if ack.AckHandle != 0x01fe0037 || !ack.IsBufferResponse || !bytes.Equal(ack.AckData, want) {
		t.Fatalf("ACK = %+v, want buffer ACK with data %x", ack, want)
	}
}
