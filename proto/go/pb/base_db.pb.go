// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: base_db.proto

package pb

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

//@db:string|Redis|DBConn:%s,UID
type DBConnect struct {
	UID                  string   `protobuf:"bytes,1,opt,name=UID,proto3" json:"UID,omitempty"`
	GatewayAddr          string   `protobuf:"bytes,2,opt,name=GatewayAddr,proto3" json:"GatewayAddr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DBConnect) Reset()         { *m = DBConnect{} }
func (m *DBConnect) String() string { return proto.CompactTextString(m) }
func (*DBConnect) ProtoMessage()    {}
func (*DBConnect) Descriptor() ([]byte, []int) {
	return fileDescriptor_6d05003c37bbec36, []int{0}
}
func (m *DBConnect) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DBConnect) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DBConnect.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DBConnect) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DBConnect.Merge(m, src)
}
func (m *DBConnect) XXX_Size() int {
	return m.Size()
}
func (m *DBConnect) XXX_DiscardUnknown() {
	xxx_messageInfo_DBConnect.DiscardUnknown(m)
}

var xxx_messageInfo_DBConnect proto.InternalMessageInfo

func (m *DBConnect) GetUID() string {
	if m != nil {
		return m.UID
	}
	return ""
}

func (m *DBConnect) GetGatewayAddr() string {
	if m != nil {
		return m.GatewayAddr
	}
	return ""
}

//@db:string|Redis|DBSessionKey:%s,SessionKey
type DBSessionKey struct {
	UID                  string          `protobuf:"bytes,1,opt,name=UID,proto3" json:"UID,omitempty"`
	SessionKey           string          `protobuf:"bytes,2,opt,name=SessionKey,proto3" json:"SessionKey,omitempty"`
	DeviceToken          string          `protobuf:"bytes,3,opt,name=DeviceToken,proto3" json:"DeviceToken,omitempty"`
	SendRspID            int64           `protobuf:"varint,4,opt,name=SendRspID,proto3" json:"SendRspID,omitempty"`
	AckPeerID            int64           `protobuf:"varint,5,opt,name=AckPeerID,proto3" json:"AckPeerID,omitempty"`
	PacketList           []*Packet       `protobuf:"bytes,6,rep,name=PacketList,proto3" json:"PacketList,omitempty"`
	RcvClientList        map[int64]int32 `protobuf:"bytes,7,rep,name=RcvClientList,proto3" json:"RcvClientList,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *DBSessionKey) Reset()         { *m = DBSessionKey{} }
func (m *DBSessionKey) String() string { return proto.CompactTextString(m) }
func (*DBSessionKey) ProtoMessage()    {}
func (*DBSessionKey) Descriptor() ([]byte, []int) {
	return fileDescriptor_6d05003c37bbec36, []int{1}
}
func (m *DBSessionKey) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DBSessionKey) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DBSessionKey.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DBSessionKey) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DBSessionKey.Merge(m, src)
}
func (m *DBSessionKey) XXX_Size() int {
	return m.Size()
}
func (m *DBSessionKey) XXX_DiscardUnknown() {
	xxx_messageInfo_DBSessionKey.DiscardUnknown(m)
}

var xxx_messageInfo_DBSessionKey proto.InternalMessageInfo

func (m *DBSessionKey) GetUID() string {
	if m != nil {
		return m.UID
	}
	return ""
}

func (m *DBSessionKey) GetSessionKey() string {
	if m != nil {
		return m.SessionKey
	}
	return ""
}

func (m *DBSessionKey) GetDeviceToken() string {
	if m != nil {
		return m.DeviceToken
	}
	return ""
}

func (m *DBSessionKey) GetSendRspID() int64 {
	if m != nil {
		return m.SendRspID
	}
	return 0
}

func (m *DBSessionKey) GetAckPeerID() int64 {
	if m != nil {
		return m.AckPeerID
	}
	return 0
}

func (m *DBSessionKey) GetPacketList() []*Packet {
	if m != nil {
		return m.PacketList
	}
	return nil
}

func (m *DBSessionKey) GetRcvClientList() map[int64]int32 {
	if m != nil {
		return m.RcvClientList
	}
	return nil
}

func init() {
	proto.RegisterType((*DBConnect)(nil), "forevernine.base.proto.DBConnect")
	proto.RegisterType((*DBSessionKey)(nil), "forevernine.base.proto.DBSessionKey")
	proto.RegisterMapType((map[int64]int32)(nil), "forevernine.base.proto.DBSessionKey.RcvClientListEntry")
}

func init() { proto.RegisterFile("base_db.proto", fileDescriptor_6d05003c37bbec36) }

