import React from 'react';
import { Spinner } from '@patternfly/react-core';

export const LoadingSpinner: React.FunctionComponent = () => (
  <div
    style={{
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      height: '100vh',
    }}
  >
    <Spinner size="xl" />
  </div>
);
