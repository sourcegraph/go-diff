// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: diff.proto

package diff

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
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
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// A FileDiff represents a unified diff for a single file.
//
// A file unified diff has a header that resembles the following:
//
//   --- oldname	2009-10-11 15:12:20.000000000 -0700
//   +++ newname	2009-10-11 15:12:30.000000000 -0700
type FileDiff struct {
	// the original name of the file
	OrigName string `protobuf:"bytes,1,opt,name=OrigName,proto3" json:"OrigName,omitempty"`
	// the original timestamp (nil if not present)
	OrigTime []byte `protobuf:"bytes,2,opt,name=OrigTime,proto3" json:"OrigTime,omitempty"`
	// the new name of the file (often same as OrigName)
	NewName string `protobuf:"bytes,3,opt,name=NewName,proto3" json:"NewName,omitempty"`
	// the new timestamp (nil if not present)
	NewTime []byte `protobuf:"bytes,4,opt,name=NewTime,proto3" json:"NewTime,omitempty"`
	// extended header lines (e.g., git's "new mode <mode>", "rename from <path>", etc.)
	Extended []string `protobuf:"bytes,5,rep,name=Extended,proto3" json:"Extended,omitempty"`
	// hunks that were changed from orig to new
	Hunks                []*Hunk  `protobuf:"bytes,6,rep,name=Hunks,proto3" json:"Hunks,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FileDiff) Reset()         { *m = FileDiff{} }
func (m *FileDiff) String() string { return proto.CompactTextString(m) }
func (*FileDiff) ProtoMessage()    {}
func (*FileDiff) Descriptor() ([]byte, []int) {
	return fileDescriptor_686521effc814b25, []int{0}
}
func (m *FileDiff) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *FileDiff) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_FileDiff.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *FileDiff) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FileDiff.Merge(m, src)
}
func (m *FileDiff) XXX_Size() int {
	return m.Size()
}
func (m *FileDiff) XXX_DiscardUnknown() {
	xxx_messageInfo_FileDiff.DiscardUnknown(m)
}

var xxx_messageInfo_FileDiff proto.InternalMessageInfo

// A Hunk represents a series of changes (additions or deletions) in a file's
// unified diff.
type Hunk struct {
	// starting line number in original file
	OrigStartLine int32 `protobuf:"varint,1,opt,name=OrigStartLine,proto3" json:"OrigStartLine,omitempty"`
	// number of lines the hunk applies to in the original file
	OrigLines int32 `protobuf:"varint,2,opt,name=OrigLines,proto3" json:"OrigLines,omitempty"`
	// if > 0, then the original file had a 'No newline at end of file' mark at this offset
	OrigNoNewlineAt int32 `protobuf:"varint,3,opt,name=OrigNoNewlineAt,proto3" json:"OrigNoNewlineAt,omitempty"`
	// starting line number in new file
	NewStartLine int32 `protobuf:"varint,4,opt,name=NewStartLine,proto3" json:"NewStartLine,omitempty"`
	// number of lines the hunk applies to in the new file
	NewLines int32 `protobuf:"varint,5,opt,name=NewLines,proto3" json:"NewLines,omitempty"`
	// optional section heading
	Section string `protobuf:"bytes,6,opt,name=Section,proto3" json:"Section,omitempty"`
	// 0-indexed line offset in unified file diff (including section headers); this is
	// only set when Hunks are read from entire file diff (i.e., when ReadAllHunks is
	// called) This accounts for hunk headers, too, so the StartPosition of the first
	// hunk will be 1.
	StartPosition int32 `protobuf:"varint,7,opt,name=StartPosition,proto3" json:"StartPosition,omitempty"`
	// hunk body (lines prefixed with '-', '+', or ' ')
	Body                 []byte   `protobuf:"bytes,8,opt,name=Body,proto3" json:"Body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Hunk) Reset()         { *m = Hunk{} }
func (m *Hunk) String() string { return proto.CompactTextString(m) }
func (*Hunk) ProtoMessage()    {}
func (*Hunk) Descriptor() ([]byte, []int) {
	return fileDescriptor_686521effc814b25, []int{1}
}
func (m *Hunk) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Hunk) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Hunk.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Hunk) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Hunk.Merge(m, src)
}
func (m *Hunk) XXX_Size() int {
	return m.Size()
}
func (m *Hunk) XXX_DiscardUnknown() {
	xxx_messageInfo_Hunk.DiscardUnknown(m)
}

var xxx_messageInfo_Hunk proto.InternalMessageInfo

