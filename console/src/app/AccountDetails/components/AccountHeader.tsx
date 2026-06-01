import { Flex, FlexItem, PageSection, Title } from '@patternfly/react-core';
import { AccountsHeaderProps } from './types';
import { ResourceLabel } from '@app/utils/renderUtils';
import { AccountDetailsDropdown } from './AccountDetailsDropdown';
import React from 'react';

export const AccountsHeader: React.FunctionComponent<AccountsHeaderProps> = ({ accountName, accountId }) => {
  return (
    <PageSection hasBodyWrapper={false}>
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
