package java

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-oracle-terraform/client"
)

const WaitForServiceInstanceReadyTimeout = time.Duration(3600 * time.Second)
const WaitForServiceInstanceDeleteTimeout = time.Duration(3600 * time.Second)
const ServiceInstanceDeleteRetry = 30

var (
	ServiceInstanceContainerPath = "/paas/service/jcs/api/v1.1/instances/%s"
	ServiceInstanceResourcePath  = "/paas/service/jcs/api/v1.1/instances/%s/%s"
)

// ServiceInstanceClient is a client for the Service functions of the Java API.
type ServiceInstanceClient struct {
	ResourceClient
	Timeout time.Duration
}

// ServiceInstanceClient obtains an ServiceInstanceClient which can be used to access to the
// Service Instance functions of the Java Cloud API
func (c *JavaClient) ServiceInstanceClient() *ServiceInstanceClient {
	return &ServiceInstanceClient{
		ResourceClient: ResourceClient{
			JavaClient:       c,
			ContainerPath:    ServiceInstanceContainerPath,
			ResourceRootPath: ServiceInstanceResourcePath,
		}}
}

type ServiceInstanceLevel string

const (
	// PAAS: Production-level service. This is the default. Supports Oracle Java Cloud Service instance creation
	// and monitoring, backup and restoration, patching, and scaling. Use PAAS if you want to enable domain partitions
	// using WebLogic Server 12.2.1, use AppToCloud artifacts to create a service instance, or create a service instance
	// for an Oracle Fusion Middleware product.
	ServiceInstanceLevelPAAS ServiceInstanceLevel = "PAAS"
	// BASIC: Development-level service. Supports Oracle Java Cloud Service instance creation and monitoring
	// but does not support backup and restoration, patching, or scaling.
	ServiceInstanceLevelBasic ServiceInstanceLevel = "BASIC"
)

type ServiceInstanceBackupDestination string

const (
	// BOTH - Enable backups. This is the default. This means automated scheduled backups are enabled,
	// and on-demand backups can be initiated. All backups are stored on disk and the Oracle Storage
	// Cloud Service container that is specified in cloudStorageContainer.
	ServiceInstanceBackupDestinationBoth ServiceInstanceBackupDestination = "BOTH"
	// NONE - Do not enable backups. This means automated scheduled backups are not enabled,
	// and on-demand backups cannot be initiated. When set to NONE, cloudStorageContainer is not required.
	ServiceInstanceBackupDestinationNone ServiceInstanceBackupDestination = "NONE"
)

type ServiceInstanceTargetDataSourceType string

const (
	// If the specified Database Cloud Service database deployment does not use Oracle RAC, the value must be Generic.
	ServiceInstanceTargetDataSourceTypeGeneric ServiceInstanceTargetDataSourceType = "Generic"
	// If the specified Database Cloud Service database deployment uses Oracle RAC and the specified edition
	// (for WebLogic Server software) is EE, the value must be Multi.
	ServiceInstanceTargetDataSourceTypeMulti ServiceInstanceTargetDataSourceType = "Multi"
	// If the specified Database Cloud Service database deployment uses Oracle RAC and the specified edition
	// (for WebLogic Server software) is SUITE, the value can be GridLink or Multi.
	ServiceInstanceTargetDataSourceTypeGridLink ServiceInstanceTargetDataSourceType = "GridLink"
)

type ServiceInstanceProtocol string

const (
	ServiceInstanceProtocolT3    ServiceInstanceProtocol = "t3"
	ServiceInstanceProtocolT3S   ServiceInstanceProtocol = "t3s"
	ServiceInstanceProtocolIIOP  ServiceInstanceProtocol = "iiop"
	ServiceInstanceProtocolIIOPS ServiceInstanceProtocol = "iiops"
)

type ServiceInstanceDomainMode string

const (
	ServiceInstanceDomainModeDev ServiceInstanceDomainMode = "DEVELOPMENT"
	ServiceInstanceDomainModePro ServiceInstanceDomainMode = "PRODUCTION"
)

type ServiceInstanceEdition string

const (
	ServiceInstanceEditionSE    ServiceInstanceEdition = "SE"
	ServiceInstanceEditionEE    ServiceInstanceEdition = "EE"
	ServiceInstanceEditionSuite ServiceInstanceEdition = "SUITE"
)

type ServiceInstanceLoadBalancingPolicy string

const (
	ServiceInstanceLoadBalancingPolicyLCC ServiceInstanceLoadBalancingPolicy = "least_connection_count"
	ServiceInstanceLoadBalancingPolicyLRT ServiceInstanceLoadBalancingPolicy = "least_response_time"
	ServiceInstanceLoadBalancingPolicyRR  ServiceInstanceLoadBalancingPolicy = "round_robin"
)

type ServiceInstanceShape string

const (
	// oc3: 1 OCPU, 7.5 GB memory
	ServiceInstanceShapeOC3 ServiceInstanceShape = "oc3"
	// oc4: 2 OCPUs, 15 GB memory
	ServiceInstanceShapeOC4 ServiceInstanceShape = "oc4"
	// oc5: 4 OCPUs, 30 GB memory
	ServiceInstanceShapeOC5 ServiceInstanceShape = "oc5"
	// oc6: 8 OCPUs, 60 GB memory
	ServiceInstanceShapeOC6 ServiceInstanceShape = "oc6"
	// oc7: 16 OCPUS, 120 GB memory
	ServiceInstanceShapeOC7 ServiceInstanceShape = "oc7"
	// oc1m: 1 OCPU, 15 GB memory
	ServiceInstanceShapeOC1M ServiceInstanceShape = "oc1m"
	// oc2m: 2 OCPUs, 30 GB memory
	ServiceInstanceShapeOC2M ServiceInstanceShape = "oc2m"
	// oc3m: 4 OCPUs, 60 GB memory
	ServiceInstanceShapeOC3M ServiceInstanceShape = "oc3m"
	// oc4m: 8 OCPUs, 120 GB memory
	ServiceInstanceShapeOC4M ServiceInstanceShape = "oc4m"
	// oc5m: 16 OCPUS, 240 GB memory
	ServiceInstanceShapeOC5M ServiceInstanceShape = "oc5m"
)

type ServiceInstanceScalingUnitName string

const (
	ServiceInstanceScalingUnitNameBasic  ServiceInstanceScalingUnitName = "BASIC"
	ServiceInstanceScalingUnitNameSmall  ServiceInstanceScalingUnitName = "SMALL"
	ServiceInstanceScalingUnitNameMedium ServiceInstanceScalingUnitName = "MEDIUM"
	ServiceInstanceScalingUnitNameLarge  ServiceInstanceScalingUnitName = "LARGE"
)

type ServiceInstanceType string

const (
	ServiceInstanceTypeWebLogic ServiceInstanceType = "weblogic"
	ServiceInstanceTypeDataGrid ServiceInstanceType = "datagrid"
	ServiceInstanceTypeOTD      ServiceInstanceType = "otd"
)

type ServiceInstanceUpperStackProductName string

const (
	ServiceInstanceUpperStackProductNameODI ServiceInstanceUpperStackProductName = "ODI"
	ServiceInstanceUpperStackProductNameWCP ServiceInstanceUpperStackProductName = "WCP"
)

type ServiceInstanceVersion string

const (
	ServiceInstanceVersion1221 ServiceInstanceVersion = "12.2.1"
	ServiceInstanceVersion1213 ServiceInstanceVersion = "12.1.3"
	ServiceInstanceVersion1036 ServiceInstanceVersion = "10.3.6"
)

type ServiceInstanceSubscriptionType string

const (
	ServiceInstanceSubscriptionTypeHourly  ServiceInstanceSubscriptionType = "HOURLY"
	ServiceInstanceSubscriptionTypeMonthly ServiceInstanceSubscriptionType = "MONTHLY"
)

type ServiceInstanceScalingUnitInstanceStatus string

const (
	ServiceInstanceScalingUnitInstanceStatusReady    ServiceInstanceScalingUnitInstanceStatus = "Ready"
	ServiceInstanceScalingUnitInstanceStatusStarting ServiceInstanceScalingUnitInstanceStatus = "Starting"
	ServiceInstanceScalingUnitInstanceStatusStopping ServiceInstanceScalingUnitInstanceStatus = "Stopping"
	ServiceInstanceScalingUnitInstanceStatusError    ServiceInstanceScalingUnitInstanceStatus = "Error"
)

type ServiceInstanceServiceComponentType string

