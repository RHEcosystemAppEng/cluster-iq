import { PageSection, Title } from '@patternfly/react-core';
import { AccountsHeaderProps } from './types';
import { ResourceLabel } from '@app/utils/renderUtils';
import React from 'react';

export const AccountsHeader: React.FunctionComponent<AccountsHeaderProps> = ({ accountName }) => {
  return (
    <PageSection hasBodyWrapper={false}>
      <Title headingLevel="h1" size="2xl">
        <ResourceLabel label="Account" color="#c9190b" /> {accountName}
      </Title>
    </PageSection>
  );
};

export default AccountsHeader;
