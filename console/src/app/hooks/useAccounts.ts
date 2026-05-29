import { useQuery } from '@tanstack/react-query';
import { api, AccountResponseApi } from '@api';

export function useAccounts() {
  return useQuery<AccountResponseApi[]>({
    queryKey: ['accounts'],
    queryFn: async ({ signal }) => {
      const { data } = await api.accounts.accountsList({ page: 1, page_size: 10000 }, { signal });
      return data.items || [];
    },
  });
}
