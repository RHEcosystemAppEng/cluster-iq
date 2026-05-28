import React from 'react';
import { CardDefinition, CardLayout, DashboardState } from '../types';
import { CLOUD_PROVIDERS, STATUSES, TOTAL_COUNT_ICONS } from '../constants';
import { SystemEventResponseApi } from '@api';
import { ActivityTable } from './ActivityTable';

export const generateCards = (
  state: DashboardState,
  events: SystemEventResponseApi[] = []
): Record<string, CardDefinition[]> => {
  const isValidTimestamp = state.lastScanTimestamp && state.lastScanTimestamp !== '0001-01-01T00:00:00Z';
  const scannerContent = isValidTimestamp
    ? `${new Date(state.lastScanTimestamp!).toLocaleString()}`
    : 'No scan data available';
  const totalClusters = (state.clustersByStatus.running || 0) + (state.clustersByStatus.stopped || 0);
  const totalInstances = state.instances || 0;

  const statusCards = [
    {
      title: 'Clusters',
      content: Object.entries(STATUSES).map(([key, status]) => ({
        icon: status.icon,
        value: state.clustersByStatus[key] || 0,
        ref: status.route,
      })),
      layout: CardLayout.MULTI_ICON,
      totalCount: {
        icon: TOTAL_COUNT_ICONS.clusters,
        value: totalClusters,
        label: 'Total',
      },
    },
    {
      title: 'Instances',
      content: Object.entries(STATUSES).map(([key, status]) => ({
        icon: status.icon,
        value: state.instancesByStatus[key] || 0,
        ref: status.route,
      })),
      layout: CardLayout.MULTI_ICON,
      totalCount: {
        icon: TOTAL_COUNT_ICONS.instances,
        value: totalInstances,
        label: 'Total',
      },
    },
    {
      title: 'Last Scan Timestamp',
      content: [{ value: scannerContent }],
      layout: CardLayout.MULTI_ICON,
    },
  ];

  const providerCards = Object.values(CLOUD_PROVIDERS).map(provider => ({
    title: provider.title,
    content: [
      {
        value: `${state.clustersByProvider[provider.key] ?? 0} Cluster(s)`,
        icon: provider.icon,
        ref: `/clusters?provider=${provider.key}`,
      },
      {
        value: `${state.accountsByProvider[provider.key] ?? 0} Account(s)`,
        icon: provider.providerIcon,
        ref: `/accounts?provider=${provider.key}`,
      },
    ],
    layout: CardLayout.MULTI_ICON,
  }));

  const activityCards = [
    {
      title: 'Recent events',
      content: [], // Empty content since we're using customComponent
      layout: CardLayout.MULTI_ICON,
      customComponent: <ActivityTable events={events} />,
    },
  ];

  return {
    statusCards,
    providerCards,
    activityCards,
  };
};
