import { ProviderApi } from '@api';

export interface AccountsToolbarProps {
  searchValue: string;
  setSearchValue: (value: string) => void;
  providerSelections: ProviderApi[] | null;
  setProviderSelections: (value: ProviderApi[] | null) => void;
}