const (
	ServiceInstanceServiceComponentTypeJDK    ServiceInstanceServiceComponentType = "JDK"
	ServiceInstanceServiceComponentTypeOTD    ServiceInstanceServiceComponentType = "OTD"
	ServiceInstanceServiceComponentTypeOTDJDK ServiceInstanceServiceComponentType = "OTD_JDK"
	ServiceInstanceServiceComponentTypeWLS    ServiceInstanceServiceComponentType = "WLS"
)

type ServiceInstanceServiceComponentVersion string

const (
	ServiceInstanceServiceComponentVersionWLS    ServiceInstanceServiceComponentVersion = "12.1.3.0.5"
	ServiceInstanceServiceComponentVersionOTD    ServiceInstanceServiceComponentVersion = "11.1.1.9.1"
	ServiceInstanceServiceComponentVersionJDK    ServiceInstanceServiceComponentVersion = "1.7.0_91"
	ServiceInstanceServiceComponentVersionOTDJDK ServiceInstanceServiceComponentVersion = "1.7.0_91"
)

type ServiceInstanceShiftStatus string

const (
	ServiceInstanceShiftStatusReady     ServiceInstanceShiftStatus = "readyToShift"
	ServiceInstanceShiftStatusCompleted ServiceInstanceShiftStatus = "shiftCompleted"
	ServiceInstanceShiftStatusFailed    ServiceInstanceShiftStatus = "shiftFailed"
)

type ServiceInstanceStatus string

const (
	// 	Failed: the service instance has failed.
	ServiceInstanceFailed ServiceInstanceStatus = "Failed"
	//	In Progress: the service instance is being created.
	ServiceInstanceInProgress ServiceInstanceStatus = "In Progress"
	//	Maintenance: the service instance is being stopped, started, restarted or scaled.
	ServiceInstanceMaintenance ServiceInstanceStatus = "Maintenance"
	//	Running: the service instance is running.
	ServiceInstanceRunning ServiceInstanceStatus = "Running"
	//	Stopped: the service instance is stopped.
	ServiceInstanceStopped ServiceInstanceStatus = "Stopped"
	//	Terminating: the service instance is being deleted.
	ServiceInstanceTerminating ServiceInstanceStatus = "Terminating"
)

type ServiceInstanceMiddlewareVersion string

const (
	ServiceInstanceMiddlewareVersion12c212 ServiceInstanceMiddlewareVersion = "12cRelease212"
	ServiceInstanceMiddlewareVersion12c2   ServiceInstanceMiddlewareVersion = "12cRelease2"
	ServiceInstanceMiddlewareVersion12cR3  ServiceInstanceMiddlewareVersion = "12cR3"
	ServiceInstanceMiddlewareVersion11g    ServiceInstanceMiddlewareVersion = "11g"
)

type ServiceInstance struct {
	// Flag that specifies whether updates to the Oracle Cloud Tools are automatically
	// applied to the Oracle Java Cloud Service instance during the maintenance window.
	// The Oracle Cloud Tools are used to manage the lifecycle of your service instance.
	AutoUpdate bool `json:"auto_update"`
	// Name of the cluster that contains the Managed Servers for the service instance.
	ClusterName string `json:"cluster_name"`
	// Status indicating whether the version of Oracle Cloud Tools is out of compliance.
	// Oracle uses Oracle Cloud Tools to manage the lifecycle of your service instance.
	// If the Oracle Cloud Tools are out of compliance, this attribute is set to one of
	// the following status values:
	// NEW_VERSION: Indicates a new version of Oracle Cloud Tools is available.
	// Oracle strongly recommends that you apply this update as soon as possible.
	// DEPRECATED: Indicates that the current version of Oracle Cloud Tools is deprecated.
	// Apply the latest Oracle Cloud Tools update to avoid any disruption of service in
	// the future.
	// UNSUPPORTED: Indicates that the current version of Oracle Cloud Tools is not supported.
	// Apply the latest patch to resume normal operations.
	// If the Oracle Cloud Tools are up-to-date, this attribute is blank.
	ComplianceStatus string `json:"compliance_status"`
	// Description that provides more details about the compliance status of the
	// Oracle Cloud Tools, used to manage the lifecycle of your Oracle Java Cloud Service
	// instance. If the service instance is out of compliance, this attribute is set to one
	// of the following descriptions, based on the status value:If the Oracle Cloud Tools
	// are out of compliance, this attribute is set to one of the following status values:
	// NEW_VERSION: A newer version of Oracle tools latestVersion is available.
	// This update includes critical fixes to Oracle Cloud Tools. Oracle uses cloud tools
	// to manage lifecycle of your service. Oracle strongly recommends that customers apply
	// this update as soon as possible.
	// DEPRECATED: This service is currently in a deprecated state because the Oracle tools
	// version deprecatedVersion is deprecated. Apply the latest Oracle tools update as this
	// version may not be supported in the future.
	// UNSUPPORTED: This service is currently in an unsupported state because the Oracle
	// tools version unsupportedVersion is obsolete. Apply the latest patch to resume normal
	// operations.
	// If the Oracle Cloud Tools are up-to-date, this attribute is blank.
	ComplianceStatusDescription string `json:"compliance_status_description"`
	// Location where the service instance is provisioned.
	ComputeSiteName string `json:"compute_site_name"`
	// Resource URL for accessing the deployed applications using HTTP.
	ContentURL string `json:"content_url"`
	// Name of the user account used to create the Oracle Java Cloud Service instance.
	CreatedBy string `json:"created_by"`
	// Job ID for the create job.
	CreationJobID string `json:"creation_job_id"`
	// Date and time the Oracle Java Cloud Service instance was created.
	CreationTime string `json:"creation_time"`
	// Groups details of Database Cloud Service database deployments and databases used.
	DBAssociations []DBAssociation `json:"db_associations"`
	// Database that is used to host the Oracle Required Schema.
	DBInfo string `json:"db_info"`
	// Name of the Database Cloud Service database deployment that is used to host
	// the Oracle Required Schema.
	DBServiceName string `json:"db_service_name"`
	// Resource URL for the Oracle Database Cloud Service database deployment for this
	// service instance.
	DBServiceURI string `json:"db_service_uri"`
	// Job ID for the delete job.
	DeletionJobID int `json:"deletion_job_id"`
	// Free-form text that provides additional information about the service instance.
	Description string `json:"description"`
	// Mode of the domain. Valid values include: DEVELOPMENT and PRODUCTION.
	DomainMode string `json:"domain_mode"`
	// Name of the WebLogic domain.
	DomainName string `json:"domain_name"`
	// Software edition. Valid values include: SE, EE, or SUITE.
	Edition string `json:"edition"`
	// Error status that describes the reason that the Oracle Java Cloud Service instance
	// is in an erroneous state. The following provides an example of the description:
	// This service is currently in an erroneous state as the tools are in an inconsistent
	// state. Reason - error details. Apply the latest tools patch to resume normal operations
	// on this service.
	ErrorStatusDesc string `json:"error_status_desc"`
	// URL to Enterprise Manager Fusion Middleware Control.
	FMWControlURL string `json:"fmw_control_url"`
	// Identity domain ID for the Oracle Java Cloud Service account (on Oracle Public Cloud).
	// Tenant name for the Oracle Java Cloud Service instance (on Oracle Cloud Machine).
	IdentityDomain string `json:"identity_domain"`
	// Groups one or more IP reservations in use on this service instance.
	// This attribute is only applicable to accounts where regions are supported.
	IPReservations []IPReservation `json:"ip_reservations"`
	// This attribute is only applicable to accounts where regions are supported.
	// The three-part name of an IP network to which the service instance is attached. For example: /Compute-identity_domain/user/object
	IPNetwork string `json:"ipNetwork"`
	// Flag that specifies whether this service instance is created with AppToCloud artifacts.
	// This attribute is displayed only if the service instance was created with AppToCloud artifacts.
	IsApp2Cloud bool `json:"isApp2Cloud"`
	// Date and time the Oracle Java Cloud Service instance was last modified.
	LastModifiedTime string `json:"last_modified_time"`
	// Service level. Valid values include:
	// PAAS: Production-level service. Supports Oracle Java Cloud Service instance
	// creation and monitoring; backup and restoration; patching; and scaling.
	// This is the default.
	// BASIC: Developer-level service. Supports Oracle Java Cloud Service instance
	// creation and monitoring. Note: This service level does not support backup and
	// restoration, patching, or scaling.
	Level ServiceInstanceLevel `json:"level"`
	// Job ID of a lifecycle control request. Is a Long value. This attribute appears only
	// if a lifecycle control request is in progress. You can use this ID to check the status
	// of the lifecycle control request (for example, a start/stop/restart operation)
	LifecycleControlJobID int `json:"lifecycle_control_job_id"`
	// Total amount of memory in GBs allocated across all nodes in the service instance.
	MemorySize int `json:"memory_size"`
	// Total number of public IP addresses reserved for the Oracle Java Cloud Service instance.
	NumIPReservations int `json:"num_ip_reservations"`
	// Number of Managed Servers in the domain.
	NumNodes int `json:"num_nodes"`
	// Total number of Oracle Compute Units (OCPUs) allocated across all nodes in the service instance.
	OCPUCount int `json"ocpu_count"`
	// Groups information about the Coherence data tier.
	Options []Option `json:"options"`
	// URL to load balancer Administration Console.
	OTDAdminURL string `json:"otd_admin_url"`
	// Flag that specifies whether the load balancer is enabled.
	OTDProvisioned string `json:"otd_provisioned"`
	// Desired compute shape for the load balancer.
	OTDShape string `json:"otd_shape"`
	// Storage size of the load balancer in GBs.
	OTDStorageSize int `json:"otd_storage_size"`
	// Version of the PaaS Service Manager.
	PSMPluginVersion string `json:"psm_plugin_version"`
	// This attribute is only applicable to accounts where regions are supported.
	// Location where the service instance is provisioned.
	Region string `json:"region"`
	// URL for accessing the sample application, if it was installed and deployed when
	// the service instance was provisioned.
	SampleAppURL string `json:"sample_app_url"`
	// URL for accessing the deployed applications using HTTPS.
	SecureContentURL string `json:"secure_content_url"`
	// Groups service component details.
	ServiceComponents []ServiceComponent `json:"service_components"`
	// Name of Oracle Java Cloud Service instance.
	ServiceName string `json:"service_name"`
	// Resource URL for the Oracle Java Cloud Service instance.
	ServiceURI string `json:"service_uri"`
	// Desired compute shape.
	Shape string `json:"shape"`
	// Job ID of the customPayload import operation with regards to AppToCloud migration.
	// This attribute is displayed only if the service instance was created with AppToCloud artifacts.
	ShiftJobID string `json:"shiftJobId"`
	// Status of the service instance with regards to AppToCloud migration. Possible values:
	// readyToShift, shiftCompleted, shiftFailed.
	// This attribute is displayed only if the service instance was created with AppToCloud artifacts.
	ShiftStatus ServiceInstanceShiftStatus `json:"shiftStatus"`
	// Flag that specifies the status of the Oracle Java Cloud Service instance.
	// Valid values include: Running, In Progress, Maintenance, Stopped, Terminating, and Failed.
	Status ServiceInstanceStatus `json:"status"`
	// Total amount of block storage in GBs allocated across all nodes in the service instance.
	StorageSize int `json:"storage_size"`
	// Billing frequency. Valid values include:
	// HOURLY: Pay only for the number of hours used during your billing period.
	// MONTHLY: Pay one price for the full month irrespective of the number of hours used.
	SubscriptionType ServiceInstanceSubscriptionType `json:"subscription_type"`
	// This attribute is not available on Oracle Cloud Machine.
	// The Oracle Fusion Middleware product installer added to this service instance.
	// For example: WCP
	UpperStackProductName string `json:"upper_stack_product_name"`
	// Oracle Fusion Middleware software version.
	// Valid values include: 12cRelease212, 12cRelease2, 12cR3 and 11g.
	Version ServiceInstanceMiddlewareVersion `json:"version"`
	// URL to the WebLogic Administration Console.
	WLSAdminURL string `json:"wls_admin_url"`
	// Port for accessing the Administration Server using WLST.
	WLSDeploymentChannelPort int `json:"wls_deployment_channel_port"`
	// Oracle WebLogic Server software version. For example: 12.2.1.2.0, 12.2.1.0.x, 12.1.3.0.x and 10.3.6.0.x.
	WLSVersion string `json:"wlsVersion"`
}

