import { Breadcrumb, BreadcrumbItem, Flex, FlexItem, PageSection, Title } from '@patternfly/react-core';
import { AccountsHeaderProps } from './types';
import { ResourceLabel } from '@app/utils/renderUtils';
import { AccountDetailsDropdown } from './AccountDetailsDropdown';
import React from 'react';
import { Link } from 'react-router-dom';

export const AccountsHeader: React.FunctionComponent<AccountsHeaderProps> = ({ accountName, accountId }) => {
  return (
    <PageSection hasBodyWrapper={false}>
      <Breadcrumb>
        <BreadcrumbItem>
          <Link to="/">Home</Link>
        </BreadcrumbItem>
        <BreadcrumbItem>
          <Link to="/accounts">Accounts</Link>
        </BreadcrumbItem>
        <BreadcrumbItem isActive>{accountName}</BreadcrumbItem>
      </Breadcrumb>
      <Flex
        spaceItems={{ default: 'spaceItemsMd' }}
        alignItems={{ default: 'alignItemsFlexStart' }}
        flexWrap={{ default: 'nowrap' }}
      >
        <FlexItem>
          <Title headingLevel="h1" size="2xl">
            <ResourceLabel label="Account" color="#c9190b" /> {accountName}
          </Title>
        </FlexItem>
        <FlexItem align={{ default: 'alignRight' }}>
          <AccountDetailsDropdown accountId={accountId} accountName={accountName} />
        </FlexItem>
      </Flex>
    </PageSection>
  );
};

export default AccountsHeader;
