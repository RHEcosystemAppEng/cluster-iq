import { ResourceStatusApi, ProviderApi } from '@api';

export interface ServersTableProps {
  searchValue: string;
  statusSelection: string | null;
  providerSelections: ProviderApi[] | null;
  showTerminated: boolean;
}

export interface ServersTableToolbarProps {
  searchValue: string;
  setSearchValue: (value: string) => void;
  statusSelection: ResourceStatusApi | null;
  setStatusSelection: (value: ResourceStatusApi | null) => void;
  providerSelections: ProviderApi[] | null;
  setProviderSelections: (value: ProviderApi[] | null) => void;
  showTerminated: boolean;
  setShowTerminated: (value: boolean) => void;
}