type DBAssociation struct {
	// The URL to use to connect to Oracle Application Express on the service instance.
	DBApexURL string `json:"db_apex_url"`
	// Flag that specifies whether this database is used to host application schemas (true).
	DBApp bool `json:"db_app"`
	// Information about the Database Cloud Service database deployment or the database
	// connection string.
	DBConnectString string `json:"db_connect_string"`
	// The URL to use to connect to Enterprise Manager on the service instance.
	DBEmURL string `json:"db_em_url"`
	// Flag that specifies whether this database is used to host the Oracle Required
	// Schema (true).
	DBInfra bool `json:"db_infra"`
	// The URL to use to connect to Oracle DBaaS Monitor on the service instance.
	DBMonitorURL string `json:"db_monitor_url"`
	// Service level for the Database Cloud Service database deployment.
	DBServiceLevel string `json:"db_service_level"`
	// Name of the Database Cloud Service database deployment.
	DBServiceName string `json:"db_service_name"`
	// Version of the Oracle database.
	DBVersion string `json:"db_version"`
	// Name of the pluggable database created for Oracle Database 12c.
	// This value does not apply to a Database Cloud Service database deployment that
	// is running Oracle Database 11g.
	PDBServiceName string `json:"pdb_service_name"`
}

type IPReservation struct {
	// Name of an IP reservation that is assigned to a node on the service instance.
	Name string `json:"name"`
}

type Option struct {
	// Groups the information about the datagrid cluster.
	Clusters []Cluster `json:"clusters"`
	// Specifies the Coherence data tier cluster.
	// This value is set to datagrid.
	Type string `json:"type"`
}

type Cluster struct {
	// Name of a Coherence data tier cluster for a service instance.
	ClusterName string `json:"clusterName"`
	// Additional heap available when a capacity unit is added as a result of a scale out
	// operation.
	HeapIncrements string `json:"heapIncrements"`
	// Heap size to configure per JVM in a capacity unit.
	HeapSize string `json:"heapSize"`
	// Number of JVMs or Managed Servers to configure per VM in a capacity unit.
	JVMCount int `json:"jvmCount"`
	// Maximum JVM heap available, based on the site.
	MaxHeap string `json:"maxHeaps"`
	// Maximum primary cache that can be provided, based on the site.
	MaxPrimary string `json:"maxPrimary"`
	// Maximum number of capacity units that can be provisioned, based on the site.
	MaxScalingUnit int `json:"maxScalingUnit"`
	// Primary cache storage that can be added when a capacity unit is added as a
	// result of a scale out operation.
	PrimaryIncrements string `json:"primaryIncrements"`
	// Number of capacity units provisioned for the service instance.
	ScalingUnitCount int `json:"scalingUnitCount"`
	// Groups information about the capacity units provisioned in the service instance.
	ScalingUnitInstances []ScalingUnitInstance `json:"scalingUnitInstances"`
	// Name of a default capacity unit, if used for the service instance.
	ScalingUnitName string `json:"scalingUnitName"`
	// Shape of the virtual machines provisioned by a capacity unit.
	Shape string `json:"shape"`
	// Total heap available, based on the number of JVMs configured per capacity unit.
	TotalHeap string `json:"totalHeap"`
	// Total primary cache storage to allocate for Coherence, based on the general
	// rule of splitting the JVM heap size into thirds, using 1/3 for primary cache storage,
	// 1/3 for backup storage, and 1/3 for scratch space.
	TotalPrimary string `json:"totalPrimary"`
	// Number of virtual machines configured per capacity unit.
	VMCount int `json:"vmCount"`
}

type ScalingUnitInstance struct {
	// Unique ID for managing a capacity unit.
	ScalingUnitInstanceId int `json:"scalingUnitInstanceId"`
	// Groups information about the Managed Servers provisioned by a capacity unit.
	Servers []Server `json:"servers"`
	// Status of the capacity unit. Valid values are:
	// Ready: Fully operational
	// Starting: Being created or initialized
	// Stopping: Being removed
	// Error: Has some error condition(s)
	Status ServiceInstanceScalingUnitInstanceStatus `json:"status"`
}

type Server struct {
	// Name of a Managed Server on the Coherence data tier.
	Name string `json:"name"`
}

type ServiceComponent struct {
	// Service component type. Valid values are JDK, OTD, OTD_JDK, or WLS.
	Type ServiceInstanceServiceComponentType `json:"type"`
	// Software version of the specified component.
	// For example, 12.1.3.0.5 for WLS, 11.1.1.9.1 for OTD, 1.7.0_91 for OTD_JDK, or 1.7.0_91 for JDK.
	Version ServiceInstanceServiceComponentVersion `json:"version"`
}

