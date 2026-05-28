import { LoadingSpinner } from '@app/components/common/LoadingSpinner';
import { renderStatusLabel } from '@app/utils/renderUtils';
import { sortItems } from '@app/utils/tableFilters';
import { api, InstanceResponseApi } from '@api';
import { ThProps, Table, Thead, Tr, Th, Tbody, Td } from '@patternfly/react-table';
import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';

const ClusterDetailsInstances: React.FunctionComponent = () => {
  const { clusterID } = useParams();
  const [data, setData] = useState<InstanceResponseApi[]>([]);
  const [loading, setLoading] = useState(true);

  // Index of the currently active column
  const [activeSortIndex, setActiveSortIndex] = useState<number | undefined>(1);
  const [activeSortDirection, setActiveSortDirection] = useState<'asc' | 'desc'>('asc');

  useEffect(() => {
    const fetchData = async () => {
      try {
        console.log('Fetching data...');
        const { data: fetchedInstancesPerCluster } = await api.clusters.instancesList(clusterID!);
        console.log('Fetched data:', fetchedInstancesPerCluster);
        setData(fetchedInstancesPerCluster);
      } catch (error) {
        console.error('Error fetching data:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [clusterID]);

  if (!clusterID) {
    return <LoadingSpinner />;
  }

  console.log('Rendered with data:', data);

  let sortedData = data;
  if (activeSortIndex !== undefined && activeSortDirection) {
    const sortFields: (keyof InstanceResponseApi)[] = [
      'instanceId',
      'instanceName',
      'instanceType',
      'availabilityZone',
    ];
    sortedData = sortItems(data, sortFields[activeSortIndex], activeSortDirection);
  }

  // set table column properties
  const getSortParams = (columnIndex: number): ThProps['sort'] => ({
    sortBy: {
      index: activeSortIndex,
      direction: activeSortDirection,
      defaultDirection: 'asc', // starting sort direction when first sorting a column. Defaults to 'asc'
    },
    onSort: (_event, index, direction) => {
      setActiveSortIndex(index);
      setActiveSortDirection(direction);
    },
    columnIndex,
  });

  return (
    <React.Fragment>
      {loading ? (
        <LoadingSpinner />
      ) : (
        <Table aria-label="Simple table">
          <Thead>
            <Tr>
              <Th sort={getSortParams(0)}>ID</Th>
              <Th sort={getSortParams(1)}>Name</Th>
              <Th sort={getSortParams(2)}>Type</Th>
              <Th>Status</Th>
              <Th sort={getSortParams(3)}>AvailabilityZone</Th>
            </Tr>
          </Thead>
          <Tbody>
            {sortedData.map(instance => (
              <Tr key={instance.instanceId}>
                <Td dataLabel={instance.instanceId}>
                  <Link to={`/instances/${instance.instanceId}`}>{instance.instanceId}</Link>
                </Td>
                <Td>{instance.instanceName}</Td>
                <Td>{instance.instanceType}</Td>
                <Td dataLabel={instance.status}>{renderStatusLabel(instance.status)}</Td>
                <Td>{instance.availabilityZone}</Td>
              </Tr>
            ))}
          </Tbody>
        </Table>
      )}
    </React.Fragment>
  );
};

export default ClusterDetailsInstances;
