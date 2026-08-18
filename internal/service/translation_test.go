package service

import (
	"errors"
	"testing"
	"translator/internal/model"
)

func TestResolveTranslation(t *testing.T) {
	svc := &TranslationService{}
	template := model.MessageTemplate{}

	t.Run("not georgian text returns sentinel error", func(t *testing.T) {
		_, err := svc.resolveTranslation("hello", template, "NOT_GEORGIAN_TEXT")
		if !errors.Is(err, ErrNotGeorgianText) {
			t.Fatalf("got err %v, want ErrNotGeorgianText", err)
		}
	})

	t.Run("passes through translation and normalizes em-dash", func(t *testing.T) {
		got, err := svc.resolveTranslation("text", template, "hello — world")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "hello - world" {
			t.Errorf("got %q, want %q", got, "hello - world")
		}
	})
}