type CreateServiceInstanceInput struct {
	// This attribute is not available on Oracle Cloud Machine.
	// This attribute is only applicable when level is set to PAAS. Specifies whether to
	// enable backups for this Oracle Java Cloud Service instance.
	// Optional.
	BackupDestination ServiceInstanceBackupDestination `json:"backupDestination,omitempty"`
	// Where to store your service instance backups.
	// Note the difference between Oracle Public Cloud and Oracle Cloud Machine.
	// On Oracle Public Cloud, this is the name of a Oracle Storage Cloud Service container. For
	// example, you can specify:
	// Storage-<identitydomainid>/<containername>
	// <storageservicename>-<identitydomainid>/<containername>
	// https://foo.storage.oraclecloud.com/v1/MyService-bar/MyContainer
	// The format to use to specify the Storage container name depends on the URL of your Oracle
	// Storage Cloud Service account. To identify the URL of your storage account, see About REST URLs for Oracle Storage Cloud Service Resources in Using Oracle Storage Cloud Service.
	// Note:
	// Do not use an Oracle Storage Cloud container that you use to back up Oracle Java Cloud Service
	// instances for any other purpose. For example, do not also use the same container to back up Oracle
	// Database Cloud Service database deployments. Using one container for multiple purposes can result in
	// billing errors.
	// You do not have to specify a storage container if you provision the service instance without enabling
	// backups.
	// On Oracle Cloud Machine, this is the NFS URI of the remote storage disk used to store service instance
	// backups. Get from your Cloud or Tenant administrator the filer ip address and export path that is
	// designated to store Oracle Java Cloud Service instance backups. Specify the NFS URI in the format:
	// ip-address:export-path
	// This is not required when provisioning an Oracle Java Cloud Service - Virtual Image instance (BASIC level).
	// Optional.
	CloudStorageContainer string `json:"cloudStorageContainer,omitempty"`
	// Password for the Oracle Storage Cloud Service administrator.
	// Must be specified if cloudStorageContainer is set.
	CloudStoragePassword string `json:"cloudStoragePassword,omitempty"`
	// Username for the Oracle Storage Cloud Service administrator.
	// Must be specified if cloudStorageContainer is set.
	CloudStorageUsername string `json:"cloudStorageUser,omitempty"`
	// Specify if the given cloudStorageContainer is to be created if it does not already exist.
	// Default value is false.
	// Optional.
	CreateStorageContainerIfMissing bool `json:"createStorageContainerIfMissing,omitempty"`
	// Free-form text that provides additional information about the service instance.
	// Optional.
	Description string `json:"description,omitempty"`
	// This attribute is not available on Oracle Cloud Machine.
	// Flag that specifies whether to enable (true) or disable (false) the access rules that control external
	// communication to the WebLogic Server Administration Console, Fusion Middleware Control, and Load Balancer Console.
	// If you do not set it to true, after the service instance is created, you have to explicitly enable the rules
	// for the administration consoles before you can gain access to them.
	// The default value is false.
	// Optional
	EnableAdminConsole bool `json:"enableAdminConsole,omitempty"`
	// This attribute is not available on Oracle Cloud Machine.
	// This attribute is only applicable to accounts where regions are supported.
	// The three-part name of a custom IP network to attach this service instance to. For example:
	// /Compute-identity_domain/user/object
	// A region name must be specified in order to use ipNetwork. Only those IP networks created in the specified region can be used.
	// If using an IP network, note that the dbServiceName for the service instance should be
	// attached to the same ipNetwork. If your Oracle Java Cloud Service and Oracle Database Cloud Service are attached to
	// different IP networks, then the two IP networks must be connected to the same IP network exchange.
	// Access rules required for the communication between the Oracle Java Cloud Service instance and Oracle Database Cloud Service
	// database deployment are created automatically.
	// See Creating an IP Network in Using Oracle Compute Cloud Service (IaaS).
	// Optional.
	IPNetwork string `json:"ipNetwork,omitempty"`
	// Service level for the service instance
	// The default is PAAS
	// Optional.
	Level ServiceInstanceLevel `json:"level,omitempty"`
	// Groups component-specific attributes in the following categories:
	// WebLogic Server ("type":"weblogic")
	// Oracle Traffic Director ("type":"otd")
	// Oracle Coherence ("type":"datagrid")
	// Required.
	Parameters []Parameter `json:"parameters"`
	// Flag that specifies whether to enable the load balancer.
	// The default value is true when you configure more than one Managed Server for the Oracle
	// Java Cloud Service instance. Otherwise, the default value is false
	// Optional.
	ProvisionOTD bool `json:"provisionOTD,omitempty"`
	// This attribute is only available on Oracle Cloud Machine.
	// Path to the network from which the Oracle Java Cloud Service REST API will be accessed.
	// You can connect to an external network via the InfiniBand-to-10 GB Ethernet gateways
	// using Ethernet over InfiniBand (EoIB).
	// Optional.
	PublicNetwork string `json:"publicNetwork,omitempty"`
	// This attribute is not available on Oracle Cloud Machine.
	// This attribute is only applicable to accounts where regions are supported.
	// Name of the region where the Oracle Java Cloud Service instance is to be provisioned.
	// If no region name is specified, the service instance is provisioned in the same compute site as the site of the Oracle Database Cloud Service database deployment specified in dbServiceName.
	// If a region name is specified, note that the dbServiceName for this service instance must be one that is provisioned in the same region.
	// A region name must be specified if you want to use ipReservations or ipNetwork.
	// Optional.
	Region string `json:"region,omitempty"`
	// Flag that specifies whether to automatically deploy and start the sample application,
	// sample-app.war, to the Managed Server in your service instance.
	// The default value is false.
	// Optional.
	SampleAppDeploymentRequested bool `json:"sampleAppDeploymentRequested,omitempty"`
	// Name of Oracle Java Cloud Service instance. The service name:
	// Must not exceed 50 characters.
	// Must start with a letter.
	// Must contain only letters, numbers, or hyphens.
	// Must not contain any other special characters.
	// Must be unique within the identity domain.
	// By default, the names of the domain and cluster in the service instance will be
	// generated from the first eight characters of the service instance name (serviceName),
	// using the following formats, respectively:
	// first8charsOfServiceInstanceName_domain
	// first8charsOfServiceInstanceName_cluster
	// Required.
	ServiceName string `json:"serviceName"`
	// Metering frequency. Valid values include:
	// HOURLY - Pay only for the number of hours used during your billing period. This is the default.
	// MONTHLY - Pay one price for the full month irrespective of the number of hours used.
	// Required.
	SubscriptionType ServiceInstanceSubscriptionType `json:"subscriptionType"`
}

