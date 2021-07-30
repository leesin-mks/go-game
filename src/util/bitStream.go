package util

import "log"

//----------------bitsream---------------
//for example
//buf := make([]byte, 256)
//var bitstream base.BitStream
//bitstream.BuildPacketStream(buf, 256)
//bitstream.WriteInt(1000, 16)
// or
//bitstream := NewBitStream(buf)
//----------------------------------------

const (
	Bit8   = 8
	Bit16  = 16
	Bit32  = 32
	Bit64  = 64
	Bit128 = 128
	MAX_PACKET = 1 * 1024 * 1024 //1MB
	MAX_CLIENT_PACKET = 10 * 1024 //10KB
)

type (
	BitStream struct {
		dataPtr        []byte
		bitNum         int
		flagNum        int
		tailFlag       bool
		bufSize        int
		bitsLimite     int
		error          bool
		maxReadBitNum  int
		maxWriteBitNum int
	}

	IBitStream interface {
		BuildPacketStream([]byte, int) bool
		setBuffer([]byte, int, int)
		GetBuffer() []byte
		GetBytePtr() []byte
		GetReadByteSize() int
		GetCurPos() int
		GetPosition() int
		GetStreamSize() int
		SetPosition(int) bool
		clear()
		resize() bool

		WriteBits([]byte, int)
		ReadBits(int)[]byte
		WriteInt(int, int)
		ReadInt(int) int
		ReadFlag() bool
		WriteFlag(bool) bool
		WriteString(string)
		ReadString() string

		WriteInt64(int64, int)
		ReadInt64(int) int64
		WriteFloat(float32)
		ReadFloat() float32
		WriteFloat64(float64)
		ReadFloat64() float64
	}
)

func Assert(x bool, y string) {
	if bool(x) == false {
		log.Printf("\nFatal :{%s}", y)
	}
}

func (this *BitStream) BuildPacketStream(buffer []byte, writeSize int) bool {
	if writeSize <= 0 {
		return false
	}

	this.setBuffer(buffer, writeSize, -1)
	this.SetPosition(0)
	return true
}

func (this *BitStream) setBuffer(bufPtr []byte, size int, maxSize int) {
	this.dataPtr = bufPtr
	this.bitNum = 0
	this.flagNum = 0
	this.tailFlag = false
	this.bufSize = size
	this.maxReadBitNum = size << 3
	if maxSize < 0 {
		maxSize = size
	}
	this.maxWriteBitNum = maxSize << 3
	this.bitsLimite = size
	this.error = false
}

func (this *BitStream) GetBuffer() []byte {
	return this.dataPtr[0:this.GetPosition()]
}

func (this *BitStream) GetBytePtr() []byte {
	return this.dataPtr[this.GetPosition():]
}

func (this *BitStream) GetReadByteSize() int {
	return (this.maxReadBitNum >> 3) - this.GetPosition()
}

func (this *BitStream) GetCurPos() int {
	return this.bitNum
}

func (this *BitStream) GetPosition() int {
	return (this.bitNum + 7) >> 3
}

func (this *BitStream) GetStreamSize() int {
	return this.bufSize
}

func (this *BitStream) SetPosition(pos int) bool {
	Assert(pos == 0 || this.flagNum == 0, "不正确的setPosition调用")
	if pos != 0 && this.flagNum != 0 {
		return false
	}

	this.bitNum = pos << 3
	this.flagNum = 0
	return true
}

func (this *BitStream) clear() {
	var buff []byte
	buff = make([]byte, this.bufSize)
	this.dataPtr = buff
}

func (this *BitStream) resize() bool{
	//fmt.Println("BitStream Resize")
	this.dataPtr = append(this.dataPtr, make([]byte, this.bitsLimite)...)
	size := this.bitsLimite * 2
	if size <= 0 || size >= MAX_PACKET * 2{
		return false
	}
	this.bufSize = size
	this.maxReadBitNum = size << 3
	this.maxWriteBitNum = size << 3
	this.bitsLimite = size
	return true
}

func (this *BitStream) WriteBits(bitPtr []byte, bitCount int) {
	if bitCount == 0 {
		return
	}

	if this.tailFlag {
		this.error = true
		Assert(false, "Out of range write")
		return
	}

	if (bitCount & 0x7) != 0 {
		bitCount = (bitCount & ^0x7) + 8
	}

	for bitCount+this.bitNum > this.maxWriteBitNum {
		if !this.resize(){
			this.error = true
			Assert(false, "Out of range write")
			return
		}
	}

	bitNum := this.bitNum >> 3
	byteCount := (bitCount + 7) >> 3
	copy(this.dataPtr[bitNum:], bitPtr[:byteCount])
	/*for i, v := range bitPtr[:byteCount] {
		this.dataPtr[bitNum+i] = v
	}*/
	this.bitNum += bitCount
}

