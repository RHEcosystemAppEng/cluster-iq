import { AccountResponseApi, ClusterResponseApi } from '@api';
import React from 'react';

export interface AccountsHeaderProps {
  accountName: string;
  label: string;
}

export interface AccountsTabsProps {
  detailsTabContent: React.ReactNode;
  clustersTabContent: React.ReactNode;
}

export interface AccountDetailsContentProps {
  loading: boolean;
  accountData: AccountResponseApi | null;
}

export interface AccountDescriptionListProps {
  account: AccountResponseApi;
}

export interface ClustersTableProps {
  clusters: ClusterResponseApi[];
}
