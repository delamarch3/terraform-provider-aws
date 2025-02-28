// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package bedrockagent

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagent"
	awstypes "github.com/aws/aws-sdk-go-v2/service/bedrockagent/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
	"github.com/hashicorp/terraform-provider-aws/internal/framework/flex"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// Function annotations are used for resource registration to the Provider. DO NOT EDIT.
// @FrameworkResource("aws_bedrockagent_flow_version", name="Flow Version")
func newResourceFlowVersion(_ context.Context) (resource.ResourceWithConfigure, error) {
	r := &resourceFlowVersion{}

	r.SetDefaultCreateTimeout(5 * time.Minute)
	r.SetDefaultUpdateTimeout(5 * time.Minute)
	r.SetDefaultDeleteTimeout(5 * time.Minute)

	return r, nil
}

const (
	ResNameFlowVersion = "Flow Version"
)

type resourceFlowVersion struct {
	framework.ResourceWithConfigure
	framework.WithTimeouts
}


// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/schemas?page=schemas
func (r *resourceFlowVersion) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"arn": framework.ARNAttributeComputedOnly(),
			"version": framework.IDAttribute(),
			"name": framework.IDAttribute(),
			"flow_identifier": framework.IDAttribute(),
			"description": schema.StringAttribute{
				Optional: true,
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
			}),
		},
	}
}

func (r *resourceFlowVersion) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// TIP: -- 1. Get a client connection to the relevant service
	conn := r.Meta().BedrockAgentClient(ctx)

	// TIP: -- 2. Fetch the plan
	var plan resourceFlowVersionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TIP: -- 3. Populate a Create input structure
	var input bedrockagent.CreateFlowVersionInput
	resp.Diagnostics.Append(flex.Expand(ctx, plan, &input)...)
	if resp.Diagnostics.HasError() {
		return
	}

    input.ClientToken = aws.String(id.UniqueId())

	// TIP: -- 4. Call the AWS Create function
	out, err := conn.CreateFlowVersion(ctx, &input)
	if err != nil {
		// TIP: Since ID has not been set yet, you cannot use plan.ID.String()
		// in error messages at this point.
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.BedrockAgent, create.ErrActionCreating, ResNameFlowVersion, plan.Name.String(), err),
			err.Error(),
		)
		return
	}
	if out == nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.BedrockAgent, create.ErrActionCreating, ResNameFlowVersion, plan.Name.String(), nil),
			errors.New("empty output").Error(),
		)
		return
	}

	// TIP: -- 5. Using the output from the create function, set attributes
	resp.Diagnostics.Append(flex.Flatten(ctx, out, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TIP: -- 6. Use a waiter to wait for create to complete
	createTimeout := r.CreateTimeout(ctx, plan.Timeouts)
	_, err = waitFlowVersionCreated(ctx, conn, plan.FlowIdentifier.ValueString(), plan.Version.ValueString(), createTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.BedrockAgent, create.ErrActionWaitingForCreation, ResNameFlowVersion, plan.Name.String(), err),
			err.Error(),
		)
		return
	}

	// TIP: -- 7. Save the request plan to response state
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *resourceFlowVersion) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// TIP: -- 1. Get a client connection to the relevant service
	conn := r.Meta().BedrockAgentClient(ctx)

	// TIP: -- 2. Fetch the state
	var state resourceFlowVersionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TIP: -- 3. Get the resource from AWS using an API Get, List, or Describe-
	// type function, or, better yet, using a finder.
	out, err := findFlowVersionByID(ctx, conn, state.FlowIdentifier.ValueString(), state.Version.ValueString())
	// TIP: -- 4. Remove resource from state if it is not found
	if tfresource.NotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.BedrockAgent, create.ErrActionSetting, ResNameFlowVersion, state.FlowIdentifier.String(), err),
			err.Error(),
		)
		return
	}

	// TIP: -- 5. Set the arguments and attributes
	resp.Diagnostics.Append(flex.Flatten(ctx, out, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TIP: -- 6. Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *resourceFlowVersion) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	return
}

func (r *resourceFlowVersion) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// TIP: -- 1. Get a client connection to the relevant service
	conn := r.Meta().BedrockAgentClient(ctx)

	// TIP: -- 2. Fetch the state
	var state resourceFlowVersionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TIP: -- 3. Populate a delete input structure
	input := bedrockagent.DeleteFlowVersionInput{
        FlowIdentifier: state.FlowIdentifier.ValueStringPointer(),
        FlowVersion: state.Version.ValueStringPointer(),
        SkipResourceInUseCheck: false,
	}

	// TIP: -- 4. Call the AWS delete function
	_, err := conn.DeleteFlowVersion(ctx, &input)
	// TIP: On rare occassions, the API returns a not found error after deleting a
	// resource. If that happens, we don't want it to show up as an error.
	if err != nil {
		if errs.IsA[*awstypes.ResourceNotFoundException](err) {
			return
		}

		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.BedrockAgent, create.ErrActionDeleting, ResNameFlowVersion, state.Version.String(), err),
			err.Error(),
		)
		return
	}

	// TIP: -- 5. Use a waiter to wait for delete to complete
	deleteTimeout := r.DeleteTimeout(ctx, state.Timeouts)
	_, err = waitFlowVersionDeleted(ctx, conn, state.FlowIdentifier.ValueString(), state.Version.ValueString(), deleteTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.BedrockAgent, create.ErrActionWaitingForDeletion, ResNameFlowVersion, state.Version.String(), err),
			err.Error(),
		)
		return
	}
}

