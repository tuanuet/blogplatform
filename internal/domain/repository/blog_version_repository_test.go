package repository_test

import (
	"github.com/aiagent/internal/domain/repository"
	"testing"
)

func TestBlogVersionRepository_Interface(t *testing.T) {
	// Just verify the type exists and can be referenced
	var _ repository.BlogVersionRepository = nil
}
