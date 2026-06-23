package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func rsId() schema.StringAttribute {
	return schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "The unique ID of this resource.",
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

func orgUnitAPIPath(id string) string {
	return "id:" + strings.TrimPrefix(id, "id:")
}

func importSplitId(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse, boolAttr, idAttr string) {
	parts := strings.SplitN(req.ID, ",", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("Expected format: %s,%s. Got: %q", boolAttr, idAttr, req.ID))
		return
	}
	boolVal, err := strconv.ParseBool(parts[0])
	if err != nil {
		resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("Cannot parse %q as bool: %v", parts[0], err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(boolAttr), boolVal)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(idAttr), parts[1])...)
}