var fileDescriptor_6d05003c37bbec36 = []byte{
	// 354 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0xcf, 0x4a, 0xc3, 0x40,
	0x10, 0xc6, 0xbb, 0x8d, 0xad, 0x74, 0x6b, 0x41, 0x17, 0x91, 0x50, 0x64, 0x09, 0x3d, 0xe5, 0x94,
	0x80, 0x1e, 0x14, 0x0f, 0x6a, 0xd3, 0x88, 0x14, 0x3d, 0x94, 0xad, 0x5e, 0x04, 0x91, 0xfc, 0x19,
	0x25, 0xa4, 0xdd, 0x0d, 0x49, 0x8c, 0xe4, 0x0d, 0x3d, 0xf6, 0x11, 0x34, 0x27, 0x8f, 0xbe, 0x81,
	0x92, 0xac, 0xd8, 0x14, 0xeb, 0x6d, 0xe7, 0xf7, 0xcd, 0x7c, 0x33, 0xf9, 0x08, 0xee, 0xb9, 0x4e,
	0x02, 0x0f, 0xbe, 0x6b, 0x44, 0xb1, 0x48, 0x05, 0xd9, 0x7b, 0x14, 0x31, 0x64, 0x10, 0xf3, 0x80,
	0x83, 0x51, 0x4a, 0x92, 0xf7, 0x77, 0xaa, 0xb6, 0xc8, 0xf1, 0x42, 0x48, 0x25, 0x1a, 0x9c, 0xe1,
	0x8e, 0x6d, 0x8d, 0x04, 0xe7, 0xe0, 0xa5, 0x64, 0x1b, 0x2b, 0xb7, 0x63, 0x5b, 0x45, 0x1a, 0xd2,
	0x3b, 0xac, 0x7c, 0x12, 0x0d, 0x77, 0x2f, 0x9d, 0x14, 0x5e, 0x9c, 0x7c, 0xe8, 0xfb, 0xb1, 0xda,
	0xac, 0x94, 0x3a, 0x1a, 0x7c, 0x35, 0xf1, 0x96, 0x6d, 0x4d, 0x21, 0x49, 0x02, 0xc1, 0xaf, 0x20,
	0x5f, 0x63, 0x42, 0x31, 0x5e, 0xea, 0x3f, 0x1e, 0x35, 0x52, 0x2e, 0xb1, 0x21, 0x0b, 0x3c, 0xb8,
	0x11, 0x21, 0x70, 0x55, 0x91, 0x4b, 0x6a, 0x88, 0xec, 0xe3, 0xce, 0x14, 0xb8, 0xcf, 0x92, 0x68,
	0x6c, 0xab, 0x1b, 0x1a, 0xd2, 0x15, 0xb6, 0x04, 0xa5, 0x3a, 0xf4, 0xc2, 0x09, 0x40, 0x3c, 0xb6,
	0xd5, 0x96, 0x54, 0x7f, 0x01, 0x39, 0xc5, 0x78, 0x52, 0x7d, 0xf1, 0x75, 0x90, 0xa4, 0x6a, 0x5b,
	0x53, 0xf4, 0xee, 0x01, 0x35, 0xd6, 0x27, 0x64, 0xc8, 0x4e, 0x56, 0x9b, 0x20, 0xf7, 0xb8, 0xc7,
	0xbc, 0x6c, 0x34, 0x0b, 0x80, 0x4b, 0x8b, 0xcd, 0xca, 0xe2, 0xe8, 0x3f, 0x8b, 0x7a, 0x18, 0xc6,
	0xca, 0xe4, 0x05, 0x4f, 0xe3, 0x9c, 0xad, 0xba, 0xf5, 0xcf, 0x31, 0xf9, 0xdb, 0x54, 0x86, 0x18,
	0x42, 0x5e, 0x85, 0xa8, 0xb0, 0xf2, 0x49, 0x76, 0x71, 0x2b, 0x73, 0x66, 0xcf, 0x50, 0xe5, 0xd7,
	0x62, 0xb2, 0x38, 0x69, 0x1e, 0x23, 0xcb, 0x5a, 0xbc, 0xd3, 0xc6, 0x6b, 0x41, 0xd1, 0xa2, 0xa0,
	0xe8, 0xad, 0xa0, 0xe8, 0xa3, 0xa0, 0x8d, 0xcf, 0x82, 0xa2, 0x3b, 0xbd, 0x7e, 0xa2, 0x27, 0xe6,
	0xe6, 0x3c, 0xf0, 0xa3, 0x99, 0x93, 0x9a, 0xf2, 0x3f, 0x28, 0xcf, 0x35, 0x9f, 0x84, 0x19, 0xb9,
	0x6e, 0xbb, 0x2a, 0x0e, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0x12, 0x8b, 0x6d, 0xe7, 0x49, 0x02,
	0x00, 0x00,
}

