package cpio

import "io"

const (
	magicNewc = "070701"
	magicCrc  = "070702"
)

// Entry represents a single CPIO newc entry.
//
// The Data field is only set for regular files.
type Entry struct {
	Name     string
	Mode     uint32
	UID      uint32
	GID      uint32
	NLink    uint32
	MTime    uint32
	FileSize uint32
	Data     []byte
}

func (e *Entry) IsTrailer() bool { return e != nil && e.Name == "TRAILER!!!" }

// Reader reads CPIO newc/crc archives.
type Reader struct {
	r    io.Reader
	done bool
}

func NewReader(r io.Reader) *Reader { return &Reader{r: r} }

// Writer writes CPIO newc archives.
type Writer struct {
	w      io.Writer
	ino    uint32
	uid    uint32
	gid    uint32
	mtime  uint32
	closed bool
}

type WriterOption func(*Writer)

func WithMTimeUnix(mtime uint32) WriterOption { return func(w *Writer) { w.mtime = mtime } }
func WithUIDGID(uid, gid uint32) WriterOption {
	return func(w *Writer) {
		w.uid = uid
		w.gid = gid
	}
}

func NewWriter(w io.Writer, opts ...WriterOption) *Writer {
	wr := &Writer{w: w, ino: 1}
	for _, opt := range opts {
		opt(wr)
	}
	return wr
}
