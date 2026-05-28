import { ResourceStatusApi, ProviderApi } from '@api';

export interface ClustersTableToolbarProps {
  clusterNameSearch: string;
  setClusterNameSearch: (value: string) => void;
  accountNameSearch: string;
  setAccountNameSearch: (value: string) => void;
  statusSelection: ResourceStatusApi | null;
  setStatusSelection: (value: ResourceStatusApi | null) => void;
  providerSelections: ProviderApi[] | null;
  setProviderSelections: (value: ProviderApi[] | null) => void;
  showTerminated: boolean;
  setShowTerminated: (value: boolean) => void;
}

export interface ClustersTableProps {
  clusterNameSearch: string;
  accountNameSearch: string;
  statusFilter: string | null;
  providerSelections: ProviderApi[] | null;
  showTerminated: boolean;
}
