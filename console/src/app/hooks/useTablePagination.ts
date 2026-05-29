import { useState, useMemo, useEffect } from 'react';

interface UseTablePaginationOptions<T> {
  data: T[];
  initialPage?: number;
  initialPerPage?: number;
  filterDeps?: unknown[];
}

interface UseTablePaginationResult<T> {
  page: number;
  perPage: number;
  setPage: (page: number) => void;
  setPerPage: (perPage: number) => void;
  paginatedData: T[];
  totalItems: number;
}

export function useTablePagination<T>({
  data,
  initialPage = 1,
  initialPerPage = 10,
  filterDeps = [],
}: UseTablePaginationOptions<T>): UseTablePaginationResult<T> {
  const [page, setPage] = useState(initialPage);
  const [perPage, setPerPage] = useState(initialPerPage);

  useEffect(() => {
    setPage(1);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, filterDeps);

  const { paginatedData, effectivePage } = useMemo(() => {
    const maxPage = Math.max(1, Math.ceil(data.length / perPage));
    const safePage = Math.min(page, maxPage);
    const startIndex = (safePage - 1) * perPage;
    const endIndex = startIndex + perPage;

    return {
      paginatedData: data.slice(startIndex, endIndex),
      effectivePage: safePage,
    };
  }, [data, page, perPage]);

  return {
    page: effectivePage,
    perPage,
    setPage,
    setPerPage,
    paginatedData,
    totalItems: data.length,
  };
}
