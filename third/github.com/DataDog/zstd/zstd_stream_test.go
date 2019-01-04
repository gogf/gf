package zstd

import (
	"bytes"
	"io"
	"io/ioutil"
	"runtime/debug"
	"testing"
)

func failOnError(t *testing.T, msg string, err error) {
	if err != nil {
		debug.PrintStack()
		t.Fatalf("%s: %s", msg, err)
	}
}

func testCompressionDecompression(t *testing.T, dict []byte, payload []byte) {
	var w bytes.Buffer
	writer := NewWriterLevelDict(&w, DefaultCompression, dict)
	_, err := writer.Write(payload)
	failOnError(t, "Failed writing to compress object", err)
	failOnError(t, "Failed to close compress object", writer.Close())
	out := w.Bytes()
	t.Logf("Compressed %v -> %v bytes", len(payload), len(out))
	failOnError(t, "Failed compressing", err)
	rr := bytes.NewReader(out)
	// Check that we can decompress with Decompress()
	decompressed, err := Decompress(nil, out)
	failOnError(t, "Failed to decompress with Decompress()", err)
	if string(payload) != string(decompressed) {
		t.Fatalf("Payload did not match, lengths: %v & %v", len(payload), len(decompressed))
	}

	// Decompress
	r := NewReaderDict(rr, dict)
	dst := make([]byte, len(payload))
	n, err := r.Read(dst)
	if err != nil {
		failOnError(t, "Failed to read for decompression", err)
	}
	dst = dst[:n]
	if string(payload) != string(dst) { // Only print if we can print
		if len(payload) < 100 && len(dst) < 100 {
			t.Fatalf("Cannot compress and decompress: %s != %s", payload, dst)
		} else {
			t.Fatalf("Cannot compress and decompress (lengths: %v bytes & %v bytes)", len(payload), len(dst))
		}
	}
	// Check EOF
	n, err = r.Read(dst)
	if err != io.EOF && len(dst) > 0 { // If we want 0 bytes, that should work
		t.Fatalf("Error should have been EOF, was %s instead: (%v bytes read: %s)", err, n, dst[:n])
	}
	failOnError(t, "Failed to close decompress object", r.Close())
}

func TestResize(t *testing.T) {
	if len(resize(nil, 129)) != 129 {
		t.Fatalf("Cannot allocate new slice")
	}
	a := make([]byte, 1, 200)
	b := resize(a, 129)
	if &a[0] != &b[0] {
		t.Fatalf("Address changed")
	}
	if len(b) != 129 {
		t.Fatalf("Wrong size")
	}
	c := make([]byte, 5, 10)
	d := resize(c, 129)
	if len(d) != 129 {
		t.Fatalf("Cannot allocate a new slice")
	}
}

func TestStreamSimpleCompressionDecompression(t *testing.T) {
	testCompressionDecompression(t, nil, []byte("Hello world!"))
}

func TestStreamEmptySlice(t *testing.T) {
	testCompressionDecompression(t, nil, []byte{})
}

func TestZstdReaderLong(t *testing.T) {
	var long bytes.Buffer
	for i := 0; i < 10000; i++ {
		long.Write([]byte("Hellow World!"))
	}
	testCompressionDecompression(t, nil, long.Bytes())
}

func TestStreamCompressionDecompression(t *testing.T) {
	payload := []byte("Hello World!")
	repeat := 10000
	var intermediate bytes.Buffer
	w := NewWriterLevel(&intermediate, 4)
	for i := 0; i < repeat; i++ {
		_, err := w.Write(payload)
		failOnError(t, "Failed writing to compress object", err)
	}
	w.Close()
	// Decompress
	r := NewReader(&intermediate)
	dst := make([]byte, len(payload))
	for i := 0; i < repeat; i++ {
		n, err := r.Read(dst)
		failOnError(t, "Failed to decompress", err)
		if n != len(payload) {
			t.Fatalf("Did not read enough bytes: %v != %v", n, len(payload))
		}
		if string(dst) != string(payload) {
			t.Fatalf("Did not read the same %s != %s", string(dst), string(payload))
		}
	}
	// Check EOF
	n, err := r.Read(dst)
	if err != io.EOF {
		t.Fatalf("Error should have been EOF, was %s instead: (%v bytes read: %s)", err, n, dst[:n])
	}
	failOnError(t, "Failed to close decompress object", r.Close())
}

func TestStreamRealPayload(t *testing.T) {
	if raw == nil {
		t.Skip(ErrNoPayloadEnv)
	}
	testCompressionDecompression(t, nil, raw)
}

func TestStreamEmptyPayload(t *testing.T) {
	w := bytes.NewBuffer(nil)
	writer := NewWriter(w)
	_, err := writer.Write(nil)
	failOnError(t, "failed to write empty slice", err)
	err = writer.Close()
	failOnError(t, "failed to close", err)
	compressed := w.Bytes()
	t.Logf("compressed buffer: 0x%x", compressed)
	// Now recheck that if we decompress, we get empty slice
	r := bytes.NewBuffer(compressed)
	reader := NewReader(r)
	decompressed, err := ioutil.ReadAll(reader)
	failOnError(t, "failed to read", err)
	err = reader.Close()
	failOnError(t, "failed to close", err)
	if string(decompressed) != "" {
		t.Fatalf("Expected empty slice as decompressed, got %v instead", decompressed)
	}
}

func BenchmarkStreamCompression(b *testing.B) {
	if raw == nil {
		b.Fatal(ErrNoPayloadEnv)
	}
	var intermediate bytes.Buffer
	w := NewWriter(&intermediate)
	defer w.Close()
	b.SetBytes(int64(len(raw)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := w.Write(raw)
		if err != nil {
			b.Fatalf("Failed writing to compress object: %s", err)
		}
		// Prevent from unbound buffer growth.
		intermediate.Reset()
	}
}

func BenchmarkStreamDecompression(b *testing.B) {
	if raw == nil {
		b.Fatal(ErrNoPayloadEnv)
	}
	compressed, err := Compress(nil, raw)
	if err != nil {
		b.Fatalf("Failed to compress: %s", err)
	}
	_, err = Decompress(nil, compressed)
	if err != nil {
		b.Fatalf("Problem: %s", err)
	}

	dst := make([]byte, len(raw))
	b.SetBytes(int64(len(raw)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := bytes.NewReader(compressed)
		r := NewReader(rr)
		_, err := r.Read(dst)
		if err != nil {
			b.Fatalf("Failed to decompress: %s", err)
		}
	}
}

type breakingReader struct {
}

func (r *breakingReader) Read(p []byte) (int, error) {
	return len(p) - 1, io.ErrUnexpectedEOF
}

func TestUnexpectedEOFHandling(t *testing.T) {
	r := NewReader(&breakingReader{})
	_, err := r.Read(make([]byte, 1024))
	if err == nil {
		t.Error("Underlying error was handled silently")
	}
}
