/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface AccountListResponseApi {
  count?: number;
  items?: AccountResponseApi[];
}

export interface AccountRequestApi {
  accountId?: string;
  accountName?: string;
  createdAt?: string;
  lastScanTimestamp?: string;
  provider?: ProviderApi;
}

export interface AccountResponseApi {
  last15DaysCost?: number;
  accountId?: string;
  accountName?: string;
  clusterCount?: number;
  createdAt?: string;
  currentMonthSoFarCost?: number;
  lastMonthCost?: number;
  lastScanTimestamp?: string;
  provider?: ProviderApi;
  totalCost?: number;
}

export interface ActionListResponseApi {
  count?: number;
  items?: ActionResponseApi[];
}

export interface ActionRequestApi {
  accountId?: string;
  clusterId?: string;
  cronExpression?: string;
  enabled?: boolean;
  id?: string;
  instances?: string[];
  operation?: string;
  region?: string;
  status?: string;
  time?: string;
  type?: string;
}

export interface ActionResponseApi {
  accountId?: string;
  clusterId?: string;
  cronExpression?: string;
  enabled?: boolean;
  id?: string;
  instances?: string[];
  operation?: string;
  region?: string;
  status?: string;
  time?: string;
  type?: string;
}

export interface ActionTargetApi {
  accountName?: string;
  clusterId?: string;
  instances?: string[];
  region?: string;
}

export interface ClusterEventListResponseApi {
  count?: number;
  items?: ClusterEventResponseApi[];
}

export interface ClusterEventResponseApi {
  action?: string;
  description?: string;
  id?: number;
  resourceId?: string;
  resourceType?: string;
  result?: ResultStatus;
  severity?: string;
  timestamp?: string;
  triggeredBy?: string;
}

export interface ClusterListResponseApi {
  count?: number;
  items?: ClusterResponseApi[];
}

export interface ClusterRequestApi {
  accountId?: string;
  age?: number;
  clusterId?: string;
  clusterName?: string;
  consoleLink?: string;
  createdAt?: string;
  infraId?: string;
  lastScanTimestamp?: string;
  owner?: string;
  provider?: ProviderApi;
  region?: string;
  status?: ResourceStatusApi;
}

export interface ClusterResponseApi {
  last15DaysCost?: number;
  accountId?: string;
  accountName?: string;
  age?: number;
  clusterId?: string;
  clusterName?: string;
  consoleLink?: string;
  createdAt?: string;
  currentMonthSoFarCost?: number;
  infraId?: string;
  instanceCount?: number;
  lastMonthCost?: number;
  lastScanTimestamp?: string;
  owner?: string;
  provider?: ProviderApi;
  region?: string;
  status?: ResourceStatusApi;
  totalCost?: number;
}

export interface ClusterSummaryApi {
  archived?: number;
  running?: number;
  stopped?: number;
}

export interface EventRequestApi {
  action?: string;
  description?: string;
  id?: number;
  resourceId?: string;
  resourceType?: string;
  result?: string;
  severity?: string;
  timestamp?: string;
  triggeredBy?: string;
}

export interface ExpenseListResponseApi {
  count?: number;
  items?: ExpenseResponseApi[];
}

export interface ExpenseRequestApi {
  amount?: number;
  date?: string;
  instanceId?: string;
}

export interface ExpenseResponseApi {
  amount?: number;
  date?: string;
  instanceId?: string;
}

export interface GenericErrorResponseApi {
  message?: string;
}

export interface GenericResponseApi {
  message?: string;
}

export interface HealthCheckResponseApi {
  apiHealth?: boolean;
  dbHealth?: boolean;
}

export interface InstanceListResponseApi {
  count?: number;
  items?: InstanceResponseApi[];
}

export interface InstanceRequestApi {
  age?: number;
  availabilityZone?: string;
  clusterId?: string;
  createdAt?: string;
  instanceId?: string;
  instanceName?: string;
  instanceType?: string;
  lastScanTimestamp?: string;
  owner?: string;
  provider?: ProviderApi;
  status?: ResourceStatusApi;
  tags?: TagRequestApi[];
}

export interface InstanceResponseApi {
  last15DaysCost?: number;
  age?: number;
  availabilityZone?: string;
  clusterId?: string;
  clusterName?: string;
  creationTimestamp?: string;
  currentMonthSoFarCost?: number;
  instanceId?: string;
  instanceName?: string;
  instanceType?: string;
  lastMonthCost?: number;
  lastScanTimestamp?: string;
  owner?: string;
  provider?: ProviderApi;
  status?: ResourceStatusApi;
  tags?: TagResponseApi[];
  totalCost?: number;
}

export interface InstancesSummaryApi {
  archived?: number;
  running?: number;
  stopped?: number;
}

export interface OverviewSummaryApi {
  clusters?: ClusterSummaryApi;
  instances?: InstancesSummaryApi;
  providers?: ProvidersSummaryApi;
  scanner?: ScannerApi;
}

export interface PostResponseApi {
  count?: number;
  status?: string;
}

export enum ProviderApi {
  /** AWSProvider - Amazon Web Services Cloud Provider */
  AWSProvider = 'AWS',
  /** AzureProvider - Microsoft Azure Cloud Provider */
  AzureProvider = 'Azure',
  /** GCPProvider - Google Cloud Platform Cloud Provider */
  GCPProvider = 'GCP',
  /** UnknownProvider - Unknown Provider */
  UnknownProvider = 'UNKNOWN_PROVIDER',
}

export interface ProviderDetailsApi {
  accountCount?: number;
  clusterCount?: number;
}

export interface ProvidersSummaryApi {
  aws?: ProviderDetailsApi;
  azure?: ProviderDetailsApi;
  gcp?: ProviderDetailsApi;
}

export enum ResourceStatusApi {
  /** Running Instance status */
  Running = 'Running',
  /** Stopped Instance status */
  Stopped = 'Stopped',
  /** Terminated Instance status */
  Terminated = 'Terminated',
}

export interface ScannerApi {
  lastScanTimestamp?: string;
}

export interface ScheduledActionApi {
  /** for cron_action */
  cronExpression?: string;
  enabled?: boolean;
  id?: string;
  operation?: string;
  status?: string;
  target?: ActionTargetApi;
  /** for scheduled_action */
  time?: string;
  type?: string;
}

export interface SystemEventListResponseApi {
  count?: number;
  items?: SystemEventResponseApi[];
}

export interface SystemEventResponseApi {
  id?: number;
  action?: string;
  resourceId?: string;
  resourceType?: string;
  timestamp?: string;
  result?: ResultStatus;
  severity?: string;
  triggeredBy?: string;
  description?: string;
  accountId?: string;
  provider?: string;
}

export interface TagRequestApi {
  key?: string;
  value?: string;
}

export interface TagResponseApi {
  instanceId?: string;
  key?: string;
  value?: string;
}
