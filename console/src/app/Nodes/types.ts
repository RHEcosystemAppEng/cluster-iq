import { ResourceStatusApi, ProviderApi } from '@api';

export interface NodesTableProps {
  searchValue: string;
  statusSelection: string | null;
  providerSelections: ProviderApi[] | null;
  showTerminated: boolean;
}

export interface NodesTableToolbarProps {
  searchValue: string;
  setSearchValue: (value: string) => void;
  statusSelection: ResourceStatusApi | null;
  setStatusSelection: (value: ResourceStatusApi | null) => void;
  providerSelections: ProviderApi[] | null;
  setProviderSelections: (value: ProviderApi[] | null) => void;
  showTerminated: boolean;
  setShowTerminated: (value: boolean) => void;
}
