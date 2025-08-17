package components

import (
	"context"
	"io"
	"testing"
)

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
