/* eslint-disable @typescript-eslint/no-explicit-any */
import React from 'react';
import {
  Card,
  CardBody,
  CardTitle,
  Gallery,
  Grid,
  GridItem,
  PageSection,
  Content,
  Alert,
  Button,
  EmptyState,
  EmptyStateBody,
  EmptyStateFooter,
  EmptyStateActions,
} from '@patternfly/react-core';
import { CubesIcon } from '@patternfly/react-icons';
import { LoadingSpinner } from '@app/components/common/LoadingSpinner';
import { generateCards } from './components/CardData';
import { ProviderApi } from '@api';
import { renderContent } from './utils/cardRendererUtils.tsx';
import { useDashboardData } from './hooks/useDashboardData';
import { useEventsData } from './hooks/useEventsData';
import { DashboardState } from './types';

const AggregateStatusCards: React.FunctionComponent = () => {
  const { inventoryData, loading, error } = useDashboardData();
  const { events, loading: eventsLoading, error: eventsError } = useEventsData();

  if (loading) {
    return <LoadingSpinner />;
  }

  if (error || !inventoryData) {
    return (
      <PageSection hasBodyWrapper={false}>
        <EmptyState variant="lg" titleText="Unable to load dashboard" headingLevel="h1" icon={CubesIcon}>
          <EmptyStateBody>Dashboard unavailable. Refresh to try again.</EmptyStateBody>
          <EmptyStateFooter>
            <EmptyStateActions>
              <Button variant="primary" onClick={() => window.location.reload()}>
                Refresh page
              </Button>
            </EmptyStateActions>
          </EmptyStateFooter>
        </EmptyState>
      </PageSection>
    );
  }

  const dashboardState: DashboardState = {
    clustersByStatus: {
      running: inventoryData?.clusters?.running || 0,
      stopped: inventoryData?.clusters?.stopped || 0,
      terminated: inventoryData?.clusters?.archived || 0,
    },
    instancesByStatus: {
      running: inventoryData?.instances?.running || 0,
      stopped: inventoryData?.instances?.stopped || 0,
      terminated: inventoryData?.instances?.archived || 0,
    },
    clustersByProvider: {
      [ProviderApi.AWSProvider]: inventoryData.providers?.aws?.clusterCount || 0,
      [ProviderApi.GCPProvider]: inventoryData.providers?.gcp?.clusterCount || 0,
      [ProviderApi.AzureProvider]: inventoryData.providers?.azure?.clusterCount || 0,
      [ProviderApi.UnknownProvider]: 0,
    },
    accountsByProvider: {
      [ProviderApi.AWSProvider]: inventoryData.providers?.aws?.accountCount || 0,
      [ProviderApi.GCPProvider]: inventoryData.providers?.gcp?.accountCount || 0,
      [ProviderApi.AzureProvider]: inventoryData.providers?.azure?.accountCount || 0,
      [ProviderApi.UnknownProvider]: 0,
    },
    instances: (inventoryData?.instances?.running || 0) + (inventoryData?.instances?.stopped || 0),
    lastScanTimestamp: inventoryData?.scanner?.lastScanTimestamp,
  };

  const cardData = generateCards(dashboardState, events);

  return (
    <React.Fragment>
      <PageSection hasBodyWrapper={false}>
        <Content>
          <Content component="h1">Overview</Content>
        </Content>
      </PageSection>
      <PageSection hasBodyWrapper={false}>
        <Grid hasGutter>
          {Object.entries(cardData).map(([groupName, cards], groupIndex) => (
            <GridItem key={groupIndex} span={groupName === 'activityCards' ? 12 : undefined}>
              {groupName === 'activityCards' ? (
                // Full width Activity card with double height
                <Card className="pf-v6-u-min-height" component="div">
                  <CardTitle className="pf-v6-u-text-align-center">{cards[0].title}</CardTitle>
                  <CardBody className="pf-v6-u-p-md">
                    {eventsLoading ? (
                      <LoadingSpinner />
                    ) : eventsError ? (
                      <Alert variant="danger" title="Unable to load events" isInline>
                        <p>{eventsError}</p>
                        <p>Check the console for more details or try refreshing the page.</p>
                      </Alert>
                    ) : cards[0].customComponent ? (
                      cards[0].customComponent
                    ) : (
                      renderContent(cards[0].content, cards[0].layout, cards[0].totalCount)
                    )}
                  </CardBody>
                </Card>
              ) : (
                // Regular cards in Gallery
                <Gallery
                  hasGutter
                  style={
                    {
                      '--pf-v6-l-gallery--GridTemplateColumns--min': '30%',
                    } as any
                  }
                >
                  {cards.map((card, cardIndex) => (
                    <Card key={`${groupIndex}${cardIndex}`} component="div" className="pf-v6-u-min-height">
                      <CardTitle
                        className="pf-v6-u-text-align-center"
                        style={{ textAlign: 'center', justifyContent: 'center' }}
                      >
                        {card.title}
                      </CardTitle>
                      <CardBody>{renderContent(card.content, card.layout, card.totalCount)}</CardBody>
                    </Card>
                  ))}
                </Gallery>
              )}
            </GridItem>
          ))}
        </Grid>
      </PageSection>
    </React.Fragment>
  );
};

export default AggregateStatusCards;
