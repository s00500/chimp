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
		_ = Notification(NotificationSuccess, "a test notification").Render(context.Background(), w)
		_ = w.Close()
	}()
	b, _ := io.ReadAll(r)
	t.Log(string(b))
}

func TestFormInput(t *testing.T) {
	ctx := context.Background()
	var buf strings.Builder

	err := FormInput("Email", "email",
		WithType("email"),
		WithModel("form.email"),
		WithPlaceholder("you@example.com"),
		WithRequired(),
		WithError("$errors.email"),
	).Render(ctx, &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	output := buf.String()
	t.Log("Output:", output)

	// Verify label
	if !strings.Contains(output, `<label class="label"`) {
		t.Error("expected label element")
	}
	if !strings.Contains(output, "Email") {
		t.Error("expected label text 'Email'")
	}

	// Verify input attributes
	if !strings.Contains(output, `type="email"`) {
		t.Error("expected type='email'")
	}
	if !strings.Contains(output, `data-model="form.email"`) {
		t.Error("expected data-model attribute")
	}
	if !strings.Contains(output, `placeholder="you@example.com"`) {
		t.Error("expected placeholder")
	}
	if !strings.Contains(output, `required`) {
		t.Error("expected required attribute")
	}

	// Verify error element
	if !strings.Contains(output, `data-show="$errors.email"`) {
		t.Error("expected error element with data-show")
	}
}

func TestFormSelect(t *testing.T) {
	ctx := context.Background()
	var buf strings.Builder

	err := FormSelect("Country", "country",
		WithModel("form.country"),
		WithOptions([]SelectOption{
			{Value: "us", Label: "United States"},
			{Value: "de", Label: "Germany"},
		}),
		WithEmptyOption("Select a country"),
	).Render(ctx, &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	output := buf.String()
	t.Log("Output:", output)

	// Verify select element
	if !strings.Contains(output, `<select`) {
		t.Error("expected select element")
	}

	// Verify options
	if !strings.Contains(output, "United States") {
		t.Error("expected 'United States' option")
	}
	if !strings.Contains(output, "Germany") {
		t.Error("expected 'Germany' option")
	}
	if !strings.Contains(output, "Select a country") {
		t.Error("expected empty option placeholder")
	}
}

func TestFormCheckbox(t *testing.T) {
	ctx := context.Background()
	var buf strings.Builder

	err := FormCheckbox("Accept terms", "terms",
		WithModel("form.acceptTerms"),
	).Render(ctx, &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	output := buf.String()
	t.Log("Output:", output)

	if !strings.Contains(output, `type="checkbox"`) {
		t.Error("expected checkbox input")
	}
	if !strings.Contains(output, "Accept terms") {
		t.Error("expected label text")
	}
	if !strings.Contains(output, `data-model="form.acceptTerms"`) {
		t.Error("expected data-model attribute")
	}
}

func TestButton(t *testing.T) {
	ctx := context.Background()
	var buf strings.Builder

	err := Button("Save",
		WithVariant(VariantPrimary),
		WithOn("click", "@post('/save')"),
	).Render(ctx, &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	output := buf.String()
	t.Log("Output:", output)

	if !strings.Contains(output, "Save") {
		t.Error("expected button text 'Save'")
	}
	if !strings.Contains(output, "btn-primary") {
		t.Error("expected btn-primary class")
	}
	if !strings.Contains(output, `data-on:click="@post(&#39;/save&#39;)"`) && !strings.Contains(output, `data-on:click="@post('/save')"`) {
		t.Error("expected data-on:click attribute")
	}
}

func TestStack(t *testing.T) {
	ctx := context.Background()
	var buf strings.Builder

	err := Stack(WithGap("lg")).Render(ctx, &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	output := buf.String()
	t.Log("Output:", output)

	if !strings.Contains(output, "flex") {
		t.Error("expected flex class")
	}
	if !strings.Contains(output, "gap-6") {
		t.Error("expected gap-6 class for lg gap")
	}
}

func TestRow(t *testing.T) {
	ctx := context.Background()
	var buf strings.Builder

	err := Row(WithGap("sm")).Render(ctx, &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	output := buf.String()
	t.Log("Output:", output)

	if !strings.Contains(output, "flex-row") {
		t.Error("expected flex-row class")
	}
	if !strings.Contains(output, "gap-2") {
		t.Error("expected gap-2 class for sm gap")
	}
}

func TestCard(t *testing.T) {
	ctx := context.Background()
	var buf strings.Builder

	err := Card(WithCardTitle("Test Card"), WithPadding("lg")).Render(ctx, &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	output := buf.String()
	t.Log("Output:", output)

	if !strings.Contains(output, "card") {
		t.Error("expected card class")
	}
	if !strings.Contains(output, "Test Card") {
		t.Error("expected card title")
	}
	if !strings.Contains(output, "p-6") {
		t.Error("expected p-6 class for lg padding")
	}
}

func TestSection(t *testing.T) {
	ctx := context.Background()
	var buf strings.Builder

	err := Section(WithSectionTitle("Test Section")).Render(ctx, &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	output := buf.String()
	t.Log("Output:", output)

	if !strings.Contains(output, "rounded-lg border") {
		t.Error("expected section border classes")
	}
	if !strings.Contains(output, "Test Section") {
		t.Error("expected section title")
	}
	if !strings.Contains(output, "border-b") {
		t.Error("expected header border-b class")
	}
}
