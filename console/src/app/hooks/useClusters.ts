import { useQuery } from '@tanstack/react-query';
import { api, ClusterResponseApi } from '@api';

export function useClusters() {
  return useQuery<ClusterResponseApi[]>({
    queryKey: ['clusters'],
    queryFn: async ({ signal }) => {
      const { data } = await api.clusters.clustersList({ page: 1, page_size: 100000 }, { signal });
      return data.items || [];
    },
  });
}
