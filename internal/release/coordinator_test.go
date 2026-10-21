package release

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"openresty-admin/internal/model"
	"openresty-admin/internal/nginx"
)

// fakeHost is an in-memory Host adapter for testing the coordinator.
type fakeHost struct {
	files map[string]string

	// Scripted Test results; each call pops the front. Falls back to testOK.
	testResults []bool
	testOK      bool
	testOut     string

	reloadOK  bool
	reloadOut string

	// writeCount is bumped on every WriteConfig; failWriteOn makes the Nth
	// write fail (used to simulate a failed restore).
	writeCount  int32
	failWriteOn int32

	readErr   error
	deleteErr error

	// When non-nil, Test blocks until the channel is closed.
	blockTest chan struct{}
	testCalls int32
}

func newFakeHost() *fakeHost {
	return &fakeHost{files: map[string]string{}, testOK: true, reloadOK: true}
}

func (f *fakeHost) ReadForRelease(rel string) (string, bool, error) {
	if f.readErr != nil {
		return "", false, f.readErr
	}
	c, ok := f.files[rel]
	return c, ok, nil
}

func (f *fakeHost) WriteConfig(rel, content string) error {
	n := atomic.AddInt32(&f.writeCount, 1)
	if f.failWriteOn > 0 && n == f.failWriteOn {
		return errors.New("injected write failure")
	}
	f.files[rel] = content
	return nil
}

func (f *fakeHost) DeleteConfig(rel string) error {
	if f.deleteErr != nil {
		return f.deleteErr
	}
	delete(f.files, rel)
	return nil
}

func (f *fakeHost) Test() (bool, string) {
	atomic.AddInt32(&f.testCalls, 1)
	if f.blockTest != nil {
		<-f.blockTest
	}
	if len(f.testResults) > 0 {
		r := f.testResults[0]
		f.testResults = f.testResults[1:]
		return r, f.testOut
	}
	return f.testOK, f.testOut
}

func (f *fakeHost) Reload() (bool, string) {
	return f.reloadOK, f.reloadOut
}

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&model.ConfigRelease{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	t.Cleanup(func() { sqlDB.Close() })
	return db
}

// newCoord builds a coordinator whose health probe always passes by default.
func newCoord(t *testing.T, h Host) *Coordinator {
	t.Helper()
	c := New(h, newTestDB(t))
	c.probe = func([]ReleaseFile) bool { return true }
	return c
}

func TestPublishSuccess(t *testing.T) {
	h := newFakeHost()
	h.files["conf.d/site.conf"] = "old"
	c := newCoord(t, h)

	rel, err := c.Publish([]nginx.Change{{Path: "conf.d/site.conf", Content: "new"}}, "test")
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if rel.Status != "succeeded" {
		t.Fatalf("status = %q, want succeeded", rel.Status)
	}
	if rel.ID == 0 {
		t.Fatal("release should be persisted with a non-zero ID")
	}
	if got := h.files["conf.d/site.conf"]; got != "new" {
		t.Fatalf("file = %q, want new", got)
	}
}

func TestPublishCreateAction(t *testing.T) {
	h := newFakeHost()
	c := newCoord(t, h)

	rel, err := c.Publish([]nginx.Change{{Path: "conf.d/brand-new.conf", Content: "server {}"}}, "test")
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if rel.Status != "succeeded" {
		t.Fatalf("status = %q, want succeeded", rel.Status)
	}
	if _, ok := h.files["conf.d/brand-new.conf"]; !ok {
		t.Fatal("expected file to be created")
	}
}

func TestPublishDeleteAction(t *testing.T) {
	h := newFakeHost()
	h.files["conf.d/site.conf"] = "old"
	c := newCoord(t, h)

	rel, err := c.Publish([]nginx.Change{{Path: "conf.d/site.conf", Delete: true}}, "test")
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if rel.Status != "succeeded" {
		t.Fatalf("status = %q, want succeeded", rel.Status)
	}
	if _, ok := h.files["conf.d/site.conf"]; ok {
		t.Fatal("expected file to be deleted")
	}
}

