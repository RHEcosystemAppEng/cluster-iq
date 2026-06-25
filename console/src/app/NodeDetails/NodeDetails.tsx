import React, { useEffect, useState } from 'react';
import { renderStatusLabel, ResourceLabel } from '@app/utils/renderUtils';
import { parseScanTimestamp, parseNumberToCurrency } from 'src/app/utils/parseFuncs';
import { useParams } from 'react-router-dom';
import {
  PageSection,
  Tabs,
  Tab,
  TabContent,
  TabContentBody,
  TabTitleText,
  Title,
  DescriptionList,
  DescriptionListGroup,
  DescriptionListTerm,
  DescriptionListDescription,
  Flex,
  FlexItem,
  Label,
  LabelGroup,
  Bullseye,
  Spinner,
} from '@patternfly/react-core';
import { api, InstanceResponseApi, TagResponseApi } from '@api';
import { Link } from 'react-router-dom';
import { useDocumentTitle } from '@app/utils/useDocumentTitle';
import { debug } from '@app/utils/debugLogs';

interface LabelGroupOverflowProps {
  labels: Array<TagResponseApi>;
}

const LabelGroupOverflow: React.FunctionComponent<LabelGroupOverflowProps> = ({ labels }) => (
  <LabelGroup>
    {labels.map(label => (
      <Label key={label.key}>
        {label.key}: {label.value}
      </Label>
    ))}
  </LabelGroup>
);

const NodeDetails: React.FunctionComponent = () => {
  const { instanceID } = useParams();
  const [activeTabKey, setActiveTabKey] = React.useState(0);
  const [instanceData, setInstanceData] = useState<InstanceResponseApi | null>(null);
  const [loading, setLoading] = useState(true);
  useDocumentTitle(`${instanceData?.instanceName || instanceID || ''} — ClusterIQ`);
  useEffect(() => {
    if (!instanceID) return;
    let cancelled = false;
    const fetchData = async () => {
      try {
        debug('Fetching Account Clusters ', instanceID);
        const { data: fetchedInstance } = await api.instances.instancesDetail(instanceID);
        if (cancelled) return;
        setInstanceData(fetchedInstance);
        debug('Fetched Account Clusters data:', instanceID);
      } catch (error) {
        if (!cancelled) console.error('Error fetching data:', error);
      } finally {
        if (!cancelled) setLoading(false);
      }
    };

    fetchData();
    return () => {
      cancelled = true;
    };
  }, [instanceID]);

  const handleTabClick = (_event, tabIndex) => {
    setActiveTabKey(tabIndex);
  };

  const detailsTabContent = (
    <React.Fragment>
      {loading ? (
        <Bullseye>
          <Spinner size="xl" />
        </Bullseye>
      ) : (
        <Flex direction={{ default: 'column' }}>
          <FlexItem spacer={{ default: 'spacerLg' }}>
            <Title headingLevel="h2" size="lg" className="pf-v6-u-mt-sm" id="open-tabs-example-tabs-list-details-title">
              Node details
            </Title>
          </FlexItem>

          <FlexItem>
            <DescriptionList
              columnModifier={{ lg: '2Col' }}
              aria-labelledby="open-tabs-example-tabs-list-details-title"
            >
              <DescriptionListGroup>
                <DescriptionListTerm>Name</DescriptionListTerm>
                <DescriptionListDescription>{instanceID}</DescriptionListDescription>
                <DescriptionListTerm>Status</DescriptionListTerm>
                <DescriptionListDescription>{renderStatusLabel(instanceData?.status)}</DescriptionListDescription>
                <DescriptionListTerm>Cluster ID</DescriptionListTerm>
                <DescriptionListDescription>
                  <Link to={`/clusters/${instanceData?.clusterId}`}>{instanceData?.clusterId}</Link>
                </DescriptionListDescription>
                <DescriptionListTerm>Cloud Provider</DescriptionListTerm>
                <DescriptionListDescription>{instanceData?.provider}</DescriptionListDescription>
              </DescriptionListGroup>

              <DescriptionListGroup>
                <DescriptionListTerm>Labels</DescriptionListTerm>
                <LabelGroupOverflow labels={instanceData?.tags || []} />
                <DescriptionListTerm>Last scan</DescriptionListTerm>
                <DescriptionListDescription>
                  {parseScanTimestamp(instanceData?.lastScanTimestamp)}
                </DescriptionListDescription>
                <DescriptionListTerm>Created at</DescriptionListTerm>
                <DescriptionListDescription>
                  {parseScanTimestamp(instanceData?.creationTimestamp)}
                </DescriptionListDescription>
              </DescriptionListGroup>

              <DescriptionListGroup></DescriptionListGroup>

              <DescriptionListGroup>
                <DescriptionListTerm>Total Cost (aprox)</DescriptionListTerm>
                <DescriptionListDescription>
                  {parseNumberToCurrency(instanceData?.totalCost)}
                </DescriptionListDescription>
              </DescriptionListGroup>

              <DescriptionListGroup></DescriptionListGroup>
            </DescriptionList>
          </FlexItem>
        </Flex>
      )}
    </React.Fragment>
  );

  return (
    <React.Fragment>
      {/* Page header */}
      <PageSection hasBodyWrapper={false}>
        <Title headingLevel="h1" size="2xl">
          <ResourceLabel label="Node" color="#4cb140" /> {instanceData?.instanceName || instanceID}
        </Title>
        {/* Page tabs */}
      </PageSection>
      <PageSection hasBodyWrapper={false} type="tabs">
        <Tabs activeKey={activeTabKey} onSelect={handleTabClick} usePageInsets id="open-tabs-example-tabs-list">
          <Tab eventKey={0} title={<TabTitleText>Details</TabTitleText>} tabContentId={`tabContent${0}`} />
        </Tabs>
      </PageSection>
      <PageSection hasBodyWrapper={false} isFilled>
        <TabContent key={0} eventKey={0} id={`tabContent${0}`} activeKey={activeTabKey} hidden={0 !== activeTabKey}>
          <TabContentBody>{detailsTabContent}</TabContentBody>
        </TabContent>
        <TabContent
          key={1}
          eventKey={1}
          id={`tabContent${1}`}
          activeKey={activeTabKey}
          hidden={1 !== activeTabKey}
        ></TabContent>
      </PageSection>
    </React.Fragment>
  );
};

export default NodeDetails;
