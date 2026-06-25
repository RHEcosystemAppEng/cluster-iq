import React from 'react';
import {
  Button,
  EmptyState,
  EmptyStateBody,
  EmptyStateFooter,
  EmptyStateActions,
  PageSection,
} from '@patternfly/react-core';
import { ExclamationTriangleIcon } from '@patternfly/react-icons';

interface ErrorBoundaryState {
  hasError: boolean;
  error: Error | null;
}

export class ErrorBoundary extends React.Component<{ children: React.ReactNode }, ErrorBoundaryState> {
  state: ErrorBoundaryState = { hasError: false, error: null };

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error };
  }

  render() {
    if (!this.state.hasError) return this.props.children;

    return (
      <PageSection hasBodyWrapper={false} isFilled>
        <EmptyState headingLevel="h1" icon={ExclamationTriangleIcon} titleText="Something went wrong" variant="lg">
          <EmptyStateBody>{this.state.error?.message || 'An unexpected error occurred.'}</EmptyStateBody>
          <EmptyStateFooter>
            <EmptyStateActions>
              <Button variant="primary" onClick={() => this.setState({ hasError: false, error: null })}>
                Try again
              </Button>
              <Button variant="link" onClick={() => (window.location.href = '/')}>
                Go to Dashboard
              </Button>
            </EmptyStateActions>
          </EmptyStateFooter>
        </EmptyState>
      </PageSection>
    );
  }
}
