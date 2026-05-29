import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import {
  Toolbar,
  ToolbarContent,
  ToolbarItem,
  Switch,
  EmptyState,
  Title,
  EmptyStateBody,
  EmptyStateVariant,
} from '@patternfly/react-core';
import { CubesIcon } from '@patternfly/react-icons';
import { LoadingSpinner } from '@app/components/common/LoadingSpinner';
import { ClustersTable } from './ClustersTable';
import { api, ClusterResponseApi } from '@api';
import { debug } from '@app/utils/debugLogs';

export const AccountClusters: React.FunctionComponent = () => {
  const [clusters, setClusters] = useState<ClusterResponseApi[]>([]);
  const [showTerminated, setShowTerminated] = useState(false);
  const [loading, setLoading] = useState(true);
  const { accountId } = useParams();

  useEffect(() => {
    const fetchData = async () => {
      try {
        debug('Fetching data...');
        const { data } = await api.accounts.clustersList(accountId);
        debug('Fetched Account data:', data);
        setClusters(data.items || []);
      } catch (error) {
        console.error('Error fetching data:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [accountId]);

  // Filter terminated clusters
  const displayClusters = showTerminated ? clusters : clusters.filter(cluster => cluster.status !== 'Terminated');

  const handleToggleChange = (_event: React.FormEvent<HTMLInputElement>, checked: boolean) => {
    setShowTerminated(checked);
  };

  if (loading) {
    return <LoadingSpinner />;
  }

  return (
    <>
      <Toolbar>
        <ToolbarContent>
          <ToolbarItem>
            <Switch
              id="show-terminated-clusters"
              label="Show terminated clusters"
              isChecked={showTerminated}
              onChange={handleToggleChange}
            />
          </ToolbarItem>
        </ToolbarContent>
      </Toolbar>

      {displayClusters.length === 0 ? (
        <EmptyState
          titleText={
            <Title headingLevel="h4" size="md">
              No clusters found
            </Title>
          }
          icon={CubesIcon}
          variant={EmptyStateVariant.sm}
        >
          <EmptyStateBody>
            {!showTerminated ? (
              <>
                There are no active clusters.
                <br />
                Toggle &apos;Show terminated clusters&apos; to view all clusters.
              </>
            ) : (
              'No clusters found for this account.'
            )}
          </EmptyStateBody>
        </EmptyState>
      ) : (
        <ClustersTable clusters={displayClusters} />
      )}
    </>
  );
};
