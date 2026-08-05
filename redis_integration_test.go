//go:build integration

package redis_test

import (
	"context"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	spicecache "github.com/spice-framework/spice/cache"
	redisstarter "github.com/spice-framework/starter-redis"
)

type integrationValue struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

func TestRedisClientIntegration(t *testing.T) {
	connectionURL := os.Getenv("SPICE_REDIS_TEST_URL")
	if connectionURL == "" {
		t.Fatal("SPICE_REDIS_TEST_URL is required for integration tests")
	}
	client, cleanup, err := redisstarter.Open(redisstarter.Options{
		URL:           connectionURL,
		PoolSize:      4,
		AllowInsecure: strings.HasPrefix(connectionURL, "redis://"),
	})
	if err != nil {
		t.Fatalf("open Redis client: %v", err)
	}
	t.Cleanup(func() {
		cleanupContext, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()
		if cleanupErr := cleanup(cleanupContext); cleanupErr != nil {
			t.Errorf("cleanup Redis client: %v", cleanupErr)
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Ping(ctx); err != nil {
		t.Fatalf("ping Redis client: %v", err)
	}

	var observations []spicecache.Observation
	store, err := redisstarter.NewJSONStore[integrationValue](
		client,
		redisstarter.StoreOptions{
			Definition: spicecache.Definition{
				ID:     "integration.orders",
				Module: "example.com/integration/orders",
			},
			Prefix:        "spice-integration-orders",
			MaxValueBytes: 1024,
		},
		func(_ context.Context, observation spicecache.Observation) {
			observations = append(observations, observation)
		},
	)
	if err != nil {
		t.Fatalf("construct Redis cache: %v", err)
	}
	want := integrationValue{ID: "41", Quantity: 3}
	if err := store.Put(ctx, "41", want, 100*time.Millisecond); err != nil {
		t.Fatalf("put Redis cache entry: %v", err)
	}
	got, found, err := store.Get(ctx, "41")
	if err != nil {
		t.Fatalf("get Redis cache entry: %v", err)
	}
	if !found || got != want {
		t.Fatalf("cache Get() = %#v, %t; want %#v, true", got, found, want)
	}
	expirationDeadline := time.Now().Add(5 * time.Second)
	for found && time.Now().Before(expirationDeadline) {
		time.Sleep(25 * time.Millisecond)
		_, found, err = store.Get(ctx, "41")
		if err != nil {
			t.Fatalf("get expiring Redis cache entry: %v", err)
		}
	}
	if found {
		t.Fatal("Redis cache entry did not expire")
	}
	if err := store.Put(ctx, "42", want, 0); err != nil {
		t.Fatalf("put persistent Redis cache entry: %v", err)
	}
	if err := store.Delete(ctx, "42"); err != nil {
		t.Fatalf("delete Redis cache entry: %v", err)
	}
	if snapshot := store.Snapshot(); snapshot.Hits == 0 ||
		snapshot.Misses == 0 ||
		snapshot.Puts != 2 ||
		snapshot.Deletes != 1 {
		t.Fatalf("Redis cache snapshot = %#v", snapshot)
	}
	if len(observations) < 5 ||
		observations[0].Definition.Module !=
			"example.com/integration/orders" {
		t.Fatalf("Redis cache observations = %#v", observations)
	}

	canceled, cancelCanceled := context.WithCancel(context.Background())
	cancelCanceled()
	if _, _, err := store.Get(canceled, "41"); err == nil {
		t.Fatal("canceled Redis cache Get() unexpectedly succeeded")
	}
}

func TestRedisPoolsRemainIndependentAndFailuresRedactCredentials(t *testing.T) {
	connectionURL := os.Getenv("SPICE_REDIS_TEST_URL")
	if connectionURL == "" {
		t.Fatal("SPICE_REDIS_TEST_URL is required for integration tests")
	}
	options := redisstarter.Options{
		URL:           connectionURL,
		PoolSize:      2,
		AllowInsecure: strings.HasPrefix(connectionURL, "redis://"),
	}
	first, firstCleanup, err := redisstarter.Open(options)
	if err != nil {
		t.Fatalf("open first Redis client: %v", err)
	}
	second, secondCleanup, err := redisstarter.Open(options)
	if err != nil {
		if cleanupErr := firstCleanup(context.Background()); cleanupErr != nil {
			t.Fatalf("open second Redis client: %v; cleanup first client: %v", err, cleanupErr)
		}
		t.Fatalf("open second Redis client: %v", err)
	}
	t.Cleanup(func() {
		if cleanupErr := firstCleanup(context.Background()); cleanupErr != nil {
			t.Errorf("cleanup first Redis client: %v", cleanupErr)
		}
		if cleanupErr := secondCleanup(context.Background()); cleanupErr != nil {
			t.Errorf("cleanup second Redis client: %v", cleanupErr)
		}
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := first.Ping(ctx); err != nil {
		t.Fatalf("ping first Redis client: %v", err)
	}
	if err := second.Ping(ctx); err != nil {
		t.Fatalf("ping second Redis client: %v", err)
	}
	if err := firstCleanup(ctx); err != nil {
		t.Fatalf("cleanup first Redis client: %v", err)
	}
	if err := second.Ping(ctx); err != nil {
		t.Fatalf("second pool failed after first cleanup: %v", err)
	}

	parsed, err := url.Parse(connectionURL)
	if err != nil {
		t.Fatalf("parse integration URL: %v", err)
	}
	const rejectedSecret = "wrong-secret-value"
	parsed.User = url.UserPassword("default", rejectedSecret)
	rejected, rejectedCleanup, err := redisstarter.Open(redisstarter.Options{
		URL:           parsed.String(),
		PoolSize:      1,
		AllowInsecure: strings.HasPrefix(connectionURL, "redis://"),
	})
	if err != nil {
		t.Fatalf("open rejected Redis client: %v", err)
	}
	defer func() {
		if cleanupErr := rejectedCleanup(context.Background()); cleanupErr != nil {
			t.Errorf("cleanup rejected Redis client: %v", cleanupErr)
		}
	}()
	if pingErr := rejected.Ping(ctx); pingErr == nil {
		t.Fatal("ping with rejected credential unexpectedly succeeded")
	} else if strings.Contains(pingErr.Error(), rejectedSecret) ||
		strings.Contains(pingErr.Error(), connectionURL) {
		t.Fatalf("authentication failure exposed credentials: %v", pingErr)
	}
}
