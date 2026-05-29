import { renderStatusLabel } from '@app/utils/renderUtils';
import { EmptyState, EmptyStateVariant, EmptyStateBody, Title } from '@patternfly/react-core';
import { ThProps, Table, Thead, Tr, Th, Tbody, Td } from '@patternfly/react-table';
import React, { useState, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { ServersTableProps } from '../types';
import { InstanceResponseApi } from '@api';
import { TablePagination } from '@app/components/common/TablesPagination';
import { searchItems, filterByStatus, filterByProvider, sortItems } from '@app/utils/tableFilters';
import { LoadingSpinner } from '@app/components/common/LoadingSpinner';
import { ServerIcon } from '@patternfly/react-icons';
import { useInstances } from '@app/hooks/useInstances';
import { useTablePagination } from '@app/hooks/useTablePagination';

export const ServersTable: React.FunctionComponent<ServersTableProps> = ({
  searchValue,
  statusSelection,
  providerSelections,
  showTerminated,
}) => {
  const { data: allInstances = [], isLoading } = useInstances();

  const [activeSortIndex, setActiveSortIndex] = useState<number | undefined>(1);
  const [activeSortDirection, setActiveSortDirection] = useState<'asc' | 'desc'>('asc');

  const filtered = useMemo(() => {
    let result = allInstances;

    if (!showTerminated) {
      result = result.filter(instance => instance.status !== 'Terminated');
    }

    result = searchItems(result, searchValue, ['instanceName']);
    result = filterByStatus(result, statusSelection);
    result = filterByProvider(result, providerSelections);

    if (activeSortIndex !== undefined && activeSortDirection) {
      const sortFields: (keyof InstanceResponseApi)[] = [
        'instanceId',
        'instanceName',
        'status',
        'provider',
        'availabilityZone',
        'instanceType',
      ];
      if (activeSortIndex !== 2) {
        result = sortItems(result, sortFields[activeSortIndex], activeSortDirection);
      }
    }

    return result;
  }, [
    allInstances,
    showTerminated,
    searchValue,
    statusSelection,
    providerSelections,
    activeSortIndex,
    activeSortDirection,
  ]);

  const { page, perPage, setPage, setPerPage, paginatedData, totalItems } = useTablePagination({
    data: filtered,
    filterDeps: [searchValue, statusSelection, providerSelections, showTerminated],
  });

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

  const columnNames = {
    id: 'ID',
    name: 'Name',
    status: 'Status',
    provider: 'Provider',
    availabilityZone: 'AZ',
    instanceType: 'Type',
  };

  if (isLoading) {
    return <LoadingSpinner />;
  }

  if (filtered.length === 0) {
    return (
      <EmptyState
        titleText={
          <Title headingLevel="h4" size="md">
            No instances found
          </Title>
        }
        icon={ServerIcon}
        variant={EmptyStateVariant.sm}
      >
        <EmptyStateBody>
          {!showTerminated ? (
            <>
              There are no active instances.
              <br />
              Toggle &apos;Show terminated instances&apos; to view all instances.
            </>
          ) : (
            'No instances found.'
          )}
        </EmptyStateBody>
      </EmptyState>
    );
  }

  return (
    <React.Fragment>
      <Table aria-label="Servers table">
        <Thead>
          <Tr>
            <Th sort={getSortParams(0)}>{columnNames.id}</Th>
            <Th sort={getSortParams(1)}>{columnNames.name}</Th>
            <Th>{columnNames.status}</Th>
            <Th sort={getSortParams(3)}>{columnNames.provider}</Th>
            <Th sort={getSortParams(4)}>{columnNames.availabilityZone}</Th>
            <Th sort={getSortParams(5)}>{columnNames.instanceType}</Th>
          </Tr>
        </Thead>
        <Tbody>
          {paginatedData.map(instance => (
            <Tr key={instance.instanceId}>
              <Td dataLabel={columnNames.id} width={15}>
                <Link to={`/instances/${instance.instanceId}`}>{instance.instanceId}</Link>
              </Td>
              <Td dataLabel={columnNames.name} width={30}>
                {instance.instanceName}
              </Td>
              <Td dataLabel={columnNames.status}>{renderStatusLabel(instance.status)}</Td>
              <Td dataLabel={columnNames.provider}>{instance.provider}</Td>
              <Td dataLabel={columnNames.availabilityZone}>{instance.availabilityZone}</Td>
              <Td dataLabel={columnNames.instanceType}>{instance.instanceType}</Td>
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

export default ServersTable;
