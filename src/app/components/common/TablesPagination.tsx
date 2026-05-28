import React from 'react';
import { Pagination } from '@patternfly/react-core';

interface PaginationProps {
  itemCount: number;
  page: number;
  perPage: number;
  onSetPage: (page: number) => void;
  onPerPageSelect: (perPage: number) => void;
}

export const TablePagination: React.FC<PaginationProps> = ({
  itemCount,
  page,
  perPage,
  onSetPage,
  onPerPageSelect,
}) => {
  return (
    <Pagination
      itemCount={itemCount}
      page={page}
      perPage={perPage}
      onSetPage={(_evt, newPage) => onSetPage(newPage)}
      onPerPageSelect={(_evt, newPerPage) => onPerPageSelect(newPerPage)}
      isLastFullPageShown
      perPageOptions={[
        { title: '10', value: 10 },
        { title: '20', value: 20 },
        { title: '50', value: 50 },
      ]}
    />
  );
};
