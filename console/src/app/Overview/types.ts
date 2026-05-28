import React from 'react';
import { ProviderApi } from '@api';

export enum CardLayout {
  SINGLE_ICON = 'icon',
  MULTI_ICON = 'multiIcon',
  WITH_SUBTITLE = 'withSubtitle',
}

export interface CardContentItem {
  icon?: React.ReactNode;
  value: string | number;
  ref?: string;
}

export interface CardTotalCount {
  icon: React.ReactNode;
  value: number;
  label: string;
}

export interface CardDefinition {
  title: string;
  content: CardContentItem[];
  layout: CardLayout;
  customComponent?: React.ReactNode;
  totalCount?: CardTotalCount;
}

export interface DashboardState {
  clustersByStatus: Record<string, number>;
  instancesByStatus: Record<string, number>;
  clustersByProvider: Record<ProviderApi, number>;
  accountsByProvider: Record<ProviderApi, number>;
  instances: number;
  lastScanTimestamp?: string;
}
