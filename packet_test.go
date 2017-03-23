package radigo

import (
	"reflect"
	"testing"
)

func TestPacketDecode(t *testing.T) {
	// sample packet taken out of RFC2865 -Section 7.2.
	encdPkt := []byte{
		0x01, 0x01, 0x00, 0x47, 0x2a, 0xee, 0x86, 0xf0, 0x8d, 0x0d, 0x55, 0x96, 0x9c, 0xa5, 0x97, 0x8e,
		0x0d, 0x33, 0x67, 0xa2, 0x01, 0x08, 0x66, 0x6c, 0x6f, 0x70, 0x73, 0x79, 0x03, 0x13, 0x16, 0xe9,
		0x75, 0x57, 0xc3, 0x16, 0x18, 0x58, 0x95, 0xf2, 0x93, 0xff, 0x63, 0x44, 0x07, 0x72, 0x75, 0x04,
		0x06, 0xc0, 0xa8, 0x01, 0x10, 0x05, 0x06, 0x00, 0x00, 0x00, 0x14, 0x06, 0x06, 0x00, 0x00, 0x00,
		0x02, 0x07, 0x06, 0x00, 0x00, 0x00, 0x01,
	}
	ePkt := &Packet{
		Code:       AccessRequest,
		Identifier: 1,
		Authenticator: [16]byte{0x2a, 0xee, 0x86, 0xf0, 0x8d, 0x0d, 0x55, 0x96, 0x9c, 0xa5, 0x97, 0x8e,
			0x0d, 0x33, 0x67, 0xa2},
		AVPs: []*AVP{
			&AVP{
				Type:  UserName,
				Value: []byte{0x66, 0x6c, 0x6f, 0x70, 0x73, 0x79}, // flopsy
			},
			&AVP{
				Type: CHAPPassword,
				Value: []byte{0x16, 0xe9,
					0x75, 0x57, 0xc3, 0x16, 0x18, 0x58, 0x95, 0xf2, 0x93, 0xff, 0x63, 0x44, 0x07, 0x72, 0x75}, // 3
			},
			&AVP{
				Type:  NASIPAddress,
				Value: []byte{0xc0, 0xa8, 0x01, 0x10}, // 192.168.1.16
			},
			&AVP{
				Type:  NASPort,
				Value: []byte{0x00, 0x00, 0x00, 0x14}, // 20
			},
			&AVP{
				Type:  ServiceType,
				Value: []byte{0x00, 0x00, 0x00, 0x02}, // 2
			},
			&AVP{
				Type:  FramedProtocol,
				Value: []byte{0x00, 0x00, 0x00, 0x01}, // 1
			},
		},
	}
	pkt := new(Packet)
	if err := pkt.Decode(encdPkt); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(ePkt, pkt) {
		t.Errorf("Expecting: %+v, received: %+v", ePkt, pkt)
	}
}

func TestPacketEncode(t *testing.T) {
	pkt := &Packet{
		Code:       AccessAccept,
		Identifier: 1,
		Authenticator: [16]byte{0x2a, 0xee, 0x86, 0xf0, 0x8d, 0x0d, 0x55, 0x96, 0x9c, 0xa5, 0x97, 0x8e,
			0x0d, 0x33, 0x67, 0xa2}, // Authenticator out of origin request
		AVPs: []*AVP{
			&AVP{
				Type:  ServiceType,
				Value: []byte{0x00, 0x00, 0x00, 0x02}, // 2
			},
			&AVP{
				Type:  FramedProtocol,
				Value: []byte{0x00, 0x00, 0x00, 0x01}, // 1
			},
			&AVP{
				Type:  FramedIPAddress,
				Value: []byte{0xff, 0xff, 0xff, 0xfe}, // 255.255.255.254
			},
			&AVP{
				Type:  FramedRouting,
				Value: []byte{0x00, 0x00, 0x00, 0x02}, // 0
			},
			&AVP{
				Type:  FramedCompression,
				Value: []byte{0x00, 0x00, 0x00, 0x01}, // 1
			},
			&AVP{
				Type:  FramedMTU,
				Value: []byte{0x00, 0x00, 0x05, 0xdc}, // 1500
			},
		},
	}
	ePktEncd := []byte{
		0x02, 0x01, 0x00, 0x38, 0x71, 0xf7, 0xe6, 0x82, 0x87, 0x23, 0xc8, 0x4a, 0xa0, 0xc3, 0x1d, 0xec,
		0x3f, 0x21, 0x43, 0xf7, 0x06, 0x06, 0x00, 0x00, 0x00, 0x02, 0x07, 0x06, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0xff, 0xff, 0xff, 0xfe, 0x0a, 0x06, 0x00, 0x00, 0x00, 0x02, 0x0d, 0x06, 0x00, 0x00,
		0x00, 0x01, 0x0c, 0x06, 0x00, 0x00, 0x05, 0xdc,
	}
	var buf [4096]byte
	n, err := pkt.Encode(buf[:])
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(ePktEncd, buf[:n]) { // except authenticator which in RFC is randomly generated, ours is hash with secret on original
		t.Errorf("Expecting: %+v, received: %+v", ePktEncd, buf[:n])
	}

}
