package quark

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type QuarkEc2InstanceInfo struct {
	InstanceId     string `json:"instance_id"`
	ImageId        string `json:"image_id"`
	InstanceType   string `json:"instance_type"`
	RootDeviceName string `json:"root_device_name"`
	ClientToken    string `json:"client_token"`
}

const (
	defaultRegion = "us-west-2"
)

func tableQuarkEc2Instance() *plugin.Table {
	return &plugin.Table{
		Name:        "quark_ec2_instance",
		Description: "To list all ec2 instances",
		List: &plugin.ListConfig{
			Hydrate: listQuarkEc2Instance,
		},
		Columns: []*plugin.Column{
			{
				Name:        "instance_id",
				Description: "To display the instance_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InstanceId"),
			},
			{
				Name:        "image_id",
				Description: "To display the image_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImageId"),
			},
			{
				Name:        "instance_type",
				Description: "To display the instance type",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InstanceType"),
			},
			{
				Name:        "root_device_name",
				Description: "To display the root_device_name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RootDeviceName"),
			},
			{
				Name:        "client_token",
				Description: "To display the client_token",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClientToken"),
			},
		},
	}
}

func listQuarkEc2Instance(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Info("listQuarkEc2Instance started by logger")
	conn, ok := d.Connection.Config.(awsConfig)
	if !ok {
		return nil, nil
	}

	// Create a custom AWS configuration
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithDefaultRegion(defaultRegion),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				conn.AwsAccessKeyID,
				conn.AwsSecretAccessKey,
				conn.AwsSessionToken),
		),
	)

	if err != nil {
		plugin.Logger(ctx).Error("unable to load SDK config, %v", err)
		return nil, nil
	}

	// Create an EC2 client
	svc := ec2.NewFromConfig(cfg)

	// Define the input parameters
	maxLimit := int32(10)
	input := &ec2.DescribeInstancesInput{
		MaxResults: aws.Int32(maxLimit),
	}

	// fetch all regions
	regions, regionsErr := svc.DescribeRegions(ctx, &ec2.DescribeRegionsInput{})

	if regionsErr != nil {
		plugin.Logger(ctx).Error("unable to fetch regions , %v", regionsErr)
		return nil, nil
	}

	for _, region := range regions.Regions {
		plugin.Logger(ctx).Info("fetching Region: ", *region.RegionName)
		cfg.Region = *region.RegionName
		ec2Svc := ec2.NewFromConfig(cfg)

		// Create a paginator
		paginator := ec2.NewDescribeInstancesPaginator(ec2Svc, input)

		// Iterate through the pages
		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				plugin.Logger(ctx).Error("failed to get page, %s %v", cfg.Region, err)
				break
			}

			// Print instance details
			for _, reservation := range page.Reservations {
				for _, instance := range reservation.Instances {
					d.StreamListItem(ctx, QuarkEc2InstanceInfo{
						InstanceId:     *instance.InstanceId,
						ImageId:        *instance.ImageId,
						InstanceType:   string(instance.InstanceType),
						RootDeviceName: *instance.RootDeviceName,
						ClientToken:    *instance.ClientToken,
					})
				}
			}
		}

	}

	return nil, nil
}
