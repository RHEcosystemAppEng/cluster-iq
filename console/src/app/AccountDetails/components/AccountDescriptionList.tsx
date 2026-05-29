import { parseNumberToCurrency, parseScanTimestamp } from '@app/utils/parseFuncs';
import {
  DescriptionList,
  DescriptionListGroup,
  DescriptionListTerm,
  DescriptionListDescription,
} from '@patternfly/react-core';

import { AccountDescriptionListProps } from './types';
import React from 'react';

export const AccountDescriptionList: React.FunctionComponent<AccountDescriptionListProps> = ({ account }) => {
  return (
    <DescriptionList columnModifier={{ lg: '2Col' }} aria-labelledby="open-tabs-example-tabs-list-details-title">
      <DescriptionListGroup>
        <DescriptionListTerm>Account Name</DescriptionListTerm>
        <DescriptionListDescription>{account.accountName}</DescriptionListDescription>
        <DescriptionListTerm>Account ID</DescriptionListTerm>
        <DescriptionListDescription>{account.accountId}</DescriptionListDescription>
        <DescriptionListTerm>Clusters count</DescriptionListTerm>
        <DescriptionListDescription>{account.clusterCount}</DescriptionListDescription>
        <DescriptionListTerm>Cloud Provider</DescriptionListTerm>
        <DescriptionListDescription>{account.provider}</DescriptionListDescription>
      </DescriptionListGroup>

      <DescriptionListGroup>
        <DescriptionListTerm>Total Cost (Estimated since the account is being scanned) </DescriptionListTerm>
        <DescriptionListDescription>{parseNumberToCurrency(account.totalCost)}</DescriptionListDescription>
        <DescriptionListTerm>Last 15 days cost</DescriptionListTerm>
        <DescriptionListDescription>{parseNumberToCurrency(account.last15DaysCost)}</DescriptionListDescription>
        <DescriptionListTerm>Current month (so far) cost</DescriptionListTerm>
        <DescriptionListDescription>{parseNumberToCurrency(account.currentMonthSoFarCost)}</DescriptionListDescription>
        <DescriptionListTerm>Last month cost</DescriptionListTerm>
        <DescriptionListDescription>{parseNumberToCurrency(account.lastMonthCost)}</DescriptionListDescription>
      </DescriptionListGroup>

      <DescriptionListGroup></DescriptionListGroup>
      <DescriptionListGroup>
        <DescriptionListTerm>Last scanned at</DescriptionListTerm>
        <DescriptionListDescription>{parseScanTimestamp(account.lastScanTimestamp)}</DescriptionListDescription>
        <DescriptionListTerm>Created at (in ClusterIQ)</DescriptionListTerm>
        <DescriptionListDescription>{parseScanTimestamp(account.createdAt)}</DescriptionListDescription>
      </DescriptionListGroup>
    </DescriptionList>
  );
};
