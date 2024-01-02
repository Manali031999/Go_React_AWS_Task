package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"gorm.io/gorm"

	"new_be/db"
	"new_be/models"
)

// AwsAccess retrieves EC2 instance details and CPU utilization metrics, then stores them in the database
func AwsAccess() error {
	dbDriver, err := db.GetDBInstance()
	if err != nil {
		log.Fatalf("Error getting database instance: %v", err)
	}
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")

	// Create AWS session
	sess, err := createAWSSession(accessKeyID, secretAccessKey, region)
	if err != nil {
		return fmt.Errorf("error creating AWS session: %v", err)
	}

	// Create EC2 service client
	ec2Svc := createEC2Service(sess)

	// Describe EC2 instances
	instances, err := describeEC2Instances(ec2Svc)
	if err != nil {
		return fmt.Errorf("error describing EC2 instances: %v", err)
	}

	// Fetch and insert CPU utilization for each instance
	for _, instance := range instances {
		if err := fetchAndInsertCPUUtilization(dbDriver, sess, instance); err != nil {
			// Handle the error accordingly
			return fmt.Errorf("error processing instance %s: %v", *instance.InstanceId, err)
		}
	}
	return nil
}

// fetchCPUUtilization retrieves CPU utilization metric data for a specific EC2 instance
func fetchCPUUtilization(sess *session.Session, instanceID string) (map[time.Time]float64, float64, error) {

	// Create a CloudWatch service client
	cloudwatchSvc := cloudwatch.New(sess)

	// Construct input for retrieving CPU utilization metric data
	metricDataInput := &cloudwatch.GetMetricDataInput{
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("cpu_utilization_query"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Namespace:  aws.String("AWS/EC2"),
						MetricName: aws.String("CPUUtilization"),
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("InstanceId"),
								Value: aws.String(instanceID),
							},
						},
					},
					Period: aws.Int64(300), // 5-minute period
					Stat:   aws.String("Average"),
				},
			},
		},
		StartTime: aws.Time(time.Now().Add(-1 * 3600 * time.Second)), // Start time is one hour ago
		EndTime:   aws.Time(time.Now()),
	}

	// Fetch CPU utilization metric data using CloudWatch
	result, err := cloudwatchSvc.GetMetricData(metricDataInput)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching CPU utilization data: %v", err)

	}

	// Extract CPU utilization data and create graphData map
	graphData := make(map[time.Time]float64)
	if len(result.MetricDataResults) > 0 {
		for i, timestamp := range result.MetricDataResults[0].Timestamps {
			value := result.MetricDataResults[0].Values[i]
			graphData[*timestamp] = *value
		}
		var cpuUtilization float64
		if len(result.MetricDataResults) > 0 && len(result.MetricDataResults[0].Values) > 0 {
			cpuUtilization = *result.MetricDataResults[0].Values[0]
		} else {
			cpuUtilization = 0
		}

		return graphData, cpuUtilization, nil
	}

	// If no graphData is available, check for a single CPU utilization value
	if len(result.MetricDataResults) > 0 && len(result.MetricDataResults[0].Values) > 0 {
		cpuUtilization := *result.MetricDataResults[0].Values[0]
		return graphData, cpuUtilization, nil
	}
	return nil, 0, fmt.Errorf("no CPU utilization data available")
}

// createAWSSession creates a new AWS session
func createAWSSession(accessKeyID, secretAccessKey, region string) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKeyID,
			secretAccessKey,
			"", // A session token is not needed for basic AWS credentials
		),
	})
	if err != nil {
		return nil, err
	}
	return sess, nil
}

// createEC2Service creates an EC2 service client
func createEC2Service(sess *session.Session) *ec2.EC2 {
	return ec2.New(sess)
}

// describeEC2Instances retrieves details of EC2 instances
func describeEC2Instances(ec2Svc *ec2.EC2) ([]*ec2.Instance, error) {
	result, err := ec2Svc.DescribeInstances(nil)
	if err != nil {
		return nil, err
	}
	var instances []*ec2.Instance
	for _, reservation := range result.Reservations {
		instances = append(instances, reservation.Instances...)
	}
	return instances, nil
}

// fetchAndInsertCPUUtilization fetches CPU utilization and inserts data into the database
func fetchAndInsertCPUUtilization(dbDriver *gorm.DB, sess *session.Session, instance *ec2.Instance) error {
	graphData, cpuUtilization, err := fetchCPUUtilization(sess, *instance.InstanceId)

	if err != nil {
		fmt.Printf("Error fetching CPU utilization for instance %s: %v\n", *instance.InstanceId, err)
		cpuUtilization = 0
	}

	fmt.Printf("CPU Utilization: %.2f%%\n", cpuUtilization)
	instanceData := models.Instance{
		InstanceID:   *instance.InstanceId,
		InstanceType: *instance.InstanceType,
		Region:       *instance.Placement.AvailabilityZone,
	}
	err = db.InsertInstance(dbDriver, instanceData)
	if err != nil {
		log.Printf("Error inserting instance data: %v", err)
	}

	jsonData, err := json.Marshal(graphData)
	if err != nil {
		fmt.Println("Error while marshalling graph data", err)
		return err
	}
	metricData := models.MetricData{
		InstanceID: instanceData.InstanceID,
		CPU:        cpuUtilization,
		GraphData:  jsonData,
	}
	err = db.InsertMetricData(dbDriver, metricData)
	if err != nil {
		log.Printf("Error inserting metric data: %v", err)
	}
	return nil
}
