import { useQuery, useQueryClient } from '@tanstack/react-query';
import { api, ActionResponseApi } from '@api';

export const SCHEDULE_ACTIONS_QUERY_KEY = ['scheduleActions'] as const;

export function useScheduleActions() {
  return useQuery<ActionResponseApi[]>({
    queryKey: SCHEDULE_ACTIONS_QUERY_KEY,
    queryFn: async ({ signal }) => {
      const { data } = await api.schedule.scheduleList({ page: 1, page_size: 10000 }, { signal });
      return data.items || [];
    },
  });
}

export function useInvalidateScheduleActions() {
  const queryClient = useQueryClient();
  return () => queryClient.invalidateQueries({ queryKey: SCHEDULE_ACTIONS_QUERY_KEY });
}
