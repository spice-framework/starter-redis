package redis

import (
	"slices"
	"testing"

	spicestarter "github.com/spice-framework/spice/annotation/sdk/starter"
)

func TestManifestDeclaresIndependentRedisStarter(t *testing.T) {
	t.Parallel()
	manifest := Manifest()
	spec := manifest.Spec()
	if spec.Schema != spicestarter.Schema ||
		spec.ID != "github.com/spice-framework/starter-redis" ||
		spec.Module != "github.com/spice-framework/starter-redis" ||
		spec.SpiceAPI != spicestarter.APIVersion ||
		spec.Review != "docs/dependency-review.md" ||
		len(spec.Activation.EntryPoints) != 1 ||
		spec.Activation.EntryPoints[0].Package != "github.com/spice-framework/starter-redis" ||
		spec.Activation.EntryPoints[0].Symbol != "Open" {
		t.Fatalf("Manifest().Spec() = %#v", spec)
	}
	if !slices.Equal(spec.Capabilities, []string{"cache.redis", "data.redis"}) ||
		len(spec.Dependencies) != 1 ||
		spec.Dependencies[0] != (spicestarter.Dependency{
			Module:  "github.com/redis/go-redis/v9",
			Version: "v9.21.0",
			License: "BSD-2-Clause",
		}) {
		t.Fatalf("capabilities=%v dependencies=%#v", spec.Capabilities, spec.Dependencies)
	}
	if err := manifest.Compatible(spicestarter.APIVersion, "go1.26.5"); err != nil {
		t.Fatalf("Compatible() error = %v", err)
	}
	content, err := manifest.JSON()
	if err != nil {
		t.Fatalf("JSON() error = %v", err)
	}
	parsed, err := spicestarter.Parse(content)
	if err != nil {
		t.Fatalf("Parse(JSON()) error = %v", err)
	}
	if parsed.Spec().ID != spec.ID {
		t.Fatalf("parsed ID = %q, want %q", parsed.Spec().ID, spec.ID)
	}
}
