package components

import (
	"context"
	"io"
	"strings"
	"testing"
)

func TestNestedSlotsWithLazyInit(t *testing.T) {
	// Test that slots work without middleware/Root - just a plain context
	ctx := context.Background()
	var buf strings.Builder

	err := NestedSlotTest().Render(ctx, &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	output := buf.String()
	t.Log("Output:", output)

	// Verify outer slot value is present
	if !strings.Contains(output, "Outer: hello") {
		t.Error("expected 'Outer: hello' in output")
	}

	// Verify inner slot value is present
	if !strings.Contains(output, "Inner: 42") {
		t.Error("expected 'Inner: 42' in output")
	}

	// Verify nested access works
	if !strings.Contains(output, "Outer says: hello") {
		t.Error("expected 'Outer says: hello' in nested content")
	}
	if !strings.Contains(output, "Inner says: 42") {
		t.Error("expected 'Inner says: 42' in nested content")
	}
}

func TestSiblingSlotsWithLazyInit(t *testing.T) {
	ctx := context.Background()
	var buf strings.Builder

	err := SiblingSlotTest().Render(ctx, &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	output := buf.String()
	t.Log("Output:", output)

	// Verify both slots render their children
	if !strings.Contains(output, "<h1>") || !strings.Contains(output, "<p>") {
		t.Error("expected DualSlot structure with h1 and p tags")
	}
}

func TestNotification(t *testing.T) {
	// Pipe the rendered template into goquery.
	r, w := io.Pipe()
	go func() {
		_ = Notification("Test", "a test notification").Render(context.Background(), w)
		_ = w.Close()
	}()
	b, _ := io.ReadAll(r)
	t.Log(string(b))
	//t.Fail()
	// doc, err := goquery.NewDocumentFromReader(r)
	// if err != nil {
	// 	t.Fatalf("failed to read template: %v", err)
	// }
	// // Expect the component to be present.
	// if doc.Find(`[data-testid="headerTemplate"]`).Length() == 0 {
	// 	t.Error("expected data-testid attribute to be rendered, but it wasn't")
	// }
	// // Expect the page name to be set correctly.
	// expectedPageName := "Posts"
	// if actualPageName := doc.Find("h1").Text(); actualPageName != expectedPageName {
	// 	t.Errorf("expected page name %q, got %q", expectedPageName, actualPageName)
	// }
}
