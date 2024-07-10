package types

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"encoding/binary"
	io "io"
)

func DenomToID(denom string) (uint64, error) {
	hash := sha256.New()
	if _, err := hash.Write([]byte(denom)); err != nil {
		return 0, err
	}

	md := hash.Sum(nil)
	return binary.LittleEndian.Uint64(md), nil
}

func (m VoteExtension2) Compress() ([]byte, error) {
	// Encode vote extension to bytes
	bz, err := m.Marshal()
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer

	// we use the best compression level as size reduction is prioritized
	w := zlib.NewWriter(&b)
	defer w.Close()

	// write and flush the buffer
	if _, err := w.Write(bz); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (m *VoteExtension2) Decompress(bz []byte) error {
	if len(bz) == 0 {
		return nil
	}
	r, err := zlib.NewReader(bytes.NewReader(bz))
	if err != nil {
		return err
	}
	r.Close()

	// read bytes and return
	veBz, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return m.Unmarshal(veBz)
}
