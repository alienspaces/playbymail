package compress

import (
	"io"

	"github.com/klauspost/compress/zstd"
)

// See https://github.com/klauspost/compress/tree/master/zstd#zstd

var (
	encoder, _ = zstd.NewWriter(nil)
	decoder, _ = zstd.NewReader(nil, zstd.WithDecoderConcurrency(0))
)

// ZStdCompressBuffer allows optionally passing an out buffer. If nil, a new buffer will be allocated.
func ZStdCompressBuffer(in, out []byte) []byte {
	return encoder.EncodeAll(in, out)
}

// ZStdDecompressBuffer allows optionally passing an out buffer. If nil, a new buffer will be allocated.
func ZStdDecompressBuffer(in, out []byte) ([]byte, error) {
	return decoder.DecodeAll(in, out)
}

func ZStdCompressStream(in io.Reader, out io.Writer) error {
	enc, err := zstd.NewWriter(out)
	if err != nil {
		return err
	}

	_, err = io.Copy(enc, in)
	if err != nil {
		enc.Close()
		return err
	}

	return enc.Close()
}

func ZStdCompressStreamPiped(in io.Reader, out io.Writer) (io.Reader, error) {
	enc, err := zstd.NewWriter(out)
	if err != nil {
		return nil, err
	}
	defer enc.Close()

	pr, pw := io.Pipe()
	go func() {
		_, err := io.Copy(enc, in)
		if err != nil {
			pw.CloseWithError(err)
			return
		}

		pw.Close()
	}()

	return pr, nil
}

func ZStdDecompressStream(in io.Reader, out io.Writer) error {
	d, err := zstd.NewReader(in, zstd.WithDecoderConcurrency(0))
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = io.Copy(out, d)
	return err
}