// A Stat is a diff stat that represents the number of lines added/changed/deleted.
type Stat struct {
	// number of lines added
	Added int32 `protobuf:"varint,1,opt,name=Added,proto3" json:""`
	// number of lines changed
	Changed int32 `protobuf:"varint,2,opt,name=Changed,proto3" json:""`
	// number of lines deleted
	Deleted              int32    `protobuf:"varint,3,opt,name=Deleted,proto3" json:""`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Stat) Reset()         { *m = Stat{} }
func (m *Stat) String() string { return proto.CompactTextString(m) }
func (*Stat) ProtoMessage()    {}
func (*Stat) Descriptor() ([]byte, []int) {
	return fileDescriptor_686521effc814b25, []int{2}
}
func (m *Stat) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Stat) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Stat.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Stat) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Stat.Merge(m, src)
}
func (m *Stat) XXX_Size() int {
	return m.Size()
}
func (m *Stat) XXX_DiscardUnknown() {
	xxx_messageInfo_Stat.DiscardUnknown(m)
}

var xxx_messageInfo_Stat proto.InternalMessageInfo

func init() {
	proto.RegisterType((*FileDiff)(nil), "diff.FileDiff")
	proto.RegisterType((*Hunk)(nil), "diff.Hunk")
	proto.RegisterType((*Stat)(nil), "diff.Stat")
}

func init() { proto.RegisterFile("diff.proto", fileDescriptor_686521effc814b25) }

var fileDescriptor_686521effc814b25 = []byte{
	// 381 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x92, 0x5f, 0x4e, 0xf2, 0x50,
	0x10, 0xc5, 0x29, 0xb4, 0xfc, 0x99, 0x8f, 0x2f, 0x9a, 0xfb, 0xd4, 0x10, 0x53, 0x9b, 0xc6, 0x87,
	0xbe, 0x08, 0x89, 0xae, 0x00, 0x44, 0xe3, 0x83, 0xa9, 0xa6, 0xb8, 0x81, 0x96, 0x4e, 0xcb, 0x8d,
	0xd0, 0x6b, 0xe0, 0x92, 0xea, 0x0e, 0x5c, 0x90, 0x8b, 0xe0, 0xd1, 0x15, 0x18, 0xe5, 0xd1, 0x55,
	0x98, 0x3b, 0xb7, 0x85, 0xe0, 0xdb, 0x9c, 0xdf, 0xe9, 0xdc, 0x99, 0x33, 0x29, 0x40, 0xc2, 0xd3,
	0xb4, 0xff, 0xbc, 0x14, 0x52, 0x30, 0x53, 0xd5, 0xbd, 0xf3, 0x8c, 0xcb, 0xd9, 0x3a, 0xee, 0x4f,
	0xc5, 0x62, 0x90, 0x89, 0x4c, 0x0c, 0xc8, 0x8c, 0xd7, 0x29, 0x29, 0x12, 0x54, 0xe9, 0x26, 0xef,
	0xdd, 0x80, 0xf6, 0x0d, 0x9f, 0xe3, 0x98, 0xa7, 0x29, 0xeb, 0x41, 0xfb, 0x7e, 0xc9, 0xb3, 0x20,
	0x5a, 0xa0, 0x6d, 0xb8, 0x86, 0xdf, 0x09, 0x77, 0xba, 0xf2, 0x1e, 0xf9, 0x02, 0xed, 0xba, 0x6b,
	0xf8, 0xdd, 0x70, 0xa7, 0x99, 0x0d, 0xad, 0x00, 0x0b, 0x6a, 0x6b, 0x50, 0x5b, 0x25, 0x4b, 0x87,
	0x9a, 0x4c, 0x6a, 0xaa, 0xa4, 0x7a, 0xef, 0xfa, 0x45, 0x62, 0x9e, 0x60, 0x62, 0x5b, 0x6e, 0x43,
	0xcd, 0xaa, 0x34, 0x73, 0xc1, 0xba, 0x5d, 0xe7, 0x4f, 0x2b, 0xbb, 0xe9, 0x36, 0xfc, 0x7f, 0x17,
	0xd0, 0xa7, 0x94, 0x0a, 0x85, 0xda, 0xf0, 0xde, 0xea, 0x60, 0xaa, 0x8a, 0x9d, 0xc1, 0x7f, 0xb5,
	0xc6, 0x44, 0x46, 0x4b, 0x79, 0xc7, 0x73, 0xbd, 0xb7, 0x15, 0x1e, 0x42, 0x76, 0x02, 0x1d, 0x05,
	0x54, 0xbd, 0xa2, 0xed, 0xad, 0x70, 0x0f, 0x98, 0x0f, 0x47, 0x14, 0x53, 0x04, 0x58, 0xcc, 0x79,
	0x8e, 0x43, 0x49, 0x31, 0xac, 0xf0, 0x2f, 0x66, 0x1e, 0x74, 0x03, 0x2c, 0xf6, 0xc3, 0x4c, 0xfa,
	0xec, 0x80, 0xa9, 0x60, 0x01, 0x16, 0x7a, 0x94, 0x45, 0xfe, 0x4e, 0xab, 0x73, 0x4c, 0x70, 0x2a,
	0xb9, 0xc8, 0xed, 0xa6, 0x3e, 0x54, 0x29, 0x55, 0x0e, 0x7a, 0xe2, 0x41, 0xac, 0x38, 0xf9, 0x2d,
	0x9d, 0xe3, 0x00, 0x32, 0x06, 0xe6, 0x48, 0x24, 0xaf, 0x76, 0x9b, 0x6e, 0x49, 0xb5, 0x17, 0x83,
	0x39, 0x91, 0x91, 0x64, 0x3d, 0xb0, 0x86, 0x89, 0xba, 0x26, 0x5d, 0x60, 0x64, 0xfe, 0x7c, 0x9e,
	0xd6, 0x42, 0x8d, 0x98, 0x03, 0xad, 0xab, 0x59, 0x94, 0x67, 0x98, 0xe8, 0xf4, 0xa5, 0x5b, 0x41,
	0xe5, 0x8f, 0x71, 0x8e, 0x12, 0x13, 0x9d, 0xbc, 0xf2, 0x4b, 0x38, 0x3a, 0xde, 0x7c, 0x3b, 0xb5,
	0xcd, 0xd6, 0x31, 0x3e, 0xb6, 0x8e, 0xf1, 0xb5, 0x75, 0x8c, 0xb8, 0x49, 0xbf, 0xcf, 0xe5, 0x6f,
	0x00, 0x00, 0x00, 0xff, 0xff, 0xb1, 0x2f, 0x14, 0x9f, 0x81, 0x02, 0x00, 0x00,
}