// TIP: ==== TERRAFORM IMPORTING ====
// If Read can get all the information it needs from the Identifier
// (i.e., path.Root("id")), you can use the PassthroughID importer. Otherwise,
// you'll need a custom import function.
//
// See more:
// https://developer.hashicorp.com/terraform/plugin/framework/resources/import
func (r *resourceFlowVersion) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}


// TIP: ==== STATUS CONSTANTS ====
// Create constants for states and statuses if the service does not
// already have suitable constants. We prefer that you use the constants
// provided in the service if available (e.g., awstypes.StatusInProgress).
const (
	statusChangePending = "Pending"
	statusDeleting      = "Deleting"
	statusNormal        = "Normal"
	statusUpdated       = "Updated"
)

// TIP: ==== WAITERS ====
// Some resources of some services have waiters provided by the AWS API.
// Unless they do not work properly, use them rather than defining new ones
// here.
//
// Sometimes we define the wait, status, and find functions in separate
// files, wait.go, status.go, and find.go. Follow the pattern set out in the
// service and define these where it makes the most sense.
//
// If these functions are used in the _test.go file, they will need to be
// exported (i.e., capitalized).
//
// You will need to adjust the parameters and names to fit the service.
func waitFlowVersionCreated(ctx context.Context, conn *bedrockagent.Client, id string, version string, timeout time.Duration) (*bedrockagent.GetFlowVersionOutput, error) {
	stateConf := &retry.StateChangeConf{
		Pending:                   []string{},
		Target:                    []string{statusNormal},
		Refresh:                   statusFlowVersion(ctx, conn, id, version),
		Timeout:                   timeout,
		NotFoundChecks:            20,
		ContinuousTargetOccurence: 2,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if out, ok := outputRaw.(*bedrockagent.GetFlowVersionOutput); ok {
		return out, err
	}

	return nil, err
}

// TIP: A deleted waiter is almost like a backwards created waiter. There may
// be additional pending states, however.
func waitFlowVersionDeleted(ctx context.Context, conn *bedrockagent.Client, id string, version string, timeout time.Duration) (*bedrockagent.GetFlowVersionOutput, error) {
	stateConf := &retry.StateChangeConf{
		Pending:                   []string{statusDeleting, statusNormal},
		Target:                    []string{},
		Refresh:                   statusFlowVersion(ctx, conn, id, version),
		Timeout:                   timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if out, ok := outputRaw.(*bedrockagent.GetFlowVersionOutput); ok {
		return out, err
	}

	return nil, err
}

// TIP: ==== STATUS ====
// The status function can return an actual status when that field is
// available from the API (e.g., out.Status). Otherwise, you can use custom
// statuses to communicate the states of the resource.
//
// Waiters consume the values returned by status functions. Design status so
// that it can be reused by a create, update, and delete waiter, if possible.
func statusFlowVersion(ctx context.Context, conn *bedrockagent.Client, id string, version string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		out, err := findFlowVersionByID(ctx, conn, id, version)
		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return out, string(out.Status), nil
	}
}

// TIP: ==== FINDERS ====
// The find function is not strictly necessary. You could do the API
// request from the status function. However, we have found that find often
// comes in handy in other places besides the status function. As a result, it
// is good practice to define it separately.
func findFlowVersionByID(ctx context.Context, conn *bedrockagent.Client, id string, version string) (*bedrockagent.GetFlowVersionOutput, error) {
	in := &bedrockagent.GetFlowVersionInput{
        FlowIdentifier: aws.String(id),
        FlowVersion: aws.String(version),
	}

	out, err := conn.GetFlowVersion(ctx, in)
	if err != nil {
		if errs.IsA[*awstypes.ResourceNotFoundException](err) {
			return nil, &retry.NotFoundError{
				LastError:   err,
				LastRequest: in,
			}
		}

		return nil, err
	}

	if out == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out, nil
}

// TIP: ==== DATA STRUCTURES ====
// With Terraform Plugin-Framework configurations are deserialized into
// Go types, providing type safety without the need for type assertions.
// These structs should match the schema definition exactly, and the `tfsdk`
// tag value should match the attribute name.
//
// Nested objects are represented in their own data struct. These will
// also have a corresponding attribute type mapping for use inside flex
// functions.
//
// See more:
// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/accessing-values
type resourceFlowVersionModel struct {
	ARN             types.String                                          `tfsdk:"arn"`
    Version         types.String                                          `tfsdk:"version"`
    Name            types.String                                          `tfsdk:"name"`
	FlowIdentifier  types.String                                          `tfsdk:"flow_identifier"`
	Description     types.String                                          `tfsdk:"description"`
	Timeouts        timeouts.Value                                        `tfsdk:"timeouts"`
}

type complexArgumentModel struct {
	NestedRequired types.String `tfsdk:"nested_required"`
	NestedOptional types.String `tfsdk:"nested_optional"`
}
