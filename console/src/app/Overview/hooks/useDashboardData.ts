import { useQuery } from '@tanstack/react-query';
import { api, OverviewSummaryApi } from '@api';

export const useDashboardData = () => {
  const { data, isLoading, error } = useQuery<OverviewSummaryApi>({
    queryKey: ['overview'],
    queryFn: async ({ signal }) => {
      const { data } = await api.overview.overviewList({ signal });
      return data;
    },
    refetchInterval: 5_000,
  });

  return {
    inventoryData: data,
    loading: isLoading,
    error: error ? 'Failed to fetch inventory data' : null,
  };
};
