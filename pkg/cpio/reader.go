package cpio

import (
	"errors"
	"io"
)

func parseHex8(b []byte) (uint32, error) {
	var v uint32
	for i := 0; i < 8 && i < len(b); i++ {
		c := b[i]
		var n uint32
		switch {
		case c >= '0' && c <= '9':
			n = uint32(c - '0')
		case c >= 'A' && c <= 'F':
			n = uint32(c-'A') + 10
		case c >= 'a' && c <= 'f':
			n = uint32(c-'a') + 10
		default:
			return 0, errors.New("cpio: invalid hex")
		}
		v = (v << 4) | n
	}
	return v, nil
}

func pad4(n uint32) uint32 {
	rem := n & 3
	if rem == 0 {
		return 0
	}
	return 4 - rem
}

func discardN(r io.Reader, n uint32) error {
	if n == 0 {
		return nil
	}
	buf := make([]byte, 128)
	for n > 0 {
		chunk := uint32(len(buf))
		if chunk > n {
			chunk = n
		}
		_, err := io.ReadFull(r, buf[:chunk])
		if err != nil {
			return err
		}
		n -= chunk
	}
	return nil
}

// Next returns the next entry. When the archive ends, it returns io.EOF.
func (rd *Reader) Next() (*Entry, error) {
	if rd == nil || rd.done {
		return nil, io.EOF
	}

	var hdr [110]byte
	n, err := io.ReadFull(rd.r, hdr[:])
	if err != nil {
		// If we didn't read anything, treat as clean EOF.
		if errors.Is(err, io.EOF) && n == 0 {
			rd.done = true
			return nil, io.EOF
		}
		return nil, err
	}

	magic := string(hdr[0:6])
	if magic != magicNewc && magic != magicCrc {
		return nil, errors.New("cpio: invalid magic")
	}

	mode, err := parseHex8(hdr[14:22])
	if err != nil {
		return nil, err
	}
	uid, err := parseHex8(hdr[22:30])
	if err != nil {
		return nil, err
	}
	gid, err := parseHex8(hdr[30:38])
	if err != nil {
		return nil, err
	}
	nlink, err := parseHex8(hdr[38:46])
	if err != nil {
		return nil, err
	}
	mtime, err := parseHex8(hdr[46:54])
	if err != nil {
		return nil, err
	}
	filesize, err := parseHex8(hdr[54:62])
	if err != nil {
		return nil, err
	}
	namesize, err := parseHex8(hdr[94:102])
	if err != nil {
		return nil, err
	}

	if namesize == 0 {
		return nil, errors.New("cpio: invalid namesize")
	}

	nameb := make([]byte, namesize)
	if _, err := io.ReadFull(rd.r, nameb); err != nil {
		return nil, err
	}
	name := ""
	for i := 0; i < len(nameb); i++ {
		if nameb[i] == 0 {
			break
		}
		name += string([]byte{nameb[i]})
	}

	// Align after header + name.
	if err := discardN(rd.r, pad4(110+namesize)); err != nil {
		return nil, err
	}

	if name == "TRAILER!!!" {
		rd.done = true
		return nil, io.EOF
	}

	data := make([]byte, filesize)
	if filesize > 0 {
		if _, err := io.ReadFull(rd.r, data); err != nil {
			return nil, err
		}
	}
	if err := discardN(rd.r, pad4(filesize)); err != nil {
		return nil, err
	}

	return &Entry{
		Name:     name,
		Mode:     mode,
		UID:      uid,
		GID:      gid,
		NLink:    nlink,
		MTime:    mtime,
		FileSize: filesize,
		Data:     data,
	}, nil
}