type Parameter struct {
	// Password for WebLogic Server or Oracle Traffic Director administrator. The password must
	// meet the following requirements:
	// Starts with a letter
	// Is between 8 and 30 characters long
	// Has one or more upper case letters
	// Has one or more lower case letters
	// Has one or more numbers
	// Has one or more of the following special characters: hyphen (-), underscore (_),
	// pound sign (#), dollar sign ($). If Exadata is the database for the service instance,
	// the password cannot contain the dollar sign ($).
	// If an administrator password is not explicitly set for Oracle Traffic Director (otd),
	// the OTD administrator password defaults to the WebLogic Server administrator password.
	// Note: This attribute is valid when component type is set to weblogic or otd only;
	// it is not valid for datagrid.
	// Optional.
	AdminPassword string `json:"adminPassword,omitempty"`
	// Port for accessing the WebLogic Server or Oracle Traffic Director using HTTP. The default values are:
	// 7001 for WebLogic Server
	// 8989 for Oracle Traffic Director
	// The adminPort, contentPort, securedAdminPort, securedContentPort, and nodeManagerPort
	// values must be unique.
	// Note: This attribute is valid when component type is set to weblogic or otd only;
	// it is not valid for datagrid.
	// Optional.
	AdminPort int `json:"adminPort,omitempty"`
	// User name for the WebLogic Server or Oracle Traffic Director administrator.
	// The name must be between 8 and 128 characters long and cannot contain any of the following characters:
	// Tab
	// Brackets
	// Parentheses
	// The following special characters: left angle bracket (<), right angle bracket (>),
	// ampersand (&), pound sign (#), pipe symbol (|), and question mark (?).
	// If a username is not explicitly set for Oracle Traffic Director (otd), the OTD user name
	// defaults to the WebLogic Server administrator user name.
	// Note: This attribute is valid when component type is set to weblogic or otd only;
	// it is not valid for datagrid.
	// Optional.
	AdminUsername string `json:"adminUserName,omitempty"`
	// This attribute is not available on Oracle Cloud Machine.
	// Note: This attribute is valid when component type is set to weblogic only;
	//  it is not valid for otd or datagrid.
	// Groups details of Database Cloud Service database deployments that host application schemas, if used.
	// Optional.
	AppDBs []AppDB `json:"addDbs"`
	// Size of the backup volume for the service. The value must be a multiple of GBs.
	// You can specify this value in bytes or GBs. If specified in GBs, use the following format:
	// nG, where n specifies the number of GBs. For example, you can express 10 GBs as bytes or GBs.
	// For example: 100000000000 or 10G. This value defaults to the system configured volume size.
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// Optional.
	BackupVolumeSize string `json:"backupVolumeSize,omitempty"`
	// For WebLogic Server, specifies the name of the cluster that contains the Managed Servers
	// for the service instance. By default, the name of the cluster will be generated from the
	// first eight characters of the Oracle Java Cloud Service instance name (serviceName), using
	// the following format: first8charsOfServiceInstanceName_cluster
	// For Coherence, specifies the name of the storage-enabled WebLogic Server cluster to add
	//  for the instance. If the cluster name is empty or null, a name is generated from the first
	// eight characters of the Oracle Java Cloud Service instance name using the following format:
	// first8charsOfServiceInstanceName_DGCluster.
	// The cluster name:
	// Must not exceed 50 characters.
	// Must start with a letter.
	// Must contain only alphabetical characters, underscores (_), or dashes (-).
	// Must not contain any other special characters.
	// Must be unique within the identity domain.
	// Note: This attribute is valid when component type is set to weblogic or datagrid only;
	// it is not valid for otd.
	// Optional.
	ClusterName string `json:"clusterName,omitempty"`
	// Connection string for the database. The connection string must be entered using one of
	// the following formats:
	// host:port:SID
	// host:port/serviceName
	// For example, foo.bar.com:1521:orcl or foo.bar.com:1521/mydbservice
	// Note the difference between Oracle Public Cloud and Oracle Cloud Machine.
	// On Oracle Public Cloud, this attribute is required only when you specify a
	// Virtual Image service level of Database Cloud Service in dbServiceName. It is used to
	//  connect to the database deployment on Database Cloud Service - Virtual Image.
	// On Oracle Cloud machine, this is the string that is used to connect to the database.
	// The database can be either an on-premises database or a Database Cloud Service database deployment.
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// Optional.
	ConnectString string `json:"connectString,omitempty"`
	// Port for accessing the deployed applications using HTTP. The default value is 8001.
	// Note: This value is overridden by privilegedContentPort unless its value is set to 0.
	//  This value has no effect if the load balancer is enabled.
	// The adminPort, contentPort, securedAdminPort, securedContentPort, and nodeManagerPort
	// values must be unique.
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// Optional.
	ContentPort int `json:"contentPort,omitempty"`
	// User name for the database administrator.
	// For service instances based on Oracle WebLogic Server 11g (10.3.6), this value must
	// be set to a database user with DBA role. You can use the default user SYSTEM or a user
	// that has been granted the DBA role.
	// For service instances based on Oracle WebLogic Server 12c (12.2.1 and 12.1.3), this value
	// must be set to a database user with SYSDBA system privileges. You can use the default user
	// SYS or a user that has been granted the SYSDBA privilege.
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// Required.
	DBAName string `json:"dbaName"`
	// The Database administrator password that was specified when the database deployment
	// on Database Cloud Service was created or the password for the database administrator.
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// Required.
	DBAPassword string `json:"dbaPassword"`
	// Path to the network through which the Oracle Java Cloud Service instance will access the database.
	// You can connect to the database network using one of the following options:
	// Connect to Exadata on the same InfiniBand fabric (IPoIB). In this case, Exalogic
	// machines use a unified 32 GB per second InfiniBand quad data rate (QDR) fabric for
	// internal communication. Exalogic machines communicate with Oracle Exadata Database Machines
	// for database connectivity via IPoIB.
	// Connect to a database via an Ethernet over InfiniBand (EoIB) network. In this case,
	// Exalogic machines can be connected to an external network, including a standard database
	// hosted on a machine outside of the Exalogic machine, via the InfiniBand-to-10 GB Ethernet
	// gateways using Ethernet over InfiniBand (EoIB).
	// This attribute is only available on Oracle Cloud Machine.
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// Optional.
	DBNetwork string `json:"dbNetwork,omitempty"`
	// Name of the database deployment on Oracle Database Cloud Service to host the
	// Oracle Infrastructure schemas required for this Oracle Java Cloud Service instance.
	// If provisioning a service instance in a specific region, specify a Database Cloud Service
	// database deployment that is in the same region.
	// The specified database deployment must be running. Only an Oracle Java Cloud Service instance
	// based on WebLogic Server 12.2.1 can use a required schema database deployment that is created
	// using the Oracle Database 12.2 version.
	// When provisioning a production-level Oracle Java Cloud Service instance, you must use a
	// production-level Database Cloud Service. On Oracle Public Cloud, you can specify a Virtual
	// Image service level of Database Cloud Service if you are provisioning an Oracle Java Cloud
	// Service - Virtual Image instance. If you specify a Virtual Image service level of Database
	// Cloud Service, you must also specify its connection string using the connectString attribute.
	// See Oracle Database Cloud Service Database Deployment in Using Oracle Java Cloud Service for
	// the backup options that you can use when you create a database deployment on Database Cloud Service.
	// Note: To ensure that you can restore the database for an Oracle Java Cloud Service instance
	// without risking data loss for other service instances, do not use the same Database Cloud Service
	// database deployment with multiple Oracle Java Cloud Service instances.
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// Required.
	DBServiceName string `json:"dbServiceName"`
	// Port for accessing the Administration Server using WLST.
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// The default value is 9001.
	// Optional.
	DeploymentChannelPort int `json:"deploymentChannelPort,omitempty"`
	// Mode of the domain. Valid values include: DEVELOPMENT and PRODUCTION. The default value is PRODUCTION.
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// Optional.
	DomainMode ServiceInstanceDomainMode `json:"domainMode,omitempty"`
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// Name of the WebLogic domain. By default, the domain name will be generated from the first
	// eight characters of the Oracle Java Cloud Service instance name (serviceName), using the
	// following format: first8charsOfServiceInstanceName_domain
	// By default, the Managed Server names will be generated from the first eight characters of
	// the domain name name (domainName), using the following format: first8charsOfDomainName_server_n,
	// where n starts with 1 and is incremented by 1 for each additional Managed Server to ensure each
	// name is unique.
	// Optional.
	DomainName string `json:"domainName,omitempty"`
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// Number of partitions to enable in the domain for WebLogic Server 12.2.1.
	// Valid values include: 0 (no partitions), 1, 2, and 4.
	// Optional.
	DomainPartitionCount int `json:"domainPartitionCount,omitempty"`
	// Note: This attribute is valid when component type is set to weblogic only; it is
	// not valid for otd or datagrid.
	// Size of the domain volume for the service. The value must be a multiple of GBs.
	// You can specify this value in bytes or GBs. If specified in GBs, use the following format:
	// nG, where n specifies the number of GBs. For example, you can express 10 GBs as bytes or GBs.
	// For example: 100000000000 or 10G.
	// This value defaults to the system configured volume size.
	// Optional.
	DomainVolumeSize string `json:"domainVolumeSize,omitempty"`
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// Software edition for WebLogic Server. Valid values include:
	// SE - Standard edition. See Oracle WebLogic Server Standard Edition.
	//  Do not use the Standard edition if you are going to enable domain partitions using
	// WebLogic Server 12.2.1. Do not use the Standard edition if you are also using
	// upperStackProductName to provision a service instance for an Oracle Fusion Middleware product.
	// EE - Enterprise Edition. This is the default for both PAAS and BASIC service levels.
	// See Oracle WebLogic Server Enterprise Edition.
	// SUITE - Suite edition. See Oracle WebLogic Suite.
	// When creating an instance that has Oracle Coherence enabled, you must set this value to SUITE.
	// Optional.
	Edition ServiceInstanceEdition `json:"edition,omitempty"`
	// Note: This attribute is valid when component type is set to otd only;
	// it is not valid for weblogic or datagrid.
	// Flag that specifies whether load balancer HA is enabled.
	// This value defaults to false (that is, HA is not enabled).
	// Optional
	HAEnabled bool `json:"haEnabled,omitempty"`
	// This attribute is not available on Oracle Cloud Machine.
	// Note: This attribute is valid when component type is set to weblogic or otd;
	// it is not valid for datagrid.
	// A single IP reservation name or multiple names separated by commas.
	// Reserved or pre-allocated IP addresses can be assigned to WebLogic Managed Server nodes
	// and load balancer nodes (if OTD is enabled).
	// For weblogic type, all Managed Servers in the cluster must be provisioned with pre-allocated
	// IP reservations, so the number of names in ipReservations must match the managedServerCount
	// in the domain.
	// For otd type, the number of names in ipReservations must match the number of load balancer
	// nodes you are provisioning.
	// Note the difference between accounts where regions are supported and not supported.
	// Where regions are supported: A region name must be specified in order to use ipReservations.
	// Only those reserved IPs created in the specified region can be used.
	// See IP Reservations REST Endpoints for information about how to find unused IP reservations and,
	// if needed, create new IP reservations.
	// Where regions are not supported: When using an Oracle Database Exadata Cloud Service database
	// deployment with your Oracle Java Cloud Service instance in an account where regions are not
	// enabled, a region name is not required in order to use ipReservations. However, you must first
	// submit a request to get the IP reservations. See the My Oracle Support document titled How to
	// Request Authorized IPs for Provisioning a Java Cloud Service with Database Exadata Cloud Service
	// (MOS Note 2163568.1).
	// Optional.
	IPReservations string `json:"ipReservations,omitempty"`
	// Note: This attribute is valid when component type is set to otd only;
	// it is not valid for weblogic or datagrid.
	// Listener port for the load balancer for accessing deployed applications using HTTP.
	// The default value is 8080.
	// Note: This value is overridden by privilegedListenerPort unless its value is set to 0.
	// This value has no effect if the load balancer is disabled.
	// Optional.
	ListenerPort int `json:"listernerPort,omitempty"`
	// Note: This attribute is valid when component type is set to otd only;
	// it is not valid for weblogic or datagrid.
	// Flag that specifies whether the non-secure listener port is enabled on the load balancer.
	// The default value is true.
	// Optional.
	ListenerPortEnabled bool `json:"listenerPortEnabled,omitempty"`
	// Note: This attribute is valid when component type is set to otd only;
	// it is not valid for weblogic or datagrid.
	// Protocol used for the load balancer listener port. The default value is http.
	// Optional.
	ListenerType string `json:"listenerType,omitempty"`
	// Note: This attribute is valid when component type is set to otd only; it is not valid
	// for weblogic or datagrid.
	// Policy to use for routing requests to the load balancer. Valid policies include:
	// least_connection_count - Passes each new request to the Managed Server with the least
	// number of connections. This policy is useful for smoothing distribution when Managed Servers
	// get bogged down. Managed Servers with greater processing power to handle requests will
	// receive more connections over time. This is the default.
	// least_response_time - Passes each new request to the Managed Server with the fastest response time.
	// This policy is useful when Managed Servers are distributed across networks.
	// round_robin - Passes each new request to the next Managed Server in line, evenly distributing
	// requests across all Managed Servers regardless of the number of connections or response time.
	// Optional.
	LoadBalancingPolicy ServiceInstanceLoadBalancingPolicy `json:"loadBalancingPolicy,omitempty"`
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// Number of Managed Servers in the domain. Valid values include: 1, 2, 4, and 8.
	// The default value is 1.
	// Optional.
	ManagedServerCount int `json:"managedServerCount,omitempty"`
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// Initial Java heap size (-Xms) for a Managed Server JVM, specified in megabytes.
	// The value must be greater than -1.
	// If you specify this initial value, a value greater than 0 (zero) must also be specified
	// for msMaxHeapMB, msMaxPermMB, and msPermMB. In addition, msInitialHeapMB must be less
	// than msMaxHeapMB, and msPermMB must be less than msMaxPermMB.
	// Optional,
	MSInitialHeapMB int `json:"msInitialHeapMB,omitempty"`
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// One or more Managed Server JVM arguments separated by a space.
	// You cannot specify any arguments that are related to JVM heap sizes and PermGen
	// spaces (for example, -Xms, -Xmx, -XX:PermSize, and -XX:MaxPermSize).
	// A typical use case would be to set Java system properties using -Dname=value
	// (for example, -Dmyproject.debugDir=/var/myproject/log).
	// You can overwrite or append the default JVM arguments, which are used to start Managed
	// Server processes. See overwriteMsJvmArgs for information on how to overwrite or append
	// the server start arguments.
	// Optional.
	MSJvmArgs string `json:"msJvmArgs,omitempty"`
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// Maximum Java heap size (-Xmx) for a Managed Server JVM, specified in megabytes.
	// The value must be greater than -1.
	// If you specify this maximum value, a value greater than 0 (zero) must also be
	// specified for msInitialHeapMB, msMaxPermMB, and msPermMB. In addition, msInitialHeapMB
	// must be less than msMaxHeapMB, and msPermMB must be less than msMaxPermMB.
	// Optional.
	MSMaxHeapMB int `json:"msMaxHeapMB,omitempty"`
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// Maximum Permanent Generation (PermGen) space in Java heap memory
	// (-XX:MaxPermSize) for a Managed Server JVM, specified in megabytes.
	// The value must be greater than -1.
	// Not applicable for a WebLogic Server 12.2.1 instance, which uses JDK 8.
	// If you specify this maximum value, a value greater than 0 (zero) must also be specified
	// for msInitialHeapMB, msMaxHeapMB, and msPermMB. In addition, msInitialHeapMB must be less
	// than msMaxHeapMB, and msPermMB must be less than msMaxPermMB.
	// Optional.
	MSMaxPermMB int `json:"msMaxPermMB,omitempty"`
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// Initial Permanent Generation (PermGen) space in Java heap memory (-XX:PermSize)
	// for a Managed Server JVM, specified in megabytes. The value must be greater than -1.
	// Not applicable for a WebLogic Server 12.2.1 instance which uses JDK 8.
	// If you specify this initial value, a value greater than 0 (zero) must also be specified
	// for msInitialHeapMB, msMaxHeapMB, and msMaxPermMB. In addition, msInitialHeapMB must be less
	// than msMaxHeapMB, and msPermMB must be less than msMaxPermMB.
	// Optional.
	MSPermMB int `json:"msPermMB,omitempty"`
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// Size of the MW_HOME disk volume for the service (/u01/app/oracle/middleware).
	// The value must be a multiple of GBs. You can specify this value in bytes or GBs.
	// If specified in GBs, use the following format: nG, where n specifies the number of GBs.
	// For example, you can express 10 GBs as bytes or GBs. For example: 100000000000 or 10G.
	//  This value defaults to the system configured volume size.
	// Optional.
	MWVolumeSize string `json:"mwVolumeSize,omitempty"`
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// Password for Node Manager. This value defaults to the WebLogic administrator password
	// (adminPassword) if no value is supplied.
	// Note that the Node Manager password cannot be changed after the Oracle Java Cloud Service
	// instance is provisioned.
	// Optional.
	NodeManagerPassword string `json:"nodeManagerPassword,omitempty"`
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// Port for the Node Manager.
	// Node Manager is a WebLogic Server utility that enables you to start, shut down,
	// and restart Administration Server and Managed Server instances from a remote location.
	// The adminPort, contentPort, securedAdminPort, securedContentPort, and nodeManagerPort
	// values must be unique.
	// The default value is 5556.
	// Optional.
	NodeManagerPort int `json:"nodeManagerPort,omitempty"`
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// User name for Node Manager. This value defaults to the WebLogic administrator user name
	// (adminUserName) if no value is supplied.
	// Optional
	NodeManagerUsername string `json:"nodeManagerUserName,omitempty"`
	// Flag that determines whether the user defined Managed Server JVM arguments specified in msJvmArgs should
	// replace the server start arguments (true), or append the server start arguments (false). Default is false.
	// The server start arguments are calculated automatically by Oracle Java Cloud Service from site default values.
	// If you append (that is, overwriteMsJvmArgs is false or is not set), the user defined arguments specified in
	// msJvmArgs are added to the end of the server start arguments. If you overwrite (that is, set overwriteMsJvmArgs
	// to true), the calculated server start arguments are replaced.
	// Optional.
	OverwriteMsJVMArgs bool `json:"overwriteMsJvmArgs,omitempty"`
	// Note: This attribute is valid when component type is set to weblogic only; it is not valid for otd or datagrid.
	// Name of the pluggable database for Oracle Database 12c. If not specified, the pluggable database name configured
	// when the database was created will be used.
	// Note: This value does not apply to Oracle Database 11g.
	// Optional.
	PDBServiceName string `json:"pdbServiceName,omitempty"`
	// Note: This attribute is valid when component type is set to weblogic only; it is not valid for otd or datagrid.
	// Privileged content port for accessing the deployed applications using HTTP.
	// Note: This value has no effect if the load balancer is enabled.
	// To disable the privileged content port, set the value to 0. In this case, if the load balancer is not
	// provisioned, the content port defaults to contentPort, if specified, or 8001.
	// The default value is 80.
	// Optional.
	PrivilegedContentPort int `json:"privilegedContentPort,omitempty"`
	// Note: This attribute is valid when component type is set to otd only; it is not valid for weblogic or datagrid.
	// Privileged listener port for accessing the deployed applications using HTTP.
	// Note: This value has no effect if the load balancer is disabled.
	// To disable the privileged listener port, set the value to 0. In this case, if the load balancer
	// is provisioned, the listener port defaults to listenerPort, if specified, or 8080.
	// The default value is 80.
	// Optional.
	PrivilegedListenerPort int `json:"privilegedListenerPort,omitempty"`
	// Note: This attribute is valid when component type is set to weblogic only; it is not valid for otd or datagrid.
	// Privileged content port for accessing the deployed applications using HTTPS. The default value is 443.
	// Note: This value has no effect if the load balancer is enabled.
	// To disable the privileged listener port, set the value to 0. In this case, if the load balancer is not provisioned,
	// this value defaults to securedContentPort, if specified, or 8002.
	// The default value is 443.
	// Optional.
	PrivilegedSecuredContentPort int `json:"privilegedSecuredContentPort,omitempty"`
	// Note: This attribute is valid when component type is set to otd only; it is not valid for weblogic or datagrid.
	// Privileged listener port for accessing the deployed applications using HTTPS. The default value is 443.
	// Note: This value has no effect if the load balancer is disabled.
	// To disable the privileged listener port, set the value to 0. In this case, if the load balancer is provisioned,
	// the listener port defaults to securedListenerPort, if specified, or 8081.
	PrivilegedSecuredListenerPort int `json:"privilegedSecuredListenerPort,omitempty"`
	// Note: This attribute is valid when component type is set to datagrid only; it is not valid for weblogic or otd.
	// Required when using a custom capacity unit only. Groups attributes for a custom capacity unit.
	// Optional.
	ScalingUnits []ScalingUnit `json:"scalingUnit,omitempty"`
	// Note: This attribute is valid when component type is set to datagrid only;
	// it is not valid for weblogic or otd.
	// The number of capacity units to add.
	// Each capacity unit provides a fixed amount of primary cache storage to allocate
	// for Coherence, based on the capacity unit's predefined properties for number of
	// VMs, number of JVMs per VM, and heap size for each JVM.
	// This value cannot be 0 (zero).
	// Optional.
	ScalingUnitCount int `json:"scalingUnitCount,omitempty"`
	// Note: This attribute is valid when component type is set to datagrid only;
	// it is not valid for weblogic or otd.
	// Each default capacity unit is one or three VMs with a predefined compute shape
	// (processing power and RAM), running one or more JVMs or Managed Coherence Servers
	// per VM to provide a predefined primary cache capacity. The amount of primary cache
	// storage to allocate for Coherence is based on the general rule of splitting the
	// JVM heap size into thirds: using 1/3 for primary cache storage, 1/3 for backup storage,
	// and 1/3 for scratch space.
	// Required when using a default capacity unit only.
	ScalingUnitName ServiceInstanceScalingUnitName `json:"scalingUnitName,omitempty"`
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// Port for accessing the Administration Server using HTTPS.
	// The adminPort, contentPort, securedAdminPort, securedContentPort, and nodeManagerPort
	// values must be unique.
	// The default value is 7002.
	// Optional.
	SecuredAdminPort int `json:"securedAdminPort,omitempty"`
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// Port for accessing the Administration Server using HTTPS.
	// This value is overridden by privilegedSecuredContentPort unless its value is set to 0.
	// This value has no effect if the load balancer is enabled.
	// The adminPort, contentPort, securedAdminPort, securedContentPort, and nodeManagerPort
	// values must be unique.
	// The default value is 8002.
	// Optional.
	SecuredContentPort int `json:"securedContentPort,omitempty"`
	// Note: This attribute is valid when component type is set to otd only; it is not valid
	// for weblogic or datagrid.
	// Secured listener port for accessing the deployed applications using HTTPS.
	// This value is overridden by privilegedSecuredContentPort unless its value is set to 0.
	// This value has no effect if the load balanced is disabled.
	// The default value is 8081.
	// Optional.
	SecuredListenerPort int `json:"securedListenerPort,omitempty"`
	// Desired compute shape. A shape defines the number of Oracle Compute Units (OCPUs)
	// and amount of memory (RAM).
	// Required.
	Shape ServiceInstanceShape `json:"shape"`
	// Component type to which the set of parameters applies.
	// Valid values include:
	// weblogic - Oracle WebLogic Server
	// datagrid - Oracle Coherence
	// otd - Oracle Traffic Director (load balancer)
	// Required.
	Type ServiceInstanceType `json:"type"`
	// This attribute is not available on Oracle Cloud Machine.
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// The Oracle Fusion Middleware product installer to add to this Oracle Java Cloud Service
	// instance. Valid values are:
	// ODI - Oracle Data Integrator
	// WCP - Oracle WebCenter Portal
	// To use upperStackProductName, you must specify 12.2.1 as the WebLogic Server software
	// version, EE or SUITE as the edition, and PAAS as the service level.
	// After the service instance is provisioned, the specified Fusion Middleware product
	// installer is available in /u01/zips/upperstack on the Administration Server virtual machine.
	// To install the product over the provisioned domain, follow the instructions provided by
	// the Oracle product's installation and configuration documentation.
	// This attribute is required only if you are provisioning an Oracle Java Cloud Service
	// instance for an Oracle Fusion Middleware product.
	// Optional.
	UpperStackProductName ServiceInstanceUpperStackProductName `json:"upperStackProductName,omitempty"`
	// Note: This attribute is valid when component type is set to weblogic only;
	// it is not valid for otd or datagrid.
	// Oracle WebLogic Server software version. Valid values are: 12.2.1, 12.1.3 and 10.3.6.
	// You cannot use 10.3.6 when creating an instance and configuring the Coherence data tier
	// at the same time.
	// Only 12.2.1 is valid if you are using upperStackProductName to provision a service
	// instance for an Oracle Fusion Middleware product.
	// Required.
	Version ServiceInstanceVersion `json:"version"`
	// Note: This attribute is valid when component type is set to weblogic or otd only;
	// it is not valid for datagrid.
	// The public key for the secure shell (SSH). This key will be used for authentication
	// when connecting to the Oracle Java Cloud Service instance using an SSH client.
	// If not specified for otd, this value defaults to the VMsPublicKey value provided for weblogic.
	// Specify only one of the public key attributes for weblogic and otd:
	// VMsPublicKey or VMsPublicKeyName
	// Optional.
	VMsPublicKey string `json:"VMsPublicKey,omitempty"`
	// Note: This attribute is valid when component type is set to weblogic or otd only;
	// it is not valid for datagrid.
	// Name of the compute SSH key object referring to the public key.
	// If not specified for otd, this value defaults to the VMsPublicKeyName value provided
	// for weblogic.
	// Specify only one of the public key attributes for weblogic and otd:
	// VMsPublicKey or VMsPublicKeyName
	// Optional.
	VMsPublicKeyName string `json:"VMsPublicKeyName,omitempty"`
}

