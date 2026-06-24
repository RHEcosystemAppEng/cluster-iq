/* eslint-disable @typescript-eslint/no-explicit-any */
import React from 'react';
import {
  Card,
  CardBody,
  CardTitle,
  Gallery,
  PageSection,
  Content,
  Alert,
  Button,
  EmptyState,
  EmptyStateBody,
  EmptyStateFooter,
  EmptyStateActions,
} from '@patternfly/react-core';
import { CubesIcon, DollarSignIcon, GlobeIcon, UserIcon, HandshakeIcon, HistoryIcon } from '@patternfly/react-icons';
import { LoadingSpinner } from '@app/components/common/LoadingSpinner';
import { generateCards } from './components/CardData';
import { PartnerDonutChart } from './components/PartnerDonutChart';
import { TopMetricCard } from './components/TopMetricCard';
import { ProviderApi, TopItemApi } from '@api';
import { renderContent } from './utils/cardRendererUtils.tsx';
import { useDashboardData } from './hooks/useDashboardData';
import { useEventsData } from './hooks/useEventsData';
import { useDocumentTitle } from '@app/utils/useDocumentTitle';
import { DashboardState } from './types';
import './Overview.css';

const AggregateStatusCards: React.FunctionComponent = () => {
  useDocumentTitle('Overview — ClusterIQ');
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
    lastScanTimestamp: inventoryData?.scanner?.lastScanTimestamp,
    topRegions: inventoryData?.topRegions || [],
    topOwners: inventoryData?.topOwners || [],
    clustersByPartner: inventoryData?.clustersByPartner || [],
    costPerAccount: inventoryData?.costPerAccount || [],
  };

  const costAsTopItems: TopItemApi[] = (dashboardState.costPerAccount || []).map(a => ({
    name: a.accountName,
    clusterCount: a.currentMonthCost,
  }));
  const formatCost = (v: number) => `$${v.toFixed(2)}`;

  const cardData = generateCards(dashboardState, events);

  return (
    <React.Fragment>
      <PageSection hasBodyWrapper={false}>
        <Content>
          <Content component="h1">Overview</Content>
        </Content>
      </PageSection>
      <PageSection hasBodyWrapper={false}>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
          {/* Row 1: Summary cards */}
          <Gallery
            hasGutter
            style={
              {
                '--pf-v6-l-gallery--GridTemplateColumns--min': '22%',
              } as any
            }
          >
            {cardData.summaryCards.map((card, cardIndex) => (
              <Card key={cardIndex} component="div" className="pf-v6-u-min-height overview-card">
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

          {/* Row 2: Partner chart + ranked lists */}
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1.5rem' }}>
            <PartnerDonutChart data={dashboardState.clustersByPartner} />
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1.5rem' }}>
              <TopMetricCard
                title="Cost per Account"
                items={costAsTopItems}
                formatValue={formatCost}
                icon={<DollarSignIcon />}
              />
              <TopMetricCard title="Top Regions" items={dashboardState.topRegions} icon={<GlobeIcon />} />
              <TopMetricCard title="Top Owners" items={dashboardState.topOwners} icon={<UserIcon />} />
              <TopMetricCard
                title="Top Partners"
                items={dashboardState.clustersByPartner.slice(0, 5)}
                icon={<HandshakeIcon />}
              />
            </div>
          </div>

          {/* Row 4: Recent Events */}
          <Card className="pf-v6-u-min-height overview-card" component="div">
            <CardTitle className="pf-v6-u-text-align-center">
              <HistoryIcon style={{ marginRight: '0.4rem' }} />
              {cardData.activityCards[0].title}
            </CardTitle>
            <CardBody className="pf-v6-u-p-md">
              {eventsLoading ? (
                <LoadingSpinner />
              ) : eventsError ? (
                <Alert variant="danger" title="Unable to load events" isInline>
                  <p>{eventsError}</p>
                  <p>Check the console for more details or try refreshing the page.</p>
                </Alert>
              ) : cardData.activityCards[0].customComponent ? (
                cardData.activityCards[0].customComponent
              ) : (
                renderContent(
                  cardData.activityCards[0].content,
                  cardData.activityCards[0].layout,
                  cardData.activityCards[0].totalCount
                )
              )}
            </CardBody>
          </Card>
        </div>
      </PageSection>
    </React.Fragment>
  );
};

export default AggregateStatusCards;
