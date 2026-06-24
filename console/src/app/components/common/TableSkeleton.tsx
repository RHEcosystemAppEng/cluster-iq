import React from 'react';
import { Skeleton } from '@patternfly/react-core';

interface TableSkeletonProps {
  rows?: number;
  columns?: number;
}

export const TableSkeleton: React.FunctionComponent<TableSkeletonProps> = ({ rows = 5, columns = 4 }) => (
  <div style={{ padding: '1rem' }}>
    <div style={{ display: 'flex', gap: '2rem', marginBottom: '1rem' }}>
      {Array.from({ length: columns }, (_, i) => (
        <Skeleton key={`head-${i}`} width="20%" height="1.2rem" />
      ))}
    </div>
    {Array.from({ length: rows }, (_, rowIndex) => (
      <div
        key={`row-${rowIndex}`}
        style={{
          display: 'flex',
          gap: '2rem',
          padding: '0.75rem 0',
          borderBottom: '1px solid var(--pf-t--global--border--color--default)',
        }}
      >
        {Array.from({ length: columns }, (_, colIndex) => (
          <Skeleton key={`cell-${rowIndex}-${colIndex}`} width={colIndex === 0 ? '30%' : '18%'} height="1rem" />
        ))}
      </div>
    ))}
  </div>
);
