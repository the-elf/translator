package util

import (
	"strings"
	"testing"
)

func TestRequireEnv(t *testing.T) {
	t.Run("returns value when set", func(t *testing.T) {
		t.Setenv("REQUIRE_ENV_TEST_KEY", "value")

		got, err := RequireEnv("REQUIRE_ENV_TEST_KEY")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "value" {
			t.Errorf("got %q, want %q", got, "value")
		}
	})

	t.Run("errors when unset", func(t *testing.T) {
		_, err := RequireEnv("REQUIRE_ENV_TEST_KEY_UNSET")
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}

func TestRequireEnvs(t *testing.T) {
	t.Run("returns all values when set", func(t *testing.T) {
		t.Setenv("REQUIRE_ENVS_TEST_A", "a")
		t.Setenv("REQUIRE_ENVS_TEST_B", "b")

		got, err := RequireEnvs("REQUIRE_ENVS_TEST_A", "REQUIRE_ENVS_TEST_B")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got["REQUIRE_ENVS_TEST_A"] != "a" || got["REQUIRE_ENVS_TEST_B"] != "b" {
			t.Errorf("got %v, want a=a b=b", got)
		}
	})

	t.Run("joins errors for every missing key", func(t *testing.T) {
		_, err := RequireEnvs("REQUIRE_ENVS_TEST_MISSING_1", "REQUIRE_ENVS_TEST_MISSING_2")
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		got := err.Error()
		if !strings.Contains(got, "REQUIRE_ENVS_TEST_MISSING_1") || !strings.Contains(got, "REQUIRE_ENVS_TEST_MISSING_2") {
			t.Errorf("error %q does not mention both missing keys", got)
		}
	})
}