func (m *DBConnect) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DBConnect) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DBConnect) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.GatewayAddr) > 0 {
		i -= len(m.GatewayAddr)
		copy(dAtA[i:], m.GatewayAddr)
		i = encodeVarintBaseDb(dAtA, i, uint64(len(m.GatewayAddr)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.UID) > 0 {
		i -= len(m.UID)
		copy(dAtA[i:], m.UID)
		i = encodeVarintBaseDb(dAtA, i, uint64(len(m.UID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *DBSessionKey) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DBSessionKey) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DBSessionKey) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.RcvClientList) > 0 {
		for k := range m.RcvClientList {
			v := m.RcvClientList[k]
			baseI := i
			i = encodeVarintBaseDb(dAtA, i, uint64(v))
			i--
			dAtA[i] = 0x10
			i = encodeVarintBaseDb(dAtA, i, uint64(k))
			i--
			dAtA[i] = 0x8
			i = encodeVarintBaseDb(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x3a
		}
	}
	if len(m.PacketList) > 0 {
		for iNdEx := len(m.PacketList) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.PacketList[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintBaseDb(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	if m.AckPeerID != 0 {
		i = encodeVarintBaseDb(dAtA, i, uint64(m.AckPeerID))
		i--
		dAtA[i] = 0x28
	}
	if m.SendRspID != 0 {
		i = encodeVarintBaseDb(dAtA, i, uint64(m.SendRspID))
		i--
		dAtA[i] = 0x20
	}
	if len(m.DeviceToken) > 0 {
		i -= len(m.DeviceToken)
		copy(dAtA[i:], m.DeviceToken)
		i = encodeVarintBaseDb(dAtA, i, uint64(len(m.DeviceToken)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.SessionKey) > 0 {
		i -= len(m.SessionKey)
		copy(dAtA[i:], m.SessionKey)
		i = encodeVarintBaseDb(dAtA, i, uint64(len(m.SessionKey)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.UID) > 0 {
		i -= len(m.UID)
		copy(dAtA[i:], m.UID)
		i = encodeVarintBaseDb(dAtA, i, uint64(len(m.UID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintBaseDb(dAtA []byte, offset int, v uint64) int {
	offset -= sovBaseDb(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *DBConnect) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.UID)
	if l > 0 {
		n += 1 + l + sovBaseDb(uint64(l))
	}
	l = len(m.GatewayAddr)
	if l > 0 {
		n += 1 + l + sovBaseDb(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *DBSessionKey) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.UID)
	if l > 0 {
		n += 1 + l + sovBaseDb(uint64(l))
	}
	l = len(m.SessionKey)
	if l > 0 {
		n += 1 + l + sovBaseDb(uint64(l))
	}
	l = len(m.DeviceToken)
	if l > 0 {
		n += 1 + l + sovBaseDb(uint64(l))
	}
	if m.SendRspID != 0 {
		n += 1 + sovBaseDb(uint64(m.SendRspID))
	}
	if m.AckPeerID != 0 {
		n += 1 + sovBaseDb(uint64(m.AckPeerID))
	}
	if len(m.PacketList) > 0 {
		for _, e := range m.PacketList {
			l = e.Size()
			n += 1 + l + sovBaseDb(uint64(l))
		}
	}
	if len(m.RcvClientList) > 0 {
		for k, v := range m.RcvClientList {
			_ = k
			_ = v
			mapEntrySize := 1 + sovBaseDb(uint64(k)) + 1 + sovBaseDb(uint64(v))
			n += mapEntrySize + 1 + sovBaseDb(uint64(mapEntrySize))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovBaseDb(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozBaseDb(x uint64) (n int) {
	return sovBaseDb(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *DBConnect) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBaseDb
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: DBConnect: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DBConnect: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBaseDb
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthBaseDb
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBaseDb
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GatewayAddr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBaseDb
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthBaseDb
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBaseDb
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GatewayAddr = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBaseDb(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBaseDb
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *DBSessionKey) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBaseDb
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: DBSessionKey: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DBSessionKey: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBaseDb
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthBaseDb
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBaseDb
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SessionKey", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBaseDb
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthBaseDb
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBaseDb
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SessionKey = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DeviceToken", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBaseDb
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthBaseDb
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBaseDb
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DeviceToken = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SendRspID", wireType)
			}
			m.SendRspID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBaseDb
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SendRspID |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AckPeerID", wireType)
			}
			m.AckPeerID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBaseDb
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AckPeerID |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PacketList", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBaseDb
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthBaseDb
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBaseDb
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PacketList = append(m.PacketList, &Packet{})
			if err := m.PacketList[len(m.PacketList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RcvClientList", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBaseDb
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthBaseDb
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBaseDb
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.RcvClientList == nil {
				m.RcvClientList = make(map[int64]int32)
			}
			var mapkey int64
			var mapvalue int32
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowBaseDb
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowBaseDb
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapkey |= int64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else if fieldNum == 2 {
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowBaseDb
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapvalue |= int32(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipBaseDb(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthBaseDb
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.RcvClientList[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBaseDb(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBaseDb
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipBaseDb(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBaseDb
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowBaseDb
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowBaseDb
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthBaseDb
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupBaseDb
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthBaseDb
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthBaseDb        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBaseDb          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupBaseDb = fmt.Errorf("proto: unexpected end of group")
)