type AppDB struct {
	// User name for the database administrator.
	// For service instances based on Oracle WebLogic Server 11g (10.3.6), this value must
	// be set to a database user with DBA role. You can use the default user SYSTEM or a user
	// that has been granted the DBA role.
	// For service instances based on Oracle WebLogic Server 12c (12.2.1 and 12.1.3), this value
	// must be set to a database user with SYSDBA system privileges. You can use the default user
	// SYS or a user that has been granted the SYSDBA privilege.
	// Required.
	DBAName string `json:"dbaName"`
	// Database administrator password that was specified when the database deployment on
	// Database Cloud Service was created.
	// Required.
	DBAPassword string `json:"dbaPassword"`
	// Name of the database deployment on Database Cloud Service to use for an application
	// schema. The specified database deployment must be running.
	// Required.
	DBServiceName string `json:"dbServiceName"`
	// Name of the pluggable database for Oracle Database 12c. If not specified,
	// the pluggable database name configured when the database was created will be used.
	// Note: This value does not apply to Oracle Database 11g.
	// Optional.
	PDBServiceName string `json:"pdbServiceName,omitempty"`
}

type ScalingUnit struct {
	//Heap size to configure with each JVM, based on the memory available from the chosen compute shape.
	// Note that the JVM heap size multiplied by the number of JVMs per VM must not exceed the available memory.
	// Consider using 75% of the memory after reserving 1500 MB for the operating system to calculate the heap size per JVM.
	// The total amount of primary cache storage to allocate for Coherence is based on the general rule of
	// splitting the JVM heap size into thirds: using 1/3 for primary cache storage, 1/3 for backup storage,
	// and 1/3 for scratch space.
	// If a custom capacity unit is configured with a single VM, there might not be space to store a backup
	// copy of the Coherence data, so the actual data available might be more than 1/3 of the JVM heap size.
	// Use a number from 1 GB to 16 GB.
	// Required.
	HeapSize string `json:"heapSize"`
	// Number of JVMs to start on each VM.
	// Use a number from 1 to 8.
	// Required.
	JVMCount int `json:"jvmCount"`
	// Desired compute shape. A shape defines the number of Oracle Compute Units (OCPUs)
	// and amount of memory (RAM).
	// Required.
	Shape ServiceInstanceShape `json:"shape"`
	// Number of VMs to configure for a custom capacity unit.
	// Use a number from 1 to 3. Use 3 VMs to achieve Coherence high availability.
	// Required.
	VMCount int `json:"vmCount"`
}

