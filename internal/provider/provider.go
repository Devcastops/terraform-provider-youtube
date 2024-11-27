// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// Ensure YoutubeProvider satisfies various provider interfaces.
var _ provider.Provider = &YoutubeProvider{}
var _ provider.ProviderWithFunctions = &YoutubeProvider{}

// YoutubeProvider defines the provider implementation.
type YoutubeProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// YoutubeProviderModel describes the provider data model.
type YoutubeProviderModel struct {
	AccessToken types.String `tfsdk:"access_token"`
}

func (p *YoutubeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "youtube"
	resp.Version = p.version
}

func (p *YoutubeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"access_token": schema.StringAttribute{
				MarkdownDescription: "api key to auth to youtube, can be gotten from `https://developers.google.com/oauthplayground/#step1&apisSelect=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fyoutube&url=https%3A%2F%2F&content_type=application%2Fjson&http_method=GET&useDefaultOauthCred=unchecked&oauthEndpointSelect=Google&oauthAuthEndpointValue=https%3A%2F%2Faccounts.google.com%2Fo%2Foauth2%2Fv2%2Fauth&oauthTokenEndpointValue=https%3A%2F%2Foauth2.googleapis.com%2Ftoken&includeCredentials=unchecked&accessTokenType=bearer&autoRefreshToken=unchecked&accessType=offline&prompt=consent&response_type=code&wrapLines=on`",
				Sensitive:           true,
				Required:            true,
			},
		},
	}
}

func (p *YoutubeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data YoutubeProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// config, _ := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
	// 	TargetPrincipal: data.ImpersonatedUserEmail.ValueString(),
	// 	Scopes:          []string{youtube.YoutubeScope},
	// }, option.WithCredentialsFile("blinkbyte-7127f7aa8685.json"))
	// service, err := youtube.NewService(ctx, option.WithTokenSource(config),option.WithScopes())
	// service, err := youtube.NewService(ctx, option.WithCredentialsFile("blinkbyte-7127f7aa8685.json"), option.WithScopes(youtube.YoutubeScope))
	service, err := youtube.NewService(ctx, option.WithTokenSource(oauth2.StaticTokenSource(&oauth2.Token{AccessToken: data.AccessToken.ValueString()})), option.WithScopes(youtube.YoutubeScope))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create youtube API Client",
			"An unexpected error occurred when creating the youtube API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Youtube Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = service
	resp.ResourceData = service
}

func (p *YoutubeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewVideoResource,
	}
}

func (p *YoutubeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewVideoDataSource,
	}
}

func (p *YoutubeProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &YoutubeProvider{
			version: version,
		}
	}
}
