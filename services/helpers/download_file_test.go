package helpers

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	sdk "github.com/Smartling/api-sdk-go"
	sdkfile "github.com/Smartling/api-sdk-go/helpers/sm_file"
)

type stubDownloadClient struct {
	sdk.APIClient
	body   io.Reader
	stream func(ctx context.Context) io.Reader
}

func (c stubDownloadClient) DownloadTranslation(ctx context.Context, _, _ string, _ sdk.FileDownloadRequest) (io.ReadCloser, error) {
	if c.stream != nil {
		return io.NopCloser(c.stream(ctx)), nil
	}
	return io.NopCloser(c.body), nil
}

type stalledReader struct {
	ctx  context.Context
	sent bool
}

func (r *stalledReader) Read(p []byte) (int, error) {
	if !r.sent {
		r.sent = true
		return copy(p, "partial"), nil
	}
	<-r.ctx.Done()
	return 0, r.ctx.Err()
}

type trickleReader struct {
	chunks int
	delay  time.Duration
}

func (r *trickleReader) Read(p []byte) (int, error) {
	if r.chunks == 0 {
		return 0, io.EOF
	}
	time.Sleep(r.delay)
	r.chunks--
	return copy(p, "x"), nil
}

func setIdleTimeout(t *testing.T, d time.Duration) {
	t.Helper()
	prev := downloadIdleTimeout
	downloadIdleTimeout = d
	t.Cleanup(func() { downloadIdleTimeout = prev })
}

func downloadStream(t *testing.T, stream func(ctx context.Context) io.Reader, path string) error {
	t.Helper()
	return DownloadFile(
		context.Background(),
		stubDownloadClient{stream: stream},
		"test-project",
		sdkfile.File{FileURI: "a.json"},
		"fr-FR",
		path,
		sdk.RetrievePublished,
	)
}

func TestDownloadFile_SymlinkKeepsLinkAndUpdatesTarget(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks require elevated privileges on Windows")
	}
	targetDir := t.TempDir()
	target := filepath.Join(targetDir, "real.json")
	if err := os.WriteFile(target, []byte("old"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	linkDir := t.TempDir()
	link := filepath.Join(linkDir, "a_fr-FR.json")
	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := download(t, strings.NewReader("new"), link); err != nil {
		t.Fatalf("DownloadFile error: %v", err)
	}

	info, err := os.Lstat(link)
	if err != nil {
		t.Fatalf("lstat link: %v", err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Errorf("symlink at %s was replaced by a regular file", link)
	}
	body, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read target: %v", err)
	}
	if string(body) != "new" {
		t.Errorf("target content = %q, want %q", body, "new")
	}
	assertNoTempFiles(t, targetDir)
	assertNoTempFiles(t, linkDir)
}

func TestDownloadFile_StalledBodyTimesOut(t *testing.T) {
	setIdleTimeout(t, 50*time.Millisecond)
	dir := t.TempDir()
	path := filepath.Join(dir, "a_fr-FR.json")

	done := make(chan error, 1)
	go func() {
		done <- downloadStream(t, func(ctx context.Context) io.Reader {
			return &stalledReader{ctx: ctx}
		}, path)
	}()

	select {
	case err := <-done:
		if !errors.Is(err, errDownloadStalled) {
			t.Fatalf("err = %v, want errDownloadStalled", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("stalled download was not aborted")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("partial file created at %s (stat err: %v)", path, err)
	}
	assertNoTempFiles(t, dir)
}

func TestDownloadFile_SteadySlowBodyDoesNotTimeOut(t *testing.T) {
	const (
		timeout = 200 * time.Millisecond
		delay   = 20 * time.Millisecond
		chunks  = 20
	)
	setIdleTimeout(t, timeout)
	path := filepath.Join(t.TempDir(), "a_fr-FR.json")

	start := time.Now()
	err := downloadStream(t, func(context.Context) io.Reader {
		return &trickleReader{chunks: chunks, delay: delay}
	}, path)
	if err != nil {
		t.Fatalf("DownloadFile error: %v", err)
	}
	if elapsed := time.Since(start); elapsed <= timeout {
		t.Fatalf("download took %s, must exceed idle timeout %s to prove the timer resets", elapsed, timeout)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read result: %v", err)
	}
	if want := strings.Repeat("x", chunks); string(body) != want {
		t.Errorf("content = %q, want %q", body, want)
	}
}

type failingReader struct {
	data string
	done bool
}

func (r *failingReader) Read(p []byte) (int, error) {
	if r.done {
		return 0, errors.New("connection reset")
	}
	r.done = true
	return copy(p, r.data), nil
}

func download(t *testing.T, body io.Reader, path string) error {
	t.Helper()
	return DownloadFile(
		context.Background(),
		stubDownloadClient{body: body},
		"test-project",
		sdkfile.File{FileURI: "a.json"},
		"fr-FR",
		path,
		sdk.RetrievePublished,
	)
}

func assertNoTempFiles(t *testing.T, dir string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".tmp") {
			t.Errorf("temp file left behind: %s", e.Name())
		}
	}
}

func TestDownloadFile_WritesContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "a_fr-FR.json")

	if err := download(t, strings.NewReader("translated"), path); err != nil {
		t.Fatalf("DownloadFile error: %v", err)
	}

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read result: %v", err)
	}
	if string(body) != "translated" {
		t.Errorf("content = %q, want %q", body, "translated")
	}
	if runtime.GOOS != "windows" {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("stat result: %v", err)
		}
		if got := info.Mode().Perm(); got != 0o644 {
			t.Errorf("mode = %o, want 644", got)
		}
	}
	assertNoTempFiles(t, filepath.Dir(path))
}

func TestDownloadFile_ReplacesExistingKeepingMode(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a_fr-FR.json")
	if err := os.WriteFile(path, []byte("old"), 0o600); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := download(t, strings.NewReader("new"), path); err != nil {
		t.Fatalf("DownloadFile error: %v", err)
	}

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read result: %v", err)
	}
	if string(body) != "new" {
		t.Errorf("content = %q, want %q", body, "new")
	}
	if runtime.GOOS != "windows" {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("stat result: %v", err)
		}
		if got := info.Mode().Perm(); got != 0o600 {
			t.Errorf("mode = %o, want 600 (preserved)", got)
		}
	}
	assertNoTempFiles(t, dir)
}

func TestDownloadFile_InterruptedDownloadLeavesNoPartialFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a_fr-FR.json")

	err := download(t, &failingReader{data: "partial"}, path)
	if err == nil {
		t.Fatal("expected error from interrupted download, got nil")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("partial file created at %s (stat err: %v)", path, err)
	}
	assertNoTempFiles(t, dir)
}

func TestDownloadFile_InterruptedDownloadKeepsExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a_fr-FR.json")
	if err := os.WriteFile(path, []byte("previous complete"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := download(t, &failingReader{data: "partial"}, path); err == nil {
		t.Fatal("expected error from interrupted download, got nil")
	}

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read existing: %v", err)
	}
	if string(body) != "previous complete" {
		t.Errorf("existing file changed to %q", body)
	}
	assertNoTempFiles(t, dir)
}