// CreateServiceInstance creates a new ServiceInstace.
func (c *ServiceInstanceClient) CreateServiceInstance(input *CreateServiceInstanceInput) (*ServiceInstance, error) {
	var (
		serviceInstance      *ServiceInstance
		serviceInstanceError error
	)
	if c.Timeout == 0 {
		c.Timeout = WaitForServiceInstanceReadyTimeout
	}

	// Since these CloudStorageUsername and CloudStoragePassword are sensitive we'll read them
	// from the environment if they aren't passed in.
	if input.CloudStorageContainer != "" && input.CloudStorageUsername == "" && input.CloudStoragePassword == "" {
		input.CloudStorageUsername = *c.ResourceClient.JavaClient.client.UserName
		input.CloudStoragePassword = *c.ResourceClient.JavaClient.client.Password
	}

	for i := 0; i < *c.JavaClient.client.MaxRetries; i++ {
		c.client.DebugLogString(fmt.Sprintf("(Iteration: %d of %d) Creating service instance with name %s\n Input: %+v", i, *c.JavaClient.client.MaxRetries, input.ServiceName, input))

		serviceInstance, serviceInstanceError = c.startServiceInstance(input.ServiceName, input)
		if serviceInstanceError == nil {
			c.client.DebugLogString(fmt.Sprintf("(Iteration: %d of %d) Finished creating service instance with name %s\n Info: %+v", i, *c.JavaClient.client.MaxRetries, input.ServiceName, serviceInstance))
			return serviceInstance, nil
		}
	}
	return nil, serviceInstanceError
}

