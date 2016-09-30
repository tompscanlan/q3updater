package q3updater

import (
	"testing"
)

var approvalTests = []struct {
	in    Approval
	valid bool
}{
	{
		Approval{Id: 99, TeamId: 99, Blob: 99, Description: "for testing", Approved: false},
		true,
	},
	{
		Approval{Id: 99, TeamId: 99, Blob: 99, Description: "for testing new!", Approved: false},
		true,
	},
}

func CompareApproval(a *Approval, b *Approval, t *testing.T) {

	if a.Id != b.Id {
		t.Errorf("expected %d, got %d", a.Id, b.Id)
	}

	if a.TeamId != b.TeamId {
		t.Errorf("expected %d, got %d", a.TeamId, b.TeamId)
	}

	if a.Blob != b.Blob {
		t.Errorf("expected %d, got %d", a.Blob, b.Blob)
	}
	if a.Description != b.Description {
		t.Errorf("expected %s, got %s", a.Description, b.Description)
	}
	if a.Approved != b.Approved {
		t.Errorf("expected %t, got %t", a.Approved, b.Approved)
	}
}
