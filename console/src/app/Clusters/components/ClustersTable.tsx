import { renderStatusLabel, ResourceBadge, renderProviderIcon } from '@app/utils/renderUtils';
import { parseNumberToCurrency } from '@app/utils/parseFuncs';
import { ThProps, Table, Thead, Tr, Th, Tbody, Td } from '@patternfly/react-table';
import React, { useState, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { ClusterResponseApi } from '@api';
import { ClustersTableProps } from '../types';
import { TableSkeleton } from '@app/components/common/TableSkeleton';
import { TablePagination } from '@app/components/common/TablesPagination';
import { searchItems, filterByStatus, filterByProvider, sortItems } from '@app/utils/tableFilters';
import { EmptyState, EmptyStateVariant, EmptyStateBody, Title } from '@patternfly/react-core';
import { CubesIcon } from '@patternfly/react-icons';
import { useClusters } from '@app/hooks/useClusters';
import { useTablePagination } from '@app/hooks/useTablePagination';

export const ClustersTable: React.FunctionComponent<ClustersTableProps> = ({
  clusterNameSearch,
  accountNameSearch,
  statusFilter,
  providerSelections,
  showTerminated,
}) => {
  const { data: allClusters = [], isLoading } = useClusters();

  const [activeSortIndex, setActiveSortIndex] = useState<number | undefined>(0);
  const [activeSortDirection, setActiveSortDirection] = useState<'asc' | 'desc'>('asc');

  const filtered = useMemo(() => {
    let processed = allClusters;

    if (!showTerminated) {
      processed = processed.filter(cluster => cluster.status !== 'Terminated');
    }

    if (clusterNameSearch) {
      processed = searchItems(processed, clusterNameSearch, ['clusterName']);
    }

    if (accountNameSearch) {
      processed = searchItems(processed, accountNameSearch, ['accountName']);
    }

    processed = filterByStatus(processed, statusFilter);
    processed = filterByProvider(processed, providerSelections);

    if (activeSortIndex !== undefined && activeSortDirection) {
      const sortFields: (keyof ClusterResponseApi)[] = [
        'clusterId',
        'clusterName',
        'status',
        'accountId',
        'provider',
        'region',
        'last15DaysCost',
        'consoleLink',
      ];
      processed = sortItems(processed, sortFields[activeSortIndex], activeSortDirection);
    }

    return processed;
  }, [
    allClusters,
    showTerminated,
    clusterNameSearch,
    accountNameSearch,
    statusFilter,
    providerSelections,
    activeSortIndex,
    activeSortDirection,
  ]);

  const { page, perPage, setPage, setPerPage, paginatedData, totalItems } = useTablePagination({
    data: filtered,
    filterDeps: [clusterNameSearch, accountNameSearch, statusFilter, providerSelections, showTerminated],
  });

  const columnNames = {
    id: 'ID',
    name: 'Name',
    status: 'Status',
    account: 'Account',
    provider: 'Provider',
    region: 'Region',
    cost15d: 'Cost (15d)',
    console: 'Web console',
  };

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

  if (isLoading) {
    return <TableSkeleton columns={8} />;
  }

  if (filtered.length === 0) {
    return (
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
            'No clusters match the current filters.'
          )}
        </EmptyStateBody>
      </EmptyState>
    );
  }

  return (
    <React.Fragment>
      <Table aria-label="Clusters table">
        <Thead>
          <Tr>
            <Th sort={getSortParams(0)}>{columnNames.id}</Th>
            <Th sort={getSortParams(1)}>{columnNames.name}</Th>
            <Th>{columnNames.status}</Th>
            <Th sort={getSortParams(3)}>{columnNames.account}</Th>
            <Th sort={getSortParams(4)}>{columnNames.provider}</Th>
            <Th sort={getSortParams(5)}>{columnNames.region}</Th>
            <Th sort={getSortParams(6)}>{columnNames.cost15d}</Th>
            <Th>{columnNames.console}</Th>
          </Tr>
        </Thead>
        <Tbody>
          {paginatedData.map(cluster => (
            <Tr key={cluster.clusterId}>
              <Td dataLabel={columnNames.id}>
                <ResourceBadge label="C" color="#0066cc" />{' '}
                <Link to={`/clusters/${cluster.clusterId}`}>{cluster.clusterId}</Link>
              </Td>
              <Td dataLabel={columnNames.name}>{cluster.clusterName}</Td>
              <Td dataLabel={columnNames.status}>{renderStatusLabel(cluster.status)}</Td>
              <Td dataLabel={columnNames.account}>
                <Link to={`/accounts/${cluster.accountId}`}>{cluster.accountName}</Link>
              </Td>
              <Td dataLabel={columnNames.provider}>{renderProviderIcon(cluster.provider)}</Td>
              <Td dataLabel={columnNames.region}>{cluster.region}</Td>
              <Td dataLabel={columnNames.cost15d}>{parseNumberToCurrency(cluster.last15DaysCost)}</Td>
              <Td dataLabel={columnNames.console}>
                <a href={cluster.consoleLink} target="_blank" rel="noopener noreferrer">
                  Console
                </a>
              </Td>
            </Tr>
          ))}
        </Tbody>
      </Table>
      <TablePagination
        itemCount={totalItems}
        page={page}
        perPage={perPage}
        onSetPage={setPage}
        onPerPageSelect={setPerPage}
      />
    </React.Fragment>
  );
};

export default ClustersTable;