func (c *ServiceInstanceClient) startServiceInstance(name string, input *CreateServiceInstanceInput) (*ServiceInstance, error) {
	if err := c.createResource(*input, nil); err != nil {
		return nil, err
	}

	// Call wait for instance ready now, as creating the instance is an eventually consistent operation
	getInput := &GetServiceInstanceInput{
		Name: name,
	}

	// Wait for the service instance to be running and return the result
	// Don't have to unqualify any objects, as the GetServiceInstance method will handle that
	serviceInstance, serviceInstanceError := c.WaitForServiceInstanceRunning(getInput, c.Timeout)
	// If the service instance enters an error state we need to delete the instance and retry
	if serviceInstanceError != nil {
		deleteInput := &DeleteServiceInstanceInput{
			Name: name,
		}
		err := c.DeleteServiceInstance(deleteInput)
		if err != nil {
			return nil, fmt.Errorf("Error deleting service instance %s: %s", name, err)
		}
		return nil, serviceInstanceError
	}
	return serviceInstance, nil
}

// WaitForServiceInstanceRunning waits for a service instance to be completely initialized and available.
func (c *ServiceInstanceClient) WaitForServiceInstanceRunning(input *GetServiceInstanceInput, timeoutSeconds time.Duration) (*ServiceInstance, error) {
	var info *ServiceInstance
	var getErr error
	err := c.client.WaitFor("service instance to be ready", timeoutSeconds, func() (bool, error) {
		info, getErr = c.GetServiceInstance(input)
		if getErr != nil {
			return false, getErr
		}
		c.client.DebugLogString(fmt.Sprintf("Service instance name is %v, Service instance info is %+v", info.ServiceName, info))
		switch s := info.Status; s {
		case ServiceInstanceRunning: // Target State
			c.client.DebugLogString("Service Instance Running")
			return true, nil
		case ServiceInstanceFailed:
			c.client.DebugLogString("Service Instance Failed")
			return false, fmt.Errorf("Service Instance Creation Failed")
		case ServiceInstanceInProgress:
			c.client.DebugLogString("Service Instance is being created")
			return false, nil
		default:
			c.client.DebugLogString(fmt.Sprintf("Unknown instance state: %s, waiting", s))
			return false, nil
		}
	})
	return info, err
}

type GetServiceInstanceInput struct {
	// Name of the Java Cloud Service instance.
	// Required.
	Name string `json:"serviceId"`
}

// GetServiceInstance retrieves the SeriveInstance with the given name.
func (c *ServiceInstanceClient) GetServiceInstance(getInput *GetServiceInstanceInput) (*ServiceInstance, error) {
	var serviceInstance ServiceInstance
	if err := c.getResource(getInput.Name, &serviceInstance); err != nil {
		return nil, err
	}

	return &serviceInstance, nil
}

type DeleteServiceInstanceInput struct {
	// Name of the Java Cloud Service instance.
	// Required.
	Name string `json:"-"`
	// User name for the database administrator.
	// Required.
	DBAUsername string `json:"dbaName"`
	// The database administrator password that was specified when the Database Cloud Service database deployment
	// was created or the password for the database administrator.
	// Required.
	DBAPassword string `json:"dbaPassword"`
	// Flag that specifies whether you want to force the removal of the service instance even if the database
	// instance cannot be reached to delete the database schemas. If set to true, you may need to delete the associated
	// database schemas manually on the database instance if they are not deleted as part of the service instance
	// delete operation.
	// The default value is false.
	// Optional.
	ForceDelete bool `json:"forceDelete,omitempty"`
	// Flag that specifies whether you want to back up the service instance or skip backing up the instance before deleting it.
	// The default value is true (that is, skip backing up).
	// Optional.
	SkipBackupOnTerminate bool `json:"skipBackupOnTerminate,omitempty"`
}

func (c *ServiceInstanceClient) DeleteServiceInstance(deleteInput *DeleteServiceInstanceInput) error {
	if c.Timeout == 0 {
		c.Timeout = WaitForServiceInstanceDeleteTimeout
	}

	deleteErr := c.deleteInstanceResource(deleteInput.Name, deleteInput)
	if deleteErr != nil {
		return deleteErr
	}

	// Call wait for instance deleted now, as deleting the instance is an eventually consistent operation
	getInput := &GetServiceInstanceInput{
		Name: deleteInput.Name,
	}

	// Wait for instance to be deleted
	return c.WaitForServiceInstanceDeleted(getInput, c.Timeout)
}

// WaitForServiceInstanceDeleted waits for a service instance to be fully deleted.
func (c *ServiceInstanceClient) WaitForServiceInstanceDeleted(input *GetServiceInstanceInput, timeoutSeconds time.Duration) error {
	return c.client.WaitFor("service instance to be deleted", timeoutSeconds, func() (bool, error) {
		info, err := c.GetServiceInstance(input)
		if err != nil {
			if client.WasNotFoundError(err) {
				// Service Instance could not be found, thus deleted
				return true, nil
			}
			// Some other error occurred trying to get instance, exit
			return false, err
		}
		switch s := info.Status; s {
		case ServiceInstanceTerminating:
			c.client.DebugLogString("Service Instance terminating")
			return false, nil
		default:
			c.client.DebugLogString(fmt.Sprintf("Unknown instance state: %s, waiting", s))
			return false, nil
		}
	})
}