func TestPublishValidationFailureRollsBack(t *testing.T) {
	h := newFakeHost()
	h.files["conf.d/site.conf"] = "old"
	// First Test (validation) fails; retest after restore succeeds.
	h.testResults = []bool{false, true}
	h.testOut = "nginx: [emerg] unknown directive"
	c := newCoord(t, h)

	rel, err := c.Publish([]nginx.Change{{Path: "conf.d/site.conf", Content: "bad"}}, "test")
	if err == nil {
		t.Fatal("expected error on validation failure")
	}
	if rel.Status != "validation_failed" {
		t.Fatalf("status = %q, want validation_failed", rel.Status)
	}
	if got := h.files["conf.d/site.conf"]; got != "old" {
		t.Fatalf("file should be restored to old, got %q", got)
	}
}

func TestPublishReloadFailureRollsBack(t *testing.T) {
	h := newFakeHost()
	h.files["conf.d/site.conf"] = "old"
	h.testOK = true
	h.reloadOK = false
	h.reloadOut = "nginx: [error] open() failed"
	c := newCoord(t, h)

	rel, err := c.Publish([]nginx.Change{{Path: "conf.d/site.conf", Content: "new"}}, "test")
	if err == nil {
		t.Fatal("expected error on reload failure")
	}
	if rel.Status != "reload_failed" {
		t.Fatalf("status = %q, want reload_failed", rel.Status)
	}
	if got := h.files["conf.d/site.conf"]; got != "old" {
		t.Fatalf("file should be restored to old, got %q", got)
	}
}

func TestPublishHealthCheckFailureRollsBack(t *testing.T) {
	h := newFakeHost()
	h.files["conf.d/site.conf"] = "old"
	c := newCoord(t, h)
	c.probe = func([]ReleaseFile) bool { return false } // health check fails

	rel, err := c.Publish([]nginx.Change{{Path: "conf.d/site.conf", Content: "new"}}, "test")
	if err == nil {
		t.Fatal("expected error on health check failure")
	}
	if rel.Status != "rolled_back" {
		t.Fatalf("status = %q, want rolled_back", rel.Status)
	}
	if got := h.files["conf.d/site.conf"]; got != "old" {
		t.Fatalf("file should be restored to old, got %q", got)
	}
}

func TestPublishRecoveryRequiredWhenRestoreFails(t *testing.T) {
	h := newFakeHost()
	h.files["conf.d/site.conf"] = "old"
	h.testOK = false // validation fails
	h.testOut = "nginx: [emerg] syntax error"
	// Write 1 = the change; write 2 = the restore, which we force to fail.
	h.failWriteOn = 2
	c := newCoord(t, h)

	rel, err := c.Publish([]nginx.Change{{Path: "conf.d/site.conf", Content: "bad"}}, "test")
	if err == nil {
		t.Fatal("expected error on recovery")
	}
	if rel.Status != "recovery_required" {
		t.Fatalf("status = %q, want recovery_required", rel.Status)
	}
	if rel.Error == nil || *rel.Error == "" {
		t.Fatal("expected error message on recovery_required release")
	}
}

func TestIsLockedDuringPublish(t *testing.T) {
	h := newFakeHost()
	h.blockTest = make(chan struct{})
	c := newCoord(t, h)

	done := make(chan struct{})
	go func() {
		c.Publish([]nginx.Change{{Path: "x.conf", Content: "y"}}, "test")
		close(done)
	}()

	// Wait until the goroutine has acquired the lock and reached Test().
	deadline := time.Now().Add(2 * time.Second)
	for atomic.LoadInt32(&h.testCalls) == 0 {
		if time.Now().After(deadline) {
			t.Fatal("publish never reached Test")
		}
		time.Sleep(time.Millisecond)
	}

	if !c.IsLocked() {
		t.Fatal("expected IsLocked == true while a publish is in flight")
	}

	close(h.blockTest)
	<-done

	if c.IsLocked() {
		t.Fatal("expected IsLocked == false after publish completes")
	}
}
