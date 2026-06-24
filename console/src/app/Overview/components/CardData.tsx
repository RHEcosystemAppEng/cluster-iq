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

  const totalAccounts = Object.values(state.accountsByProvider).reduce((sum, count) => sum + count, 0);
  const totalClustersByProvider =
    Object.values(state.clustersByProvider).reduce((sum, count) => sum + count, 0) -
    (state.clustersByStatus.terminated || 0);
  const totalClustersByStatus = (state.clustersByStatus.running || 0) + (state.clustersByStatus.stopped || 0);

  const summaryCards: CardDefinition[] = [
    {
      title: 'Accounts',
      content: Object.values(CLOUD_PROVIDERS).map(provider => ({
        icon: provider.providerIcon,
        value: state.accountsByProvider[provider.key] ?? 0,
        ref: `/accounts?provider=${provider.key}`,
      })),
      layout: CardLayout.MULTI_ICON,
      totalCount: {
        icon: TOTAL_COUNT_ICONS.clusters,
        value: totalAccounts,
        label: 'Total',
      },
    },
    {
      title: 'Clusters by Provider',
      content: Object.values(CLOUD_PROVIDERS).map(provider => ({
        icon: provider.icon,
        value: state.clustersByProvider[provider.key] ?? 0,
        ref: `/clusters?provider=${provider.key}`,
      })),
      layout: CardLayout.MULTI_ICON,
      totalCount: {
        icon: TOTAL_COUNT_ICONS.clusters,
        value: totalClustersByProvider,
        label: 'Total',
      },
    },
    {
      title: 'Clusters by Status',
      content: Object.entries(STATUSES).map(([key, status]) => ({
        icon: status.icon,
        value: state.clustersByStatus[key] || 0,
        ref: status.route,
      })),
      layout: CardLayout.MULTI_ICON,
      totalCount: {
        icon: TOTAL_COUNT_ICONS.clusters,
        value: totalClustersByStatus,
        label: 'Total',
      },
    },
    {
      title: 'Last Scan Timestamp',
      content: [{ value: scannerContent }],
      layout: CardLayout.MULTI_ICON,
    },
  ];

  const activityCards: CardDefinition[] = [
    {
      title: 'Recent events',
      content: [],
      layout: CardLayout.MULTI_ICON,
      customComponent: <ActivityTable events={events} />,
    },
  ];

  return {
    summaryCards,
    activityCards,
  };
};
