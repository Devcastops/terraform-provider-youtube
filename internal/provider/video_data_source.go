// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"google.golang.org/api/youtube/v3"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &VideoDataSource{}

func NewVideoDataSource() datasource.DataSource {
	return &VideoDataSource{}
}

// VideoDataSource defines the data source implementation.
type VideoDataSource struct {
	client *youtube.Service
}

// VideoDataSourceModel describes the data source data model.
type VideoDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Res         types.String `tfsdk:"res"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
}

func (d *VideoDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_video"
}

func (d *VideoDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language video.
		MarkdownDescription: "Video data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the youtube video",
				Required:            true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "Title of the youtube video",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Title of the youtube video",
				Computed:            true,
			},
			"res": schema.StringAttribute{
				MarkdownDescription: "test for what gets returned",
				Computed:            true,
			},
		},
	}
}

func (d *VideoDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VideoDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data VideoDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	call := d.client.Videos.List([]string{
		"snippet"}).Id(data.Id.ValueString())

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
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
