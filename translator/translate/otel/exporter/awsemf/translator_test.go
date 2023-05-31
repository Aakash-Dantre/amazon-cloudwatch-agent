// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package awsemf

import (
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awsemfexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/resourcetotelemetry"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/confmap"

	legacytranslator "github.com/aws/private-amazon-cloudwatch-agent-staging/translator"
)

var nilSlice []string
var nilMetricDescriptorsSlice []awsemfexporter.MetricDescriptor

func TestTranslator(t *testing.T) {
	tt := NewTranslator()
	require.EqualValues(t, "awsemf", tt.ID().String())
	testCases := map[string]struct {
		env     map[string]string
		input   map[string]interface{}
		want    map[string]interface{} // Can't construct & use awsemfexporter.Config as it uses internal only types
		wantErr error
	}{
		"GenerateAwsEmfExporterConfigEcs": {
			input: map[string]interface{}{
				"logs": map[string]interface{}{
					"metrics_collected": map[string]interface{}{
						"ecs": map[string]interface{}{},
					},
				},
			},
			want: map[string]interface{}{
				"namespace":                              "ECS/ContainerInsights",
				"log_group_name":                         "/aws/ecs/containerinsights/{ClusterName}/performance",
				"log_stream_name":                        "NodeTelemetry-{ContainerInstanceId}",
				"dimension_rollup_option":                "NoDimensionRollup",
				"parse_json_encoded_attr_values":         []string{"Sources"},
				"output_destination":                     "cloudwatch",
				"eks_fargate_container_insights_enabled": false,
				"resource_to_telemetry_conversion": resourcetotelemetry.Settings{
					Enabled: true,
				},
				"metric_declarations": []*awsemfexporter.MetricDeclaration{
					{
						Dimensions: [][]string{{"ContainerInstanceId", "InstanceId", "ClusterName"}},
						MetricNameSelectors: []string{"instance_cpu_reserved_capacity", "instance_cpu_utilization",
							"instance_filesystem_utilization", "instance_memory_reserved_capacity",
							"instance_memory_utilization", "instance_network_total_bytes", "instance_number_of_running_tasks"},
					},
					{
						Dimensions: [][]string{{"ClusterName"}},
						MetricNameSelectors: []string{"instance_cpu_limit", "instance_cpu_reserved_capacity",
							"instance_cpu_usage_total", "instance_cpu_utilization", "instance_filesystem_utilization",
							"instance_memory_limit", "instance_memory_reserved_capacity", "instance_memory_utilization",
							"instance_memory_working_set", "instance_network_total_bytes", "instance_number_of_running_tasks"},
					},
				},
				"metric_descriptors": nilMetricDescriptorsSlice,
			},
		},
		"GenerateAwsEmfExporterConfigKubernetes": {
			input: map[string]interface{}{
				"logs": map[string]interface{}{
					"metrics_collected": map[string]interface{}{
						"kubernetes": map[string]interface{}{},
					},
				},
			},
			want: map[string]interface{}{
				"namespace":                              "ContainerInsights",
				"log_group_name":                         "/aws/containerinsights/{ClusterName}/performance",
				"log_stream_name":                        "{NodeName}",
				"dimension_rollup_option":                "NoDimensionRollup",
				"parse_json_encoded_attr_values":         []string{"Sources", "kubernetes"},
				"output_destination":                     "cloudwatch",
				"eks_fargate_container_insights_enabled": false,
				"resource_to_telemetry_conversion": resourcetotelemetry.Settings{
					Enabled: true,
				},
				"metric_declarations": []*awsemfexporter.MetricDeclaration{
					{
						Dimensions: [][]string{{"PodName", "Namespace", "ClusterName"}, {"Service", "Namespace", "ClusterName"}, {"Namespace", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"pod_cpu_utilization", "pod_memory_utilization",
							"pod_network_rx_bytes", "pod_network_tx_bytes", "pod_cpu_utilization_over_pod_limit",
							"pod_memory_utilization_over_pod_limit"},
					},
					{
						Dimensions:          [][]string{{"PodName", "Namespace", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"pod_cpu_reserved_capacity", "pod_memory_reserved_capacity"},
					},
					{
						Dimensions:          [][]string{{"PodName", "Namespace", "ClusterName"}},
						MetricNameSelectors: []string{"pod_number_of_container_restarts"},
					},
					{
						Dimensions: [][]string{{"NodeName", "InstanceId", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"node_cpu_utilization", "node_memory_utilization",
							"node_network_total_bytes", "node_cpu_reserved_capacity",
							"node_memory_reserved_capacity", "node_number_of_running_pods", "node_number_of_running_containers"},
					},
					{
						Dimensions:          [][]string{{"ClusterName"}},
						MetricNameSelectors: []string{"node_cpu_usage_total", "node_cpu_limit", "node_memory_working_set", "node_memory_limit"},
					},
					{
						Dimensions:          [][]string{{"NodeName", "InstanceId", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"node_filesystem_utilization"},
					},
					{
						Dimensions:          [][]string{{"Service", "Namespace", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"service_number_of_running_pods"},
					},
					{
						Dimensions:          [][]string{{"Namespace", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"namespace_number_of_running_pods"},
					},
					{
						Dimensions:          [][]string{{"ClusterName"}},
						MetricNameSelectors: []string{"cluster_node_count", "cluster_failed_node_count"},
					},
				},
				"metric_descriptors": nilMetricDescriptorsSlice,
			},
		},
		"GenerateAwsEmfExporterConfigKubernetesWithEnableFullPodMetrics": {
			input: map[string]interface{}{
				"logs": map[string]interface{}{
					"metrics_collected": map[string]interface{}{
						"kubernetes": map[string]interface{}{
							"enable_full_pod_metrics": true,
						},
					},
				},
			},
			want: map[string]interface{}{
				"namespace":                              "ContainerInsights",
				"log_group_name":                         "/aws/containerinsights/{ClusterName}/performance",
				"log_stream_name":                        "{NodeName}",
				"dimension_rollup_option":                "NoDimensionRollup",
				"parse_json_encoded_attr_values":         []string{"Sources", "kubernetes"},
				"output_destination":                     "cloudwatch",
				"eks_fargate_container_insights_enabled": false,
				"resource_to_telemetry_conversion": resourcetotelemetry.Settings{
					Enabled: true,
				},
				"metric_declarations": []*awsemfexporter.MetricDeclaration{
					{
						Dimensions: [][]string{{"FullPodName", "PodName", "Namespace", "ClusterName"}, {"PodName", "Namespace", "ClusterName"}, {"Service", "Namespace", "ClusterName"}, {"Namespace", "ClusterName"},
							{"ClusterName"}},
						MetricNameSelectors: []string{"pod_cpu_utilization", "pod_memory_utilization",
							"pod_network_rx_bytes", "pod_network_tx_bytes", "pod_cpu_utilization_over_pod_limit",
							"pod_memory_utilization_over_pod_limit"},
					},
					{
						Dimensions: [][]string{{"FullPodName", "PodName", "Namespace", "ClusterName"}, {"PodName", "Namespace", "ClusterName"}, {"ClusterName"}, {"Service", "Namespace", "ClusterName"}},
						MetricNameSelectors: []string{"pod_cpu_reserved_capacity", "pod_memory_reserved_capacity", "pod_number_of_container_restarts",
							"pod_number_of_containers", "pod_number_of_running_containers"},
					},
					{
						Dimensions: [][]string{{"NodeName", "InstanceId", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"node_cpu_utilization", "node_memory_utilization",
							"node_network_total_bytes", "node_cpu_reserved_capacity",
							"node_memory_reserved_capacity", "node_number_of_running_pods", "node_number_of_running_containers"},
					},
					{
						Dimensions:          [][]string{{"ClusterName"}},
						MetricNameSelectors: []string{"node_cpu_usage_total", "node_cpu_limit", "node_memory_working_set", "node_memory_limit"},
					},
					{
						Dimensions:          [][]string{{"NodeName", "InstanceId", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"node_filesystem_utilization"},
					},
					{
						Dimensions:          [][]string{{"Service", "Namespace", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"service_number_of_running_pods"},
					},
					{
						Dimensions:          [][]string{{"Namespace", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"namespace_number_of_running_pods"},
					},
					{
						Dimensions:          [][]string{{"ClusterName"}},
						MetricNameSelectors: []string{"cluster_node_count", "cluster_failed_node_count"},
					},
				},
				"metric_descriptors": nilMetricDescriptorsSlice,
			},
		},
		"GenerateAwsEmfExporterConfigKubernetesWithEnableContainerMetrics": {
			input: map[string]interface{}{
				"logs": map[string]interface{}{
					"metrics_collected": map[string]interface{}{
						"kubernetes": map[string]interface{}{
							"enable_container_metrics": true,
						},
					},
				},
			},
			want: map[string]interface{}{
				"namespace":                              "ContainerInsights",
				"log_group_name":                         "/aws/containerinsights/{ClusterName}/performance",
				"log_stream_name":                        "{NodeName}",
				"dimension_rollup_option":                "NoDimensionRollup",
				"parse_json_encoded_attr_values":         []string{"Sources", "kubernetes"},
				"output_destination":                     "cloudwatch",
				"eks_fargate_container_insights_enabled": false,
				"resource_to_telemetry_conversion": resourcetotelemetry.Settings{
					Enabled: true,
				},
				"metric_declarations": []*awsemfexporter.MetricDeclaration{
					{
						Dimensions:          [][]string{{"ContainerName", "FullPodName", "Namespace", "ClusterName"}, {"ContainerName", "Namespace", "ClusterName"}},
						MetricNameSelectors: []string{"container_cpu_utilization", "container_memory_utilization", "container_filesystem_usage"},
					},
					{
						Dimensions: [][]string{{"PodName", "Namespace", "ClusterName"}, {"Service", "Namespace", "ClusterName"}, {"Namespace", "ClusterName"},
							{"ClusterName"}},
						MetricNameSelectors: []string{"pod_cpu_utilization", "pod_memory_utilization",
							"pod_network_rx_bytes", "pod_network_tx_bytes", "pod_cpu_utilization_over_pod_limit",
							"pod_memory_utilization_over_pod_limit"},
					},
					{
						Dimensions:          [][]string{{"PodName", "Namespace", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"pod_cpu_reserved_capacity", "pod_memory_reserved_capacity"},
					},
					{
						Dimensions:          [][]string{{"PodName", "Namespace", "ClusterName"}},
						MetricNameSelectors: []string{"pod_number_of_container_restarts"},
					},
					{
						Dimensions: [][]string{{"NodeName", "InstanceId", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"node_cpu_utilization", "node_memory_utilization",
							"node_network_total_bytes", "node_cpu_reserved_capacity",
							"node_memory_reserved_capacity", "node_number_of_running_pods", "node_number_of_running_containers"},
					},
					{
						Dimensions:          [][]string{{"ClusterName"}},
						MetricNameSelectors: []string{"node_cpu_usage_total", "node_cpu_limit", "node_memory_working_set", "node_memory_limit"},
					},
					{
						Dimensions:          [][]string{{"NodeName", "InstanceId", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"node_filesystem_utilization"},
					},
					{
						Dimensions:          [][]string{{"Service", "Namespace", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"service_number_of_running_pods"},
					},
					{
						Dimensions:          [][]string{{"PodName", "Namespace", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"deployment_spec_replicas", "deployment_status_replicas", "deployment_status_replicas_available", "deployment_status_replicas_unavailable"},
					},
					{
						Dimensions: [][]string{{"PodName", "Namespace", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"daemonset_status_number_available", "daemonset_status_number_unavailable",
							"daemonset_status_desired_number_scheduled", "daemonset_status_current_number_scheduled"},
					},
					{
						Dimensions:          [][]string{{"Namespace", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"namespace_number_of_running_pods"},
					},
					{
						Dimensions:          [][]string{{"ClusterName"}},
						MetricNameSelectors: []string{"cluster_node_count", "cluster_failed_node_count"},
					},
				},
				"metric_descriptors": nilMetricDescriptorsSlice,
			},
		},
		"GenerateAwsEmfExporterConfigKubernetesWithEnableFullPodAndContainerMetrics": {
			input: map[string]interface{}{
				"logs": map[string]interface{}{
					"metrics_collected": map[string]interface{}{
						"kubernetes": map[string]interface{}{
							"enable_full_pod_metrics":  true,
							"enable_container_metrics": true,
						},
					},
				},
			},
			want: map[string]interface{}{
				"namespace":                              "ContainerInsights",
				"log_group_name":                         "/aws/containerinsights/{ClusterName}/performance",
				"log_stream_name":                        "{NodeName}",
				"dimension_rollup_option":                "NoDimensionRollup",
				"parse_json_encoded_attr_values":         []string{"Sources", "kubernetes"},
				"output_destination":                     "cloudwatch",
				"eks_fargate_container_insights_enabled": false,
				"resource_to_telemetry_conversion": resourcetotelemetry.Settings{
					Enabled: true,
				},
				"metric_declarations": []*awsemfexporter.MetricDeclaration{
					{
						Dimensions:          [][]string{{"ContainerName", "FullPodName", "Namespace", "ClusterName"}, {"ContainerName", "Namespace", "ClusterName"}},
						MetricNameSelectors: []string{"container_cpu_utilization", "container_memory_utilization", "container_filesystem_usage"},
					},
					{
						Dimensions: [][]string{{"FullPodName", "PodName", "Namespace", "ClusterName"}, {"PodName", "Namespace", "ClusterName"}, {"Service", "Namespace", "ClusterName"}, {"Namespace", "ClusterName"},
							{"ClusterName"}},
						MetricNameSelectors: []string{"pod_cpu_utilization", "pod_memory_utilization",
							"pod_network_rx_bytes", "pod_network_tx_bytes", "pod_cpu_utilization_over_pod_limit",
							"pod_memory_utilization_over_pod_limit"},
					},
					{
						Dimensions: [][]string{{"FullPodName", "PodName", "Namespace", "ClusterName"}, {"PodName", "Namespace", "ClusterName"}, {"ClusterName"}, {"Service", "Namespace", "ClusterName"}},
						MetricNameSelectors: []string{"pod_cpu_reserved_capacity", "pod_memory_reserved_capacity", "pod_number_of_container_restarts",
							"pod_number_of_containers", "pod_number_of_running_containers"},
					},
					{
						Dimensions: [][]string{{"NodeName", "InstanceId", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"node_cpu_utilization", "node_memory_utilization",
							"node_network_total_bytes", "node_cpu_reserved_capacity",
							"node_memory_reserved_capacity", "node_number_of_running_pods", "node_number_of_running_containers"},
					},
					{
						Dimensions:          [][]string{{"ClusterName"}},
						MetricNameSelectors: []string{"node_cpu_usage_total", "node_cpu_limit", "node_memory_working_set", "node_memory_limit"},
					},
					{
						Dimensions:          [][]string{{"NodeName", "InstanceId", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"node_filesystem_utilization"},
					},
					{
						Dimensions:          [][]string{{"Service", "Namespace", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"service_number_of_running_pods"},
					},
					{
						Dimensions:          [][]string{{"PodName", "Namespace", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"deployment_spec_replicas", "deployment_status_replicas", "deployment_status_replicas_available", "deployment_status_replicas_unavailable"},
					},
					{
						Dimensions: [][]string{{"PodName", "Namespace", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"daemonset_status_number_available", "daemonset_status_number_unavailable",
							"daemonset_status_desired_number_scheduled", "daemonset_status_current_number_scheduled"},
					},
					{
						Dimensions:          [][]string{{"Namespace", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"namespace_number_of_running_pods"},
					},
					{
						Dimensions:          [][]string{{"ClusterName"}},
						MetricNameSelectors: []string{"cluster_node_count", "cluster_failed_node_count"},
					},
				},
				"metric_descriptors": nilMetricDescriptorsSlice,
			},
		},
		"GenerateAwsEmfExporterConfigKubernetesWithEnableNodeDetailedMetrics": {
			input: map[string]interface{}{
				"logs": map[string]interface{}{
					"metrics_collected": map[string]interface{}{
						"kubernetes": map[string]interface{}{
							"enable_node_detailed_metrics": true,
						},
					},
				},
			},
			want: map[string]interface{}{
				"namespace":                              "ContainerInsights",
				"log_group_name":                         "/aws/containerinsights/{ClusterName}/performance",
				"log_stream_name":                        "{NodeName}",
				"dimension_rollup_option":                "NoDimensionRollup",
				"parse_json_encoded_attr_values":         []string{"Sources", "kubernetes"},
				"output_destination":                     "cloudwatch",
				"eks_fargate_container_insights_enabled": false,
				"resource_to_telemetry_conversion": resourcetotelemetry.Settings{
					Enabled: true,
				},
				"metric_declarations": []*awsemfexporter.MetricDeclaration{
					{
						Dimensions: [][]string{{"PodName", "Namespace", "ClusterName"}, {"Service", "Namespace", "ClusterName"}, {"Namespace", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"pod_cpu_utilization", "pod_memory_utilization",
							"pod_network_rx_bytes", "pod_network_tx_bytes", "pod_cpu_utilization_over_pod_limit",
							"pod_memory_utilization_over_pod_limit"},
					},
					{
						Dimensions:          [][]string{{"PodName", "Namespace", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"pod_cpu_reserved_capacity", "pod_memory_reserved_capacity"},
					},
					{
						Dimensions:          [][]string{{"PodName", "Namespace", "ClusterName"}},
						MetricNameSelectors: []string{"pod_number_of_container_restarts"},
					},
					{
						Dimensions: [][]string{{"NodeName", "InstanceId", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"node_cpu_utilization", "node_memory_utilization",
							"node_network_total_bytes", "node_cpu_reserved_capacity",
							"node_memory_reserved_capacity", "node_number_of_running_pods", "node_number_of_running_containers",
							"node_cpu_usage_total", "node_cpu_limit", "node_memory_working_set", "node_memory_limit",
							"node_status_condition_ready", "node_status_condition_disk_pressure", "node_status_condition_memory_pressure",
							"node_status_condition_pid_pressure", "node_status_condition_network_unavailable",
							"node_status_capacity_pods", "node_status_allocatable_pods"},
					},
					{
						Dimensions:          [][]string{{"NodeName", "InstanceId", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"node_filesystem_utilization", "node_filesystem_inodes", "node_filesystem_inodes_free"},
					},
					{
						Dimensions:          [][]string{{"Service", "Namespace", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"service_number_of_running_pods"},
					},
					{
						Dimensions:          [][]string{{"Namespace", "ClusterName"}, {"ClusterName"}},
						MetricNameSelectors: []string{"namespace_number_of_running_pods"},
					},
					{
						Dimensions:          [][]string{{"ClusterName"}},
						MetricNameSelectors: []string{"cluster_node_count", "cluster_failed_node_count"},
					},
				},
				"metric_descriptors": nilMetricDescriptorsSlice,
			},
		},
		"GenerateAwsEmfExporterConfigPrometheus": {
			input: map[string]interface{}{
				"logs": map[string]interface{}{
					"metrics_collected": map[string]interface{}{
						"prometheus": map[string]interface{}{
							"log_group_name":  "/test/log/group",
							"log_stream_name": "{ServiceName}",
							"emf_processor": map[string]interface{}{
								"metric_declaration": []interface{}{
									map[string]interface{}{
										"source_labels":    []string{"Service", "Namespace"},
										"label_matcher":    "(.*node-exporter.*|.*kube-dns.*);kube-system$",
										"dimensions":       [][]string{{"Service", "Namespace"}},
										"metric_selectors": []string{"^coredns_dns_request_type_count_total$"},
									},
								},
								"metric_unit": map[string]interface{}{
									"jvm_gc_collection_seconds_sum": "Milliseconds",
								},
							},
						},
					},
				},
			},
			want: map[string]interface{}{
				"namespace":                              "CWAgent/Prometheus",
				"log_group_name":                         "/test/log/group",
				"log_stream_name":                        "{ServiceName}",
				"dimension_rollup_option":                "NoDimensionRollup",
				"parse_json_encoded_attr_values":         nilSlice,
				"output_destination":                     "cloudwatch",
				"eks_fargate_container_insights_enabled": false,
				"resource_to_telemetry_conversion": resourcetotelemetry.Settings{
					Enabled: true,
				},
				"metric_declarations": []*awsemfexporter.MetricDeclaration{
					{
						Dimensions:          [][]string{{"Service", "Namespace"}},
						MetricNameSelectors: []string{"^coredns_dns_request_type_count_total$"},
						LabelMatchers: []*awsemfexporter.LabelMatcher{
							{
								LabelNames: []string{"Service", "Namespace"},
								Regex:      "(.*node-exporter.*|.*kube-dns.*);kube-system$",
							},
						},
					},
				},
				"metric_descriptors": []awsemfexporter.MetricDescriptor{
					{
						MetricName: "jvm_gc_collection_seconds_sum",
						Unit:       "Milliseconds",
					},
				},
			},
		},
		"GenerateAwsEmfExporterConfigPrometheusNoDeclarations": {
			input: map[string]interface{}{
				"logs": map[string]interface{}{
					"metrics_collected": map[string]interface{}{
						"prometheus": map[string]interface{}{
							"log_group_name":  "/test/log/group",
							"log_stream_name": "{ServiceName}",
							"emf_processor": map[string]interface{}{
								"metric_unit": map[string]interface{}{
									"jvm_gc_collection_seconds_sum": "Milliseconds",
								},
							},
						},
					},
				},
			},
			want: map[string]interface{}{
				"namespace":                              "CWAgent/Prometheus",
				"log_group_name":                         "/test/log/group",
				"log_stream_name":                        "{ServiceName}",
				"dimension_rollup_option":                "NoDimensionRollup",
				"parse_json_encoded_attr_values":         nilSlice,
				"output_destination":                     "cloudwatch",
				"eks_fargate_container_insights_enabled": false,
				"resource_to_telemetry_conversion": resourcetotelemetry.Settings{
					Enabled: true,
				},
				"metric_declarations": []*awsemfexporter.MetricDeclaration{
					{
						MetricNameSelectors: []string{"$^"},
					},
				},
				"metric_descriptors": []awsemfexporter.MetricDescriptor{
					{
						MetricName: "jvm_gc_collection_seconds_sum",
						Unit:       "Milliseconds",
					},
				},
			},
		},
		"GenerateAwsEmfExporterConfigPrometheusNoEmfProcessor": {
			input: map[string]interface{}{
				"logs": map[string]interface{}{
					"metrics_collected": map[string]interface{}{
						"prometheus": map[string]interface{}{
							"log_group_name":  "/test/log/group",
							"log_stream_name": "{ServiceName}",
						},
					},
				},
			},
			want: map[string]interface{}{
				"namespace":                              "",
				"log_group_name":                         "/test/log/group",
				"log_stream_name":                        "{ServiceName}",
				"dimension_rollup_option":                "NoDimensionRollup",
				"parse_json_encoded_attr_values":         nilSlice,
				"output_destination":                     "cloudwatch",
				"eks_fargate_container_insights_enabled": false,
				"resource_to_telemetry_conversion": resourcetotelemetry.Settings{
					Enabled: true,
				},
				"metric_declarations": []*awsemfexporter.MetricDeclaration{
					{
						MetricNameSelectors: []string{"$^"},
					},
				},
				"metric_descriptors": nilMetricDescriptorsSlice,
			},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			conf := confmap.NewFromStringMap(testCase.input)
			got, err := tt.Translate(conf)
			require.Equal(t, testCase.wantErr, err)
			require.Truef(t, legacytranslator.IsTranslateSuccess(), "Error in legacy translation rules: %v", legacytranslator.ErrorMessages)
			if err == nil {
				require.NotNil(t, got)
				gotCfg, ok := got.(*awsemfexporter.Config)
				require.True(t, ok)
				require.Equal(t, testCase.want["namespace"], gotCfg.Namespace)
				require.Equal(t, testCase.want["log_group_name"], gotCfg.LogGroupName)
				require.Equal(t, testCase.want["log_stream_name"], gotCfg.LogStreamName)
				require.Equal(t, testCase.want["dimension_rollup_option"], gotCfg.DimensionRollupOption)
				require.Equal(t, testCase.want["parse_json_encoded_attr_values"], gotCfg.ParseJSONEncodedAttributeValues)
				require.Equal(t, testCase.want["output_destination"], gotCfg.OutputDestination)
				require.Equal(t, testCase.want["eks_fargate_container_insights_enabled"], gotCfg.EKSFargateContainerInsightsEnabled)
				require.Equal(t, testCase.want["resource_to_telemetry_conversion"], gotCfg.ResourceToTelemetrySettings)
				require.Equal(t, testCase.want["metric_declarations"], gotCfg.MetricDeclarations)
				require.Equal(t, testCase.want["metric_descriptors"], gotCfg.MetricDescriptors)
			}
		})
	}
}
