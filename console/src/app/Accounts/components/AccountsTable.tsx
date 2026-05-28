import { ThProps, Table, Thead, Tr, Th, Tbody, Td } from '@patternfly/react-table';
import React, { useState, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { AccountResponseApi, ProviderApi } from '@api';
import { LoadingSpinner } from '@app/components/common/LoadingSpinner';
import { TablePagination } from '@app/components/common/TablesPagination';
import { searchItems, filterByProvider, sortItems } from '@app/utils/tableFilters';
import { useAccounts } from '@app/hooks/useAccounts';
import { useTablePagination } from '@app/hooks/useTablePagination';

export const AccountsTable: React.FunctionComponent<{
  searchValue: string;
  providerSelections: ProviderApi[] | null;
}> = ({ searchValue, providerSelections }) => {
  const { data: allAccounts = [], isLoading } = useAccounts();

  const [activeSortIndex, setActiveSortIndex] = useState<number | undefined>(0);
  const [activeSortDirection, setActiveSortDirection] = useState<'asc' | 'desc'>('asc');

  const filtered = useMemo(() => {
    let result = allAccounts;
    result = searchItems(result, searchValue, ['accountName']);
    result = filterByProvider(result, providerSelections);

    if (activeSortIndex !== undefined && activeSortDirection) {
      const sortFields: (keyof AccountResponseApi)[] = ['accountName', 'provider', 'clusterCount'];
      result = sortItems(result, sortFields[activeSortIndex], activeSortDirection);
    }

    return result;
  }, [allAccounts, searchValue, providerSelections, activeSortIndex, activeSortDirection]);

  const { page, perPage, setPage, setPerPage, paginatedData, totalItems } = useTablePagination({
    data: filtered,
    initialPerPage: 20,
    filterDeps: [searchValue, providerSelections],
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
    name: 'Name',
    cloudProvider: 'Cloud Provider',
    clusterCount: 'Cluster Count',
  };

  return (
    <>
      {isLoading ? (
        <LoadingSpinner />
      ) : (
        <Table aria-label="Accounts table">
          <Thead>
            <Tr>
              <Th sort={getSortParams(0)}>{columnNames.name}</Th>
              <Th sort={getSortParams(1)}>{columnNames.cloudProvider}</Th>
              <Th sort={getSortParams(2)}>{columnNames.clusterCount}</Th>
            </Tr>
          </Thead>
          <Tbody>
            {paginatedData.map(account => (
              <Tr key={account.accountId}>
                <Td dataLabel={columnNames.name}>
                  <Link to={`/accounts/${account.accountId}`}>{account.accountName}</Link>
                </Td>
                <Td dataLabel={columnNames.cloudProvider}>{account.provider}</Td>
                <Td dataLabel={columnNames.clusterCount}>{account.clusterCount}</Td>
              </Tr>
            ))}
          </Tbody>
        </Table>
      )}
      <TablePagination
        itemCount={totalItems}
        page={page}
        perPage={perPage}
        onSetPage={setPage}
        onPerPageSelect={setPerPage}
      />
    </>
  );
};

export default AccountsTable;
