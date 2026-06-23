package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func TestUserStateUpgradeV0(t *testing.T) {
	ctx := context.Background()
	r := &userResource{}
	upgraders := r.UpgradeState(ctx)

	upgrader, ok := upgraders[0]
	if !ok {
		t.Fatal("expected state upgrader for version 0")
	}

	oldState := []byte(`{
		"id": "user-old-123",
		"primary_email": "test@way.cloud",
		"org_unit_path": "/Engineering",
		"suspended": false,
		"archived": false,
		"name": [{"given_name": "Test", "family_name": "User"}]
	}`)

	req := resource.UpgradeStateRequest{
		RawState: &tfprotov6.RawState{JSON: oldState},
	}
	resp := &resource.UpgradeStateResponse{
		State: tfsdk.State{Schema: testUserSchema(ctx)},
	}

	upgrader.StateUpgrader(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("state upgrade failed: %s", resp.Diagnostics.Errors())
	}

	var state userResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to read upgraded state: %s", resp.Diagnostics.Errors())
	}

	if state.Id.ValueString() != "user-old-123" {
		t.Errorf("expected id=user-old-123, got %s", state.Id.ValueString())
	}
	if state.PrimaryEmail.ValueString() != "test@way.cloud" {
		t.Errorf("expected primary_email=test@way.cloud, got %s", state.PrimaryEmail.ValueString())
	}
	if state.Name == nil || state.Name.GivenName.ValueString() != "Test" {
		t.Errorf("expected name.given_name=Test, got %v", state.Name)
	}
	if state.Name.FamilyName.ValueString() != "User" {
		t.Errorf("expected name.family_name=User, got %s", state.Name.FamilyName.ValueString())
	}
	if state.Suspended.ValueBool() != false {
		t.Error("expected suspended=false")
	}
	if state.Archived.ValueBool() != false {
		t.Error("expected archived=false")
	}
}

func testUserSchema(ctx context.Context) fwschema.Schema {
	r := &userResource{}
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}
	r.Schema(ctx, req, resp)
	return resp.Schema
}
