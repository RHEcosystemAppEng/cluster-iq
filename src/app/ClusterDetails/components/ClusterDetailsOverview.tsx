import { LoadingSpinner } from '@app/components/common/LoadingSpinner';
import { parseNumberToCurrency, parseScanTimestamp } from '@app/utils/parseFuncs';
import { renderStatusLabel } from '@app/utils/renderUtils';
import { ClusterResponseApi, TagResponseApi } from '@api';
import {
  Flex,
  FlexItem,
  Title,
  DescriptionList,
  DescriptionListGroup,
  DescriptionListTerm,
  DescriptionListDescription,
  TabContentBody,
  PageSection,
  Label,
  Divider,
  Tabs,
  Tab,
  TabTitleText,
  TabContent,
} from '@patternfly/react-core';
import React, { useState, useEffect, useMemo } from 'react';
import { useParams } from 'react-router-dom';
import { ClusterDetailsDropdown } from './ClusterDetailsDropdown';
import { ClusterDetailsEvents } from './ClusterDetailsEvents';
import { api } from '@api';
import ClusterDetailsInstances from './ClusterDetailsInstances';
import { LabelGroupOverflow } from '@app/components/common/LabelGroupOverflow';

const ClusterDetailsOverview: React.FunctionComponent = () => {
  const { clusterID } = useParams();
  const [activeTabKey, setActiveTabKey] = React.useState(0);
  const [tags, setTagData] = useState<TagResponseApi[]>([]);
  const [cluster, setClusterData] = useState<ClusterResponseApi | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!clusterID) return;

    const fetchData = async () => {
      try {
        const { data: fetchedCluster } = await api.clusters.clustersDetail(clusterID!);
        setClusterData(fetchedCluster);
        const { data: fetchedTags } = await api.clusters.tagsList(clusterID!);
        setTagData(fetchedTags);
      } catch (error) {
        console.error('Error fetching data:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [clusterID]);

  const filterTagsByKey = key => {
    const result = tags.filter(tag => tag.key == key);
    if (result[0] !== undefined && result[0] != null) {
      return result[0].value;
    }
    return 'unknown';
  };

  const ownerTag = filterTagsByKey('Owner');
  const partnerTag = filterTagsByKey('Partner');

  const handleTabClick = (_, tabIndex) => {
    setActiveTabKey(tabIndex);
  };

  const detailsTabContent = (
    <React.Fragment>
      {loading ? (
        <LoadingSpinner />
      ) : (
        <Flex direction={{ default: 'column' }}>
          <FlexItem spacer={{ default: 'spacerLg' }}>
            <Title headingLevel="h2" size="lg" className="pf-v6-u-mt-sm" id="open-tabs-example-tabs-list-details-title">
              Cluster details
            </Title>
          </FlexItem>

          <FlexItem>
            <DescriptionList
              columnModifier={{ lg: '2Col' }}
              aria-labelledby="open-tabs-example-tabs-list-details-title"
            >
              <DescriptionListGroup name="Basic Info">
                <DescriptionListTerm>Name</DescriptionListTerm>
                <DescriptionListDescription>{cluster?.clusterName}</DescriptionListDescription>
                <DescriptionListTerm>Infrastructure ID</DescriptionListTerm>
                <DescriptionListDescription>{cluster?.infraId}</DescriptionListDescription>
                <DescriptionListTerm>Status</DescriptionListTerm>
                <DescriptionListDescription>{renderStatusLabel(cluster?.status)}</DescriptionListDescription>
              </DescriptionListGroup>

              <DescriptionListGroup name="Cluster links">
                <DescriptionListTerm>Web console</DescriptionListTerm>
                <DescriptionListDescription>
                  <a href={cluster?.consoleLink} target="_blank" rel="noopener noreferrer">
                    Console
                  </a>
                </DescriptionListDescription>
                <DescriptionListTerm>Number of nodes</DescriptionListTerm>
                <DescriptionListDescription>{String(cluster?.instanceCount)}</DescriptionListDescription>
              </DescriptionListGroup>

              <DescriptionListGroup name="Cloud Properties">
                <DescriptionListTerm>Cloud Provider</DescriptionListTerm>
                <DescriptionListDescription>{cluster?.provider}</DescriptionListDescription>
                <DescriptionListTerm>Account</DescriptionListTerm>
                <DescriptionListDescription>{cluster?.accountId || 'unknown'}</DescriptionListDescription>
                <DescriptionListTerm>Region</DescriptionListTerm>
                <DescriptionListDescription>{cluster?.region || 'unknown'}</DescriptionListDescription>
              </DescriptionListGroup>

              <DescriptionListGroup name="Timestamps">
                <DescriptionListTerm>Created at</DescriptionListTerm>
                <DescriptionListDescription>{parseScanTimestamp(cluster?.createdAt)}</DescriptionListDescription>
                <DescriptionListTerm>Last scanned at</DescriptionListTerm>
                <DescriptionListDescription>
                  {parseScanTimestamp(cluster?.lastScanTimestamp)}
                </DescriptionListDescription>
                <DescriptionListTerm>Age (days)</DescriptionListTerm>
                <DescriptionListDescription>{cluster?.age}</DescriptionListDescription>
              </DescriptionListGroup>

              <DescriptionListGroup name="Extra metadata">
                <DescriptionListTerm>Labels</DescriptionListTerm>
                <LabelGroupOverflow labels={tags} />
                <DescriptionListTerm>Partner</DescriptionListTerm>
                <DescriptionListDescription>{partnerTag}</DescriptionListDescription>
                <DescriptionListTerm>Owner</DescriptionListTerm>
                <DescriptionListDescription>{ownerTag}</DescriptionListDescription>
              </DescriptionListGroup>

              <DescriptionListGroup name="Costs">
                <DescriptionListTerm>
                  Cluster Total Cost (Estimated since the cluster is being scanned)
                </DescriptionListTerm>
                <DescriptionListDescription>{parseNumberToCurrency(cluster?.totalCost)}</DescriptionListDescription>
                <DescriptionListTerm>Cluster Total (Current month so far)</DescriptionListTerm>
                <DescriptionListDescription>
                  {parseNumberToCurrency(cluster?.currentMonthSoFarCost)}
                </DescriptionListDescription>
                <DescriptionListTerm>Cluster Total (Last 15 days)</DescriptionListTerm>
                <DescriptionListDescription>
                  {parseNumberToCurrency(cluster?.last15DaysCost)}
                </DescriptionListDescription>
                <DescriptionListTerm>Cluster Total (Last Month)</DescriptionListTerm>
                <DescriptionListDescription>{parseNumberToCurrency(cluster?.lastMonthCost)}</DescriptionListDescription>
              </DescriptionListGroup>
            </DescriptionList>
          </FlexItem>
        </Flex>
      )}
    </React.Fragment>
  );

  const serversTabContent = useMemo(
    () => (
      <TabContentBody>
        <ClusterDetailsInstances />
      </TabContentBody>
    ),
    []
  );

  const eventsTabContent = useMemo(
    () => (
      <TabContentBody>
        <ClusterDetailsEvents />
      </TabContentBody>
    ),
    []
  );

  return (
    <React.Fragment>
      <PageSection hasBodyWrapper={false}>
        <Flex
          spaceItems={{ default: 'spaceItemsMd' }}
          alignItems={{ default: 'alignItemsFlexStart' }}
          flexWrap={{ default: 'nowrap' }}
        >
          <FlexItem>
            <Label color="blue">Cluster</Label>
          </FlexItem>

          <FlexItem>
            <Title headingLevel="h1" size="2xl">
              {clusterID}
            </Title>
          </FlexItem>

          <FlexItem align={{ default: 'alignRight' }}>
            <ClusterDetailsDropdown clusterStatus={cluster?.status ?? null} />
          </FlexItem>
        </Flex>
        {/* Page tabs */}
      </PageSection>
      <PageSection hasBodyWrapper={false} type="tabs">
        <Divider />
        <Tabs activeKey={activeTabKey} onSelect={handleTabClick} usePageInsets id="open-tabs-example-tabs-list">
          <Tab eventKey={0} title={<TabTitleText>Details</TabTitleText>} tabContentId={`tabContent${0}`} />
          <Tab eventKey={1} title={<TabTitleText>Servers</TabTitleText>} tabContentId={`tabContent${1}`} />
          <Tab eventKey={2} title={<TabTitleText>Events</TabTitleText>} tabContentId={`tabContent${2}`} />
        </Tabs>
      </PageSection>
      <PageSection hasBodyWrapper={false} isFilled>
        <TabContent key={0} eventKey={0} id={`tabContent${0}`} activeKey={activeTabKey} hidden={0 !== activeTabKey}>
          <TabContentBody>{detailsTabContent}</TabContentBody>
        </TabContent>
        <TabContent key={1} eventKey={1} id={`tabContent${1}`} activeKey={activeTabKey} hidden={1 !== activeTabKey}>
          <TabContentBody>{serversTabContent}</TabContentBody>
        </TabContent>
        <TabContent key={2} eventKey={2} id={`tabContent${2}`} activeKey={activeTabKey} hidden={2 !== activeTabKey}>
          <TabContentBody>{eventsTabContent}</TabContentBody>
        </TabContent>
      </PageSection>
    </React.Fragment>
  );
};
export default ClusterDetailsOverview;
