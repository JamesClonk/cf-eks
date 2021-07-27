package main

import (
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/awseks"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type ClusterStackProps struct {
	awscdk.StackProps
}

func NewClusterStack(scope constructs.Construct, id string, props *ClusterStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	workerRole := awsiam.NewRole(stack, jsii.String("EKSWorkerRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("ec2.amazonaws.com"), nil),
	})

	eksCluster := awseks.NewCluster(stack, jsii.String("EKSCluster"), &awseks.ClusterProps{
		DefaultCapacity: jsii.Number(0),
		Version:         awseks.KubernetesVersion_V1_20(),
		EndpointAccess:  awseks.EndpointAccess_PUBLIC_AND_PRIVATE(),
	})
	eksCluster.AddNodegroupCapacity(jsii.String("EKSNodeGroup"), &awseks.NodegroupOptions{
		NodeRole:      workerRole,
		MinSize:       jsii.Number(1),
		MaxSize:       jsii.Number(5),
		InstanceTypes: &[]awsec2.InstanceType{awsec2.NewInstanceType(jsii.String("t3.medium"))},
	})
	//eksCluster.AwsAuth().AddUserMapping(user, mapping)

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewClusterStack(app, "EKSClusterStack", &ClusterStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Region: jsii.String("eu-central-1"),
	}
}