func (this *BitStream) ReadBits(bitCount int) []byte{
	if bitCount == 0 {
		return []byte{}
	}

	if this.tailFlag {
		this.error = true
		Assert(false, "Out of range read")
		return []byte{}
	}

	if (bitCount & 0x7) != 0 {
		bitCount = (bitCount & ^0x7) + 8
	}

	for bitCount+this.bitNum > this.maxReadBitNum {
		if !this.resize(){
			this.error = true
			Assert(false, "Out of range read")
			return []byte{}
		}
	}

	byteCount := (bitCount + 7) >> 3
	bitNum := this.bitNum >> 3
	stPtr := this.dataPtr[bitNum : bitNum + byteCount]
	this.bitNum += bitCount
	return stPtr
}

func (this *BitStream) WriteInt(value int, bitCount int) {
	this.WriteBits(IntToBytes(value), bitCount)
}

func (this *BitStream) ReadInt(bitCount int) int {
	var ret int
	buf := this.ReadBits(bitCount)
	ret = BytesToInt(buf)
	if bitCount == Bit32 {
		return int(ret)
	} else {
		ret &= (1 << uint32(bitCount)) - 1
	}

	return int(ret)
}

func (this *BitStream) ReadFlag() bool {
	if ((this.flagNum - (this.flagNum>>3)<<3) == 0) && !this.tailFlag {
		this.flagNum = this.bitNum
		if this.bitNum+8 < this.maxReadBitNum {
			this.bitNum += 8
		} else {
			if !this.resize(){
				this.tailFlag = true
			}else{
				this.bitNum += 8
			}
		}
	}

	if this.flagNum+1 > this.maxReadBitNum {
		this.error = true
		Assert(false, "Out of range read")
		return false
	}

	mask := 1 << uint32(this.flagNum&0x7)
	ret := (int(this.dataPtr[(this.flagNum>>3)]) & mask) != 0
	this.flagNum++
	return ret
}

func (this *BitStream) WriteFlag(value bool) bool {
	if ((this.flagNum - (this.flagNum>>3)<<3) == 0) && !this.tailFlag {
		this.flagNum = this.bitNum

		if this.bitNum+8 < this.maxWriteBitNum {
			this.bitNum += 8 //跳开8个用于写flag
		} else {
			if !this.resize(){
				this.tailFlag = true
			}else {
				this.bitNum += 8 //跳开8个用于写flag
			}
		}
	}

	if this.flagNum+1 > this.maxWriteBitNum {
		this.error = true
		Assert(false, "Out of range write")
		return false
	}

	if value {
		this.dataPtr[(this.flagNum >> 3)] |= (1 << uint32(this.flagNum&0x7))
	} else {
		this.dataPtr[(this.flagNum >> 3)] &= ^(1 << uint32(this.flagNum&0x7))
	}
	this.flagNum++
	return (value)
}

func (this *BitStream) ReadString() string {
	if this.ReadFlag() {
		nLen := this.ReadInt(Bit16)
		buf := this.ReadBits(nLen<<3)
		return string(buf)
	}
	return string("")
}

func (this *BitStream) WriteString(value string) {
	buf := []byte(value)
	nLen := len(buf)

	if this.WriteFlag(nLen > 0) {
		this.WriteInt(nLen, Bit16)
		this.WriteBits(buf, nLen<<3,)
	}
}

func (this *BitStream) WriteInt64(value int64, bitCount int) {
	this.WriteBits(Int64ToBytes(value), bitCount)
}

func (this *BitStream) ReadInt64(bitCount int) int64 {
	var ret int64
	buf := this.ReadBits(bitCount)
	ret = BytesToInt64(buf)
	if bitCount == Bit64 {
		return int64(ret)
	} else {
		ret &= (1 << uint64(bitCount)) - 1
	}

	return int64(ret)
}

func (this *BitStream) WriteFloat(value float32) {
	this.WriteBits(Float32ToByte(value), Bit32)
}

func (this *BitStream) ReadFloat() float32 {
	var ret float32
	buf := this.ReadBits(Bit32)
	ret = BytesToFloat32(buf)

	return float32(ret)
}

func (this *BitStream) WriteFloat64(value float64) {
	this.WriteBits(Float64ToByte(value), Bit64)
}

func (this *BitStream) ReadFloat64() float64 {
	var ret float64
	buf := this.ReadBits(Bit64)
	ret = BytesToFloat64(buf)

	return float64(ret)
}

func NewBitStream(buf []byte, nLen int) *BitStream {
	var bitstream BitStream
	bitstream.BuildPacketStream(buf, nLen)
	return &bitstream
}
