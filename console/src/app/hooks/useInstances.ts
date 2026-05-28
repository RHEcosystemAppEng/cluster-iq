import { useQuery } from '@tanstack/react-query';
import { api, InstanceResponseApi } from '@api';

export function useInstances() {
  return useQuery<InstanceResponseApi[]>({
    queryKey: ['instances'],
    queryFn: async ({ signal }) => {
      const { data } = await api.instances.instancesList({ page: 1, page_size: 100000 }, { signal });
      return data.items || [];
    },
  });
}
