import { useState, useEffect } from 'react';
import { api, OverviewSummaryApi } from '@api';

export const useDashboardData = () => {
  const [inventoryData, setInventoryData] = useState<OverviewSummaryApi | undefined>();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const inventoryOverview = async () => {
      try {
        setLoading(true);
        setError(null);
        const { data } = await api.overview.overviewList();
        setInventoryData(data);
      } catch {
        setError('Failed to fetch inventory data');
        console.error('Failed to fetch inventory data.');
      } finally {
        setLoading(false);
      }
    };
    inventoryOverview();
  }, []);

  return { inventoryData, loading, error };
};
