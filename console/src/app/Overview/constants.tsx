/* eslint-disable react-refresh/only-export-components */
import React from 'react';
import {
  CheckCircleIcon,
  ErrorCircleOIcon,
  OpenshiftIcon,
  AwsIcon,
  GoogleIcon,
  AzureIcon,
  ArchiveIcon,
  DatabaseIcon,
  RegistryIcon,
} from '@patternfly/react-icons';
import { ResourceStatusApi, ProviderApi } from '@api';

const PATTERNFLY_COLORS = {
  success: 'var(--pf-t--global--color--status--success--default)',
  danger: 'var(--pf-t--global--color--status--danger--default)',
  warning: 'var(--pf-t--global--color--status--warning--default)',
  disabled: 'var(--pf-t--global--text--color--disabled)',
} as const;

const CLUSTER_ICON = <OpenshiftIcon color={PATTERNFLY_COLORS.danger} />;

const PROVIDER_ICONS = {
  [ProviderApi.AWSProvider]: <AwsIcon color={PATTERNFLY_COLORS.danger} />,
  [ProviderApi.GCPProvider]: <GoogleIcon color={PATTERNFLY_COLORS.danger} />,
  [ProviderApi.AzureProvider]: <AzureIcon color={PATTERNFLY_COLORS.danger} />,
} as const;

export const STATUSES = {
  running: {
    key: ResourceStatusApi.Running,
    icon: <CheckCircleIcon color={PATTERNFLY_COLORS.success} />,
    route: '/clusters?status=Running',
  },
  stopped: {
    key: ResourceStatusApi.Stopped,
    icon: <ErrorCircleOIcon color={PATTERNFLY_COLORS.danger} />,
    route: '/clusters?status=Stopped',
  },
  terminated: {
    key: ResourceStatusApi.Terminated,
    icon: <ArchiveIcon color={PATTERNFLY_COLORS.disabled} />,
    route: '/clusters?status=Terminated',
  },
} as const;

export const CLOUD_PROVIDERS = {
  [ProviderApi.AWSProvider]: {
    key: ProviderApi.AWSProvider,
    title: 'AWS Clusters',
    icon: CLUSTER_ICON,
    providerIcon: PROVIDER_ICONS[ProviderApi.AWSProvider],
  },
  [ProviderApi.GCPProvider]: {
    key: ProviderApi.GCPProvider,
    title: 'GCP Clusters',
    icon: CLUSTER_ICON,
    providerIcon: PROVIDER_ICONS[ProviderApi.GCPProvider],
  },
  [ProviderApi.AzureProvider]: {
    key: ProviderApi.AzureProvider,
    title: 'Azure Clusters',
    icon: CLUSTER_ICON,
    providerIcon: PROVIDER_ICONS[ProviderApi.AzureProvider],
  },
} as const;

export const TOTAL_COUNT_ICONS = {
  clusters: <DatabaseIcon color={PATTERNFLY_COLORS.disabled} />,
  instances: <RegistryIcon color={PATTERNFLY_COLORS.disabled} />,
} as const;
