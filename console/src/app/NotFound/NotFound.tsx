import React from 'react';
import { Button, EmptyState, EmptyStateBody, EmptyStateVariant, PageSection, Title } from '@patternfly/react-core';
import { useNavigate } from 'react-router-dom';

const NotFound: React.FunctionComponent = () => {
  const navigate = useNavigate();

  return (
    <PageSection hasBodyWrapper={false} isFilled>
      <EmptyState
        titleText={
          <Title headingLevel="h1" size="lg">
            404: Page not found
          </Title>
        }
        variant={EmptyStateVariant.full}
      >
        <EmptyStateBody>The page you are looking for does not exist.</EmptyStateBody>
        <Button variant="primary" onClick={() => navigate('/')}>
          Go to Dashboard
        </Button>
      </EmptyState>
    </PageSection>
  );
};

export default NotFound;