func (m *FileDiff) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *FileDiff) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.OrigName) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintDiff(dAtA, i, uint64(len(m.OrigName)))
		i += copy(dAtA[i:], m.OrigName)
	}
	if len(m.OrigTime) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintDiff(dAtA, i, uint64(len(m.OrigTime)))
		i += copy(dAtA[i:], m.OrigTime)
	}
	if len(m.NewName) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintDiff(dAtA, i, uint64(len(m.NewName)))
		i += copy(dAtA[i:], m.NewName)
	}
	if len(m.NewTime) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintDiff(dAtA, i, uint64(len(m.NewTime)))
		i += copy(dAtA[i:], m.NewTime)
	}
	if len(m.Extended) > 0 {
		for _, s := range m.Extended {
			dAtA[i] = 0x2a
			i++
			l = len(s)
			for l >= 1<<7 {
				dAtA[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			dAtA[i] = uint8(l)
			i++
			i += copy(dAtA[i:], s)
		}
	}
	if len(m.Hunks) > 0 {
		for _, msg := range m.Hunks {
			dAtA[i] = 0x32
			i++
			i = encodeVarintDiff(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *Hunk) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Hunk) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.OrigStartLine != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintDiff(dAtA, i, uint64(m.OrigStartLine))
	}
	if m.OrigLines != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintDiff(dAtA, i, uint64(m.OrigLines))
	}
	if m.OrigNoNewlineAt != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintDiff(dAtA, i, uint64(m.OrigNoNewlineAt))
	}
	if m.NewStartLine != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintDiff(dAtA, i, uint64(m.NewStartLine))
	}
	if m.NewLines != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintDiff(dAtA, i, uint64(m.NewLines))
	}
	if len(m.Section) > 0 {
		dAtA[i] = 0x32
		i++
		i = encodeVarintDiff(dAtA, i, uint64(len(m.Section)))
		i += copy(dAtA[i:], m.Section)
	}
	if m.StartPosition != 0 {
		dAtA[i] = 0x38
		i++
		i = encodeVarintDiff(dAtA, i, uint64(m.StartPosition))
	}
	if len(m.Body) > 0 {
		dAtA[i] = 0x42
		i++
		i = encodeVarintDiff(dAtA, i, uint64(len(m.Body)))
		i += copy(dAtA[i:], m.Body)
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *Stat) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Stat) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Added != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintDiff(dAtA, i, uint64(m.Added))
	}
	if m.Changed != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintDiff(dAtA, i, uint64(m.Changed))
	}
	if m.Deleted != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintDiff(dAtA, i, uint64(m.Deleted))
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintDiff(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *FileDiff) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.OrigName)
	if l > 0 {
		n += 1 + l + sovDiff(uint64(l))
	}
	l = len(m.OrigTime)
	if l > 0 {
		n += 1 + l + sovDiff(uint64(l))
	}
	l = len(m.NewName)
	if l > 0 {
		n += 1 + l + sovDiff(uint64(l))
	}
	l = len(m.NewTime)
	if l > 0 {
		n += 1 + l + sovDiff(uint64(l))
	}
	if len(m.Extended) > 0 {
		for _, s := range m.Extended {
			l = len(s)
			n += 1 + l + sovDiff(uint64(l))
		}
	}
	if len(m.Hunks) > 0 {
		for _, e := range m.Hunks {
			l = e.Size()
			n += 1 + l + sovDiff(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Hunk) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.OrigStartLine != 0 {
		n += 1 + sovDiff(uint64(m.OrigStartLine))
	}
	if m.OrigLines != 0 {
		n += 1 + sovDiff(uint64(m.OrigLines))
	}
	if m.OrigNoNewlineAt != 0 {
		n += 1 + sovDiff(uint64(m.OrigNoNewlineAt))
	}
	if m.NewStartLine != 0 {
		n += 1 + sovDiff(uint64(m.NewStartLine))
	}
	if m.NewLines != 0 {
		n += 1 + sovDiff(uint64(m.NewLines))
	}
	l = len(m.Section)
	if l > 0 {
		n += 1 + l + sovDiff(uint64(l))
	}
	if m.StartPosition != 0 {
		n += 1 + sovDiff(uint64(m.StartPosition))
	}
	l = len(m.Body)
	if l > 0 {
		n += 1 + l + sovDiff(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Stat) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Added != 0 {
		n += 1 + sovDiff(uint64(m.Added))
	}
	if m.Changed != 0 {
		n += 1 + sovDiff(uint64(m.Changed))
	}
	if m.Deleted != 0 {
		n += 1 + sovDiff(uint64(m.Deleted))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovDiff(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozDiff(x uint64) (n int) {
	return sovDiff(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *FileDiff) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDiff
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
			return fmt.Errorf("proto: FileDiff: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FileDiff: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrigName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDiff
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
				return ErrInvalidLengthDiff
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDiff
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OrigName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrigTime", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDiff
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthDiff
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthDiff
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OrigTime = append(m.OrigTime[:0], dAtA[iNdEx:postIndex]...)
			if m.OrigTime == nil {
				m.OrigTime = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NewName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDiff
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
				return ErrInvalidLengthDiff
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDiff
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NewName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NewTime", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDiff
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthDiff
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthDiff
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NewTime = append(m.NewTime[:0], dAtA[iNdEx:postIndex]...)
			if m.NewTime == nil {
				m.NewTime = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Extended", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDiff
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
				return ErrInvalidLengthDiff
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDiff
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Extended = append(m.Extended, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hunks", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDiff
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
				return ErrInvalidLengthDiff
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthDiff
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hunks = append(m.Hunks, &Hunk{})
			if err := m.Hunks[len(m.Hunks)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDiff(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDiff
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthDiff
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
func (m *Hunk) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDiff
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
			return fmt.Errorf("proto: Hunk: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Hunk: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrigStartLine", wireType)
			}
			m.OrigStartLine = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDiff
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OrigStartLine |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrigLines", wireType)
			}
			m.OrigLines = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDiff
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OrigLines |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrigNoNewlineAt", wireType)
			}
			m.OrigNoNewlineAt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDiff
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OrigNoNewlineAt |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NewStartLine", wireType)
			}
			m.NewStartLine = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDiff
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NewStartLine |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NewLines", wireType)
			}
			m.NewLines = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDiff
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NewLines |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Section", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDiff
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
				return ErrInvalidLengthDiff
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDiff
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Section = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StartPosition", wireType)
			}
			m.StartPosition = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDiff
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.StartPosition |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Body", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDiff
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthDiff
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthDiff
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Body = append(m.Body[:0], dAtA[iNdEx:postIndex]...)
			if m.Body == nil {
				m.Body = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDiff(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDiff
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthDiff
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
func (m *Stat) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDiff
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
			return fmt.Errorf("proto: Stat: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Stat: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Added", wireType)
			}
			m.Added = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDiff
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Added |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Changed", wireType)
			}
			m.Changed = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDiff
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Changed |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Deleted", wireType)
			}
			m.Deleted = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDiff
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Deleted |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipDiff(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDiff
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthDiff
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
func skipDiff(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDiff
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
					return 0, ErrIntOverflowDiff
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowDiff
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
				return 0, ErrInvalidLengthDiff
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthDiff
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowDiff
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipDiff(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthDiff
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthDiff = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDiff   = fmt.Errorf("proto: integer overflow")
)
