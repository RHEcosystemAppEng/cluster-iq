import { ActionOperations, ResultStatus } from '@app/types/types';
import { ProviderApi } from '@api';

export interface AuditLogsTableToolbarProps {
  searchValue: string;
  setSearchValue: (value: string) => void;
  action: ActionOperations[] | null;
  setAction: (value: ActionOperations[] | null) => void;
  result: ResultStatus[] | null;
  setResult: (value: ResultStatus[] | null) => void;
  providerSelections: ProviderApi[] | null;
  setProviderSelections: (value: ProviderApi[] | null) => void;
  triggered_by: string;
  setTriggeredBy: (value: string) => void;
}

export interface AuditLogsTableProps {
  accountName?: string;
  action?: ActionOperations[];
  provider?: ProviderApi[];
  result?: ResultStatus[];
  triggered_by?: string;
}
