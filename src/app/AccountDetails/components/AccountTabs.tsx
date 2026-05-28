import { PageSection, Tab, TabContent, TabContentBody, Tabs, TabTitleText } from '@patternfly/react-core';
import React from 'react';
import { AccountsTabsProps } from './types';

export const AccountsTabs: React.FunctionComponent<AccountsTabsProps> = ({ detailsTabContent, clustersTabContent }) => {
  const [activeTabKey, setActiveTabKey] = React.useState(0);
  const handleTabClick = (_event: React.MouseEvent, eventKey: string | number) => {
    setActiveTabKey(eventKey as number);
  };

  return (
    <>
      <PageSection hasBodyWrapper={false} type="tabs">
        <Tabs activeKey={activeTabKey} onSelect={handleTabClick} usePageInsets id="open-tabs-example-tabs-list">
          <Tab eventKey={0} title={<TabTitleText>Details</TabTitleText>} tabContentId={`tabContent${0}`} />
          <Tab eventKey={1} title={<TabTitleText>Clusters</TabTitleText>} tabContentId={`tabContent${1}`} />
        </Tabs>
      </PageSection>
      <PageSection hasBodyWrapper={false} isFilled>
        <TabContent key={0} eventKey={0} id={`tabContent${0}`} activeKey={activeTabKey} hidden={0 !== activeTabKey}>
          <TabContentBody>{detailsTabContent}</TabContentBody>
        </TabContent>
        <TabContent key={1} eventKey={1} id={`tabContent${1}`} activeKey={activeTabKey} hidden={1 !== activeTabKey}>
          <TabContentBody>{clustersTabContent}</TabContentBody>
        </TabContent>
      </PageSection>
    </>
  );
};

export default AccountsTabs;
