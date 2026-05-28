import { renderStatusLabel } from '@app/utils/renderUtils';
import { ThProps, Table, Thead, Tr, Th, Tbody, Td } from '@patternfly/react-table';
import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { ClustersTableProps } from './types';
import { sortItems } from '@app/utils/tableFilters';
import { ClusterResponseApi } from '@api';

const columnNames = {
  id: 'ID',
  name: 'Name',
  status: 'Status',
  cloudProvider: 'Cloud Provider',
  instanceCount: 'Instance Count',
};

export const ClustersTable: React.FunctionComponent<ClustersTableProps> = ({ clusters }) => {
  const [activeSortIndex, setActiveSortIndex] = useState<number | undefined>(1);
  const [activeSortDirection, setActiveSortDirection] = useState<'asc' | 'desc'>('asc');

  let sortedClusters = clusters;
  if (activeSortIndex !== undefined && activeSortDirection) {
    const sortFields: (keyof ClusterResponseApi)[] = [
      'clusterId',
      'clusterName',
      'status',
      'provider',
      'instanceCount',
    ];
    // status column (index 2) is not sortable
    if (activeSortIndex !== 2) {
      sortedClusters = sortItems(clusters, sortFields[activeSortIndex], activeSortDirection);
    }
  }

  const getSortParams = (columnIndex: number): ThProps['sort'] => ({
    sortBy: {
      index: activeSortIndex,
      direction: activeSortDirection,
      defaultDirection: 'asc',
    },
    onSort: (_event, index, direction) => {
      setActiveSortIndex(index);
      setActiveSortDirection(direction);
    },
    columnIndex,
  });

  return (
    <Table aria-label="Simple table">
      <Thead>
        <Tr>
          <Th sort={getSortParams(0)}>{columnNames.id}</Th>
          <Th sort={getSortParams(1)}>{columnNames.name}</Th>
          <Th>{columnNames.status}</Th>
          <Th sort={getSortParams(3)}>{columnNames.cloudProvider}</Th>
          <Th sort={getSortParams(4)}>{columnNames.instanceCount}</Th>
        </Tr>
      </Thead>
      <Tbody>
        {sortedClusters.map(cluster => (
          <Tr key={cluster.clusterId}>
            <Td dataLabel={cluster.clusterName}>
              <Link to={`/clusters/${cluster.clusterId}`}>{cluster.clusterId}</Link>
            </Td>
            <Td>{cluster.clusterName}</Td>
            <Td dataLabel={cluster.status}>{renderStatusLabel(cluster.status)}</Td>
            <Td>{cluster.provider}</Td>
            <Td>{cluster.instanceCount}</Td>
          </Tr>
        ))}
      </Tbody>
    </Table>
  );
};
