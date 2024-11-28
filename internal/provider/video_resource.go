// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"google.golang.org/api/youtube/v3"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &VideoResource{}
var _ resource.ResourceWithImportState = &VideoResource{}

func NewVideoResource() resource.Resource {
	return &VideoResource{}
}

// VideoResource defines the resource implementation.
type VideoResource struct {
	client *youtube.Service
}

// VideoResourceModel describes the resource data model.
type VideoResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Res         types.String `tfsdk:"res"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
}

func (d *VideoResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_video"
}

func (d *VideoResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language video.
		MarkdownDescription: "Video resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the youtube video",
				Required:            true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "Title of the youtube video",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Title of the youtube video",
				Required:            true,
			},
			"res": schema.StringAttribute{
				MarkdownDescription: "test for what gets returned",
				Computed:            true,
			},
		},
	}
}

func (d *VideoResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*youtube.Service)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *VideoResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data VideoResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	call := d.client.Videos.List([]string{
		"contentDetails",
		// "fileDetails",
		"id",
		"liveStreamingDetails",
		"localizations",
		// "paidProductPlacementDetails",
		"player",
		// "processingDetails",
		"recordingDetails",
		"snippet",
		"statistics",
		"status",
		// "suggestions",
		"topicDetails",
	}).Id(data.Id.ValueString())

	response, err := call.Do()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get Video",
			err.Error(),
		)
		return
	}
	if len(response.Items) < 1 {
		resp.Diagnostics.AddError(
			"Unable to get Video",
			"no videos found for ID",
		)
		return
	}
	data.Res = types.StringValue(fmt.Sprintf("%+v", response.Items[0].Snippet))
	data.Title = types.StringValue(response.Items[0].Snippet.Title)
	data.Description = types.StringValue(response.Items[0].Snippet.Description)
	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Create implements resource.Resource.
func (d *VideoResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddError(
		"Creation not supported",
		"Please use an import block, as youtube videos cant be ",
	)
}

// Delete implements resource.Resource.
func (d *VideoResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

// Update implements resource.Resource.
func (d *VideoResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data VideoResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	call := d.client.Videos.List([]string{
		"contentDetails",
		// "fileDetails",
		"id",
		"liveStreamingDetails",
		"localizations",
		// "paidProductPlacementDetails",
		"player",
		// "processingDetails",
		"recordingDetails",
		"snippet",
		"statistics",
		"status",
		// "suggestions",
		"topicDetails",
	}).Id(data.Id.ValueString())

	response, err := call.Do()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get Video",
			err.Error(),
		)
		return
	}
	if len(response.Items) < 1 {
		resp.Diagnostics.AddError(
			"Unable to get Video",
			"no videos found for ID",
		)
		return
	}
	response.Items[0].Snippet.Description = data.Description.ValueString()

	updateCall := d.client.Videos.Update([]string{
		"snippet"}, response.Items[0])

	updateRes, err := updateCall.Do()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get Video",
			err.Error(),
		)
		return
	}
	data.Res = types.StringValue(fmt.Sprintf("%+v", updateRes))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ImportState implements resource.ResourceWithImportState.
func (d *VideoResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
