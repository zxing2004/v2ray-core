package buf_test

import (
	"crypto/rand"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"

	"v2ray.com/core/common"
	. "v2ray.com/core/common/buf"
)

func TestMultiBufferRead(t *testing.T) {
	b1 := New()
	common.Must2(b1.WriteString("ab"))

	b2 := New()
	common.Must2(b2.WriteString("cd"))
	mb := MultiBuffer{b1, b2}

	bs := make([]byte, 32)
	_, nBytes := SplitBytes(mb, bs)
	if nBytes != 4 {
		t.Error("expect 4 bytes split, but got ", nBytes)
	}
	if r := cmp.Diff(bs[:nBytes], []byte("abcd")); r != "" {
		t.Error(r)
	}
}

func TestMultiBufferAppend(t *testing.T) {
	var mb MultiBuffer
	b := New()
	common.Must2(b.WriteString("ab"))
	mb = append(mb, b)
	if mb.Len() != 2 {
		t.Error("expected length 2, but got ", mb.Len())
	}
}

func TestMultiBufferSliceBySizeLarge(t *testing.T) {
	lb := make([]byte, 8*1024)
	common.Must2(io.ReadFull(rand.Reader, lb))

	mb := MergeBytes(nil, lb)

	mb, mb2 := SplitSize(mb, 1024)
	if mb2.Len() != 1024 {
		t.Error("expect length 1024, but got ", mb2.Len())
	}
	if mb.Len() != 7*1024 {
		t.Error("expect length 7*1024, but got ", mb.Len())
	}

	mb, mb3 := SplitSize(mb, 7*1024)
	if mb3.Len() != 7*1024 {
		t.Error("expect length 7*1024, but got", mb.Len())
	}

	if !mb.IsEmpty() {
		t.Error("expect empty buffer, but got ", mb.Len())
	}
}

func TestMultiBufferSplitFirst(t *testing.T) {
	b1 := New()
	b1.WriteString("b1")

	b2 := New()
	b2.WriteString("b2")

	b3 := New()
	b3.WriteString("b3")

	var mb MultiBuffer
	mb = append(mb, b1, b2, b3)

	mb, c1 := SplitFirst(mb)
	if diff := cmp.Diff(b1.String(), c1.String()); diff != "" {
		t.Error(diff)
	}

	mb, c2 := SplitFirst(mb)
	if diff := cmp.Diff(b2.String(), c2.String()); diff != "" {
		t.Error(diff)
	}

	mb, c3 := SplitFirst(mb)
	if diff := cmp.Diff(b3.String(), c3.String()); diff != "" {
		t.Error(diff)
	}

	if !mb.IsEmpty() {
		t.Error("expect empty buffer, but got ", mb.String())
	}
}

func BenchmarkSplitBytes(b *testing.B) {
	var mb MultiBuffer
	raw := make([]byte, Size)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer := StackNew()
		buffer.Extend(Size)
		mb = append(mb, &buffer)
		mb, _ = SplitBytes(mb, raw)
	}
}
